package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"golang/internal/infrastructure/database/models"
	"golang/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentRepository struct {
	DB *pgxpool.Pool
}

func (r *DocumentRepository) UpdateComment(ctx context.Context, commentId int, content string) (any, any) {
	panic("unimplemented")
}

func (r *DocumentRepository) CheckDocumentAccess(ctx context.Context, userId int, documentId int) (bool, error) {
	var hasAccess bool
	err := r.DB.QueryRow(ctx, `
        SELECT EXISTS(
			SELECT 1 
			FROM documents d
			JOIN documents_users du ON du.document_id = d.id
			WHERE d.id = $1 AND du.user_id = $2 AND d.is_public = true
		)
    `, documentId, userId).Scan(&hasAccess)
	return hasAccess, err
}

func (r *DocumentRepository) CheckIsOwner(ctx context.Context, documentId int, userId int) (bool, error) {
	var hasAccess bool
	err := r.DB.QueryRow(ctx, `
        SELECT EXISTS(
			SELECT 1 
			FROM documents d
			JOIN documents_users du ON du.document_id = d.id
			WHERE d.id = $1 AND du.owner_id = $2
		)
    `, documentId, userId).Scan(&hasAccess)
	return hasAccess, err
}

func (r *DocumentRepository) CreateDocument(
	ctx context.Context,
	form models.CreateDocumentModel,
	userId int,
) (*models.BaseDocumentModel, error) {
	var document models.BaseDocumentModel

	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO documents (title, owner_id, is_public)
		VALUES ($1, $2, $3)
		RETURNING id, title, content, is_public, created_at, updated_at
	`
	insertDocErr := tx.QueryRow(ctx, query, form.Title, userId, form.IsPublic).Scan(
		&document.Id, &document.Title, &document.Content,
		&document.IsPublic, &document.CreatedAt, &document.UpdatedAt,
	)
	if insertDocErr != nil {
		return nil, insertDocErr
	}

	query = `
		INSERT INTO documents_users (document_id, user_id)
		VALUES ($1, $2)
	`
	_, insertDocUserErr := tx.Exec(ctx, query, document.Id, userId)
	if insertDocUserErr != nil {
		return nil, insertDocUserErr
	}

	return &document, nil
}

func (r *DocumentRepository) GetDocumentById(
	ctx context.Context,
	documentId int,
	userId int,
) (*models.DocumentModel, error) {
	var document models.DocumentModel
	var members []byte

	query := `
		SELECT 
			d.id, d.title, d.content, d.created_at, d.updated_at, owner.id, owner.username, owner.email,
			json_agg(
				json_build_object(
					'id', u.id, 
					'username', u.username, 
					'email', u.email
				)
			) AS members
			FROM documents AS d
			JOIN users owner ON d.owner_id = owner.id
			JOIN documents_users AS d_u ON d.id = d_u.document_id
			JOIN users AS u ON u.id = d_u.user_id
		WHERE d.id = $1
	`
	err := r.DB.QueryRow(ctx, query, documentId, userId).Scan(
		&document.Id, &document.Title, &document.Content, &document.CreatedAt, &document.UpdatedAt,
		&document.Owner.Id, &document.Owner.Username, &document.Owner.Email, &members,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(members, &document.Members); err != nil {
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
	var document models.BaseDocumentModel
	clauses, args := utils.GetSetParams(documentFormEncoded)

	query := fmt.Sprintf(`
		UPDATE documents 
		SET %s 
		WHERE id = $%d AND owner_id = $%d
		RETURNING id, title, content, is_public, created_at, updated_at
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

func (r *DocumentRepository) DeleteDocument(ctx context.Context, documentId int, userId int) error {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Exec(ctx, "DELETE FROM documents WHERE id = $1 AND owner_id = $2", documentId, userId)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	query := `
		DELETE FROM documents_users
		WHERE document_id = $1 AND user_id = $2
	`
	_, err = tx.Exec(ctx, query, documentId, userId)
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
	var document models.DocumentModel

	query := `
		UPDATE documents 
		SET content = $1
		WHERE id = $2 AND owner_id = $3
		RETURNING id, title, content, is_public, created_at, updated_at
	`
	err := r.DB.QueryRow(ctx, query, content, documentId, userId).Scan(
		&document.Id, &document.Title, &document.Content,
		&document.IsPublic, &document.CreatedAt, &document.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (r *DocumentRepository) GetDocumentsByUserId(ctx context.Context, userId int) ([]models.DocumentModel, error) {
	var documents []models.DocumentModel

	query := `
		SELECT d.id, d.title, d.content, d.is_public, d.created_at, d.updated_at
		FROM documents AS d
		JOIN documents_users AS d_u ON d_u.user_id = $1
		WHERE d_u.user_id = $1 AND d.is_public = true
	`
	rows, err := r.DB.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var document models.DocumentModel
		err := rows.Scan(
			&document.Id, &document.Title, &document.Content,
			&document.IsPublic, &document.CreatedAt, &document.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}

	return documents, nil
}

func (r *DocumentRepository) GetUserDocuments(
	ctx context.Context,
	userId int,
	limit int,
	offset int,
) ([]*models.BaseDocumentModel, error) {
	query := `
		SELECT d.id, d.title, d.content, d.is_public, d.created_at, d.updated_at
		FROM documents AS d
		JOIN documents_users AS d_u ON d_u.document_id = d.id
		WHERE d_u.user_id = $1
		LIMIT $2 OFFSET $3
	`
	rows, err := r.DB.Query(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var documents []*models.BaseDocumentModel
	for rows.Next() {
		var document models.BaseDocumentModel
		err := rows.Scan(
			&document.Id, &document.Title, &document.Content,
			&document.IsPublic, &document.CreatedAt, &document.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}
	return documents, nil
}

func (r *DocumentRepository) AddDocumentSnapshot(
	ctx context.Context,
	documentId int,
	userId int,
) (*models.BaseSnapshotModel, error) {
	query := `
		INSERT INTO document_snapshots (document_id, user_id, content)
		SELECT d.id, d_u.user_id, d.content
		FROM documents AS d
		JOIN documents_users AS d_u ON d_u.document_id = d.id
		WHERE d.id = $1 AND d_u.user_id = $2
		RETURNING id, document_id, user_id, created_at
	`

	var snapshot models.BaseSnapshotModel

	err := r.DB.QueryRow(ctx, query, documentId, userId).Scan(
		&snapshot.Id, &snapshot.DocumentId, &snapshot.UserId, &snapshot.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}
