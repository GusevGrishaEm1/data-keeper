package handlers

import (
	"context"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

// mockCtxConverter mock ctx converter
type mockCtxConverter struct {
	mock.Mock
}

// ConvertEchoCtxToCtx mock method
func (m *mockCtxConverter) ConvertEchoCtxToCtx(ctx echo.Context) (context.Context, error) {
	m.Called(ctx)
	return nil, nil
}

type mockLogPassService struct {
	mock.Mock
}

func (m *mockLogPassService) Create(ctx context.Context, r CreateLogPassRequest) (*CreateLogPassResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*CreateLogPassResponse), args.Error(1)
}

func (m *mockLogPassService) Update(ctx context.Context, r UpdateLogPassRequest) (*UpdateLogPassResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*UpdateLogPassResponse), args.Error(1)
}

func (m *mockLogPassService) Delete(ctx context.Context, r DeleteLogPassRequest) (*DeleteLogPassResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*DeleteLogPassResponse), args.Error(1)
}

func (m *mockLogPassService) GetAll(ctx context.Context, r GetAllLogPassesRequest) (*GetAllLogPassesResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*GetAllLogPassesResponse), args.Error(1)
}

type mockFileService struct {
	mock.Mock
}

func (m *mockFileService) UploadFile(ctx context.Context, r UploadFileRequest) (*UploadFileResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*UploadFileResponse), args.Error(1)
}

func (m *mockFileService) DeleteFile(ctx context.Context, r DeleteFileRequest) (*DeleteFileResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*DeleteFileResponse), args.Error(1)
}

func (m *mockFileService) GetAllFiles(ctx context.Context, r GetAllFilesRequest) (*GetAllFilesResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*GetAllFilesResponse), args.Error(1)
}

func (m *mockFileService) DownloadFile(ctx context.Context, r DownloadFileRequest) (*DownloadFileResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*DownloadFileResponse), args.Error(1)
}

// Mock service
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) SignIn(ctx context.Context, r LoginRequest) (*LoginResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*LoginResponse), args.Error(1)
}

func (m *mockAuthService) SignUp(ctx context.Context, r RegisterRequest) (*RegisterResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*RegisterResponse), args.Error(1)
}
