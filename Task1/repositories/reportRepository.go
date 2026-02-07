package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db}
}

func (repo *ReportRepository) GetReport() (*models.Report, error) {
	// Query total revenue hari ini
	var totalSelling *int
	err := repo.db.QueryRow(
		"SELECT SUM(total_amount) FROM transactions WHERE DATE(created_at) = CURRENT_DATE",
	).Scan(&totalSelling)
	if err != nil {
		return nil, err
	}

	totalRevenue := 0
	if totalSelling != nil {
		totalRevenue = *totalSelling
	}

	// Query total transaksi hari ini
	var totalTransaction int
	err = repo.db.QueryRow(
		"SELECT COUNT(*) FROM transactions WHERE DATE(created_at) = CURRENT_DATE",
	).Scan(&totalTransaction)
	if err != nil {
		return nil, err
	}

	// Query produk terlaris hari ini
	var name string
	var qtyTerjual int
	err = repo.db.QueryRow(`
		SELECT p.name, SUM(td.quantity) AS qty_terjual FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
		`).
		Scan(&name, &qtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &models.Report{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaction,
		ProdukTerlaris: models.ProdukTerlaris{
			Name:       name,
			QtyTerjual: qtyTerjual,
		},
	}, nil
}
