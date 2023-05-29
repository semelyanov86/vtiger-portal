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
	repository "github.com/semelyanov86/vtiger-portal/internal/repository"
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

// FindContactsInAccount mocks base method.
func (m *MockUsersCrm) FindContactsInAccount(ctx context.Context, filter repository.PaginationQueryFilter) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindContactsInAccount", ctx, filter)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindContactsInAccount indicates an expected call of FindContactsInAccount.
func (mr *MockUsersCrmMockRecorder) FindContactsInAccount(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindContactsInAccount", reflect.TypeOf((*MockUsersCrm)(nil).FindContactsInAccount), ctx, filter)
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

// Update mocks base method.
func (m *MockUsersCrm) Update(ctx context.Context, id string, user domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, user)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUsersCrmMockRecorder) Update(ctx, id, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUsersCrm)(nil).Update), ctx, id, user)
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

// MockCompany is a mock of Company interface.
type MockCompany struct {
	ctrl     *gomock.Controller
	recorder *MockCompanyMockRecorder
}

// MockCompanyMockRecorder is the mock recorder for MockCompany.
type MockCompanyMockRecorder struct {
	mock *MockCompany
}

// NewMockCompany creates a new mock instance.
func NewMockCompany(ctrl *gomock.Controller) *MockCompany {
	mock := &MockCompany{ctrl: ctrl}
	mock.recorder = &MockCompanyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompany) EXPECT() *MockCompanyMockRecorder {
	return m.recorder
}

// GetCompanyInfo mocks base method.
func (m *MockCompany) GetCompanyInfo(ctx context.Context) (domain.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyInfo", ctx)
	ret0, _ := ret[0].(domain.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyInfo indicates an expected call of GetCompanyInfo.
func (mr *MockCompanyMockRecorder) GetCompanyInfo(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyInfo", reflect.TypeOf((*MockCompany)(nil).GetCompanyInfo), ctx)
}

// MockHelpDesk is a mock of HelpDesk interface.
type MockHelpDesk struct {
	ctrl     *gomock.Controller
	recorder *MockHelpDeskMockRecorder
}

// MockHelpDeskMockRecorder is the mock recorder for MockHelpDesk.
type MockHelpDeskMockRecorder struct {
	mock *MockHelpDesk
}

// NewMockHelpDesk creates a new mock instance.
func NewMockHelpDesk(ctrl *gomock.Controller) *MockHelpDesk {
	mock := &MockHelpDesk{ctrl: ctrl}
	mock.recorder = &MockHelpDeskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHelpDesk) EXPECT() *MockHelpDeskMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockHelpDesk) Count(ctx context.Context, client string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, client)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockHelpDeskMockRecorder) Count(ctx, client interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockHelpDesk)(nil).Count), ctx, client)
}

// Create mocks base method.
func (m *MockHelpDesk) Create(ctx context.Context, ticket domain.HelpDesk) (domain.HelpDesk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, ticket)
	ret0, _ := ret[0].(domain.HelpDesk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockHelpDeskMockRecorder) Create(ctx, ticket interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockHelpDesk)(nil).Create), ctx, ticket)
}

// GetAll mocks base method.
func (m *MockHelpDesk) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.HelpDesk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, filter)
	ret0, _ := ret[0].([]domain.HelpDesk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockHelpDeskMockRecorder) GetAll(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockHelpDesk)(nil).GetAll), ctx, filter)
}

// RetrieveById mocks base method.
func (m *MockHelpDesk) RetrieveById(ctx context.Context, id string) (domain.HelpDesk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.HelpDesk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockHelpDeskMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockHelpDesk)(nil).RetrieveById), ctx, id)
}

// Revise mocks base method.
func (m *MockHelpDesk) Revise(ctx context.Context, ticket map[string]any) (domain.HelpDesk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Revise", ctx, ticket)
	ret0, _ := ret[0].(domain.HelpDesk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Revise indicates an expected call of Revise.
func (mr *MockHelpDeskMockRecorder) Revise(ctx, ticket interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Revise", reflect.TypeOf((*MockHelpDesk)(nil).Revise), ctx, ticket)
}

// Update mocks base method.
func (m *MockHelpDesk) Update(ctx context.Context, ticket domain.HelpDesk) (domain.HelpDesk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, ticket)
	ret0, _ := ret[0].(domain.HelpDesk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockHelpDeskMockRecorder) Update(ctx, ticket interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockHelpDesk)(nil).Update), ctx, ticket)
}

