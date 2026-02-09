package transaction

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	CreateTransaction(items []CheckoutItem) (*Transaction, error)
	GetDailySalesReport() (*DailySalesReport, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateTransaction(items []CheckoutItem) (*Transaction, error) {
	// Begin database transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]TransactionDetail, 0)

	// Process each item
	for _, item := range items {
		var productPrice, stock int
		var productName string

		// Get product info and check stock
		err := tx.QueryRow("SELECT nama, harga, stok FROM products WHERE id = $1", item.ProductID).
			Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Validate stock
		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s (available: %d, requested: %d)",
				productName, stock, item.Quantity)
		}

		// Calculate subtotal
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// Update product stock
		_, err = tx.Exec("UPDATE products SET stok = stok - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// Prepare detail
		details = append(details, TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Insert transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).
		Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Insert transaction details
	for i := range details {
		details[i].TransactionID = transactionID
		var detailID int
		err = tx.QueryRow(
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal,
		).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Get created_at timestamp
	var createdAt sql.NullTime
	err = r.db.QueryRow("SELECT created_at FROM transactions WHERE id = $1", transactionID).Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt.Time,
		Details:     details,
	}, nil
}

func (r *repository) GetDailySalesReport() (*DailySalesReport, error) {
	report := &DailySalesReport{}

	// Get total revenue and transaction count for today
	err := r.db.QueryRow(`
		SELECT
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as total_transaksi
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// Get top selling product for today
	var productName sql.NullString
	var qtyTerjual sql.NullInt64

	err = r.db.QueryRow(`
		SELECT
			p.nama,
			COALESCE(SUM(td.quantity), 0) as qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.id, p.nama
		ORDER BY qty DESC
		LIMIT 1
	`).Scan(&productName, &qtyTerjual)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Set top product if exists
	if productName.Valid {
		report.ProdukTerlaris = &TopProduct{
			Nama:       productName.String,
			QtyTerjual: int(qtyTerjual.Int64),
		}
	}

	return report, nil
}
