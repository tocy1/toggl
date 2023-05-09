package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func DefaultSequentialCards() string {
	result := "AS"
	count := 0
	for _, suit := range SuitsOrder {
		for _, rank := range RanksOrder {
			if count > 0 {
				result = fmt.Sprintf("%s,%s%s", result, rank, suit)
			}
			count++
		}

	}
	return result
}

func GetCardsFromString(str string) (cards []Card, err error) {
	given_cards := strings.Split(str, ",")
	if len(given_cards) < 1 {
		return nil, fmt.Errorf("Invalid string passed: Length of string less than 1.")
	}
	for _, card := range given_cards {
		curr_card := Card{}
		suit, exists := Suits[string(card[1])]
		if exists {
			curr_card.Suit = suit
		} else {
			return nil, fmt.Errorf("Invalid string passed: The given suit code: %s, is not recognized", string(card[0]))
		}
		rank, exists := Ranks[string(card[0])]
		if exists {
			curr_card.Value = rank
			curr_card.Code = card
		} else {
			return nil, fmt.Errorf("Invalid string passed: The given rank code: %s, is not recognized", string(card[0]))
		}
		cards = append(cards, curr_card)
	}
	return cards, nil

}

func ShuffledCards(cards []Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func GetStringFromCard(cards []Card) string {
	result := ""
	count := 0
	for _, card := range cards {
		if count > 0 {
			result = fmt.Sprintf("%s,%s", result, card.Code)
		} else {
			result = result + card.Code
		}

		count++
	}
	return result
}
