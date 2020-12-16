// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasapplication_test

import (
	"time"

	"gopkg.in/yaml.v2"

	"github.com/juju/charm/v8"
	"github.com/juju/clock/testclock"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/facades/agent/caasapplication"
	"github.com/juju/juju/apiserver/params"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/state"
	coretesting "github.com/juju/juju/testing"
)

var _ = gc.Suite(&CAASApplicationSuite{})

type CAASApplicationSuite struct {
	coretesting.BaseSuite

	resources  *common.Resources
	authorizer *apiservertesting.FakeAuthorizer
	facade     *caasapplication.Facade
	st         *mockState
	clock      *testclock.Clock
}

func (s *CAASApplicationSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)

	s.clock = testclock.NewClock(time.Now())

	s.resources = common.NewResources()
	s.AddCleanup(func(_ *gc.C) { s.resources.StopAll() })

	s.authorizer = &apiservertesting.FakeAuthorizer{
		Tag: names.NewApplicationTag("gitlab"),
	}

	s.st = newMockState()

	facade, err := caasapplication.NewFacade(s.resources, s.authorizer, s.st, s.st, s.clock)
	c.Assert(err, jc.ErrorIsNil)
	s.facade = facade
}

func (s *CAASApplicationSuite) TestAddUnit(c *gc.C) {
	args := params.CAASUnitIntroductionArgs{
		PodName: "gitlab-0",
		PodUUID: "gitlab-uuid",
	}

	s.st.app.newUnit = &mockUnit{
		life: state.Alive,
		containerInfo: &mockCloudContainer{
			providerID: "gitlab-0",
			unit:       "gitlab/0",
		},
		updateOp: nil,
	}

	results, err := s.facade.UnitIntroduction(args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results.Error, gc.IsNil)
	c.Assert(results.Result.UnitName, gc.Equals, "gitlab/0")
	c.Assert(results.Result.AgentConf, gc.NotNil)

	s.st.CheckCallNames(c, "Model", "Application", "Unit", "ControllerConfig", "APIHostPortsForAgents")
	s.st.CheckCall(c, 1, "Application", "gitlab")
	s.st.app.CheckCallNames(c, "Life", "Name", "AddUnit")
	c.Assert(s.st.app.Calls()[2].Args[0], gc.DeepEquals, state.AddUnitParams{
		ProviderId: strPtr("gitlab-0"),
		UnitName:   strPtr("gitlab/0"),
	})
}

func (s *CAASApplicationSuite) TestReuseUnitByName(c *gc.C) {
	args := params.CAASUnitIntroductionArgs{
		PodName: "gitlab-0",
		PodUUID: "gitlab-uuid",
	}

	s.st.units["gitlab/0"] = &mockUnit{
		life: state.Alive,
		containerInfo: &mockCloudContainer{
			providerID: "gitlab-0",
			unit:       "gitlab/0",
		},
	}

	results, err := s.facade.UnitIntroduction(args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results.Error, gc.IsNil)
	c.Assert(results.Result.UnitName, gc.Equals, "gitlab/0")
	c.Assert(results.Result.AgentConf, gc.NotNil)

	s.st.CheckCallNames(c, "Model", "Application", "Unit", "ControllerConfig", "APIHostPortsForAgents")
	s.st.CheckCall(c, 1, "Application", "gitlab")
	s.st.app.CheckCallNames(c, "Life", "Name", "UpdateUnits")
	c.Assert(s.st.app.Calls()[2].Args[0], gc.DeepEquals, &state.UpdateUnitsOperation{
		Updates: []*state.UpdateUnitOperation{nil},
	})
}

