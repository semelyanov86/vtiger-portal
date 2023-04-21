// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/semelyanov86/vtiger-portal/internal/domain"
	vtiger "github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

// MockUsers is a mock of Users interface.
type MockUsers struct {
	ctrl     *gomock.Controller
	recorder *MockUsersMockRecorder
}

// MockUsersMockRecorder is the mock recorder for MockUsers.
type MockUsersMockRecorder struct {
	mock *MockUsers
}

// NewMockUsers creates a new mock instance.
func NewMockUsers(ctrl *gomock.Controller) *MockUsers {
	mock := &MockUsers{ctrl: ctrl}
	mock.recorder = &MockUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsers) EXPECT() *MockUsersMockRecorder {
	return m.recorder
}

// GetByEmail mocks base method.
func (m *MockUsers) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", ctx, email)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockUsersMockRecorder) GetByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUsers)(nil).GetByEmail), ctx, email)
}

// GetById mocks base method.
func (m *MockUsers) GetById(ctx context.Context, id int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUsersMockRecorder) GetById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUsers)(nil).GetById), ctx, id)
}

// GetForToken mocks base method.
func (m *MockUsers) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForToken", ctx, tokenScope, tokenPlaintext)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForToken indicates an expected call of GetForToken.
func (mr *MockUsersMockRecorder) GetForToken(ctx, tokenScope, tokenPlaintext interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForToken", reflect.TypeOf((*MockUsers)(nil).GetForToken), ctx, tokenScope, tokenPlaintext)
}

// Insert mocks base method.
func (m *MockUsers) Insert(ctx context.Context, user *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockUsersMockRecorder) Insert(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUsers)(nil).Insert), ctx, user)
}

// Update mocks base method.
func (m *MockUsers) Update(ctx context.Context, user *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUsersMockRecorder) Update(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUsers)(nil).Update), ctx, user)
}

// MockUsersCrm is a mock of UsersCrm interface.
type MockUsersCrm struct {
	ctrl     *gomock.Controller
	recorder *MockUsersCrmMockRecorder
}

// MockUsersCrmMockRecorder is the mock recorder for MockUsersCrm.
type MockUsersCrmMockRecorder struct {
	mock *MockUsersCrm
}

// NewMockUsersCrm creates a new mock instance.
func NewMockUsersCrm(ctrl *gomock.Controller) *MockUsersCrm {
	mock := &MockUsersCrm{ctrl: ctrl}
	mock.recorder = &MockUsersCrmMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsersCrm) EXPECT() *MockUsersCrmMockRecorder {
	return m.recorder
}

// ClearUserCodeField mocks base method.
func (m *MockUsersCrm) ClearUserCodeField(ctx context.Context, id string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearUserCodeField", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClearUserCodeField indicates an expected call of ClearUserCodeField.
func (mr *MockUsersCrmMockRecorder) ClearUserCodeField(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearUserCodeField", reflect.TypeOf((*MockUsersCrm)(nil).ClearUserCodeField), ctx, id)
}

// FindByEmail mocks base method.
func (m *MockUsersCrm) FindByEmail(ctx context.Context, email string) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUsersCrmMockRecorder) FindByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUsersCrm)(nil).FindByEmail), ctx, email)
}

// RetrieveById mocks base method.
func (m *MockUsersCrm) RetrieveById(ctx context.Context, id string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockUsersCrmMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockUsersCrm)(nil).RetrieveById), ctx, id)
}

// MockTokens is a mock of Tokens interface.
type MockTokens struct {
	ctrl     *gomock.Controller
	recorder *MockTokensMockRecorder
}

// MockTokensMockRecorder is the mock recorder for MockTokens.
type MockTokensMockRecorder struct {
	mock *MockTokens
}

// NewMockTokens creates a new mock instance.
func NewMockTokens(ctrl *gomock.Controller) *MockTokens {
	mock := &MockTokens{ctrl: ctrl}
	mock.recorder = &MockTokensMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokens) EXPECT() *MockTokensMockRecorder {
	return m.recorder
}

// DeleteAllForUser mocks base method.
func (m *MockTokens) DeleteAllForUser(ctx context.Context, scope string, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllForUser", ctx, scope, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllForUser indicates an expected call of DeleteAllForUser.
func (mr *MockTokensMockRecorder) DeleteAllForUser(ctx, scope, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllForUser", reflect.TypeOf((*MockTokens)(nil).DeleteAllForUser), ctx, scope, userId)
}

// Insert mocks base method.
func (m *MockTokens) Insert(ctx context.Context, token *domain.Token) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockTokensMockRecorder) Insert(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockTokens)(nil).Insert), ctx, token)
}

// New mocks base method.
func (m *MockTokens) New(ctx context.Context, userId int64, ttl time.Duration, scope string) (*domain.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", ctx, userId, ttl, scope)
	ret0, _ := ret[0].(*domain.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New.
func (mr *MockTokensMockRecorder) New(ctx, userId, ttl, scope interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockTokens)(nil).New), ctx, userId, ttl, scope)
}

// MockManagers is a mock of Managers interface.
type MockManagers struct {
	ctrl     *gomock.Controller
	recorder *MockManagersMockRecorder
}

// MockManagersMockRecorder is the mock recorder for MockManagers.
type MockManagersMockRecorder struct {
	mock *MockManagers
}

// NewMockManagers creates a new mock instance.
func NewMockManagers(ctrl *gomock.Controller) *MockManagers {
	mock := &MockManagers{ctrl: ctrl}
	mock.recorder = &MockManagersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManagers) EXPECT() *MockManagersMockRecorder {
	return m.recorder
}

// RetrieveById mocks base method.
func (m *MockManagers) RetrieveById(ctx context.Context, id string) (domain.Manager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.Manager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockManagersMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockManagers)(nil).RetrieveById), ctx, id)
}

// MockModules is a mock of Modules interface.
type MockModules struct {
	ctrl     *gomock.Controller
	recorder *MockModulesMockRecorder
}

// MockModulesMockRecorder is the mock recorder for MockModules.
type MockModulesMockRecorder struct {
	mock *MockModules
}

// NewMockModules creates a new mock instance.
func NewMockModules(ctrl *gomock.Controller) *MockModules {
	mock := &MockModules{ctrl: ctrl}
	mock.recorder = &MockModulesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModules) EXPECT() *MockModulesMockRecorder {
	return m.recorder
}

// GetModuleInfo mocks base method.
func (m *MockModules) GetModuleInfo(ctx context.Context, module string) (vtiger.Module, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleInfo", ctx, module)
	ret0, _ := ret[0].(vtiger.Module)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetModuleInfo indicates an expected call of GetModuleInfo.
func (mr *MockModulesMockRecorder) GetModuleInfo(ctx, module interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleInfo", reflect.TypeOf((*MockModules)(nil).GetModuleInfo), ctx, module)
}
