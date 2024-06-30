package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

func TestCreateLogPass(t *testing.T) {
	mockService := new(mockLogPassService)
	mockConverter := new(mockCtxConverter)
	handler := NewLogPassHandler(mockService, mockConverter)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler.CreateLogPass(echo.New().NewContext(r, w))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	e := httpexpect.Default(t, server.URL)
	uuidStr := uuid.NewString()
	expectedResponse := &CreateLogPassResponse{UUID: uuidStr}
	mockConverter.On("ConvertEchoCtxToCtx", mock.Anything).Return(nil, nil)
	mockService.On("Create", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	reqBody := map[string]string{
		"name":     "test",
		"login":    "user",
		"password": "pass",
	}

	e.POST("/logpass").
		WithJSON(reqBody).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		HasValue("uuid", uuidStr)

	mockService.AssertExpectations(t)
	mockConverter.AssertExpectations(t)
}
