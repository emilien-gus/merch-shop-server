package models

type InfoResponse struct {
	CoinsCount  int         `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
