package product

import (
	"database/sql"
	"time"
)

type Product struct {
	ID         int       `json:"id"`
	Nama       string    `json:"nama"`
	Harga      int       `json:"harga"`
	Stok       int       `json:"stok"`
	CategoryID *int      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProductDetail struct {
	ID           int       `json:"id"`
	Nama         string    `json:"nama"`
	Harga        int       `json:"harga"`
	Stok         int       `json:"stok"`
	CategoryID   *int      `json:"category_id"`
	CategoryName *string   `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Nama       string `json:"nama"`
	Harga      int    `json:"harga"`
	Stok       int    `json:"stok"`
	CategoryID *int   `json:"category_id"`
}

type UpdateProductRequest struct {
	Nama       string `json:"nama"`
	Harga      int    `json:"harga"`
	Stok       int    `json:"stok"`
	CategoryID *int   `json:"category_id"`
}

func (p *ProductDetail) ScanRow(rows *sql.Rows) error {
	var categoryName sql.NullString
	err := rows.Scan(
		&p.ID,
		&p.Nama,
		&p.Harga,
		&p.Stok,
		&p.CategoryID,
		&categoryName,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if categoryName.Valid {
		p.CategoryName = &categoryName.String
	}

	return nil
}
