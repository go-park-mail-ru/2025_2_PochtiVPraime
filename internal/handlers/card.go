package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

type CardHandler struct {
	CardService services.CardService
	AuthService services.AuthService
}

func NewCardHandler(cardService *services.CardService, authService *services.AuthService) *CardHandler {
	return &CardHandler{
		CardService: *cardService,
		AuthService: *authService,
	}
}

// GetUserFromRequest извлекает пользователя из токена
func (ch *CardHandler) GetUserFromRequest(ctx context.Context, r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	tokenString := cookie.Value
	return ch.AuthService.GetUserFromToken(ctx, tokenString)
}

// CreateCard обрабатывает POST /board/{boardId}/list/{listId}/tasks
func (ch *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	/*
		user, err := ch.GetUserFromRequest(ctx, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	*/

	// Получаем listID из URL параметров
	listId, err := strconv.ParseInt(r.PathValue("listId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	newCard := new(models.Card)
	err = decoder.Decode(newCard)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Println(newCard)
	// Устанавливаем создателя карточки
	zaglushka := 1
	newCard.AuthorBoardMemberId = int64(zaglushka) //реализовано пока так, тк нет реализации board_member

	createdCard, err := ch.CardService.CreateCard(ctx, newCard, listId)
	if err != nil {
		log.Printf("Error creating card: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdCard); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetCard обрабатывает GET /board/{boardId}/list/{listId}/task/{taskId}
func (ch *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := ch.GetUserFromRequest(ctx, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cardID, err := strconv.ParseInt(r.PathValue("taskId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	card, err := ch.CardService.GetCard(ctx, cardID, user.ID)
	if err != nil {
		log.Printf("Error getting card: %v", err)
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(card); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetListCards обрабатывает GET /board/{boardId}/list/{listId}/tasks
func (ch *CardHandler) GetListCards(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	listID, err := strconv.ParseInt(r.PathValue("listId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	cards, err := ch.CardService.GetListCards(ctx, listID)
	if err != nil {
		log.Printf("Error getting list cards: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cards); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateCard обрабатывает PUT /board/{boardId}/list/{listId}/task/{taskId}
func (ch *CardHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := ch.GetUserFromRequest(ctx, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cardID, err := strconv.ParseInt(r.PathValue("taskId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	card := new(models.Card)
	err = decoder.Decode(card)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//log.Println(card)

	// Проверяем, что пользователь имеет доступ к карточке
	_, err = ch.CardService.GetCard(ctx, cardID, user.ID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	updatedCard, err := ch.CardService.UpdateCard(ctx, card, cardID)
	_ = updatedCard
	if err != nil {
		log.Printf("Error updating card: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	card.ID = cardID
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(card); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteCard обрабатывает DELETE /board/{boardId}/list/{listId}/task/{taskId}
func (ch *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := ch.GetUserFromRequest(ctx, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cardID, err := strconv.ParseInt(r.PathValue("taskId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Проверяем, что карточка существует и пользователь имеет к ней доступ
	_, err = ch.CardService.GetCard(ctx, cardID, user.ID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	if err := ch.CardService.DeleteCard(ctx, cardID); err != nil {
		log.Printf("Error deleting card: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ch *CardHandler) Card(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ch.GetCard(w, r)
	case http.MethodDelete:
		ch.DeleteCard(w, r)
	case http.MethodPut:
		ch.UpdateCard(w, r)
	}
}

func (ch *CardHandler) CreateOrGetCards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ch.GetListCards(w, r)
	case http.MethodPost:
		ch.CreateCard(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
