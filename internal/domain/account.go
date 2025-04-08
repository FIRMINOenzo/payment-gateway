package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id        string
	Name      string
	Email     string
	ApiKey    string
	Balance   float64
	mu        sync.RWMutex
	CreatedAt time.Time
	UpdatedAt time.Time
}

func generateApiKey() string {
	byteArray := make([]byte, 16)
	rand.Read(byteArray)
	return hex.EncodeToString(byteArray)
}

func NewAccount(name string, email string) *Account {
	account := &Account{
		Id:        uuid.New().String(),
		Name:      name,
		Email:     email,
		ApiKey:    generateApiKey(),
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account
}

func (account *Account) UpdateBalance(amount float64) {
	account.mu.Lock()
	defer account.mu.Unlock()

	account.Balance += amount
	account.UpdatedAt = time.Now()
}
