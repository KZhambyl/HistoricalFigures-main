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

type Figure struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	YearsOfLife string    `json:"years_of_life"`
	Description string    `json:"description"`
	Version     int32     `json:"version"`
}

func ValidateFigure(v *validator.Validator, f *Figure) {
	v.Check(f.Name != "", "name", "must be provided")
	v.Check(len(f.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(len(f.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	v.Check(len(f.YearsOfLife) <= 9, "years_of_life", "must not be more than 10 bytes long")
	v.Check(len(f.YearsOfLife) > 0, "years_of_life", "must not be empty")
}

type FigureModel struct {
	DB *sql.DB
}

func (m FigureModel) Insert(figure *Figure) error {
	query := `
	INSERT INTO figures (name, years_of_life, description)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version`

	args := []interface{}{figure.Name, figure.YearsOfLife, figure.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&figure.ID, &figure.CreatedAt, &figure.Version)
}

func (m FigureModel) Get(id int64) (*Figure, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, name, years_of_life, description, version
	FROM figures WHERE id = $1`

	var figure Figure

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&figure.ID,
		&figure.CreatedAt,
		&figure.Name,
		&figure.YearsOfLife,
		&figure.Description,
		&figure.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &figure, nil
}

func (m FigureModel) Update(figure *Figure) error {
	query := `
	UPDATE figures
	SET name = $1, years_of_life = $2, description = $3, version = version + 1 WHERE id = $4 AND version = $5
	RETURNING version`
	args := []interface{}{
		figure.Name,
		figure.YearsOfLife,
		figure.Description,
		figure.ID,
		figure.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&figure.Version)
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

func (m FigureModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM figures
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

func (m FigureModel) GetAll(name string, years_of_life string, filters Filters) ([]*Figure, Metadata, error) {
	query := fmt.Sprintf(`SELECT count(*) OVER(), id, created_at, name, years_of_life, description, version
	FROM figures
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) or $1 = '') 
	AND (years_of_life = $2 OR $2 = '')
	ORDER BY %s %s, id ASC 
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, years_of_life, filters.limit(), filters.offset()}

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
