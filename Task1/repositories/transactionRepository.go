package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// inisiasi subtotal -> jumlah total transaksi keseluruhan
	totalAmount := 0
	//  inisiasi modeling transactionDetails -> nanti insert ke db
	details := make([]models.TransactionDetail, 0)

	// loop setiap item
	for _, item := range items {
		var productPrice, stock int
		var productName string
		// get product dapet pricing
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1 FOR UPDATE", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}
		// hitung current total * pricing
		// ditambahin ke subtotal
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock product")
		}

		// kurangi jumlah stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// itemnya dimasukin transactionDetails
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID, detailID int
	var createdAt time.Time
	// insert transaction
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at", totalAmount).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	// insert transaction details
	for x := range details {
		details[x].TransactionID = transactionID

		err = tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[x].ProductID, details[x].Quantity, details[x].Subtotal).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		// detailID = details[x].ID
		details[x].ID = detailID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt,
		Details:     details,
	}, nil
}
