package services

import (
	"context"
	"encoding/json"
	"golang/internal/core/repositories"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"golang/internal/utils"
	"io"

	"github.com/jackc/pgx/v5"
)


type CommentService struct {
	Repository *repositories.CommentRepository
}


func (s *CommentService) GetCommentsReplies(context context.Context, commentId int) (any, any) {
	comments, err := s.Repository.GetCommentsReplies(context, commentId)
	if err != nil {
		return nil, apierrors.CheckDBError(err, "document")
	}
	return comments, nil
}


func (s *CommentService) DeleteComment(
	ctx context.Context,
	userId int,
	commentId int,
	documentId int,
) *apierrors.APIError {
	err := s.Repository.DeleteComment(ctx, userId, commentId, documentId)
	if err != nil && err != pgx.ErrNoRows {
		return apierrors.CheckDBError(err, "comment")
	}
	return nil
}


func (s *CommentService) GetCommentsByDocument(
	ctx context.Context,
	documentId int,
	limit int,
	offset int,
) ([]*models.CommentModel, *apierrors.APIError) {
	comments, err := s.Repository.GetCommentsByDocument(ctx, documentId, limit, offset)
	if err != nil {
		return nil, apierrors.CheckDBError(err, "comments")
	}
	return comments, nil
}


func (s *CommentService) CreateComment(
	context context.Context,
	user *models.BaseUserModel,
	documentId int,
	form io.ReadCloser,
) (*models.CommentModel, error) {
	var commentForm models.CreateCommentModel

	if err := json.NewDecoder(form).Decode(&commentForm); err != nil {
		return nil, err
	}

	validateErr := utils.ValidateForm(commentForm)
	if validateErr != nil {
		return nil, validateErr
	}

	comment, dbErr := s.Repository.CreateComment(context, user.Id, documentId, commentForm)
	if dbErr != nil {
		return nil, apierrors.CheckDBError(dbErr, "comment")
	}
	return comment, nil
}


func (s *CommentService) UpdateComment(
	ctx context.Context,
	commentId int,
	userId int,
	form io.ReadCloser,
) (*models.CommentModel, error) {
	var commentForm models.UpdateCommentModel

	if err := json.NewDecoder(form).Decode(&commentForm); err != nil {
		return nil, err
	}

	validateErr := utils.ValidateForm(commentForm)
	if validateErr != nil {
		return nil, validateErr
	}

	comment, dbErr := s.Repository.UpdateComment(ctx, commentId, commentForm.Content)
	if dbErr != nil {
		return nil, apierrors.CheckDBError(dbErr, "comment")
	}
	return comment, nil
}
