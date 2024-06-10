package handlers

import (
	"context"
	"net/http"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/error"
	"github.com/labstack/echo"
)

type LogPassService interface {
	CreateLogPass(ctx context.Context, r CreateLogPassRequest) (*CreateLogPassResponse, error)
	UpdateLogPass(ctx context.Context, r UpdateLogPassRequest) (*UpdateLogPassResponse, error)
	DeleteLogPass(ctx context.Context, r DeleteLogPassRequest) (*DeleteLogPassResponse, error)
	GetAllLogPasses(ctx context.Context, r GetAllLogPassesRequest) (*GetAllLogPassesResponse, error)
}

type CreateLogPassRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateLogPassRequest struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DeleteLogPassRequest struct {
	UUID string `json:"uuid"`
}

type GetAllLogPassesRequest struct{}

type CreateLogPassResponse struct {
	UUID string `json:"uuid"`
}

type UpdateLogPassResponse struct {
	UUID string `json:"uuid"`
}

type DeleteLogPassResponse struct {
	UUID string `json:"uuid"`
}

type GetAllLogPassesResponse struct {
	Items []GetAllLogPassResponceItem `json:"items"`
}

type GetAllLogPassResponceItem struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LogPassHandler struct {
	service LogPassService
}

func NewLogPassHandler(service LogPassService) *LogPassHandler {
	return &LogPassHandler{
		service: service,
	}
}

func (h *LogPassHandler) CreateLogPass(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(CreateLogPassRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.service.CreateLogPass(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.JSON(http.StatusCreated, res)
}

func (h *LogPassHandler) UpdateLogPass(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(UpdateLogPassRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.service.UpdateLogPass(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *LogPassHandler) DeleteLogPass(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(DeleteLogPassRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.service.DeleteLogPass(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *LogPassHandler) GetAllLogPasses(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(GetAllLogPassesRequest)
	res, err := h.service.GetAllLogPasses(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}
