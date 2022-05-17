// Code generated by MockGen. DO NOT EDIT.
// Source: charts.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCharts is a mock of Charts interface.
type MockCharts struct {
	ctrl     *gomock.Controller
	recorder *MockChartsMockRecorder
}

// MockChartsMockRecorder is the mock recorder for MockCharts.
type MockChartsMockRecorder struct {
	mock *MockCharts
}

// NewMockCharts creates a new mock instance.
func NewMockCharts(ctrl *gomock.Controller) *MockCharts {
	mock := &MockCharts{ctrl: ctrl}
	mock.recorder = &MockChartsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharts) EXPECT() *MockChartsMockRecorder {
	return m.recorder
}

// Sale mocks base method.
func (m *MockCharts) Sale(title, artist string, copies int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Sale", title, artist, copies)
}

// Sale indicates an expected call of Sale.
func (mr *MockChartsMockRecorder) Sale(title, artist, copies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sale", reflect.TypeOf((*MockCharts)(nil).Sale), title, artist, copies)
}
