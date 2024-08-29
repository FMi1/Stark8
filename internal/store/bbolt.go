package store

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Bbolt struct {
	db *bolt.DB
}

func NewBbolt() (*Bbolt, error) {
	db, err := bolt.Open("stark8.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	b := &Bbolt{db: db}
	err = b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Users"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Bbolt) Connect() error {
	return nil
}

func (b *Bbolt) Close() error {
	return b.db.Close()
}

func (b *Bbolt) GetUser(username string) (User, error) {
	var user User
	err := b.db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte("Users"))
		u := users.Get([]byte(username))
		if u == nil {
			return fmt.Errorf("user not found")
		}
		err := json.Unmarshal(u, &user)
		if err != nil {
			return err
		}
		return nil
	})
	return user, err
}

func (b *Bbolt) CreateUser(up UserParams) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte("Users"))
		u := &User{
			Username:  up.Username,
			Password:  up.Password,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		key := []byte(u.Username)
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}
		if exists := users.Get(key); exists != nil {
			return fmt.Errorf("user already exists")
		}
		return users.Put(key, buf)
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (b *Bbolt) DeleteUser(username string) error {
	return nil // TODO
}
