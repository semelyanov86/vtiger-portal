// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	multipart "mime/multipart"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	domain "github.com/semelyanov86/vtiger-portal/internal/domain"
	vtiger "github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

// MockContextServiceInterface is a mock of ContextServiceInterface interface.
type MockContextServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockContextServiceInterfaceMockRecorder
}

// MockContextServiceInterfaceMockRecorder is the mock recorder for MockContextServiceInterface.
type MockContextServiceInterfaceMockRecorder struct {
	mock *MockContextServiceInterface
}

// NewMockContextServiceInterface creates a new mock instance.
func NewMockContextServiceInterface(ctrl *gomock.Controller) *MockContextServiceInterface {
	mock := &MockContextServiceInterface{ctrl: ctrl}
	mock.recorder = &MockContextServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContextServiceInterface) EXPECT() *MockContextServiceInterfaceMockRecorder {
	return m.recorder
}

// ContextGetUser mocks base method.
func (m *MockContextServiceInterface) ContextGetUser(c *gin.Context) *domain.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContextGetUser", c)
	ret0, _ := ret[0].(*domain.User)
	return ret0
}

// ContextGetUser indicates an expected call of ContextGetUser.
func (mr *MockContextServiceInterfaceMockRecorder) ContextGetUser(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContextGetUser", reflect.TypeOf((*MockContextServiceInterface)(nil).ContextGetUser), c)
}

// ContextSetUser mocks base method.
func (m *MockContextServiceInterface) ContextSetUser(c *gin.Context, user *domain.User) *gin.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContextSetUser", c, user)
	ret0, _ := ret[0].(*gin.Context)
	return ret0
}

// ContextSetUser indicates an expected call of ContextSetUser.
func (mr *MockContextServiceInterfaceMockRecorder) ContextSetUser(c, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContextSetUser", reflect.TypeOf((*MockContextServiceInterface)(nil).ContextSetUser), c, user)
}

// MockCommentServiceInterface is a mock of CommentServiceInterface interface.
type MockCommentServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCommentServiceInterfaceMockRecorder
}

// MockCommentServiceInterfaceMockRecorder is the mock recorder for MockCommentServiceInterface.
type MockCommentServiceInterfaceMockRecorder struct {
	mock *MockCommentServiceInterface
}

// NewMockCommentServiceInterface creates a new mock instance.
func NewMockCommentServiceInterface(ctrl *gomock.Controller) *MockCommentServiceInterface {
	mock := &MockCommentServiceInterface{ctrl: ctrl}
	mock.recorder = &MockCommentServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentServiceInterface) EXPECT() *MockCommentServiceInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCommentServiceInterface) Create(ctx context.Context, content, related, userId string) (domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, content, related, userId)
	ret0, _ := ret[0].(domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCommentServiceInterfaceMockRecorder) Create(ctx, content, related, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCommentServiceInterface)(nil).Create), ctx, content, related, userId)
}

// GetRelated mocks base method.
func (m *MockCommentServiceInterface) GetRelated(ctx context.Context, id string) ([]domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRelated", ctx, id)
	ret0, _ := ret[0].([]domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRelated indicates an expected call of GetRelated.
func (mr *MockCommentServiceInterfaceMockRecorder) GetRelated(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRelated", reflect.TypeOf((*MockCommentServiceInterface)(nil).GetRelated), ctx, id)
}

// MockDocumentServiceInterface is a mock of DocumentServiceInterface interface.
type MockDocumentServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDocumentServiceInterfaceMockRecorder
}

// MockDocumentServiceInterfaceMockRecorder is the mock recorder for MockDocumentServiceInterface.
type MockDocumentServiceInterfaceMockRecorder struct {
	mock *MockDocumentServiceInterface
}

// NewMockDocumentServiceInterface creates a new mock instance.
func NewMockDocumentServiceInterface(ctrl *gomock.Controller) *MockDocumentServiceInterface {
	mock := &MockDocumentServiceInterface{ctrl: ctrl}
	mock.recorder = &MockDocumentServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDocumentServiceInterface) EXPECT() *MockDocumentServiceInterfaceMockRecorder {
	return m.recorder
}

// AttachFile mocks base method.
func (m *MockDocumentServiceInterface) AttachFile(ctx context.Context, file multipart.File, id string, userModel domain.User, header *multipart.FileHeader) (domain.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachFile", ctx, file, id, userModel, header)
	ret0, _ := ret[0].(domain.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AttachFile indicates an expected call of AttachFile.
func (mr *MockDocumentServiceInterfaceMockRecorder) AttachFile(ctx, file, id, userModel, header interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachFile", reflect.TypeOf((*MockDocumentServiceInterface)(nil).AttachFile), ctx, file, id, userModel, header)
}

// DeleteFile mocks base method.
func (m *MockDocumentServiceInterface) DeleteFile(ctx context.Context, id, related string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", ctx, id, related)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockDocumentServiceInterfaceMockRecorder) DeleteFile(ctx, id, related interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockDocumentServiceInterface)(nil).DeleteFile), ctx, id, related)
}

// GetFile mocks base method.
func (m *MockDocumentServiceInterface) GetFile(ctx context.Context, id, relatedId string) (vtiger.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile", ctx, id, relatedId)
	ret0, _ := ret[0].(vtiger.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFile indicates an expected call of GetFile.
func (mr *MockDocumentServiceInterfaceMockRecorder) GetFile(ctx, id, relatedId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockDocumentServiceInterface)(nil).GetFile), ctx, id, relatedId)
}

// GetRelated mocks base method.
func (m *MockDocumentServiceInterface) GetRelated(ctx context.Context, id string) ([]domain.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRelated", ctx, id)
	ret0, _ := ret[0].([]domain.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRelated indicates an expected call of GetRelated.
func (mr *MockDocumentServiceInterfaceMockRecorder) GetRelated(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRelated", reflect.TypeOf((*MockDocumentServiceInterface)(nil).GetRelated), ctx, id)
}
