package handlers

import (
	"context"
	"net/http"
	"strings"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/error"
	"github.com/labstack/echo"
)

// File service
type FileService interface {
	// Upload file
	UploadFile(ctx context.Context, r UploadFileRequest) (*UploadFileResponse, error)
	// Delete file
	DeleteFile(ctx context.Context, r DeleteFileRequest) (*DeleteFileResponse, error)
	// Get all files
	GetAllFiles(ctx context.Context, r GetAllFilesRequest) (*GetAllFilesResponse, error)
	// Download file
	DownloadFile(ctx context.Context, r DownloadFileRequest) (*DownloadFileResponse, error)
}

// Upload file request
type UploadFileRequest struct {
	Name   string
	Format string
	File   []byte
}

// Upload file response
type UploadFileResponse struct {
	UUID string `json:"uuid"`
}

// Delete file request
type DeleteFileRequest struct {
	UUID string `json:"uuid"`
}

// Delete file response
type DeleteFileResponse struct {
	UUID string `json:"uuid"`
}

// Get all files request
type GetAllFilesRequest struct{}

// Get all files response
type GetAllFilesResponse struct {
	Items []GetAllFilesResponceItem `json:"items"`
}

// Get all files responce item
type GetAllFilesResponceItem struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Format string `json:"format"`
	Size   int    `json:"size"`
}

// Download file request
type DownloadFileRequest struct {
	UUID string
}

// Download file response
type DownloadFileResponse struct {
	Name   string
	Format string
	File   []byte
}

// File handler
type FileHandler struct {
	fileService FileService
}

// NewFileHandler create new file handler
func NewFileHandler(fileService FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
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
	defer src.Close()

	if file.Size > 5*1024*1024 {
		return c.JSON(http.StatusBadRequest, customerr.ToJson("File is too large"))
	}

	buf := make([]byte, file.Size)
	_, err = src.Read(buf)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	strs := strings.Split(file.Filename, ".")
	req := UploadFileRequest{Name: strs[0], Format: strs[1], File: buf}
	res, err := h.fileService.UploadFile(c.Request().Context(), req)
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

	res, err := h.fileService.DeleteFile(c.Request().Context(), *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// GetAllFiles get all files for user
func (h *FileHandler) GetAllFiles(c echo.Context) error {
	res, err := h.fileService.GetAllFiles(c.Request().Context(), GetAllFilesRequest{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// DownloadFile downloads file for user
func (h *FileHandler) DownloadFile(c echo.Context) error {
	uuid := c.Param("uuid")

	res, err := h.fileService.DownloadFile(c.Request().Context(), DownloadFileRequest{UUID: uuid})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.Blob(http.StatusOK, "application/octet-stream", res.File)
}
