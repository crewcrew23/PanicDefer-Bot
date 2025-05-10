package requestmodel

type RequestCommand struct {
	Command string `json:"command"`
	Value   string `json:"value"`
	ChatID  int64  `json:"chat_id"`
}
