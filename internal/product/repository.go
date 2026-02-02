package product

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	GetAll() ([]ProductDetail, error)
	GetByID(id int) (*ProductDetail, error)
	Create(req CreateProductRequest) (*Product, error)
	Update(id int, req UpdateProductRequest) (*Product, error)
	Delete(id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]ProductDetail, error) {
	query := `
		SELECT
			p.id,
			p.nama,
			p.harga,
			p.stok,
			p.category_id,
			c.name as category_name,
			p.created_at,
			p.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []ProductDetail
	for rows.Next() {
		var prod ProductDetail
		if err := prod.ScanRow(rows); err != nil {
			return nil, err
		}
		products = append(products, prod)
	}

	return products, nil
}

func (r *repository) GetByID(id int) (*ProductDetail, error) {
	query := `
		SELECT
			p.id,
			p.nama,
			p.harga,
			p.stok,
			p.category_id,
			c.name as category_name,
			p.created_at,
			p.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("product not found")
	}

	var prod ProductDetail
	if err := prod.ScanRow(rows); err != nil {
		return nil, err
	}

	return &prod, nil
}

func (r *repository) Create(req CreateProductRequest) (*Product, error) {
	query := `
		INSERT INTO products (nama, harga, stok, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, nama, harga, stok, category_id, created_at, updated_at
	`

	var prod Product
	err := r.db.QueryRow(query, req.Nama, req.Harga, req.Stok, req.CategoryID).Scan(
		&prod.ID, &prod.Nama, &prod.Harga, &prod.Stok, &prod.CategoryID, &prod.CreatedAt, &prod.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &prod, nil
}

func (r *repository) Update(id int, req UpdateProductRequest) (*Product, error) {
	query := `
		UPDATE products
		SET nama = $1, harga = $2, stok = $3, category_id = $4
		WHERE id = $5
		RETURNING id, nama, harga, stok, category_id, created_at, updated_at
	`

	var prod Product
	err := r.db.QueryRow(query, req.Nama, req.Harga, req.Stok, req.CategoryID, id).Scan(
		&prod.ID, &prod.Nama, &prod.Harga, &prod.Stok, &prod.CategoryID, &prod.CreatedAt, &prod.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, err
	}

	return &prod, nil
}

func (r *repository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
