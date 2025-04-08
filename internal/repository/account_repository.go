package repository

import (
	"database/sql"
	"time"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

type AccountRepository struct {
	connection *sql.DB
}

func NewAccountRepository(connection *sql.DB) *AccountRepository {
	return &AccountRepository{connection: connection}
}

func (repository *AccountRepository) Save(account *domain.Account) error {
	stmt, err := repository.connection.Prepare(`
	INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(account.Id, account.Name, account.Email, account.ApiKey, account.Balance, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (repository *AccountRepository) FindByApiKey(apiKey string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := repository.connection.QueryRow(`
	SELECT id, name, email, api_key, balance, created_at, updated_at
	FROM accounts
	WHERE api_key = $1
	`, apiKey).Scan(&account.Id, &account.Name, &account.Email, &account.ApiKey, &account.Balance, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt

	return &account, nil
}

func (repository *AccountRepository) FindById(id string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := repository.connection.QueryRow(`
	SELECT id, name, email, api_key, balance, created_at, updated_at
	FROM accounts
	WHERE id = $1
	`, id).Scan(&account.Id, &account.Name, &account.Email, &account.ApiKey, &account.Balance, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt

	return &account, nil
}

func (repository *AccountRepository) UpdateBalance(account *domain.Account) error {
	tx, err := repository.connection.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	var currentBalance float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, account.Id).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		return domain.ErrAccountNotFound
	}
	if err != nil {
		return err
	}

	newBalance := currentBalance + account.Balance

	_, err = tx.Exec(`UPDATE accounts SET balance = $1, updated_at = $2 WHERE id = $3`, newBalance, time.Now(), account.Id)

	if err != nil {
		return err
	}

	return tx.Commit()
}
