package category

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	GetAll() ([]Category, error)
	GetByID(id int) (*Category, error)
	Create(req CreateCategoryRequest) (*Category, error)
	Update(id int, req UpdateCategoryRequest) (*Category, error)
	Delete(id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY id ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *repository) GetByID(id int) (*Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`

	var cat Category
	err := r.db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (r *repository) Create(req CreateCategoryRequest) (*Category, error) {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at`

	var cat Category
	err := r.db.QueryRow(query, req.Name, req.Description).Scan(
		&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (r *repository) Update(id int, req UpdateCategoryRequest) (*Category, error) {
	query := `UPDATE categories SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description, created_at, updated_at`

	var cat Category
	err := r.db.QueryRow(query, req.Name, req.Description, id).Scan(
		&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (r *repository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
