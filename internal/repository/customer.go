package repository

import (
	"context"
	"errors"

	"github.com/antoniooreany/subscription-billing-simulator/internal/model"
	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepository struct{ db *pgxpool.Pool }

func NewCustomerRepository(db *pgxpool.Pool) *CustomerRepository { return &CustomerRepository{db: db} }

func (r *CustomerRepository) Create(ctx context.Context, c model.Customer) (model.Customer, error) {
	_, err := r.db.Exec(ctx, `INSERT INTO customers (id, email, name) VALUES ($1, $2, $3)`, c.ID, c.Email, c.Name)
	return c, err
}

func (r *CustomerRepository) GetByID(ctx context.Context, id string) (model.Customer, error) {
	var c model.Customer
	err := r.db.QueryRow(ctx, `SELECT id, email, name, created_at FROM customers WHERE id = $1`, id).Scan(&c.ID, &c.Email, &c.Name, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Customer{}, shared.ErrNotFound
		}
		return model.Customer{}, err
	}
	return c, nil
}
