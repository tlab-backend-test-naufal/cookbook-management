// Code generated by MockGen. DO NOT EDIT.
// Source: category_usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
	usecase "github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
)

// MockCategoryRepository is a mock of CategoryRepository interface.
type MockCategoryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryRepositoryMockRecorder
}

// MockCategoryRepositoryMockRecorder is the mock recorder for MockCategoryRepository.
type MockCategoryRepositoryMockRecorder struct {
	mock *MockCategoryRepository
}

// NewMockCategoryRepository creates a new mock instance.
func NewMockCategoryRepository(ctrl *gomock.Controller) *MockCategoryRepository {
	mock := &MockCategoryRepository{ctrl: ctrl}
	mock.recorder = &MockCategoryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCategoryRepository) EXPECT() *MockCategoryRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCategoryRepository) Create(ctx context.Context, params usecase.CategoryParams) (*entity.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, params)
	ret0, _ := ret[0].(*entity.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCategoryRepositoryMockRecorder) Create(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCategoryRepository)(nil).Create), ctx, params)
}

// Delete mocks base method.
func (m *MockCategoryRepository) Delete(ctx context.Context, id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCategoryRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCategoryRepository)(nil).Delete), ctx, id)
}

// List mocks base method.
func (m *MockCategoryRepository) List(ctx context.Context, limit, offset int) (entity.Categories, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset)
	ret0, _ := ret[0].(entity.Categories)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCategoryRepositoryMockRecorder) List(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCategoryRepository)(nil).List), ctx, limit, offset)
}

// Update mocks base method.
func (m *MockCategoryRepository) Update(ctx context.Context, id uint64, params usecase.CategoryParams) (*entity.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, params)
	ret0, _ := ret[0].(*entity.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockCategoryRepositoryMockRecorder) Update(ctx, id, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCategoryRepository)(nil).Update), ctx, id, params)
}
