// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package secrets

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names/v4"

	"github.com/juju/juju/cloud"
	"github.com/juju/juju/core/leadership"
	"github.com/juju/juju/environs/cloudspec"
	"github.com/juju/juju/secrets/provider"
	"github.com/juju/juju/secrets/provider/juju"
	"github.com/juju/juju/secrets/provider/kubernetes"
	"github.com/juju/juju/state"
)

var logger = loggo.GetLogger("juju.apiserver.common.secrets")

// For testing.
var (
	GetProvider            = provider.Provider
	GetSecretsState        = getSecretsState
	GetSecretBackendsState = getSecretBackendsState
)

func getSecretsState(m Model) state.SecretsStore {
	return state.NewSecrets(m.State())
}

func getSecretBackendsState(m Model) state.SecretBackendsStorage {
	return state.NewSecretBackends(m.State())
}

// BackendConfigGetter is a func used to get secret backend config.
type BackendConfigGetter func() (*provider.ModelBackendConfigInfo, error)

// backendsForModel returns the secret backends is use by the specified model.
// If no backend is configured, the "internal" backend is used for machine models and
// a k8s backend with the same namespace is used for k8s models.
func backendsForModel(model Model) (configs map[string]provider.BackendConfig, activeID string, _ error) {
	cfg, err := model.Config()
	if err != nil {
		return nil, "", errors.Trace(err)
	}
	backendName := cfg.SecretBackend()

	var backendType string
	switch backendName {
	case provider.Auto:
		backendType = juju.BackendType
		if model.Type() == state.ModelTypeCAAS {
			backendType = kubernetes.BackendType
		}
	case provider.Internal:
		backendType = juju.BackendType
	}

	configs = make(map[string]provider.BackendConfig)
	if backendType != "" {
		if backendType == juju.BackendType {
			activeID = model.ControllerUUID()
			configs[activeID] = juju.BuiltInConfig()
		} else {
			spec, err := cloudSpecForModel(model)
			if err != nil {
				return nil, "", errors.Trace(err)
			}
			k8sConfig, err := kubernetes.BuiltInConfig(spec)
			if err != nil {
				return nil, "", errors.Trace(err)
			}
			activeID = model.UUID()
			configs[activeID] = *k8sConfig
		}
	}
	// TODO(secrets) - only use those in use by model
	// For now, we'll return all backends on the controller.
	backendState := GetSecretBackendsState(model)
	backends, err := backendState.ListSecretBackends()
	if err != nil {
		return nil, "", errors.Trace(err)
	}
	for _, b := range backends {
		if b.Name == backendName {
			activeID = b.ID
		}
		configs[b.ID] = provider.BackendConfig{
			BackendType: b.BackendType,
			Config:      b.Config,
		}
	}
	return configs, activeID, nil
}

// AdminBackendConfigInfo returns admin config for secret backends used by the model.
func AdminBackendConfigInfo(model Model) (*provider.ModelBackendConfigInfo, string, error) {
	configs, activeID, err := backendsForModel(model)
	if err != nil {
		return nil, "", errors.Annotate(err, "getting configured secrets providers")
	}
	return &provider.ModelBackendConfigInfo{
		ControllerUUID: model.ControllerUUID(),
		ModelUUID:      model.UUID(),
		ModelName:      model.Name(),
		Configs:        configs,
	}, activeID, nil
}

