// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	ctx "context"
	"regexp"

	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cloud"
	"github.com/juju/juju/database/testing"
	dbcloud "github.com/juju/juju/domain/cloud/state"
	schematesting "github.com/juju/juju/domain/schema/testing"
)

type credentialSuite struct {
	schematesting.ControllerSuite
}

var _ = gc.Suite(&credentialSuite{})

func (s *credentialSuite) SetUpTest(c *gc.C) {
	s.ControllerSuite.SetUpTest(c)

	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))
	s.addCloud(c, st, cloud.Cloud{
		Name:      "stratus",
		Type:      "ec2",
		AuthTypes: cloud.AuthTypes{cloud.AccessKeyAuthType, cloud.UserPassAuthType},
	})
}

func (s *credentialSuite) addCloud(c *gc.C, st *State, cloud cloud.Cloud) string {
	cloudSt := dbcloud.NewState(testing.TxnRunnerFactory(s.TxnRunner()))
	ctx := ctx.Background()
	err := cloudSt.UpsertCloud(ctx, cloud)
	c.Assert(err, jc.ErrorIsNil)

	db := s.DB()
	row := db.QueryRow("SELECT uuid FROM cloud WHERE name = ?", cloud.Name)
	c.Assert(row.Err(), jc.ErrorIsNil)

	var dbCloud dbcloud.Cloud
	err = row.Scan(&dbCloud.ID)
	c.Assert(err, jc.ErrorIsNil)
	return dbCloud.ID
}

func (s *credentialSuite) TestUpdateCloudCredentialNew(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	cred := cloud.NewNamedCredential("foobar", cloud.AccessKeyAuthType, map[string]string{
		"foo": "foo val",
		"bar": "bar val",
	}, true)
	ctx := ctx.Background()
	err := st.UpsertCloudCredential(ctx, "foobar", "stratus", "bob", cred)
	c.Assert(err, jc.ErrorIsNil)

	out, err := st.CloudCredential(ctx, "foobar", "stratus", "bob")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(out, jc.DeepEquals, cred)
}

func (s *credentialSuite) TestCreateInvalidCredential(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	cred := cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
		"foo": "foo val",
		"bar": "bar val",
	})
	// Setting of these properties should have no effect when creating a new credential.
	cred.Invalid = true
	cred.InvalidReason = "because am testing you"
	ctx := ctx.Background()
	err := st.UpsertCloudCredential(ctx, "foobar", "stratus", "bob", cred)
	c.Assert(err, gc.ErrorMatches, "adding invalid credential not supported")
}

func (s *credentialSuite) TestUpdateCloudCredentialExisting(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	cred := cloud.NewNamedCredential("foobar", cloud.AccessKeyAuthType, map[string]string{
		"foo": "foo val",
		"bar": "bar val",
	}, false)
	ctx := ctx.Background()
	err := st.UpsertCloudCredential(ctx, "foobar", "stratus", "bob", cred)
	c.Assert(err, jc.ErrorIsNil)

	cred = cloud.NewNamedCredential("foobar", cloud.UserPassAuthType, map[string]string{
		"user":     "bob's nephew",
		"password": "simple",
	}, true)
	err = st.UpsertCloudCredential(ctx, "foobar", "stratus", "bob", cred)
	c.Assert(err, jc.ErrorIsNil)

	out, err := st.CloudCredential(ctx, "foobar", "stratus", "bob")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(out, jc.DeepEquals, cred)
}

func (s *credentialSuite) TestUpdateCloudCredentialInvalidAuthType(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	cred := cloud.NewNamedCredential("foobar", cloud.OAuth2AuthType, map[string]string{
		"foo": "foo val",
		"bar": "bar val",
	}, false)
	ctx := ctx.Background()
	err := st.UpsertCloudCredential(ctx, "foobar", "stratus", "bob", cred)
	c.Assert(err, gc.ErrorMatches, regexp.QuoteMeta(
		`updating credential: validating credential "foobar" owned by "bob" for cloud "stratus": supported auth-types ["access-key" "userpass"], "oauth2" not supported`))
}

func (s *credentialSuite) TestCloudCredentialsEmpty(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	creds, err := st.CloudCredentials(ctx.Background(), "bob", "dummy")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(creds, gc.HasLen, 0)
}

func (s *credentialSuite) TestCloudCredentials(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	cred1 := cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
		"foo": "foo val",
		"bar": "bar val",
	})
	ctx := ctx.Background()
	err := st.UpsertCloudCredential(ctx, "bobcred1", "stratus", "bob", cred1)
	c.Assert(err, jc.ErrorIsNil)

	cred2 := cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
		"baz": "baz val",
		"qux": "qux val",
	})
	err = st.UpsertCloudCredential(ctx, "bobcred2", "stratus", "bob", cred2)
	c.Assert(err, jc.ErrorIsNil)
	err = st.UpsertCloudCredential(ctx, "foobar", "stratus", "mary", cred2)
	c.Assert(err, jc.ErrorIsNil)

	cred1.Label = "bobcred1"
	cred2.Label = "bobcred2"

	for _, userName := range []string{"bob", "bob"} {
		creds, err := st.CloudCredentials(ctx, userName, "stratus")
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(creds, jc.DeepEquals, map[string]cloud.Credential{
			"bobcred1": cred1,
			"bobcred2": cred2,
		})
	}
}

