package card

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http/handlers"
	"github.com/GusevGrishaEm1/data-keeper/internal/lib"
	"github.com/google/uuid"
)

type CardContent struct {
	Key     string `json:"key"`
	Number  string `json:"number"`
	CVV     string `json:"cvv"`
	Name    string `json:"name"`
	Expires string `json:"expires"`
}

type DataRepo interface {
	Insert(ctx context.Context, data entity.Data) error
	Update(ctx context.Context, data entity.Data) error
	Delete(ctx context.Context, user string, uuid string) error
	GetByUUID(ctx context.Context, user string, uuid string, contentType entity.ContentType) (*entity.Data, error)
	GetByUser(ctx context.Context, user string, contentType entity.ContentType) ([]*entity.Data, error)
}

type AuthService interface {
	GetUserFromContext(ctx context.Context) (string, error)
}

type KeyService interface {
	GetKeyForUser(user string) (string, error)
}

type cardService struct {
	repo        DataRepo
	authService AuthService
	keyService  KeyService
	logger      *slog.Logger
}

func NewCardService(repo DataRepo, authService AuthService, keyService KeyService, slog *slog.Logger) *cardService {
	return &cardService{
		repo:        repo,
		authService: authService,
		keyService:  keyService,
		logger:      slog,
	}
}

func (s *cardService) CreateCard(ctx context.Context, r handlers.CreateCardRequest) (*handlers.CreateCardResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	key, err := s.keyService.GetKeyForUser(user)
	if err != nil {
		return nil, err
	}

	newuuid := uuid.New().String()

	cardContent := CardContent{
		Key:     r.Key,
		Number:  r.Number,
		CVV:     r.CVV,
		Name:    r.Name,
		Expires: r.Expires,
	}
	jsonContent, err := json.Marshal(cardContent)
	if err != nil {
		return nil, err
	}

	encryptedContent, err := lib.Encrypt(key, jsonContent)
	if err != nil {
		return nil, err
	}

	entity := entity.Data{
		UUID:        newuuid,
		Content:     encryptedContent,
		ContentType: entity.CARD,
		CreatedAt:   time.Now(),
		CreatedBy:   user,
	}
	err = s.repo.Insert(ctx, entity)
	if err != nil {
		return nil, err
	}

	return &handlers.CreateCardResponse{UUID: newuuid}, nil
}

func (s *cardService) UpdateCard(ctx context.Context, r handlers.UpdateCardRequest) (*handlers.UpdateCardResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	key, err := s.keyService.GetKeyForUser(user)
	if err != nil {
		return nil, err
	}

	cardContent, err := s.getUpdatedCard(ctx, user, key, r)
	if err != nil {
		return nil, err
	}
	jsonContent, err := json.Marshal(cardContent)
	if err != nil {
		return nil, err
	}

	encryptedContent, err := lib.Encrypt(key, jsonContent)
	if err != nil {
		return nil, err
	}

	entity := entity.Data{
		UUID:        r.UUID,
		Content:     encryptedContent,
		ContentType: entity.CARD,
		CreatedAt:   time.Now(),
		CreatedBy:   user,
	}
	err = s.repo.Update(ctx, entity)
	if err != nil {
		return nil, err
	}

	return &handlers.UpdateCardResponse{UUID: r.UUID}, nil
}

func (s *cardService) getUpdatedCard(ctx context.Context, user string, key string, r handlers.UpdateCardRequest) (*CardContent, error) {
	card, err := s.repo.GetByUUID(ctx, user, r.UUID, entity.CARD)
	if err != nil {
		return nil, err
	}
	decryptedContent, err := lib.Decrypt(key, card.Content)
	if err != nil {
		return nil, err
	}
	cardDB := &CardContent{}
	err = json.Unmarshal(decryptedContent, cardDB)
	if err != nil {
		return nil, err
	}
	if r.CVV != nil {
		cardDB.CVV = *r.CVV
	}
	if r.Name != nil {
		cardDB.Name = *r.Name
	}
	if r.Number != nil {
		cardDB.Number = *r.Number
	}
	if r.Expires != nil {
		cardDB.Expires = *r.Expires
	}
	if r.Key != nil {
		cardDB.Key = *r.Key
	}
	return cardDB, nil
}

func (s *cardService) DeleteCard(ctx context.Context, r handlers.DeleteCardRequest) (*handlers.DeleteCardResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = s.repo.Delete(ctx, user, r.UUID)
	if err != nil {
		return nil, err
	}
	return &handlers.DeleteCardResponse{UUID: r.UUID}, nil
}

func (s *cardService) GetCardsByUser(ctx context.Context, r handlers.GetAllCardsRequest) (*handlers.GetAllCardsResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	key, err := s.keyService.GetKeyForUser(user)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	cards, err := s.repo.GetByUser(ctx, user, entity.CARD)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	cardsItems := make([]handlers.GetAllCardsResponceItem, 0, len(cards))
	for _, card := range cards {
		decryptedContent, err := lib.Decrypt(key, card.Content)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}
		cardDB := &CardContent{}
		err = json.Unmarshal(decryptedContent, cardDB)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}
		cardsItems = append(cardsItems, handlers.GetAllCardsResponceItem{
			UUID:    card.UUID,
			Key:     cardDB.Key,
			Number:  cardDB.Number,
			CVV:     cardDB.CVV,
			Name:    cardDB.Name,
			Expires: cardDB.Expires,
		})
	}
	return &handlers.GetAllCardsResponse{Items: cardsItems}, nil
}
