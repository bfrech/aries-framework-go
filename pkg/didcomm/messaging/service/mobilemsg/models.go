package mobilemsg

import "time"

// Message is message model for mobile messaging protocol
type Message struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	CreatedTime time.Time `json:"created_time"`
	Goal        string    `json:"goal"`
	To          string    `json:"to"`
	Body        struct {
		Content string `json:"content"`
	} `json:"body"`
}
