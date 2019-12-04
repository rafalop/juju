// Code generated by MockGen. DO NOT EDIT.
// Source: firewallrules.go

// Package migrations is a generated GoMock package.
package migrations

import (
	gomock "github.com/golang/mock/gomock"
	description "github.com/juju/description"
	firewall "github.com/juju/juju/core/firewall"
	reflect "reflect"
)

// MockMigrationFirewallRule is a mock of MigrationFirewallRule interface
type MockMigrationFirewallRule struct {
	ctrl     *gomock.Controller
	recorder *MockMigrationFirewallRuleMockRecorder
}

// MockMigrationFirewallRuleMockRecorder is the mock recorder for MockMigrationFirewallRule
type MockMigrationFirewallRuleMockRecorder struct {
	mock *MockMigrationFirewallRule
}

// NewMockMigrationFirewallRule creates a new mock instance
func NewMockMigrationFirewallRule(ctrl *gomock.Controller) *MockMigrationFirewallRule {
	mock := &MockMigrationFirewallRule{ctrl: ctrl}
	mock.recorder = &MockMigrationFirewallRuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMigrationFirewallRule) EXPECT() *MockMigrationFirewallRuleMockRecorder {
	return m.recorder
}

// ID mocks base method
func (m *MockMigrationFirewallRule) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockMigrationFirewallRuleMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockMigrationFirewallRule)(nil).ID))
}

// WellKnownService mocks base method
func (m *MockMigrationFirewallRule) WellKnownService() firewall.WellKnownServiceType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WellKnownService")
	ret0, _ := ret[0].(firewall.WellKnownServiceType)
	return ret0
}

// WellKnownService indicates an expected call of WellKnownService
func (mr *MockMigrationFirewallRuleMockRecorder) WellKnownService() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WellKnownService", reflect.TypeOf((*MockMigrationFirewallRule)(nil).WellKnownService))
}

// WhitelistCIDRs mocks base method
func (m *MockMigrationFirewallRule) WhitelistCIDRs() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WhitelistCIDRs")
	ret0, _ := ret[0].([]string)
	return ret0
}

// WhitelistCIDRs indicates an expected call of WhitelistCIDRs
func (mr *MockMigrationFirewallRuleMockRecorder) WhitelistCIDRs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WhitelistCIDRs", reflect.TypeOf((*MockMigrationFirewallRule)(nil).WhitelistCIDRs))
}

// MockFirewallRuleSource is a mock of FirewallRuleSource interface
type MockFirewallRuleSource struct {
	ctrl     *gomock.Controller
	recorder *MockFirewallRuleSourceMockRecorder
}

// MockFirewallRuleSourceMockRecorder is the mock recorder for MockFirewallRuleSource
type MockFirewallRuleSourceMockRecorder struct {
	mock *MockFirewallRuleSource
}

// NewMockFirewallRuleSource creates a new mock instance
func NewMockFirewallRuleSource(ctrl *gomock.Controller) *MockFirewallRuleSource {
	mock := &MockFirewallRuleSource{ctrl: ctrl}
	mock.recorder = &MockFirewallRuleSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFirewallRuleSource) EXPECT() *MockFirewallRuleSourceMockRecorder {
	return m.recorder
}

// AllFirewallRules mocks base method
func (m *MockFirewallRuleSource) AllFirewallRules() ([]MigrationFirewallRule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllFirewallRules")
	ret0, _ := ret[0].([]MigrationFirewallRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllFirewallRules indicates an expected call of AllFirewallRules
func (mr *MockFirewallRuleSourceMockRecorder) AllFirewallRules() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllFirewallRules", reflect.TypeOf((*MockFirewallRuleSource)(nil).AllFirewallRules))
}

// MockFirewallRulesModel is a mock of FirewallRulesModel interface
type MockFirewallRulesModel struct {
	ctrl     *gomock.Controller
	recorder *MockFirewallRulesModelMockRecorder
}

// MockFirewallRulesModelMockRecorder is the mock recorder for MockFirewallRulesModel
type MockFirewallRulesModelMockRecorder struct {
	mock *MockFirewallRulesModel
}

// NewMockFirewallRulesModel creates a new mock instance
func NewMockFirewallRulesModel(ctrl *gomock.Controller) *MockFirewallRulesModel {
	mock := &MockFirewallRulesModel{ctrl: ctrl}
	mock.recorder = &MockFirewallRulesModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFirewallRulesModel) EXPECT() *MockFirewallRulesModelMockRecorder {
	return m.recorder
}

// AddFirewallRule mocks base method
func (m *MockFirewallRulesModel) AddFirewallRule(args description.FirewallRuleArgs) description.FirewallRule {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFirewallRule", args)
	ret0, _ := ret[0].(description.FirewallRule)
	return ret0
}

// AddFirewallRule indicates an expected call of AddFirewallRule
func (mr *MockFirewallRulesModelMockRecorder) AddFirewallRule(args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFirewallRule", reflect.TypeOf((*MockFirewallRulesModel)(nil).AddFirewallRule), args)
}
