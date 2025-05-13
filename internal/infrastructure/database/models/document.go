package models

import (
	"encoding/json"
	"time"
)


type DocumentOwner struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}


type CreateDocumentModel struct {
	Title   string 
	IsPublic bool
}


type UpdateDocumentModel struct {
	Title 		*string 	`json:"title,omitempty" db:"title" validate:"omitempty"`
	IsPublic 	*bool 		`json:"is_public,omitempty" db:"is_public" validate:"omitempty"`
}


type UpdateDocumentContent struct {
	DocumentId 	string 		`json:"documentId"`
	Content 	string 		`json:"content"`
}


type CursorMove struct {
	DocumentId 	string          `json:"doc_id"`
	Position    json.RawMessage `json:"position"`
}


type BaseDocumentModel struct {
	Id        int       	`json:"id"`
	Title     string    	`json:"title"`
	Content   string    	`json:"content"`
	CreatedAt time.Time 	`json:"createdAt"`
	IsPublic  bool			`json:"isPublic"`
	UpdatedAt time.Time 	`json:"updatedAt"`
}


type DocumentModel struct {
	BaseDocumentModel 
	Owner	  BaseUserModel 	`json:"owner"`
	Members   []BaseUserModel 	`json:"members"`
}


type BaseSnapshotModel struct {
	Id        int       	`json:"id"`
	DocumentId int       	`json:"documentId"`
	UserId   int       		`json:"userId"`
	CreatedAt time.Time 	`json:"createdAt"`
}