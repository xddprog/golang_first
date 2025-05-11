package services

import (
	"context"
	"encoding/json"
	"golang/internal/core/repositories"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"golang/internal/utils"
	"io"
	"strconv"
)


type DocumentService struct {
	Repository *repositories.DocumentRepository
}


func (s *DocumentService) CreateDocument(ctx context.Context, userId int, documentForm io.ReadCloser) (*models.BaseDocumentModel, *apierrors.APIError) {
	var documentFormEncoded models.CreateDocumentModel

	if err := json.NewDecoder(documentForm).Decode(&documentFormEncoded); err != nil {
		return nil, &apierrors.ErrInvalidRequestBody
	}

	if err := utils.ValidateForm(documentFormEncoded); err != nil {
		return nil, &apierrors.ErrEncodingError
	}

	document, err := s.Repository.CreateDocument(ctx, documentFormEncoded, userId)
	if err != nil {
		return nil, apierrors.CheckDBError(err)
	}
	return document, nil
}


func (s *DocumentService) GetDocumentById(
	ctx context.Context, 
	documentId int, 
	userId int,
) (*models.DocumentModel, *apierrors.APIError) {
	document, err := s.Repository.GetDocumentById(ctx, documentId, userId)
	if err != nil {
		return nil, apierrors.CheckDBError(err)
	}
	return document, nil
}


func (s *DocumentService) UpdateDocument(ctx context.Context, userId int, documentId int, documentForm io.ReadCloser,) (*models.BaseDocumentModel, *apierrors.APIError) {
	var documentFormEncoded models.UpdateDocumentModel

	if err := json.NewDecoder(documentForm).Decode(&documentFormEncoded); err != nil {
		return nil, &apierrors.ErrInvalidRequestBody
	}

	if err := utils.ValidateForm(documentFormEncoded); err != nil {
		return nil, &apierrors.ErrEncodingError
	}

	document, err := s.Repository.UpdateDocument(ctx, userId, documentId, documentFormEncoded)
	if err != nil {
		return nil, apierrors.CheckDBError(err)
	}
	return document, nil
}


func (s *DocumentService) DeleteDocument(ctx context.Context, documentId int, userId int) *apierrors.APIError {
	err := s.Repository.DeleteDocument(ctx, documentId, userId)
	if err != nil {
		return apierrors.CheckDBError(err)
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
		return nil, apierrors.CheckDBError(err)
	}
	return document, nil
}