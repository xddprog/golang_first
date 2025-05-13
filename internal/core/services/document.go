package services

import (
	"context"
	"encoding/json"
	"fmt"
	"golang/internal/core/repositories"
	"golang/internal/infrastructure/clients"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"golang/internal/utils"
	"io"
	"strconv"
	"time"
)

type DocumentService struct {
	Repository  *repositories.DocumentRepository
	RedisClient *clients.RedisClient
	SMTPClient  *clients.SmtpClient
}


func (s *DocumentService) CheckDocumentAccess(ctx context.Context, userId int, documentId int) *apierrors.APIError {
	_, err := s.Repository.CheckDocumentAccess(ctx, userId, documentId)
	if err != nil {
		return &apierrors.ErrDocumentAccessDenied
	}
	return nil
}

func (s *DocumentService) CreateDocument(
	ctx context.Context,
	userId int,
	documentForm io.ReadCloser,
) (*models.BaseDocumentModel, *apierrors.APIError) {
	var documentFormEncoded models.CreateDocumentModel

	if err := json.NewDecoder(documentForm).Decode(&documentFormEncoded); err != nil {
		return nil, &apierrors.ErrInvalidRequestBody
	}

	if err := utils.ValidateForm(documentFormEncoded); err != nil {
		return nil, &apierrors.ErrEncodingError
	}

	document, err := s.Repository.CreateDocument(ctx, documentFormEncoded, userId)
	if err != nil {
		return nil, apierrors.CheckDBError(err, "document")
	}
	return document, nil
}

func (s *DocumentService) GetDocumentById(
	ctx context.Context,
	documentId int,
	userId int,
) (*models.DocumentModel, *apierrors.APIError) {
	_, err := s.Repository.CheckDocumentAccess(ctx, userId, documentId)
	if err != nil {
		return nil, &apierrors.ErrDocumentAccessDenied
	}

	document, err := s.Repository.GetDocumentById(ctx, documentId, userId)
	if err != nil {
		return nil, apierrors.CheckDBError(err, "document")
	}
	return document, nil
}

func (s *DocumentService) UpdateDocument(
	ctx context.Context,
	userId int,
	documentId int,
	documentForm io.ReadCloser,
) (*models.BaseDocumentModel, *apierrors.APIError) {
	var documentFormEncoded models.UpdateDocumentModel

	if err := json.NewDecoder(documentForm).Decode(&documentFormEncoded); err != nil {
		return nil, &apierrors.ErrInvalidRequestBody
	}

	if err := utils.ValidateForm(documentFormEncoded); err != nil {
		return nil, &apierrors.ErrEncodingError
	}

	document, err := s.Repository.UpdateDocument(ctx, userId, documentId, documentFormEncoded)
	if err != nil {
		return nil, apierrors.CheckDBError(err, "document")
	}
	return document, nil
}

func (s *DocumentService) DeleteDocument(ctx context.Context, documentId int, userId int) *apierrors.APIError {
	err := s.Repository.DeleteDocument(ctx, documentId, userId)
	if err != nil {
		return apierrors.CheckDBError(err, "document")
	}
	return nil
}

func (s *DocumentService) UpdateDocumentContent(
	ctx context.Context,
	userId int,
	documentId string,
	content string,
) (*models.DocumentModel, *apierrors.APIError) {
	documentIdInt, err := strconv.Atoi(documentId)
	if err != nil {
		return nil, &apierrors.ErrInvalidRequestBody
	}

	document, err := s.Repository.UpdateDocumentContent(ctx, userId, documentIdInt, content)
	if err != nil {
		return nil, apierrors.CheckDBError(err, "document")
	}
	return document, nil
}

func (s *DocumentService) SendInvite(
	ctx context.Context,
	userId int,
	userEmail string,
	documentId int,
) *apierrors.APIError {
	if _, err := s.Repository.CheckIsOwner(ctx, userId, documentId); err != nil {
		return &apierrors.ErrDocumentAccessDenied
	}

	key := fmt.Sprintf("document:%d:user:%d:invite", documentId, userId)
	code := utils.RandSeq(6)
	err := s.RedisClient.Set(ctx, key, code, time.Hour*24)
	if err != nil {
		return &apierrors.ErrInvalidRequestBody
	}

	redisErr := s.RedisClient.Set(ctx, key, code, time.Hour*24)
	if redisErr != nil {
		return &apierrors.ErrInternalServerError
	}

	document, err := s.Repository.GetDocumentById(ctx, documentId, userId)
	if err != nil {
		return apierrors.CheckDBError(err, "document")
	}

	smtpErr := s.SMTPClient.SendInviteToDocument(
		userEmail,
		"Invite to Document",
		code,
		strconv.Itoa(documentId),
		document.Title,
	)
	if smtpErr != nil {
		return &apierrors.ErrInternalServerError
	}
	return nil
}

func (s *DocumentService) GetUserDocuments(
	ctx context.Context,
	userId int,
	limit int,
	offset int,
) ([]*models.BaseDocumentModel, *apierrors.APIError) {
	documents, err := s.Repository.GetUserDocuments(ctx, userId, limit, offset)
	if err != nil {
		return nil, apierrors.CheckDBError(err, "document")
	}
	return documents, nil
}

func (s *DocumentService) AddDocumentSnapshot(
	ctx context.Context,
	userId int,
	documentId int,
) (*models.BaseSnapshotModel, *apierrors.APIError) {
	snapshot, err := s.Repository.AddDocumentSnapshot(ctx, userId, documentId)
	if err != nil {
		return nil, apierrors.CheckDBError(err, "document")
	}
	return snapshot, nil
}
