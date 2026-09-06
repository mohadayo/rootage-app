package repository

import (
	"context"
	"database/sql"

	"github.com/rootage-ses-quiz/backend/internal/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) DB() *sql.DB {
	return r.db
}

func (r *CategoryRepository) List(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, created_at, updated_at FROM categories WHERE is_active ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*model.Category, error) {
	var c model.Category
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(ctx context.Context, cat *model.Category) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		cat.Name, cat.Description,
	).Scan(&cat.ID, &cat.CreatedAt, &cat.UpdatedAt)
}

func (r *CategoryRepository) Update(ctx context.Context, cat *model.Category) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE categories SET name = $1, description = $2, updated_at = NOW() WHERE id = $3`,
		cat.Name, cat.Description, cat.ID,
	)
	return err
}

// Deactivate はカテゴリを論理削除する（出題・一覧から外すが履歴は残す）。
func (r *CategoryRepository) Deactivate(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE categories SET is_active = FALSE, updated_at = NOW() WHERE id = $1`, id)
	return err
}
