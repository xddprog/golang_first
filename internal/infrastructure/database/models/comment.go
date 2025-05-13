package models

import "time"


type CreateCommentModel struct {
    Content     string  `json:"content" validate:"required"`
    ParentId    int     `json:"parent_id"`
}


type CommentModel struct {
    Id          int             `json:"id"`
    User        BaseUserModel   `json:"user"`
    ParentId    int             `json:"parent_id"`
    Content     string          `json:"content"`
    UserId      int             `json:"user_id"`
    DocumentId  int             `json:"document_id"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
}


type UpdateCommentModel struct {
    Content string `json:"content" validate:"required"`
}