func (s *CAASApplicationSuite) TestFindByProviderID(c *gc.C) {
	c.Skip("skip for now, because of the TODO in UnitIntroduction facade: hardcoded deploymentType := caas.DeploymentStateful")

	args := params.CAASUnitIntroductionArgs{
		PodName: "gitlab-0",
		PodUUID: "gitlab-uuid",
	}

	s.st.app.charm.meta.Deployment.DeploymentType = charm.DeploymentStateless
	s.st.units["gitlab/0"] = &mockUnit{
		life: state.Alive,
	}
	s.st.units["gitlab/0"].SetErrors(errors.NotFoundf("cloud container"))

	results, err := s.facade.UnitIntroduction(args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results.Error, gc.IsNil)
	c.Assert(results.Result.UnitName, gc.Equals, "gitlab/0")
	c.Assert(results.Result.AgentConf, gc.NotNil)

	s.st.CheckCallNames(c, "Model", "Application", "ControllerConfig", "APIHostPortsForAgents")
	s.st.CheckCall(c, 1, "Application", "gitlab")
	s.st.app.CheckCallNames(c, "Life", "Charm", "AllUnits", "UpdateUnits")
	c.Assert(s.st.app.Calls()[3].Args[0], gc.DeepEquals, &state.UpdateUnitsOperation{
		Updates: []*state.UpdateUnitOperation{nil},
	})
}

func (s *CAASApplicationSuite) TestAgentConf(c *gc.C) {
	args := params.CAASUnitIntroductionArgs{
		PodName: "gitlab-0",
		PodUUID: "gitlab-uuid",
	}

	s.st.app.newUnit = &mockUnit{
		life: state.Alive,
		containerInfo: &mockCloudContainer{
			providerID: "gitlab-0",
			unit:       "gitlab/0",
		},
		updateOp: nil,
	}

	results, err := s.facade.UnitIntroduction(args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results.Error, gc.IsNil)
	c.Assert(results.Result.UnitName, gc.Equals, "gitlab/0")
	c.Assert(results.Result.AgentConf, gc.NotNil)

	conf := map[string]interface{}{}
	err = yaml.Unmarshal(results.Result.AgentConf, conf)
	c.Assert(err, jc.ErrorIsNil)

	check := jc.NewMultiChecker()
	check.AddExpr(`_["cacert"]`, jc.Ignore)
	check.AddExpr(`_["oldpassword"]`, jc.Ignore)
	check.AddExpr(`_["values"]`, jc.Ignore)
	c.Assert(conf, check, map[string]interface{}{
		"tag":               "unit-gitlab-0",
		"datadir":           "/var/lib/juju",
		"transient-datadir": "/var/run/juju",
		"logdir":            "/var/log",
		"metricsspooldir":   "/var/lib/juju/metricspool",
		"upgradedToVersion": "1.9.99",
		"cacert":            "ignore",
		"controller":        "controller-ffffffff-ffff-ffff-ffff-ffffffffffff",
		"model":             "model-ffffffff-ffff-ffff-ffff-ffffffffffff",
		"apiaddresses": []interface{}{
			"10.0.2.1:17070",
			"52.7.1.1:17070",
		},
		"oldpassword": "ignore",
		"values":      nil,
	})
}

func (s *CAASApplicationSuite) TestDyingApplication(c *gc.C) {
	args := params.CAASUnitIntroductionArgs{
		PodName: "gitlab-0",
		PodUUID: "gitlab-uuid",
	}

	s.st.app.life = state.Dying

	results, err := s.facade.UnitIntroduction(args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results.Error, gc.ErrorMatches, `application not provisioned`)
}

func (s *CAASApplicationSuite) TestMissingArgUUID(c *gc.C) {
	args := params.CAASUnitIntroductionArgs{
		PodName: "gitlab-0",
	}

	s.st.app.life = state.Dying

	results, err := s.facade.UnitIntroduction(args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results.Error, gc.ErrorMatches, `pod-uuid not valid`)
}

func (s *CAASApplicationSuite) TestMissingArgName(c *gc.C) {
	args := params.CAASUnitIntroductionArgs{
		PodUUID: "gitlab-uuid",
	}

	s.st.app.life = state.Dying

	results, err := s.facade.UnitIntroduction(args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results.Error, gc.ErrorMatches, `pod-name not valid`)
}

func strPtr(s string) *string {
	return &s
}
