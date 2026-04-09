package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rootage-ses-quiz/backend/internal/model"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, doc *model.Document) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO documents (title, filename, content) VALUES ($1, $2, $3) RETURNING id, uploaded_at`,
		doc.Title, doc.Filename, doc.Content,
	).Scan(&doc.ID, &doc.UploadedAt)
}

func (r *DocumentRepository) List(ctx context.Context) ([]model.Document, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, filename, uploaded_at FROM documents ORDER BY uploaded_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []model.Document
	for rows.Next() {
		var d model.Document
		if err := rows.Scan(&d.ID, &d.Title, &d.Filename, &d.UploadedAt); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

func (r *DocumentRepository) GetByID(ctx context.Context, id string) (*model.Document, error) {
	var d model.Document
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, filename, content, uploaded_at FROM documents WHERE id = $1`, id,
	).Scan(&d.ID, &d.Title, &d.Filename, &d.Content, &d.UploadedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DocumentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM documents WHERE id = $1`, id)
	return err
}

func (r *DocumentRepository) CreateChunks(ctx context.Context, chunks []model.DocumentChunk) error {
	if len(chunks) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(chunks))
	valueArgs := make([]any, 0, len(chunks)*4)

	for i, chunk := range chunks {
		base := i * 4
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d::vector)", base+1, base+2, base+3, base+4))
		valueArgs = append(valueArgs, chunk.DocumentID, chunk.ChunkIndex, chunk.Content, float64SliceToString(chunk.Embedding))
	}

	query := fmt.Sprintf(
		`INSERT INTO document_chunks (document_id, chunk_index, content, embedding) VALUES %s`,
		strings.Join(valueStrings, ", "),
	)

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	return err
}

func (r *DocumentRepository) DeleteChunksByDocument(ctx context.Context, docID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM document_chunks WHERE document_id = $1`, docID)
	return err
}

func (r *DocumentRepository) SearchSimilar(ctx context.Context, embedding []float64, limit int) ([]model.DocumentChunk, error) {
	embStr := float64SliceToString(embedding)
	rows, err := r.db.QueryContext(ctx,
		`SELECT dc.id, dc.document_id, dc.chunk_index, dc.content, dc.created_at, d.title
		 FROM document_chunks dc
		 JOIN documents d ON dc.document_id = d.id
		 WHERE dc.embedding IS NOT NULL
		 ORDER BY dc.embedding <=> $1::vector
		 LIMIT $2`,
		embStr, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []model.DocumentChunk
	for rows.Next() {
		var c model.DocumentChunk
		var title string
		if err := rows.Scan(&c.ID, &c.DocumentID, &c.ChunkIndex, &c.Content, &c.CreatedAt, &title); err != nil {
			return nil, err
		}
		chunks = append(chunks, c)
	}
	return chunks, rows.Err()
}

type ChunkWithTitle struct {
	model.DocumentChunk
	DocumentTitle string
}

func (r *DocumentRepository) SearchSimilarWithTitle(ctx context.Context, embedding []float64, limit int) ([]ChunkWithTitle, error) {
	embStr := float64SliceToString(embedding)
	rows, err := r.db.QueryContext(ctx,
		`SELECT dc.id, dc.document_id, dc.chunk_index, dc.content, dc.created_at, d.title
		 FROM document_chunks dc
		 JOIN documents d ON dc.document_id = d.id
		 WHERE dc.embedding IS NOT NULL
		 ORDER BY dc.embedding <=> $1::vector
		 LIMIT $2`,
		embStr, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []ChunkWithTitle
	for rows.Next() {
		var c ChunkWithTitle
		if err := rows.Scan(&c.ID, &c.DocumentID, &c.ChunkIndex, &c.Content, &c.CreatedAt, &c.DocumentTitle); err != nil {
			return nil, err
		}
		chunks = append(chunks, c)
	}
	return chunks, rows.Err()
}

func float64SliceToString(v []float64) string {
	if len(v) == 0 {
		return "[]"
	}
	parts := make([]string, len(v))
	for i, f := range v {
		parts[i] = fmt.Sprintf("%f", f)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
