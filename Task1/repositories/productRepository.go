package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
	"log"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GET all products
func (repo *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := "SELECT products.id, products.name, price, stock, category_id, categories.name as category_name FROM products JOIN categories on products.category_id = categories.id"

	args := []interface{}{}
	if name != "" {
		query += " WHERE products.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}
	rows, err := repo.db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.Category,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// CREATE product
func (repo *ProductRepository) Create(product *models.Product) error {
	query := `INSERT INTO products (name, price, stock, category_id)
	VALUES ($1, $2, $3, $4)
	RETURNING
		id,
		(SELECT name FROM categories WHERE id = $4) AS category_name`
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID, &product.Category)

	if err != nil {
		log.Println("insert product failed:", err)
		return err
	}
	return nil
}

// GET product by ID
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `SELECT products.id, products.name, price, stock, category_id, categories.name as category_name
	FROM products
	FULL JOIN categories on products.category_id = categories.id
	WHERE products.id = $1`

	product := models.Product{}
	err := repo.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.Category,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// UPDATE product
func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no product updated")
	}

	return nil
}

// DELETE product
func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no product deleted")
	}

	return nil
}
