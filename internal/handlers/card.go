package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/websocket"
)


type CardHandler struct {
	CardService services.CardService
	AuthService services.AuthService
	WSHub       *websocket.Hub
}

func NewCardHandler(cardService *services.CardService, authService *services.AuthService, wsHub *websocket.Hub) *CardHandler {
	return &CardHandler{
		CardService: *cardService,
		AuthService: *authService,
		WSHub:       wsHub,
	}
}

func (ch *CardHandler) broadcastTaskEvent(eventType string, payload interface{}) {
    if ch.WSHub == nil {
        return
    }
    message := map[string]interface{}{
        "type":    eventType,
        "payload": payload,
    }
    
    msgBytes, err := json.Marshal(message)
    if err != nil {
        return
    }
    
    ch.WSHub.BroadcastMessage(msgBytes)
}

// CreateCard обрабатывает POST /board/{boardId}/list/{listId}/tasks
func (ch *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()

    // ПОЛУЧИ ТЕКУЩЕГО ПОЛЬЗОВАТЕЛЯ
    user, err := ch.GetUserFromRequest(ctx, r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    listId, err := strconv.ParseInt(r.PathValue("listId"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid list ID", http.StatusBadRequest)
        return
    }

    // ПОЛУЧИ boardId из URL
    boardId, err := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid board ID", http.StatusBadRequest)
        return
    }

    // НАЙДИ board_member_id ДЛЯ ЭТОГО ПОЛЬЗОВАТЕЛЯ И ДОСКИ
    boardMemberId, err := ch.findBoardMemberId(ctx, user.ID, boardId)
    if err != nil {
        http.Error(w, "User is not a member of this board", http.StatusForbidden)
        return
    }

    decoder := json.NewDecoder(r.Body)
    newCard := new(models.Card)
    err = decoder.Decode(newCard)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // ИСПОЛЬЗУЙ НАЙДЕННЫЙ board_member_id
    newCard.AuthorBoardMemberId = boardMemberId

    createdCard, err := ch.CardService.CreateCard(ctx, newCard, listId)
    if err != nil {
        log.Printf("Error creating card: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    ch.broadcastTaskEvent("TASK_CREATED", map[string]interface{}{
        "id":          createdCard.ID,
        "content":     createdCard.Content,
        "listId":      createdCard.ListId,
        "position":    createdCard.Position,
        "isCompleted": createdCard.Completed,
        "boardId":     boardId,
    })

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdCard)
}

// ДОБАВЬ МЕТОД ДЛЯ ПОИСКА board_member_id
func (ch *CardHandler) findBoardMemberId(ctx context.Context, userId, boardId int64) (int64, error) {
    // Здесь нужно реализовать логику поиска board_member_id
    // по userId и boardId через репозиторий
    
    // ВРЕМЕННАЯ ЗАГЛУШКА
    if userId == 1 && boardId == 1 {
        return 3, nil // Sergey на доске 1
    }
    if userId == 2 && boardId == 1 {
        return 4, nil // Sergey2 на доске 1  
    }
    
    return 0, nil
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

	_, err = ch.CardService.GetCard(ctx, cardID, user.ID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	updatedCard, err := ch.CardService.UpdateCard(ctx, card, cardID)
	if err != nil {
		log.Printf("Error updating card: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ch.broadcastTaskEvent("TASK_UPDATED", map[string]interface{}{
		"id":          updatedCard.ID,
		"content":     updatedCard.Content,
		"listId":      updatedCard.ListId,
		"isCompleted": updatedCard.Completed,
		"boardId":     r.PathValue("boardId"),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCard)
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

	existingCard, err := ch.CardService.GetCard(ctx, cardID, user.ID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	if err := ch.CardService.DeleteCard(ctx, cardID); err != nil {
		log.Printf("Error deleting card: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ch.broadcastTaskEvent("TASK_DELETED", map[string]interface{}{
		"taskId":  cardID,
		"listId":  existingCard.ListId,
		"boardId": r.PathValue("boardId"),
	})

	w.WriteHeader(http.StatusNoContent)
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
	json.NewEncoder(w).Encode(card)
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
	json.NewEncoder(w).Encode(cards)
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