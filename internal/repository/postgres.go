package repository

import (
	"context"
	"database/sql"

	"github.com/Xenios7/Trade-executor/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(p *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: p,
	}
}

func (p *PostgresRepository) Save(order domain.Order) error {
    _, err := p.db.ExecContext(
        context.Background(),
        `INSERT INTO orders (id, asset, side, quantity, price, status, created_at, executed_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
        order.ID,
        order.Asset,
        order.Side,
        order.Quantity,
        order.Price,
        order.Status,
        order.CreatedAt,
        order.ExecutedAt,
    )
    return err
}

func (p *PostgresRepository) GetByID(id string) (domain.Order, error) {
    var order domain.Order
    row := p.db.QueryRowContext(
        context.Background(),
        `SELECT id, asset, side, quantity, price, status, created_at, executed_at 
         FROM orders WHERE id = $1`,
        id,
    )
    err := row.Scan(
        &order.ID,
        &order.Asset,
        &order.Side,
        &order.Quantity,
        &order.Price,
        &order.Status,
        &order.CreatedAt,
        &order.ExecutedAt,
    )
    return order, err
}

func (p *PostgresRepository) GetAll() ([]domain.Order, error) {
    // Query all rows from the orders table
    rows, err := p.db.QueryContext(
        context.Background(),
        `SELECT id, asset, side, quantity, price, status, created_at, executed_at 
         FROM orders`,
    )
    if err != nil {
        return nil, err
    }
    // Close the cursor when done, frees the database connection
    defer rows.Close()

    orders := []domain.Order{}
    // Iterate over each row returned by the query
    for rows.Next() {
        var order domain.Order
        // Scan maps each column value into the corresponding struct field
        if err := rows.Scan(
            &order.ID,
            &order.Asset,
            &order.Side,
            &order.Quantity,
            &order.Price,
            &order.Status,
            &order.CreatedAt,
            &order.ExecutedAt,
        ); err != nil {
            return nil, err
        }
        // Append each scanned order to the result slice
        orders = append(orders, order)
    }
    return orders, nil
}


