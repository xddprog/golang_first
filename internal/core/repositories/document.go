package repositories

import (
	"context"
	"fmt"
	"golang/internal/infrastructure/database/models"
	"golang/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)


type DocumentRepository struct {
	DB *pgxpool.Pool
}


func (r *DocumentRepository) CreateDocument(ctx context.Context, form models.CreateDocumentModel, userId int) (*models.DocumentModel, error) {
	var document models.DocumentModel

	query := `
		SELECT 
			d.id, d.title, d.content, d.is_public, d.created_at, d.updated_at,
			u.id, u.username, u.email
		FROM 
			(
				INSERT INTO documents (title, owner_id, is_public)
				VALUES
					($1, $2, $3)
				RETURNING id, title, content, is_public
			) as d
		JOIN users u ON u.id = d.id
	`

	err := r.DB.QueryRow(ctx, query, form.Title, userId, form.Content).Scan(
		&document.Id, &document.Title, &document.Content, &document.IsPublic, &document.CreatedAt,
		&document.UpdatedAt, &document.Owner.Id, &document.Owner.Username, &document.Owner.Email,
	)
	if err != nil {
		return nil, err
	}

	return &document, nil
}


func (r *DocumentRepository) GetDocumentById(ctx context.Context, id int) (*models.DocumentModel, error) {
	var document models.DocumentModel

	err := r.DB.QueryRow(ctx, "SELECT * FROM document WHERE id = $1", id).Scan(
		&document.Id, &document.Title, &document.Content, &document.CreatedAt, &document.UpdatedAt,
		&document.Owner.Id, &document.Owner.Username, &document.Owner.Email,
	)
	if err != nil {
		return nil, err
	}

	return &document, nil
}


func (r *DocumentRepository) UpdateDocument(
	ctx context.Context, 
	documentId int, 
	documentFormEncoded models.UpdateDocumentModel,
) (*models.DocumentModel, error) {
	var document models.DocumentModel
	clauses, args := utils.GetSetParams(documentFormEncoded)

	query := fmt.Sprintf(`
		SELECT 
			d.id, d.title, d.content, d.is_public, d.created_at, d.updated_at,
			u.id, u.username, u.email
		FROM 
			(
				UPDATE documents SET %s WHERE id = $%d
				RETURNING id, title, content, is_public
			) as d
		JOIN users u ON u.id = d.id
	`, clauses, documentId)

	err := r.DB.QueryRow(ctx, query, args...).Scan(
		&document.Id, &document.Title, &document.Content, &document.CreatedAt, &document.UpdatedAt,
		&document.Owner.Id, &document.Owner.Username, &document.Owner.Email,
	)

	if err != nil {
		return nil, err 
	}
	return &document, nil
}


func (r *DocumentRepository) DeleteDocument(ctx context.Context, documentId int) {
	r.DB.QueryRow(ctx, "DELETE FROM users WHERE id = $1", documentId)
}
