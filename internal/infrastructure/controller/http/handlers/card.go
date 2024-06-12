package handlers

import (
	"context"
	"net/http"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/error"
	"github.com/labstack/echo"
)

// Card service
type CardService interface {
	// Create card
	CreateCard(ctx context.Context, r CreateCardRequest) (*CreateCardResponse, error)
	// Update card
	UpdateCard(ctx context.Context, r UpdateCardRequest) (*UpdateCardResponse, error)
	// Delete card
	DeleteCard(ctx context.Context, r DeleteCardRequest) (*DeleteCardResponse, error)
	// Get cards by user
	GetCardsByUser(ctx context.Context, r GetAllCardsRequest) (*GetAllCardsResponse, error)
}

// Create card request
type CreateCardRequest struct {
	Key     string `json:"key"`
	Number  string `json:"number"`
	CVV     string `json:"cvv"`
	Name    string `json:"name"`
	Expires string `json:"expires"`
}

// Create card response
type CreateCardResponse struct {
	UUID string `json:"uuid"`
}

// Update card request
type UpdateCardRequest struct {
	UUID    string  `json:"uuid"`
	Key     *string `json:"key,omitempty"`
	Number  *string `json:"number,omitempty"`
	CVV     *string `json:"cvv,omitempty"`
	Name    *string `json:"name,omitempty"`
	Expires *string `json:"expires,omitempty"`
}

// Update card response
type UpdateCardResponse struct {
	UUID string `json:"uuid"`
}

// Delete card request
type DeleteCardRequest struct {
	UUID string `json:"uuid"`
}

// Delete card response
type DeleteCardResponse struct {
	UUID string `json:"uuid"`
}

// Get all cards request
type GetAllCardsRequest struct{}

// Get all cards response
type GetAllCardsResponse struct {
	Items []GetAllCardsResponceItem `json:"cards"`
}

// Get all cards responce item
type GetAllCardsResponceItem struct {
	UUID    string `json:"uuid"`
	Key     string `json:"key"`
	Number  string `json:"number"`
	CVV     string `json:"cvv"`
	Name    string `json:"name"`
	Expires string `json:"expires"`
}

// Card handler
type CardHandler struct {
	cardService CardService
}

// NewCardHandler creates new card handler
func NewCardHandler(cardService CardService) *CardHandler {
	return &CardHandler{
		cardService: cardService,
	}
}

// Create card
func (h *CardHandler) CreateCard(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, customerr.ToJson("unauthorized"))
	}
	ctx := context.WithValue(c.Request().Context(), "User", user)
	req := new(CreateCardRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.cardService.CreateCard(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.JSON(http.StatusCreated, res)
}

// Update card
func (h *CardHandler) UpdateCard(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, customerr.ToJson("unauthorized"))
	}
	ctx := context.WithValue(c.Request().Context(), "User", user)
	req := new(UpdateCardRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.cardService.UpdateCard(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

// Delete card
func (h *CardHandler) DeleteCard(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, customerr.ToJson("unauthorized"))
	}
	ctx := context.WithValue(c.Request().Context(), "User", user)
	req := new(DeleteCardRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.cardService.DeleteCard(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

// Get cards by user
func (h *CardHandler) GetCardsByUser(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, customerr.ToJson("unauthorized"))
	}
	ctx := context.WithValue(c.Request().Context(), "User", user)
	req := new(GetAllCardsRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.cardService.GetCardsByUser(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}
