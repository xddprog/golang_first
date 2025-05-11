package repositories

import (
	"context"
	"fmt"
	"golang/internal/infrastructure/database/models"
	"golang/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type DocumentRepository struct {
	DB *pgxpool.Pool
}


func (r *DocumentRepository) CheckDocumentAccess(ctx context.Context, transaction pgx.Tx, userID, docID int) (bool, error) {
    var hasAccess bool
    err := r.DB.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM documents 
            WHERE id = $1 AND (
                owner_id = $2 OR 
                id IN (SELECT document_id FROM documents_users WHERE user_id = $2)
            )
        )`, docID, userID).Scan(&hasAccess)
    return hasAccess, err
}


func (r *DocumentRepository) CreateDocument(ctx context.Context, form models.CreateDocumentModel, userId int) (*models.BaseDocumentModel, error) {
	var document models.BaseDocumentModel

	query := `
		SELECT 
			d.id, d.title, d.content, d.is_public, d.created_at, d.updated_at
		FROM 
			(
				INSERT INTO documents (title, owner_id, is_public)
				VALUES
					($1, $2, $3)
				RETURNING id, title, content, is_public, created_at, updated_at
			) as d
	`
	err := r.DB.QueryRow(ctx, query, form.Title, userId, form.IsPublic).Scan(
		&document.Id, &document.Title, &document.Content, 
		&document.IsPublic, &document.CreatedAt, &document.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &document, nil
}


func (r *DocumentRepository) GetDocumentById(
	ctx context.Context, 
	documentId int, 
	userId int,
) (*models.DocumentModel, error) {
	var document models.DocumentModel

	query := `
		SELECT d.id, d.title, d.content, d.created_at, d.updated_at, u.id, u.username, u.email
			FROM documents AS d
			JOIN documents_users AS d_u ON d.id = d_u.document_id
			JOIN users AS u ON u.id = d_u.user_id
		WHERE d.id = $1 AND u.id = $2
	`
	err := r.DB.QueryRow(ctx, query, documentId, userId).Scan(
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
	userId int,
	documentId int, 
	documentFormEncoded models.UpdateDocumentModel,
) (*models.BaseDocumentModel, error) {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx)

	_, accessErr := r.CheckDocumentAccess(ctx, tx, userId, documentId)
	if accessErr != nil {
		return nil, accessErr
	}

	var document models.BaseDocumentModel
	clauses, args := utils.GetSetParams(documentFormEncoded)

	query := fmt.Sprintf(`
		SELECT 
			d.id, d.title, d.content, d.is_public, d.created_at, d.updated_at,
		FROM 
			(
				UPDATE documents 
				SET %s 
				WHERE id = $%d AND owner_id = $%d
				RETURNING id, title, content, is_public, created_at, updated_at
			) as d
	`, clauses, documentId, userId)

	documentErr := r.DB.QueryRow(ctx, query, args...).Scan(
		&document.Id, &document.Title, &document.Content, 
		&document.IsPublic, &document.CreatedAt, &document.UpdatedAt,
	)

	if documentErr != nil {
		return nil, documentErr 
	}
	return &document, nil
}


func (r *DocumentRepository) DeleteDocument(ctx context.Context, documentId int, userId int) error{
	_, err := r.DB.Exec(ctx, "DELETE FROM documents WHERE id = $1 AND owner_id = $2", documentId, userId)
	if err != nil {
		return err
	}
	return nil
}


func (r *DocumentRepository) UpdateDocumentContent(
	ctx context.Context, 
	userId int, 
	documentId int, 
	content string,
) (*models.DocumentModel, error) {
	transaction, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	_, accessErr := r.CheckDocumentAccess(ctx, transaction, userId, documentId)
	if accessErr != nil {
		return nil, accessErr
	}

	var document models.DocumentModel
	
	query := `
		SELECT 
			d.id, d.title, d.content, d.is_public, d.created_at, d.updated_at,
		FROM 
			(
				UPDATE documents 
				SET content = $1
				WHERE id = $2 AND owner_id = $3
				RETURNING id, title, content, is_public
			) as d
	`
	deleteErr := r.DB.QueryRow(ctx, query, content, documentId, userId).Scan(
		&document.Id, &document.Title, &document.Content,
		&document.CreatedAt, &document.UpdatedAt,
	)
	
	if deleteErr != nil {
		return nil, deleteErr
	}
	return &document, nil
}