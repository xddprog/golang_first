package models

import "time"


type DocumentOwner struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}


type CreateDocumentModel struct {
	Title   string 
	Content string
}


type UpdateDocumentModel struct {
	Title *string `json:"title,omitempty" db:"title" validate:"omitempty"`
	IsPublic *bool `json:"is_public,omitempty" db:"is_public" validate:"omitempty"`
}


type DocumentModel struct {
	Id        int       	`json:"id"`
	Title     string    	`json:"title"`
	Content   string    	`json:"content"`
	Owner	  BaseUserModel `json:"owner"`
	CreatedAt time.Time 	`json:"createdAt"`
	IsPublic  bool			`json:"isPublic"`
	UpdatedAt time.Time 	`json:"updatedAt"`
}