// MockComment is a mock of Comment interface.
type MockComment struct {
	ctrl     *gomock.Controller
	recorder *MockCommentMockRecorder
}

// MockCommentMockRecorder is the mock recorder for MockComment.
type MockCommentMockRecorder struct {
	mock *MockComment
}

// NewMockComment creates a new mock instance.
func NewMockComment(ctrl *gomock.Controller) *MockComment {
	mock := &MockComment{ctrl: ctrl}
	mock.recorder = &MockCommentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComment) EXPECT() *MockCommentMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockComment) Create(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, comment)
	ret0, _ := ret[0].(domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCommentMockRecorder) Create(ctx, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockComment)(nil).Create), ctx, comment)
}

// RetrieveFromModule mocks base method.
func (m *MockComment) RetrieveFromModule(ctx context.Context, id string) ([]domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveFromModule", ctx, id)
	ret0, _ := ret[0].([]domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveFromModule indicates an expected call of RetrieveFromModule.
func (mr *MockCommentMockRecorder) RetrieveFromModule(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveFromModule", reflect.TypeOf((*MockComment)(nil).RetrieveFromModule), ctx, id)
}

// MockDocument is a mock of Document interface.
type MockDocument struct {
	ctrl     *gomock.Controller
	recorder *MockDocumentMockRecorder
}

// MockDocumentMockRecorder is the mock recorder for MockDocument.
type MockDocumentMockRecorder struct {
	mock *MockDocument
}

// NewMockDocument creates a new mock instance.
func NewMockDocument(ctrl *gomock.Controller) *MockDocument {
	mock := &MockDocument{ctrl: ctrl}
	mock.recorder = &MockDocumentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDocument) EXPECT() *MockDocumentMockRecorder {
	return m.recorder
}

// RetrieveFile mocks base method.
func (m *MockDocument) RetrieveFile(ctx context.Context, id string) (vtiger.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveFile", ctx, id)
	ret0, _ := ret[0].(vtiger.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveFile indicates an expected call of RetrieveFile.
func (mr *MockDocumentMockRecorder) RetrieveFile(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveFile", reflect.TypeOf((*MockDocument)(nil).RetrieveFile), ctx, id)
}

// RetrieveFromModule mocks base method.
func (m *MockDocument) RetrieveFromModule(ctx context.Context, id string) ([]domain.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveFromModule", ctx, id)
	ret0, _ := ret[0].([]domain.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveFromModule indicates an expected call of RetrieveFromModule.
func (mr *MockDocumentMockRecorder) RetrieveFromModule(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveFromModule", reflect.TypeOf((*MockDocument)(nil).RetrieveFromModule), ctx, id)
}

// MockFaq is a mock of Faq interface.
type MockFaq struct {
	ctrl     *gomock.Controller
	recorder *MockFaqMockRecorder
}

// MockFaqMockRecorder is the mock recorder for MockFaq.
type MockFaqMockRecorder struct {
	mock *MockFaq
}

// NewMockFaq creates a new mock instance.
func NewMockFaq(ctrl *gomock.Controller) *MockFaq {
	mock := &MockFaq{ctrl: ctrl}
	mock.recorder = &MockFaqMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFaq) EXPECT() *MockFaqMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockFaq) Count(ctx context.Context, client string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, client)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockFaqMockRecorder) Count(ctx, client interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockFaq)(nil).Count), ctx, client)
}

// GetAllFaqs mocks base method.
func (m *MockFaq) GetAllFaqs(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Faq, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllFaqs", ctx, filter)
	ret0, _ := ret[0].([]domain.Faq)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllFaqs indicates an expected call of GetAllFaqs.
func (mr *MockFaqMockRecorder) GetAllFaqs(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllFaqs", reflect.TypeOf((*MockFaq)(nil).GetAllFaqs), ctx, filter)
}

// MockInvoice is a mock of Invoice interface.
type MockInvoice struct {
	ctrl     *gomock.Controller
	recorder *MockInvoiceMockRecorder
}

// MockInvoiceMockRecorder is the mock recorder for MockInvoice.
type MockInvoiceMockRecorder struct {
	mock *MockInvoice
}

// NewMockInvoice creates a new mock instance.
func NewMockInvoice(ctrl *gomock.Controller) *MockInvoice {
	mock := &MockInvoice{ctrl: ctrl}
	mock.recorder = &MockInvoiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInvoice) EXPECT() *MockInvoiceMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockInvoice) Count(ctx context.Context, client string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, client)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockInvoiceMockRecorder) Count(ctx, client interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockInvoice)(nil).Count), ctx, client)
}

