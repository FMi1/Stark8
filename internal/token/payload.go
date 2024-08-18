package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return errors.New("token is expired")
	}
	return nil
}

func (p *Payload) GetID() string {
	return p.ID.String()
}

func (p *Payload) GetUsername() string {
	return p.Username
}

func (p *Payload) GetIssuedAt() time.Time {
	return p.IssuedAt
}

func (p *Payload) GetExpiredAt() time.Time {
	return p.ExpiredAt
}
