package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/tocy1/toggl/api/model"
)

func (h *APIHandler) CreateDeck(w http.ResponseWriter, r *http.Request) {

	id := uuid.New()
	deck := model.Deck{}
	cards_string := r.URL.Query().Get("cards")
	shuffled_string := r.URL.Query().Get("shuffled")
	var shuffled bool
	var err error
	if shuffled_string == "" {
		shuffled = false
	} else {
		shuffled, err = strconv.ParseBool(shuffled_string)
		if err != nil {
			h.Logger.Error("Unable to convert shuffled param to bool", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	deck.Shuffled = shuffled
	deck.Id = id
	cards := []model.Card{}
	input_cards := ""
	if cards_string == "" {
		input_cards += model.DefaultSequentialCards()
	} else {
		input_cards += cards_string
	}
	curr_cards, err := model.GetCardsFromString(input_cards)
	if err != nil {
		h.Logger.Error("Unable to retrieve card object from given param", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cards = append(cards, curr_cards...)
	if shuffled == true {
		model.ShuffledCards(cards)
	}
	deck.Cards = model.GetStringFromCard(cards)
	deck.Remaining = len(cards)
	h.db.CreateDeck(r.Context(), deck)

	responseBody, _ := json.Marshal(deck)
	_, err = w.Write(responseBody)
	if err != nil {
		h.Logger.Error("CreateDeck: error writing response", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (h *APIHandler) OpenDeck(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		h.Logger.Error("No UUID was passed", errors.New("Unknown UUID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	deck, err := h.db.GetDeck(r.Context(), uuid)
	if err != nil {
		h.Logger.Error("OpenDeck: error retrieving deck from db.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cards, err := model.GetCardsFromString(deck.Cards)
	if err != nil {
		h.Logger.Error("Unable to retrieve card object from given param", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := model.OpenDeckResponse{
		Id:        deck.Id,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
		Cards:     cards[len(cards)-deck.Remaining:],
	}
	if deck.Remaining != len(cards) {
		deck.Remaining = len(cards)
		err = h.db.UpdateDeck(r.Context(), deck)
		if err != nil {
			h.Logger.Error("OpenDeck: error updating deck", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	responseBody, _ := json.Marshal(result)
	_, err = w.Write(responseBody)

	if err != nil {
		h.Logger.Error("OpenDeck: error writing response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *APIHandler) DrawCards(w http.ResponseWriter, r *http.Request) {

	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		h.Logger.Error("DrawCards: No count parameter was passed or unexpected value", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		h.Logger.Error("DrawCards: No UUID was passed", errors.New("Unknown UUID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	deck, err := h.db.GetDeck(r.Context(), uuid)
	if err != nil {
		h.Logger.Error("DrawCards: error retrieving deck from db.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cards, err := model.GetCardsFromString(deck.Cards)
	if err != nil {
		h.Logger.Error("DrawCards: Unable to retrieve card object from given param", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	start := len(cards) - deck.Remaining
	end := start + count
	if end > len(cards) {
		deck.Remaining = 0
	} else {
		deck.Remaining = len(cards) - end
	}
	err = h.db.UpdateDeck(r.Context(), deck)
	if err != nil {
		h.Logger.Error("DrawCards: error updating deck", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := model.DrawCardsResponse{Cards: cards[start:end]}
	responseBody, _ := json.Marshal(result)
	_, err = w.Write(responseBody)

	if err != nil {
		h.Logger.Error("DrawCards: error writing response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