// GetAll mocks base method.
func (m *MockInvoice) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, filter)
	ret0, _ := ret[0].([]domain.Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockInvoiceMockRecorder) GetAll(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockInvoice)(nil).GetAll), ctx, filter)
}

// RetrieveById mocks base method.
func (m *MockInvoice) RetrieveById(ctx context.Context, id string) (domain.Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockInvoiceMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockInvoice)(nil).RetrieveById), ctx, id)
}

// MockServiceContract is a mock of ServiceContract interface.
type MockServiceContract struct {
	ctrl     *gomock.Controller
	recorder *MockServiceContractMockRecorder
}

// MockServiceContractMockRecorder is the mock recorder for MockServiceContract.
type MockServiceContractMockRecorder struct {
	mock *MockServiceContract
}

// NewMockServiceContract creates a new mock instance.
func NewMockServiceContract(ctrl *gomock.Controller) *MockServiceContract {
	mock := &MockServiceContract{ctrl: ctrl}
	mock.recorder = &MockServiceContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceContract) EXPECT() *MockServiceContractMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockServiceContract) Count(ctx context.Context, client, contact string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, client, contact)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockServiceContractMockRecorder) Count(ctx, client, contact interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockServiceContract)(nil).Count), ctx, client, contact)
}

// GetAll mocks base method.
func (m *MockServiceContract) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.ServiceContract, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, filter)
	ret0, _ := ret[0].([]domain.ServiceContract)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockServiceContractMockRecorder) GetAll(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockServiceContract)(nil).GetAll), ctx, filter)
}

// RetrieveById mocks base method.
func (m *MockServiceContract) RetrieveById(ctx context.Context, id string) (domain.ServiceContract, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.ServiceContract)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockServiceContractMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockServiceContract)(nil).RetrieveById), ctx, id)
}

// MockCurrency is a mock of Currency interface.
type MockCurrency struct {
	ctrl     *gomock.Controller
	recorder *MockCurrencyMockRecorder
}

// MockCurrencyMockRecorder is the mock recorder for MockCurrency.
type MockCurrencyMockRecorder struct {
	mock *MockCurrency
}

// NewMockCurrency creates a new mock instance.
func NewMockCurrency(ctrl *gomock.Controller) *MockCurrency {
	mock := &MockCurrency{ctrl: ctrl}
	mock.recorder = &MockCurrencyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCurrency) EXPECT() *MockCurrencyMockRecorder {
	return m.recorder
}

// RetrieveById mocks base method.
func (m *MockCurrency) RetrieveById(ctx context.Context, id string) (domain.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.Currency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockCurrencyMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockCurrency)(nil).RetrieveById), ctx, id)
}

// MockProduct is a mock of Product interface.
type MockProduct struct {
	ctrl     *gomock.Controller
	recorder *MockProductMockRecorder
}

// MockProductMockRecorder is the mock recorder for MockProduct.
type MockProductMockRecorder struct {
	mock *MockProduct
}