func (s *credentialSuite) assertCredentialInvalidated(c *gc.C, st *State, cloudName, userName, credentialName string) {
	cred := cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
		"foo": "foo val",
		"bar": "bar val",
	})
	ctx := ctx.Background()
	err := st.UpsertCloudCredential(ctx, credentialName, cloudName, userName, cred)
	c.Assert(err, jc.ErrorIsNil)

	cred = cloud.NewCredential(cloud.UserPassAuthType, map[string]string{
		"user":     "bob's nephew",
		"password": "simple",
	})
	cred.Invalid = true
	cred.InvalidReason = "because it is really really invalid"
	err = st.UpsertCloudCredential(ctx, credentialName, cloudName, userName, cred)
	c.Assert(err, jc.ErrorIsNil)

	out, err := st.CloudCredential(ctx, credentialName, cloudName, userName)
	c.Assert(err, jc.ErrorIsNil)
	cred.Label = "foobar"
	c.Assert(out, jc.DeepEquals, cred)
}

func (s *credentialSuite) TestInvalidateCredential(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))
	s.assertCredentialInvalidated(c, st, "stratus", "bob", "foobar")
}

func (s *credentialSuite) assertCredentialMarkedValid(c *gc.C, st *State, cloudName, userName, credentialName string, credential cloud.Credential) {
	ctx := ctx.Background()
	err := st.UpsertCloudCredential(ctx, credentialName, cloudName, userName, credential)
	c.Assert(err, jc.ErrorIsNil)

	out, err := st.CloudCredential(ctx, credentialName, cloudName, userName)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(out.Invalid, jc.IsFalse)
}

func (s *credentialSuite) TestMarkInvalidCredentialAsValidExplicitly(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))
	// This call will ensure that there is an invalid credential to test with.
	s.assertCredentialInvalidated(c, st, "stratus", "bob", "foobar")

	cred := cloud.NewCredential(cloud.UserPassAuthType, map[string]string{
		"user":     "bob's nephew",
		"password": "simple",
	})
	cred.Invalid = false
	s.assertCredentialMarkedValid(c, st, "stratus", "bob", "foobar", cred)
}

func (s *credentialSuite) TestMarkInvalidCredentialAsValidImplicitly(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))
	// This call will ensure that there is an invalid credential to test with.
	s.assertCredentialInvalidated(c, st, "stratus", "bob", "foobar")

	cred := cloud.NewCredential(cloud.UserPassAuthType, map[string]string{
		"user":     "bob's nephew",
		"password": "simple",
	})
	s.assertCredentialMarkedValid(c, st, "stratus", "bob", "foobar", cred)
}

func (s *credentialSuite) TestRemoveCredentials(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	cred1 := cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
		"foo": "foo val",
		"bar": "bar val",
	})
	ctx := ctx.Background()
	err := st.UpsertCloudCredential(ctx, "bobcred1", "stratus", "bob", cred1)
	c.Assert(err, jc.ErrorIsNil)

	err = st.RemoveCloudCredential(ctx, "bobcred1", "stratus", "bob")
	c.Assert(err, jc.ErrorIsNil)

	_, err = st.CloudCredential(ctx, "bobcred1", "stratus", "bob")
	c.Assert(err, jc.Satisfies, errors.IsNotFound)
}

func (s *credentialSuite) TestAllCloudCredentialsNotFound(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	out, err := st.AllCloudCredentials(ctx.Background(), "bob")
	c.Assert(err, gc.ErrorMatches, "cloud credentials for \"bob\" not found")
	c.Assert(out, gc.IsNil)
}

func (s *credentialSuite) createCloudCredential(c *gc.C, st *State, cloudName, userName, credentialName string) cloud.Credential {
	authType := cloud.AccessKeyAuthType
	attributes := map[string]string{
		"foo": "foo val",
		"bar": "bar val",
	}

	s.addCloud(c, st, cloud.Cloud{
		Name:      cloudName,
		Type:      "ec2",
		AuthTypes: cloud.AuthTypes{cloud.AccessKeyAuthType, cloud.UserPassAuthType},
	})

	cred := cloud.NewNamedCredential(credentialName, authType, attributes, false)
	err := st.UpsertCloudCredential(ctx.Background(), credentialName, cloudName, userName, cred)
	c.Assert(err, jc.ErrorIsNil)
	return cred
}

func (s *credentialSuite) TestAllCloudCredentials(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	one := s.createCloudCredential(c, st, "cirrus", "bob", "foobar")
	two := s.createCloudCredential(c, st, "stratus", "bob", "foobar")

	// Added to make sure it is not returned.
	s.createCloudCredential(c, st, "cumulus", "mary", "foobar")

	out, err := st.AllCloudCredentials(ctx.Background(), "bob")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(out, jc.DeepEquals, []cloud.Credential{one, two})
}

func (s *credentialSuite) TestInvalidateCloudCredential(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	one := s.createCloudCredential(c, st, "cirrus", "bob", "foobar")
	c.Assert(one.Invalid, jc.IsFalse)

	ctx := ctx.Background()
	reason := "testing, testing 1,2,3"
	err := st.InvalidateCloudCredential(ctx, "foobar", "cirrus", "bob", reason)
	c.Assert(err, jc.ErrorIsNil)

	updated, err := st.CloudCredential(ctx, "foobar", "cirrus", "bob")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(updated.Invalid, jc.IsTrue)
	c.Assert(updated.InvalidReason, gc.Equals, reason)
}

func (s *credentialSuite) TestInvalidateCloudCredentialNotFound(c *gc.C) {
	st := NewState(testing.TxnRunnerFactory(s.TxnRunner()))

	ctx := ctx.Background()
	err := st.InvalidateCloudCredential(ctx, "foobar", "cirrus", "bob", "reason")
	c.Assert(err, jc.Satisfies, errors.IsNotFound)
}
