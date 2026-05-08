package events

type UserCreatedEvent struct {
	ID    string `json:"user_id"`
	Email string `json:"email"`
}
