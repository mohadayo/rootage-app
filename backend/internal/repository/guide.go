package repository

import (
	"context"
	"database/sql"

	"github.com/rootage-ses-quiz/backend/internal/model"
)

type GuideRepository struct {
	db *sql.DB
}

func NewGuideRepository(db *sql.DB) *GuideRepository {
	return &GuideRepository{db: db}
}

// GuideCategory CRUD

func (r *GuideRepository) ListCategories(ctx context.Context) ([]model.GuideCategory, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, sort_order, created_at FROM guide_categories ORDER BY sort_order, created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.GuideCategory
	for rows.Next() {
		var c model.GuideCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.SortOrder, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *GuideRepository) GetCategoryByID(ctx context.Context, id string) (*model.GuideCategory, error) {
	var c model.GuideCategory
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, sort_order, created_at FROM guide_categories WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.SortOrder, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *GuideRepository) CreateCategory(ctx context.Context, c *model.GuideCategory) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO guide_categories (name, sort_order) VALUES ($1, $2) RETURNING id, created_at`,
		c.Name, c.SortOrder,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *GuideRepository) UpdateCategory(ctx context.Context, c *model.GuideCategory) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE guide_categories SET name = $1, sort_order = $2 WHERE id = $3`,
		c.Name, c.SortOrder, c.ID,
	)
	return err
}

func (r *GuideRepository) DeleteCategory(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM guide_categories WHERE id = $1`, id)
	return err
}

// Guide CRUD

func (r *GuideRepository) ListAll(ctx context.Context) ([]model.Guide, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT g.id, g.category_id, gc.name, g.title, g.content, g.is_published, g.sort_order, g.created_at, g.updated_at
		 FROM guides g JOIN guide_categories gc ON g.category_id = gc.id
		 ORDER BY gc.sort_order, g.sort_order, g.created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGuides(rows)
}

func (r *GuideRepository) ListPublished(ctx context.Context) ([]model.Guide, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT g.id, g.category_id, gc.name, g.title, g.content, g.is_published, g.sort_order, g.created_at, g.updated_at
		 FROM guides g JOIN guide_categories gc ON g.category_id = gc.id
		 WHERE g.is_published = true
		 ORDER BY gc.sort_order, g.sort_order, g.created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGuides(rows)
}

func (r *GuideRepository) GetByID(ctx context.Context, id string) (*model.Guide, error) {
	var g model.Guide
	err := r.db.QueryRowContext(ctx,
		`SELECT g.id, g.category_id, gc.name, g.title, g.content, g.is_published, g.sort_order, g.created_at, g.updated_at
		 FROM guides g JOIN guide_categories gc ON g.category_id = gc.id
		 WHERE g.id = $1`, id,
	).Scan(&g.ID, &g.CategoryID, &g.CategoryName, &g.Title, &g.Content, &g.IsPublished, &g.SortOrder, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GuideRepository) Create(ctx context.Context, g *model.Guide) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO guides (category_id, title, content, is_published, sort_order)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`,
		g.CategoryID, g.Title, g.Content, g.IsPublished, g.SortOrder,
	).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)
}

func (r *GuideRepository) Update(ctx context.Context, g *model.Guide) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE guides SET category_id = $1, title = $2, content = $3, is_published = $4, sort_order = $5, updated_at = NOW()
		 WHERE id = $6`,
		g.CategoryID, g.Title, g.Content, g.IsPublished, g.SortOrder, g.ID,
	)
	return err
}

func (r *GuideRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM guides WHERE id = $1`, id)
	return err
}

func scanGuides(rows *sql.Rows) ([]model.Guide, error) {
	var guides []model.Guide
	for rows.Next() {
		var g model.Guide
		if err := rows.Scan(&g.ID, &g.CategoryID, &g.CategoryName, &g.Title, &g.Content, &g.IsPublished, &g.SortOrder, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		guides = append(guides, g)
	}
	return guides, rows.Err()
}
