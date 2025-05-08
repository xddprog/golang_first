package services

import (
	"context"
	"encoding/json"
	"golang/internal/core/repositories"
	"golang/internal/infrastructure/database/models"
	"golang/internal/infrastructure/errors"
	"golang/internal/utils"
	"io"
)


type DocumentService struct {
	Repository *repositories.DocumentRepository
}


func (s *DocumentService) CreateDocument(ctx context.Context, userId int, documentForm io.ReadCloser) (*models.DocumentModel, *apierrors.APIError) {
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


func (s *DocumentService) GetDocumentById(ctx context.Context, documentId int) (*models.DocumentModel, *apierrors.APIError) {
	document, err := s.Repository.GetDocumentById(ctx, documentId)
	if err != nil {
		return nil, apierrors.CheckDBError(err)
	}
	return document, nil
}


func (s *DocumentService) UpdateDocument(ctx context.Context, documentId int, documentForm io.ReadCloser) (*models.DocumentModel, *apierrors.APIError) {
	var documentFormEncoded models.UpdateDocumentModel

	if err := json.NewDecoder(documentForm).Decode(&documentFormEncoded); err != nil {
		return nil, &apierrors.ErrInvalidRequestBody
	}

	if err := utils.ValidateForm(documentFormEncoded); err != nil {
		return nil, &apierrors.ErrEncodingError
	}

	document, err := s.Repository.UpdateDocument(ctx, documentId, documentFormEncoded)
	if err != nil {
		return nil, apierrors.CheckDBError(err)
	}
	return document, nil
}