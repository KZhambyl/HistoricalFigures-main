package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/KZhambyl/HistoricalFigures/internal/validator"
	_ "github.com/lib/pq"
	"time"
)

type Category struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Version   int32     `json:"version"`
}

func ValidateCategory(v *validator.Validator, f *Category) {
	v.Check(f.Name != "", "name", "must be provided")
	v.Check(len(f.Name) <= 15, "name", "must not be more than 15 bytes long")
}

type CategoryModel struct {
	DB *sql.DB
}

func (m CategoryModel) Insert(category *Category) error {
	query := `
	INSERT INTO categories (name)
	VALUES ($1)
	RETURNING id, created_at, version`

	args := []interface{}{category.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&category.ID, &category.CreatedAt, &category.Version)
}

func (m CategoryModel) Get(id int64) (*Category, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, name, version
	FROM categories WHERE id = $1`

	var category Category

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Name,
		&category.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (m CategoryModel) Update(category *Category) error {
	query := `
	UPDATE categories
	SET name = $1, version = version + 1 WHERE id = $2 AND version = $3
	RETURNING version`
	args := []interface{}{
		category.Name,
		category.ID,
		category.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&category.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m CategoryModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM categories
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

func (m CategoryModel) GetAll(name string, filters Filters) ([]*Category, Metadata, error) {
	query := fmt.Sprintf(`SELECT count(*) OVER(), id, created_at, name, version
	FROM categories
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) or $1 = '') 
	ORDER BY %s %s, id ASC 
	LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	categories := []*Category{}

	for rows.Next() {
		var category Category

		err := rows.Scan(
			&totalRecords,
			&category.ID,
			&category.CreatedAt,
			&category.Name,
			&category.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		categories = append(categories, &category)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return categories, metadata, nil

}

func (m CategoryModel) GetCategoryFigures(name string, years_of_life string, filters Filters, id int64) ([]*Figure, Metadata, error) {
	query := fmt.Sprintf(`SELECT count(*) OVER(), figures.id, figures.created_at, figures.name, figures.years_of_life, figures.description, figures.version
	FROM figures_categories 
	join figures on figures.id = figures_categories.figure_id 
	WHERE (to_tsvector('simple', figures.name) @@ plainto_tsquery('simple', $1) or $1 = '') 
	AND (figures.years_of_life = $2 OR $2 = '')
	AND figures_categories.category_id=$5
	ORDER BY %s %s, id ASC 
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, years_of_life, filters.limit(), filters.offset(), id}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	figures := []*Figure{}

	for rows.Next() {
		var figure Figure

		err := rows.Scan(
			&totalRecords,
			&figure.ID,
			&figure.CreatedAt,
			&figure.Name,
			&figure.YearsOfLife,
			&figure.Description,
			&figure.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		figures = append(figures, &figure)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return figures, metadata, nil
}
