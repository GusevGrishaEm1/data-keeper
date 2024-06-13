package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

type mockLogPassService struct {
	mock.Mock
}

func (m *mockLogPassService) CreateLogPass(ctx context.Context, r CreateLogPassRequest) (*CreateLogPassResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*CreateLogPassResponse), args.Error(1)
}

func (m *mockLogPassService) UpdateLogPass(ctx context.Context, r UpdateLogPassRequest) (*UpdateLogPassResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*UpdateLogPassResponse), args.Error(1)
}

func (m *mockLogPassService) DeleteLogPass(ctx context.Context, r DeleteLogPassRequest) (*DeleteLogPassResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*DeleteLogPassResponse), args.Error(1)
}

func (m *mockLogPassService) GetAllLogPasses(ctx context.Context, r GetAllLogPassesRequest) (*GetAllLogPassesResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*GetAllLogPassesResponse), args.Error(1)
}

func TestCreateLogPass(t *testing.T) {
	mockService := new(mockLogPassService)
	handler := NewLogPassHandler(mockService)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.CreateLogPass(echo.New().NewContext(r, w))
	}))
	defer server.Close()

	e := httpexpect.Default(t, server.URL)
	uuid := uuid.NewString()
	expectedResponse := &CreateLogPassResponse{UUID: uuid}
	mockService.On("CreateLogPass", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	reqBody := map[string]string{
		"name":     "test",
		"login":    "user",
		"password": "pass",
	}

	e.POST("/create").
		WithJSON(reqBody).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		HasValue("uuid", uuid)

	mockService.AssertExpectations(t)
}
