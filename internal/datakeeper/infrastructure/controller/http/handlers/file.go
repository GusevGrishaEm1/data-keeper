package handlers

import (
	"context"
	"mime/multipart"
	"net/http"
	"strings"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/error"
	"github.com/labstack/echo"
)

// FileService File service
type FileService interface {
	// UploadFile Upload file
	UploadFile(ctx context.Context, r UploadFileRequest) (*UploadFileResponse, error)
	// DeleteFile Delete file
	DeleteFile(ctx context.Context, r DeleteFileRequest) (*DeleteFileResponse, error)
	// GetAllFiles Get all files for user
	GetAllFiles(ctx context.Context, r GetAllFilesRequest) (*GetAllFilesResponse, error)
	// DownloadFile Download file
	DownloadFile(ctx context.Context, r DownloadFileRequest) (*DownloadFileResponse, error)
}

// UploadFileRequest Upload file request
type UploadFileRequest struct {
	Name   string
	Format string
	File   []byte
}

// UploadFileResponse Upload file response
type UploadFileResponse struct {
	UUID string `json:"uuid"`
}

// DeleteFileRequest Delete file request
type DeleteFileRequest struct {
	UUID string `json:"uuid"`
}

// DeleteFileResponse Delete file response
type DeleteFileResponse struct {
	UUID string `json:"uuid"`
}

// GetAllFilesRequest Get all files request
type GetAllFilesRequest struct{}

// GetAllFilesResponse Get all files response
type GetAllFilesResponse struct {
	Items []GetAllFilesResponceItem `json:"items"`
}

// GetAllFilesResponceItem Get all files responce item
type GetAllFilesResponceItem struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Format string `json:"format"`
	Size   int    `json:"size"`
}

// DownloadFileRequest Download file request
type DownloadFileRequest struct {
	UUID string `json:"uuid"`
}

// DownloadFileResponse Download file response
type DownloadFileResponse struct {
	Name   string
	Format string
	File   []byte
}

// FileHandler File handler
type FileHandler struct {
	fileService  FileService
	ctxConverter ctxConverter
}

// NewFileHandler create new file handler
func NewFileHandler(fileService FileService, ctxConverter ctxConverter) *FileHandler {
	return &FileHandler{fileService: fileService, ctxConverter: ctxConverter}
}

// UploadFile upload file for user
func (h *FileHandler) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
		}
	}(src)

	if file.Size > 5*1024*1024 {
		return c.JSON(http.StatusBadRequest, customerr.ToJson("File is too large"))
	}

	buf := make([]byte, file.Size)
	_, err = src.Read(buf)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	ctx, err := h.ctxConverter.ConvertEchoCtxToCtx(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	strs := strings.Split(file.Filename, ".")
	req := UploadFileRequest{Name: strs[0], Format: strs[1], File: buf}
	res, err := h.fileService.UploadFile(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusCreated, res)
}

// DeleteFile delete file for user
func (h *FileHandler) DeleteFile(c echo.Context) error {
	req := new(DeleteFileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	ctx, err := h.ctxConverter.ConvertEchoCtxToCtx(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.fileService.DeleteFile(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// GetAllFiles get all files for user
func (h *FileHandler) GetAllFiles(c echo.Context) error {
	ctx, err := h.ctxConverter.ConvertEchoCtxToCtx(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.fileService.GetAllFiles(ctx, GetAllFilesRequest{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// DownloadFile download file for user
func (h *FileHandler) DownloadFile(c echo.Context) error {
	req := new(DownloadFileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	ctx, err := h.ctxConverter.ConvertEchoCtxToCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	res, err := h.fileService.DownloadFile(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.Blob(http.StatusOK, "application/octet-stream", res.File)
}
