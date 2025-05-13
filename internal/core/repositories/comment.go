package repositories

import (
	"context"
	"golang/internal/infrastructure/database/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	DB *pgxpool.Pool
}


func (r *CommentRepository) GetCommentsReplies(context context.Context, commentId int) ([]*models.CommentModel, error) {
	var comments []*models.CommentModel

	query := `
		SELECT 
			c.id, c.content, c.user_id, c.parent_id, c.document_id 
			c.created_at, c.updated_at, u.id, u.username, u.email
		FROM comments c
		JOIN user u ON user.id = c.user_id
		WHERE parent_id = $1
	`
	rows, err := r.DB.Query(context, query, commentId)
	if err != nil {
	    return nil, err
	}

	for rows.Next() {
		var comment models.CommentModel
		err := rows.Scan(
			&comment.Id, &comment.Content, &comment.UserId, &comment.ParentId,
			&comment.DocumentId, &comment.CreatedAt, &comment.UpdatedAt, 
			&comment.User.Id, &comment.User.Username, &comment.User.Email,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}


func (r *CommentRepository) DeleteComment(
	context context.Context,
	userId int,
	commentId int,
	documentId int,
) error {
	query := `DELETE FROM comments WHERE id = $1 AND user_id = $2`
	rows, err := r.DB.Exec(context, query, commentId, userId)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}


func (r *CommentRepository) UpdateComment(ctx context.Context, userId int, content string) (*models.CommentModel, error) {
	var comment models.CommentModel

	query := `
		INSERT INTO comments (userId, content) VALUES ($1, $2) 
		RETURNING id, userId, content
	`
	err := r.DB.QueryRow(ctx, query, userId, content).Scan(&comment.Id, &comment.UserId, &comment.Content)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}


func (r *CommentRepository) CreateComment(
	context context.Context,
	userId int,
	documentId int,
	commentForm models.CreateCommentModel,
) (*models.CommentModel, error) {
	var comment models.CommentModel

	query := `
		SELECT 
			c.id, c.userId, c.documentId, c.content, c.createdAt, c.updatedAt,
			u.id, u.username, u.email
		FROM 
			(
				INSERT INTO comments (user_id, document_id, content)
				VALUES ($1, $2, $3)
				RETURNING id, user_id, document_id, content, created_at, updatedAt
			) as c
		JOIN users u ON u.id = c.user_id
	`
	err := r.DB.QueryRow(context, query, userId, documentId, commentForm.Content).Scan(
		&comment.Id, &comment.UserId, &comment.DocumentId,
		&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
		&comment.User.Id, &comment.User.Username, &comment.User.Email,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}


func (r *CommentRepository) GetCommentsByDocument(
	ctx context.Context,
	documentId int,
	limit int,
	offset int,
) ([]*models.CommentModel, error) {
	var comments []*models.CommentModel

	query := `
		SELECT 
			c.id, c.user_id, c.document_id, c.content, c.created_at, c.updated_at,
			u.id, u.username, u.email
		FROM comments c
		JOIN documents d ON c.document_id = d.id
		JOIN users u ON c.user_id = u.id
		WHERE c.document_id = $1
		LIMIT $2 OFFSET $3
	`
	rows, err := r.DB.Query(ctx, query, documentId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.CommentModel
		err := rows.Scan(
			&comment.Id, &comment.UserId, &comment.DocumentId,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
			&comment.User.Id, &comment.User.Username, &comment.User.Email,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}
