package token

import (
	"errors"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// Maker is a struct that implements the TokenMaker interface
type PasetoMaker struct {
	paseto       *paseto.V2
	simmetricKey []byte
}

func NewPasetoMaker(simmetricKey string) (Maker, error) {
	if len(simmetricKey) != chacha20poly1305.KeySize {
		return nil, errors.New("invalid key size: must be exactly 32 bytes")
	}
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		simmetricKey: []byte(simmetricKey),
	}, nil
}

func (m *PasetoMaker) CreateToken(user string, duration time.Duration) (string, error) {
	Payload, err := NewPayload(user, duration)
	if err != nil {
		return "", err
	}
	return m.paseto.Encrypt(m.simmetricKey, Payload, nil)
}
func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	Payload := &Payload{}
	err := m.paseto.Decrypt(token, m.simmetricKey, Payload, nil)
	if err != nil {
		return nil, err
	}
	err = Payload.Valid()
	if err != nil {
		return nil, err
	}
	return Payload, nil
}
