package handlers

import (
	"context"
	"net/http"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/error"
	"github.com/labstack/echo/v4"
)

// LogPassService store log/pass for user
type LogPassService interface {
	// Create save log/pass
	Create(ctx context.Context, r CreateLogPassRequest) (*CreateLogPassResponse, error)
	// Update update log/pass
	Update(ctx context.Context, r UpdateLogPassRequest) (*UpdateLogPassResponse, error)
	// Delete delete log/pass
	Delete(ctx context.Context, r DeleteLogPassRequest) (*DeleteLogPassResponse, error)
	// GetAll get all log/pass for user
	GetAll(ctx context.Context, r GetAllLogPassesRequest) (*GetAllLogPassesResponse, error)
}

type CreateLogPassRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateLogPassRequest struct {
	UUID     string  `json:"uuid"`
	Name     *string `json:"name"`
	Login    *string `json:"login"`
	Password *string `json:"password"`
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
	Items []GetAllLogPassResponseItem `json:"items"`
}

type GetAllLogPassResponseItem struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LogPassHandler struct {
	service      LogPassService
	ctxConverter ctxConverter
}

func NewLogPassHandler(service LogPassService, ctxConverter ctxConverter) *LogPassHandler {
	return &LogPassHandler{
		service:      service,
		ctxConverter: ctxConverter,
	}
}

// CreateLogPass create a new log/pass
// @Summary Create a new log/pass
// @Description Create a new log/pass for the user
// @Tags logpass
// @Accept json
// @Produce json
// @Param logpass body CreateLogPassRequest true "LogPass request body"
// @Success 201 {object} CreateLogPassResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /logpass/create [post]
func (h *LogPassHandler) CreateLogPass(c echo.Context) error {
	req := new(CreateLogPassRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	ctx, err := h.ctxConverter.ConvertEchoCtxToCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	res, err := h.service.Create(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusCreated, res)
}

// UpdateLogPass update an existing log/pass
// @Summary Update an existing log/pass
// @Description Update an existing log/pass for the user
// @Tags logpass
// @Accept json
// @Produce json
// @Param logpass body UpdateLogPassRequest true "LogPass request body"
// @Success 200 {object} UpdateLogPassResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /logpass/update [post]
func (h *LogPassHandler) UpdateLogPass(c echo.Context) error {
	req := new(UpdateLogPassRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	ctx, err := h.ctxConverter.ConvertEchoCtxToCtx(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.service.Update(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteLogPass delete an existing log/pass
// @Summary Delete an existing log/pass
// @Description Delete an existing log/pass for the user
// @Tags logpass
// @Accept json
// @Produce json
// @Param logpass body DeleteLogPassRequest true "LogPass request body"
// @Success 200 {object} DeleteLogPassResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /logpass/delete [post]
func (h *LogPassHandler) DeleteLogPass(c echo.Context) error {
	req := new(DeleteLogPassRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	ctx, err := h.ctxConverter.ConvertEchoCtxToCtx(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.service.Delete(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}

// GetAllLogPasses get all log/pass for user
// @Summary Get all log/pass
// @Description Get all log/pass for the user
// @Tags logpass
// @Produce json
// @Success 200 {object} GetAllLogPassesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /logpass [get]
func (h *LogPassHandler) GetAllLogPasses(c echo.Context) error {
	req := new(GetAllLogPassesRequest)

	ctx, err := h.ctxConverter.ConvertEchoCtxToCtx(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.service.GetAll(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}
