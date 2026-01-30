package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
	"log"
)


type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository{
	return &CategoryRepository{db: db}
}

// GET all categories
func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := `SELECT id, name, description FROM categories`
	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		category := models.Category{}
		err := rows.Scan(
			&category.ID, &category.Name, &category.Description,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// CREATE category
func (repo *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING ID"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)

	if err != nil {
		log.Println("insert category failed:", err)
		return err
	}
	return nil
}

// GET category By ID
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	category := models.Category{}
	err := repo.db.QueryRow(query, id).Scan(
		&category.ID, &category.Name, &category.Description,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// UPDATE category
func (repo *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return  errors.New("no category updated")
	}

	return  nil
}

// DELETE category
func (repo *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := repo.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no category deleted")
	}

	return nil
}
