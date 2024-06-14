package models

type Message struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Text     *string `json:"text,omitempty"`
}

type StatsRequest struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

type Response struct {
	Data    ResData `json:"data"`
	Message string  `json:"message"`
}

type ResData struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Count    int64  `json:"count"`
}