// BackendConfigInfo returns the config to create a secret backend.
// This is called to provide config to a client like a unit agent which
// needs to access secrets. The authTag is the agent which needs access.
// The client is expected to be restricted to write only those secrets
// owned by the agent, and read only those secrets shared with the agent.
// The result includes config for all relevant backends, including the id
// of the current active backend.
func BackendConfigInfo(model Model, authTag names.Tag, leadershipChecker leadership.Checker) (*provider.ModelBackendConfigInfo, error) {
	configs, activeID, err := AdminBackendConfigInfo(model)
	if err != nil {
		return nil, errors.Annotate(err, "getting configured secrets providers")
	}
	activeCfg := configs.Configs[activeID]
	adminModelCfg := &provider.ModelBackendConfig{
		ControllerUUID: model.ControllerUUID(),
		ModelUUID:      model.UUID(),
		ModelName:      model.Name(),
		BackendConfig:  activeCfg,
	}
	p, err := GetProvider(activeCfg.BackendType)
	if err != nil {
		return nil, errors.Trace(err)
	}
	err = p.Initialise(adminModelCfg)
	if err != nil {
		return nil, errors.Annotate(err, "initialising secrets provider")
	}
	secretsState := GetSecretsState(model)

	// Find secrets owned by the agent
	// (or its app if the agent is a leader).
	ownedFilter := state.SecretsFilter{
		OwnerTags: []names.Tag{authTag},
	}
	// Find secrets shared with the agent.
	// We include secrets shared with the app or just the specified unit.
	readFilter := state.SecretsFilter{
		ConsumerTags: []names.Tag{authTag},
	}
	// Find secrets owned by the application that should be readable for non leader units.
	readAppOwnedFilter := state.SecretsFilter{}
	switch t := authTag.(type) {
	case names.UnitTag:
		appName, _ := names.UnitApplication(t.Id())
		authApp := names.NewApplicationTag(appName)
		token := leadershipChecker.LeadershipCheck(appName, t.Id())
		err := token.Check()
		if err != nil && !leadership.IsNotLeaderError(err) {
			return nil, errors.Trace(err)
		}
		if err == nil {
			// Leader unit owns application level secrets.
			ownedFilter.OwnerTags = append(ownedFilter.OwnerTags, authApp)
		} else {
			// Non leader units can read application level secrets.
			// Find secrets owned by the application.
			readAppOwnedFilter.OwnerTags = append(readAppOwnedFilter.OwnerTags, authApp)
		}
		// Granted secrets can be consumed in application level for all units.
		readFilter.ConsumerTags = append(readFilter.ConsumerTags, authApp)
	case names.ApplicationTag:
	default:
		return nil, errors.NotSupportedf("login as %q", authTag)
	}

	ownedRevisions := map[string]provider.SecretRevisions{}
	if err := getExternalRevisions(secretsState, ownedFilter, ownedRevisions); err != nil {
		return nil, errors.Trace(err)
	}

	readRevisions := map[string]provider.SecretRevisions{}
	if err := getExternalRevisions(secretsState, readFilter, readRevisions); err != nil {
		return nil, errors.Trace(err)
	}

	if len(readAppOwnedFilter.OwnerTags) > 0 {
		if err := getExternalRevisions(secretsState, readAppOwnedFilter, readRevisions); err != nil {
			return nil, errors.Trace(err)
		}
	}

	logger.Debugf("secrets for %v:\nowned: %v\nconsumed:%v", authTag.String(), ownedRevisions, readRevisions)
	clientConfigs := make(map[string]provider.BackendConfig)
	for id, adminCfg := range configs.Configs {
		modelCfg := &provider.ModelBackendConfig{
			ControllerUUID: model.ControllerUUID(),
			ModelUUID:      model.UUID(),
			ModelName:      model.Name(),
			BackendConfig:  adminCfg,
		}
		p, err := GetProvider(adminCfg.BackendType)
		if err != nil {
			return nil, errors.Trace(err)
		}
		cfg, err := p.RestrictedConfig(modelCfg, authTag, ownedRevisions[id], readRevisions[id])
		if err != nil {
			return nil, errors.Trace(err)
		}
		clientConfigs[id] = *cfg
	}
	info := &provider.ModelBackendConfigInfo{
		ControllerUUID: model.ControllerUUID(),
		ModelUUID:      model.UUID(),
		ModelName:      model.Name(),
		ActiveID:       activeID,
		Configs:        clientConfigs,
	}
	return info, nil
}

func getExternalRevisions(backend state.SecretsStore, filter state.SecretsFilter, revisions map[string]provider.SecretRevisions) error {
	secrets, err := backend.ListSecrets(filter)
	if err != nil {
		return errors.Trace(err)
	}
	for _, md := range secrets {
		revs, err := backend.ListSecretRevisions(md.URI)
		if err != nil {
			return errors.Annotatef(err, "cannot get revisions for secret %q", md.URI)
		}
		for _, rev := range revs {
			if rev.ValueRef == nil {
				continue
			}
			revs, ok := revisions[rev.ValueRef.BackendID]
			if !ok {
				revs = provider.SecretRevisions{}
			}
			revs.Add(md.URI, rev.ValueRef.RevisionID)
			revisions[rev.ValueRef.BackendID] = revs
		}
	}
	return nil
}

// BackendForInspect returns the config to create a secret backend client able
// to read any secrets for that model.
// This is called by the show-secret facade for admin users.
func BackendForInspect(model Model, backendID string) (provider.SecretsBackend, error) {
	backendStorage := GetSecretBackendsState(model)
	backend, err := backendStorage.GetSecretBackendByID(backendID)
	if err != nil {
		return nil, errors.Annotatef(err, "getting secrets backend %q", backendID)
	}
	p, err := GetProvider(backend.BackendType)
	if err != nil {
		return nil, errors.Annotatef(err, "getting %q secrets provider for backend %q", backend.BackendType, backendID)
	}
	adminCfg := provider.ModelBackendConfig{
		ControllerUUID: model.ControllerUUID(),
		ModelUUID:      model.UUID(),
		ModelName:      model.Name(),
		BackendConfig: provider.BackendConfig{
			BackendType: backend.BackendType,
			Config:      backend.Config,
		},
	}

	cfg, err := p.RestrictedConfig(&adminCfg, nil, nil, nil)
	if err != nil {
		return nil, errors.Annotate(err, "creating secrets backend config")
	}
	return p.NewBackend(&provider.ModelBackendConfig{
		ControllerUUID: model.ControllerUUID(),
		ModelUUID:      model.UUID(),
		ModelName:      model.Name(),
		BackendConfig:  *cfg,
	})
}

func cloudSpecForModel(m Model) (cloudspec.CloudSpec, error) {
	c, err := m.Cloud()
	if err != nil {
		return cloudspec.CloudSpec{}, errors.Trace(err)
	}
	cred, err := m.CloudCredential()
	if err != nil {
		return cloudspec.CloudSpec{}, errors.Trace(err)
	}
	if cred == nil {
		return cloudspec.CloudSpec{}, errors.NotValidf("cloud credential for %s is empty", m.UUID())
	}
	cloudCredential := cloud.NewCredential(
		cloud.AuthType(cred.AuthType()),
		cred.Attributes(),
	)
	return cloudspec.MakeCloudSpec(c, "", &cloudCredential)
}
