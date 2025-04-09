package repository

import (
	"database/sql"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

type InvoiceRepository struct {
	connection *sql.DB
}

func NewInvoiceRepository(connection *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{connection: connection}
}

func (repository *InvoiceRepository) Save(invoice *domain.Invoice) error {
	stmt, err := repository.connection.Prepare(`
		INSERT INTO invoices (id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(invoice.Id, invoice.AccountId, invoice.Amount, invoice.Status, invoice.Description, invoice.PaymentType, invoice.CardLastDigits, invoice.CreatedAt, invoice.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (repository *InvoiceRepository) FindById(id string) (*domain.Invoice, error) {
	stmt, err := repository.connection.Prepare(`
		SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE id = $1
	`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(id)

	var invoice domain.Invoice

	err = row.Scan(&invoice.Id, &invoice.AccountId, &invoice.Amount, &invoice.Status, &invoice.Description, &invoice.PaymentType, &invoice.CardLastDigits, &invoice.CreatedAt, &invoice.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrInvoiceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (repository *InvoiceRepository) FindByAccountId(accountId string) ([]*domain.Invoice, error) {
	stmt, err := repository.connection.Prepare(`
		SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE account_id = $1
	`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(accountId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var invoices []*domain.Invoice

	for rows.Next() {
		var invoice domain.Invoice

		err = rows.Scan(&invoice.Id, &invoice.AccountId, &invoice.Amount, &invoice.Status, &invoice.Description, &invoice.PaymentType, &invoice.CardLastDigits, &invoice.CreatedAt, &invoice.UpdatedAt)

		if err != nil {
			return nil, err
		}

		invoices = append(invoices, &invoice)
	}

	return invoices, nil
}

func (repository *InvoiceRepository) UpdateStatus(invoice *domain.Invoice) error {
	stmt, err := repository.connection.Prepare(`
		UPDATE invoices SET status = $1, updated_at = $2 WHERE id = $3
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	rows, err := stmt.Exec(invoice.Status, invoice.UpdatedAt, invoice.Id)

	if err != nil {
		return err
	}

	affectedRows, err := rows.RowsAffected()

	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return domain.ErrInvoiceNotFound
	}

	return nil
}
