package data

import (
	"api.fruitbasket/internals/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Fruit struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	FruitName string    `json:"fruit_name"`
	Price     float64   `json:"price"`
}

func ValidateFruit(v *validator.Validator, fruit *Fruit) {
	v.Check(fruit.FruitName != "", "fruit_name", "must be provided")
	v.Check(len(fruit.FruitName) <= 200, "fruit_name", "must not be more than 200 bytes long")
}

type FruitModel struct {
	DB *sql.DB
}

func (m FruitModel) Insert(fruit *Fruit) error {
	query := `
        INSERT INTO fruits (fruit_name, price) 
        VALUES ($1, $2)
        RETURNING id, created_at`

	args := []interface{}{fruit.FruitName, fruit.Price}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	d := m.DB.QueryRowContext(ctx, query, args...).Scan(&fruit.ID, &fruit.CreatedAt)
	fmt.Println(d)
	return d
}

func (m FruitModel) Get(id int64) (*Fruit, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, fruit_name, price
        FROM fruits
        WHERE id = $1`

	var fruit Fruit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&fruit.ID,
		&fruit.CreatedAt,
		&fruit.FruitName,
		&fruit.Price,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &fruit, nil
}

func (m FruitModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM fruits
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
