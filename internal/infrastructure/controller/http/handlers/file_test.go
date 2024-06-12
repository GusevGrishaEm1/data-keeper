package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) UploadFile(ctx context.Context, r UploadFileRequest) (*UploadFileResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*UploadFileResponse), args.Error(1)
}

func (m *MockFileService) DeleteFile(ctx context.Context, r DeleteFileRequest) (*DeleteFileResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*DeleteFileResponse), args.Error(1)
}

func (m *MockFileService) GetAllFiles(ctx context.Context, r GetAllFilesRequest) (*GetAllFilesResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*GetAllFilesResponse), args.Error(1)
}

func (m *MockFileService) DownloadFile(ctx context.Context, r DownloadFileRequest) (*DownloadFileResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*DownloadFileResponse), args.Error(1)
}

func setupFileServer(mockFileService *MockFileService) *echo.Echo {
	e := echo.New()
	handler := NewFileHandler(mockFileService)

	e.POST("/upload", handler.UploadFile)
	e.DELETE("/files/:uuid", handler.DeleteFile)
	e.GET("/files", handler.GetAllFiles)
	e.GET("/files/:uuid", handler.DownloadFile)

	return e
}

func TestFileHandler_UploadFile(t *testing.T) {
	mockFileService := new(MockFileService)
	uploadResponse := &UploadFileResponse{UUID: uuid.New().String()}

	mockFileService.On("UploadFile", mock.Anything, mock.Anything).Return(uploadResponse, nil)

	e := setupFileServer(mockFileService)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("CreateFormFile: %s", err)
	}
	_, err = fw.Write([]byte("test content"))
	if err != nil {
		t.Fatalf("Write: %s", err)
	}
	w.Close()

	expect.POST("/upload").
		WithHeader("Content-Type", w.FormDataContentType()).
		WithBytes(b.Bytes()).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().ContainsKey("uuid")

	mockFileService.AssertExpectations(t)
}

func TestFileHandler_DeleteFile(t *testing.T) {
	mockFileService := new(MockFileService)
	deleteResponse := &DeleteFileResponse{UUID: "123"}

	mockFileService.On("DeleteFile", mock.Anything, mock.Anything).Return(deleteResponse, nil)

	e := setupFileServer(mockFileService)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	expect.DELETE("/files/123").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("uuid").HasValue("uuid", "123")

	mockFileService.AssertExpectations(t)
}

func TestFileHandler_GetAllFiles(t *testing.T) {
	mockFileService := new(MockFileService)
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

	mockFileService.On("GetAllFiles", mock.Anything, mock.Anything).Return(getAllFilesResponse, nil)

	e := setupFileServer(mockFileService)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	expect.GET("/files").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("items")

	mockFileService.AssertExpectations(t)
}

func TestFileHandler_DownloadFile(t *testing.T) {
	mockFileService := new(MockFileService)
	downloadResponse := &DownloadFileResponse{File: []byte("test content")}

	mockFileService.On("DownloadFile", mock.Anything, mock.Anything).Return(downloadResponse, nil)

	e := setupFileServer(mockFileService)

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
}
