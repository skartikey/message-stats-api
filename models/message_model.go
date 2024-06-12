package models

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Text     string `json:"text"`
}

type Response struct {
	Status  string        `json:"status"`
	Data    ResData       `json:"data"`
	Message string        `json:"message"`
	Error   ErrorResponse `json:"error"`
}

type ResData struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Count    int    `json:"count"`
}

type ErrorResponse struct {
	Code        int    `json:"code"`
	Description string `json:"message"`
}
