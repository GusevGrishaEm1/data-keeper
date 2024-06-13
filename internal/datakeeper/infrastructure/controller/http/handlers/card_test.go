package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

// Mock service
type mockCardService struct {
	mock.Mock
}

// Mock method
func (m *mockCardService) CreateCard(ctx context.Context, r CreateCardRequest) (*CreateCardResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*CreateCardResponse), args.Error(1)
}

// Mock method
func (m *mockCardService) UpdateCard(ctx context.Context, r UpdateCardRequest) (*UpdateCardResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*UpdateCardResponse), args.Error(1)
}

// Mock method
func (m *mockCardService) DeleteCard(ctx context.Context, r DeleteCardRequest) (*DeleteCardResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*DeleteCardResponse), args.Error(1)
}

// Mock method
func (m *mockCardService) GetCardsByUser(ctx context.Context, r GetAllCardsRequest) (*GetAllCardsResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*GetAllCardsResponse), args.Error(1)
}

// Mock method
func MockMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("User", "test")
		return next(c)
	}
}

// Test Create card
func TestCreateCard(t *testing.T) {
	e := echo.New()
	mockService := new(mockCardService)
	handler := NewCardHandler(mockService)
	e.POST("/cards", handler.CreateCard)

	server := httptest.NewServer(e)
	defer server.Close()

	client := httpexpect.Default(t, server.URL)

	req := &CreateCardRequest{
		Key:     "1234",
		Number:  "1111222233334444",
		CVV:     "123",
		Name:    "Test Bank",
		Expires: "12/25",
	}

	uuid := uuid.New().String()

	mockService.On("CreateCard", mock.Anything, *req).Return(&CreateCardResponse{UUID: uuid}, nil)
	response := client.POST("/cards").WithJSON(req).Expect().Status(http.StatusCreated).JSON().Object()
	response.ContainsKey("uuid").ContainsValue(uuid)

	mockService.AssertExpectations(t)
}

// Test Update card
func TestUpdateCard(t *testing.T) {
	e := echo.New()
	mockService := new(mockCardService)
	handler := NewCardHandler(mockService)

	e.PATCH("/cards", handler.UpdateCard)

	server := httptest.NewServer(e)
	defer server.Close()

	client := httpexpect.Default(t, server.URL)

	key := `1234`
	number := `1111222233334444`
	cvv := `123`
	name := `Test Bank`
	expires := `12/25`

	uuid := uuid.New().String()

	req := UpdateCardRequest{
		UUID:    uuid,
		Key:     &key,
		Number:  &number,
		CVV:     &cvv,
		Name:    &name,
		Expires: &expires,
	}

	mockService.On("UpdateCard", mock.Anything, req).Return(&UpdateCardResponse{UUID: uuid}, nil)

	response := client.PATCH("/cards").WithJSON(req).Expect().Status(http.StatusOK).JSON().Object()

	response.ContainsKey("uuid").ContainsValue(uuid)

	mockService.AssertExpectations(t)
}

// Test Delete card
func TestDeleteCard(t *testing.T) {
	e := echo.New()
	mockService := new(mockCardService)
	handler := NewCardHandler(mockService)

	e.DELETE("/cards", handler.DeleteCard)

	server := httptest.NewServer(e)
	defer server.Close()

	client := httpexpect.Default(t, server.URL)
	uuid := uuid.New().String()

	req := DeleteCardRequest{UUID: uuid}

	mockService.On("DeleteCard", mock.Anything, req).Return(&DeleteCardResponse{UUID: uuid}, nil)

	response := client.DELETE("/cards").WithJSON(req).Expect().Status(http.StatusOK).JSON().Object()

	response.ContainsKey("uuid").ContainsValue(uuid)

	mockService.AssertExpectations(t)
}

// Test Get all cards
func TestGetAllCards(t *testing.T) {
	e := echo.New()
	mockService := new(mockCardService)
	handler := NewCardHandler(mockService)

	e.GET("/cards", handler.GetCardsByUser)

	server := httptest.NewServer(e)
	defer server.Close()

	client := httpexpect.Default(t, server.URL)

	res := GetAllCardsResponse{
		Items: []GetAllCardsResponceItem{
			{
				UUID:    "fac1b0c2-9b0b-11ec-9b6c-0a0027000001",
				Key:     "1234",
				Number:  "1111222233334444",
				CVV:     "123",
				Name:    "Test Bank",
				Expires: "12/25",
			},
		},
	}

	mockService.On("GetCardsByUser", mock.Anything, GetAllCardsRequest{}).Return(&res, nil)

	response := client.GET("/cards").Expect().Status(http.StatusOK).JSON().Object()

	cards := response.Value("cards").Array()
	cards.Length().IsEqual(1)
	card := cards.Value(0).Object()
	card.Value("uuid").String().IsEqual("fac1b0c2-9b0b-11ec-9b6c-0a0027000001")
	card.Value("key").String().IsEqual("1234")
	card.Value("number").String().IsEqual("1111222233334444")
	card.Value("cvv").String().IsEqual("123")
	card.Value("name").String().IsEqual("Test Bank")

	mockService.AssertExpectations(t)
}