// NewMockProduct creates a new mock instance.
func NewMockProduct(ctrl *gomock.Controller) *MockProduct {
	mock := &MockProduct{ctrl: ctrl}
	mock.recorder = &MockProductMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProduct) EXPECT() *MockProductMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockProduct) Count(ctx context.Context, filters map[string]any) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, filters)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockProductMockRecorder) Count(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockProduct)(nil).Count), ctx, filters)
}

// GetAll mocks base method.
func (m *MockProduct) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, filter)
	ret0, _ := ret[0].([]domain.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockProductMockRecorder) GetAll(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockProduct)(nil).GetAll), ctx, filter)
}

// RetrieveById mocks base method.
func (m *MockProduct) RetrieveById(ctx context.Context, id string) (domain.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockProductMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockProduct)(nil).RetrieveById), ctx, id)
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockService) Count(ctx context.Context, filters map[string]any) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, filters)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockServiceMockRecorder) Count(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockService)(nil).Count), ctx, filters)
}

// GetAll mocks base method.
func (m *MockService) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, filter)
	ret0, _ := ret[0].([]domain.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockServiceMockRecorder) GetAll(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockService)(nil).GetAll), ctx, filter)
}

// RetrieveById mocks base method.
func (m *MockService) RetrieveById(ctx context.Context, id string) (domain.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockServiceMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockService)(nil).RetrieveById), ctx, id)
}

// MockProject is a mock of Project interface.
type MockProject struct {
	ctrl     *gomock.Controller
	recorder *MockProjectMockRecorder
}

// MockProjectMockRecorder is the mock recorder for MockProject.
type MockProjectMockRecorder struct {
	mock *MockProject
}

// NewMockProject creates a new mock instance.
func NewMockProject(ctrl *gomock.Controller) *MockProject {
	mock := &MockProject{ctrl: ctrl}
	mock.recorder = &MockProjectMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProject) EXPECT() *MockProjectMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockProject) Count(ctx context.Context, client, contact string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, client, contact)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockProjectMockRecorder) Count(ctx, client, contact interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockProject)(nil).Count), ctx, client, contact)
}

// GetAll mocks base method.
func (m *MockProject) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, filter)
	ret0, _ := ret[0].([]domain.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockProjectMockRecorder) GetAll(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockProject)(nil).GetAll), ctx, filter)
}

// RetrieveById mocks base method.
func (m *MockProject) RetrieveById(ctx context.Context, id string) (domain.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockProjectMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockProject)(nil).RetrieveById), ctx, id)
}

// MockProjectTask is a mock of ProjectTask interface.
type MockProjectTask struct {
	ctrl     *gomock.Controller
	recorder *MockProjectTaskMockRecorder
}

// MockProjectTaskMockRecorder is the mock recorder for MockProjectTask.
type MockProjectTaskMockRecorder struct {
	mock *MockProjectTask
}

// NewMockProjectTask creates a new mock instance.
func NewMockProjectTask(ctrl *gomock.Controller) *MockProjectTask {
	mock := &MockProjectTask{ctrl: ctrl}
	mock.recorder = &MockProjectTaskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectTask) EXPECT() *MockProjectTaskMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockProjectTask) Count(ctx context.Context, parent string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, parent)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockProjectTaskMockRecorder) Count(ctx, parent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockProjectTask)(nil).Count), ctx, parent)
}

// GetFromProject mocks base method.
func (m *MockProjectTask) GetFromProject(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.ProjectTask, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFromProject", ctx, filter)
	ret0, _ := ret[0].([]domain.ProjectTask)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFromProject indicates an expected call of GetFromProject.
func (mr *MockProjectTaskMockRecorder) GetFromProject(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFromProject", reflect.TypeOf((*MockProjectTask)(nil).GetFromProject), ctx, filter)
}

// RetrieveById mocks base method.
func (m *MockProjectTask) RetrieveById(ctx context.Context, id string) (domain.ProjectTask, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveById", ctx, id)
	ret0, _ := ret[0].(domain.ProjectTask)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveById indicates an expected call of RetrieveById.
func (mr *MockProjectTaskMockRecorder) RetrieveById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveById", reflect.TypeOf((*MockProjectTask)(nil).RetrieveById), ctx, id)
}
