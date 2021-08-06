// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/globalclockupdater (interfaces: RaftApplier,Logger,Sleeper,Timer)

// Package globalclockupdater is a generated GoMock package.
package globalclockupdater

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	raft "github.com/hashicorp/raft"
)

// MockRaftApplier is a mock of RaftApplier interface.
type MockRaftApplier struct {
	ctrl     *gomock.Controller
	recorder *MockRaftApplierMockRecorder
}

// MockRaftApplierMockRecorder is the mock recorder for MockRaftApplier.
type MockRaftApplierMockRecorder struct {
	mock *MockRaftApplier
}

// NewMockRaftApplier creates a new mock instance.
func NewMockRaftApplier(ctrl *gomock.Controller) *MockRaftApplier {
	mock := &MockRaftApplier{ctrl: ctrl}
	mock.recorder = &MockRaftApplierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRaftApplier) EXPECT() *MockRaftApplierMockRecorder {
	return m.recorder
}

// Apply mocks base method.
func (m *MockRaftApplier) Apply(arg0 []byte, arg1 time.Duration) raft.ApplyFuture {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Apply", arg0, arg1)
	ret0, _ := ret[0].(raft.ApplyFuture)
	return ret0
}

// Apply indicates an expected call of Apply.
func (mr *MockRaftApplierMockRecorder) Apply(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockRaftApplier)(nil).Apply), arg0, arg1)
}

// State mocks base method.
func (m *MockRaftApplier) State() raft.RaftState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(raft.RaftState)
	return ret0
}

// State indicates an expected call of State.
func (mr *MockRaftApplierMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockRaftApplier)(nil).State))
}

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Errorf mocks base method.
func (m *MockLogger) Errorf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockLoggerMockRecorder) Errorf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockLogger)(nil).Errorf), varargs...)
}

// Infof mocks base method.
func (m *MockLogger) Infof(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockLoggerMockRecorder) Infof(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*MockLogger)(nil).Infof), varargs...)
}

// Tracef mocks base method.
func (m *MockLogger) Tracef(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Tracef", varargs...)
}

// Tracef indicates an expected call of Tracef.
func (mr *MockLoggerMockRecorder) Tracef(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tracef", reflect.TypeOf((*MockLogger)(nil).Tracef), varargs...)
}

// Warningf mocks base method.
func (m *MockLogger) Warningf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warningf", varargs...)
}

// Warningf indicates an expected call of Warningf.
func (mr *MockLoggerMockRecorder) Warningf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warningf", reflect.TypeOf((*MockLogger)(nil).Warningf), varargs...)
}

// MockSleeper is a mock of Sleeper interface.
type MockSleeper struct {
	ctrl     *gomock.Controller
	recorder *MockSleeperMockRecorder
}

// MockSleeperMockRecorder is the mock recorder for MockSleeper.
type MockSleeperMockRecorder struct {
	mock *MockSleeper
}

// NewMockSleeper creates a new mock instance.
func NewMockSleeper(ctrl *gomock.Controller) *MockSleeper {
	mock := &MockSleeper{ctrl: ctrl}
	mock.recorder = &MockSleeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSleeper) EXPECT() *MockSleeperMockRecorder {
	return m.recorder
}

// Sleep mocks base method.
func (m *MockSleeper) Sleep(arg0 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Sleep", arg0)
}

// Sleep indicates an expected call of Sleep.
func (mr *MockSleeperMockRecorder) Sleep(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sleep", reflect.TypeOf((*MockSleeper)(nil).Sleep), arg0)
}

// MockTimer is a mock of Timer interface.
type MockTimer struct {
	ctrl     *gomock.Controller
	recorder *MockTimerMockRecorder
}

// MockTimerMockRecorder is the mock recorder for MockTimer.
type MockTimerMockRecorder struct {
	mock *MockTimer
}

// NewMockTimer creates a new mock instance.
func NewMockTimer(ctrl *gomock.Controller) *MockTimer {
	mock := &MockTimer{ctrl: ctrl}
	mock.recorder = &MockTimerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimer) EXPECT() *MockTimerMockRecorder {
	return m.recorder
}

// After mocks base method.
func (m *MockTimer) After(arg0 time.Duration) <-chan time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "After", arg0)
	ret0, _ := ret[0].(<-chan time.Time)
	return ret0
}

// After indicates an expected call of After.
func (mr *MockTimerMockRecorder) After(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "After", reflect.TypeOf((*MockTimer)(nil).After), arg0)
}
