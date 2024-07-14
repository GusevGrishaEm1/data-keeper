package handlers

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

func setupFileServer(mockFileService *mockFileService, mockConverter *mockCtxConverter) *echo.Echo {
	e := echo.New()

	handler := NewFileHandler(mockFileService, mockConverter)

	e.POST("/files", handler.UploadFile)
	e.DELETE("/files", handler.DeleteFile)
	e.GET("/files", handler.GetAllFiles)
	e.GET("/files/:uuid", handler.DownloadFile)

	return e
}

func TestFileHandler_UploadFile(t *testing.T) {
	mockFileService := new(mockFileService)
	mockCtxConverter := new(mockCtxConverter)
	uploadResponse := &UploadFileResponse{UUID: uuid.New().String()}

	mockCtxConverter.On("ConvertEchoCtxToCtx", mock.Anything).Return(context.TODO(), nil)
	mockFileService.On("UploadFile", mock.Anything, mock.Anything).Return(uploadResponse, nil)

	e := setupFileServer(mockFileService, mockCtxConverter)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", "test.txt")
	require.NoError(t, err)
	_, err = fw.Write([]byte("test content"))
	require.NoError(t, err)
	err = w.Close()
	require.NoError(t, err)

	expect.POST("/files").
		WithHeader("Content-Type", w.FormDataContentType()).
		WithBytes(b.Bytes()).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().ContainsKey("uuid")

	mockFileService.AssertExpectations(t)
	mockCtxConverter.AssertExpectations(t)
}

func TestFileHandler_DeleteFile(t *testing.T) {
	mockFileService := new(mockFileService)
	mockCtxConverter := new(mockCtxConverter)
	deleteResponse := &DeleteFileResponse{UUID: "123"}

	mockCtxConverter.On("ConvertEchoCtxToCtx", mock.Anything).Return(nil, nil)
	mockFileService.On("DeleteFile", mock.Anything, mock.Anything).Return(deleteResponse, nil)

	e := setupFileServer(mockFileService, mockCtxConverter)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	expect.DELETE("/files").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("uuid").HasValue("uuid", "123")

	mockFileService.AssertExpectations(t)
	mockCtxConverter.AssertExpectations(t)
}

func TestFileHandler_GetAllFiles(t *testing.T) {
	mockFileService := new(mockFileService)
	mockCtxConverter := new(mockCtxConverter)
	getAllFilesResponse := &GetAllFilesResponse{
		Items: []GetAllFilesResponceItem{
			{
				UUID:   "123",
				Name:   "test.txt",
				Format: "text/plain",
				Size:   10,
			},
		},
	}

	mockCtxConverter.On("ConvertEchoCtxToCtx", mock.Anything).Return(nil, nil)
	mockFileService.On("GetAllFiles", mock.Anything, mock.Anything).Return(getAllFilesResponse, nil)

	e := setupFileServer(mockFileService, mockCtxConverter)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	expect.GET("/files").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("items")

	mockFileService.AssertExpectations(t)
	mockCtxConverter.AssertExpectations(t)
}

func TestFileHandler_DownloadFile(t *testing.T) {
	mockFileService := new(mockFileService)
	mockCtxConverter := new(mockCtxConverter)
	downloadResponse := &DownloadFileResponse{File: []byte("test content")}

	mockCtxConverter.On("ConvertEchoCtxToCtx", mock.Anything).Return(nil, nil)
	mockFileService.On("DownloadFile", mock.Anything, mock.Anything).Return(downloadResponse, nil)

	e := setupFileServer(mockFileService, mockCtxConverter)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	muid := uuid.New().String()

	expect.GET("/files/" + muid).
		Expect().
		Status(http.StatusOK).
		HasContentType("application/octet-stream").
		Body().IsEqual("test content")

	mockFileService.AssertExpectations(t)
	mockCtxConverter.AssertExpectations(t)
}
