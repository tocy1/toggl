package model

import (
	"github.com/google/uuid"
)

type Deck struct {
	Id        uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     string    `json:"-"`
}

type Card struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type OpenDeckResponse struct {
	Id        uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     []Card    `json:"cards"`
}

type DrawCardsResponse struct {
	Cards []Card `json:"cards"`
}

var (
	Ranks = map[string]string{
		"A": "ACE",
		"2": "2",
		"3": "3",
		"4": "4",
		"5": "5",
		"6": "6",
		"7": "7",
		"8": "8",
		"9": "9",
		"T": "10",
		"K": "KING",
		"Q": "QUEEN",
		"J": "JACK",
	}
	RanksOrder = []string{
		"A",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"T",
		"K",
		"Q",
		"J",
	}
	Suits = map[string]string{
		"S": "SPADES",
		"D": "DIAMONDS",
		"C": "CLUBS",
		"H": "HEARTS",
	}
	SuitsOrder = []string{
		"S",
		"D",
		"C",
		"H",
	}
)
