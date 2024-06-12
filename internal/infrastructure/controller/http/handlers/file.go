package handlers

import (
	"context"
	"net/http"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/error"
	"github.com/labstack/echo"
)

type FileService interface {
	UploadFile(ctx context.Context, r UploadFileRequest) (*UploadFileResponse, error)
	DeleteFile(ctx context.Context, r DeleteFileRequest) (*DeleteFileResponse, error)
	GetAllFiles(ctx context.Context, r GetAllFilesRequest) (*GetAllFilesResponse, error)
	DownloadFile(ctx context.Context, r DownloadFileRequest) (*DownloadFileResponse, error)
}

type UploadFileRequest struct {
	Name   string
	Format string
	File   []byte
}

type UploadFileResponse struct {
	UUID string `json:"uuid"`
}

type DeleteFileRequest struct {
	UUID string `json:"uuid"`
}

type DeleteFileResponse struct {
	UUID string `json:"uuid"`
}

type GetAllFilesRequest struct{}

type GetAllFilesResponse struct {
	Items []GetAllFilesResponceItem `json:"items"`
}

type GetAllFilesResponceItem struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Format string `json:"format"`
	Size   int    `json:"size"`
}

type DownloadFileRequest struct {
	UUID string
}

type DownloadFileResponse struct {
	Name   string
	Format string
	File   []byte
}

type FileHandler struct {
	fileService FileService
}

func NewFileHandler(fileService FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

func (h *FileHandler) UploadFile(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, customerr.ToJson("unauthorized"))
	}
	ctx := context.WithValue(c.Request().Context(), "User", user)
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

	fileName := file.Filename
	fileFormat := file.Header.Get("Content-Type")
	req := UploadFileRequest{Name: fileName, Format: fileFormat, File: buf}

	res, err := h.fileService.UploadFile(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *FileHandler) DeleteFile(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, customerr.ToJson("unauthorized"))
	}
	ctx := context.WithValue(c.Request().Context(), "User", user)
	req := new(DeleteFileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.fileService.DeleteFile(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

func (h *FileHandler) GetAllFiles(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, customerr.ToJson("unauthorized"))
	}
	ctx := context.WithValue(c.Request().Context(), "User", user)
	req := GetAllFilesRequest{}

	res, err := h.fileService.GetAllFiles(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *FileHandler) DownloadFile(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, customerr.ToJson("unauthorized"))
	}
	ctx := context.WithValue(c.Request().Context(), "User", user)
	uuid := c.Param("uuid")
	req := DownloadFileRequest{UUID: uuid}

	res, err := h.fileService.DownloadFile(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.Blob(http.StatusOK, "application/octet-stream", res.File)
}
