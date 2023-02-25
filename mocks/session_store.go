// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	sessions "github.com/dghubble/sessions"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Destroy mocks base method
func (m *MockStore) Destroy(arg0 http.ResponseWriter, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Destroy", arg0, arg1)
}

// Destroy indicates an expected call of Destroy
func (mr *MockStoreMockRecorder) Destroy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockStore)(nil).Destroy), arg0, arg1)
}

// Get mocks base method
func (m *MockStore) Get(arg0 *http.Request, arg1 string) (*sessions.Session[string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*sessions.Session[string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockStoreMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), arg0, arg1)
}

// New mocks base method
func (m *MockStore) New(arg0 string) *sessions.Session[string] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", arg0)
	ret0, _ := ret[0].(*sessions.Session[string])
	return ret0
}

// New indicates an expected call of New
func (mr *MockStoreMockRecorder) New(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockStore)(nil).New), arg0)
}

// Save mocks base method
func (m *MockStore) Save(arg0 http.ResponseWriter, arg1 *sessions.Session[string]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockStoreMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStore)(nil).Save), arg0, arg1)
}
