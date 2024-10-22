package dto

import (
	"mime/multipart"
)


type FileFormRequest struct {
	File *multipart.FileHeader `json:"file" form:"file" binding:"required" swaggerignore:"true"`
}

type UploadFileRequest struct {
	FileFormRequest
	Description string `json:"description" form:"description" binding:"required"`
	DocCategory string `json:"docCategory" form:"docCategory" binding:"required"`
}

type CreateDocumentRequest struct {
	DocCategory string `json:"docCategory"`	
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`	
}

type UpdateDocumentRequest struct {
	DocCategory string `json:"docCategory"`	
	Description string `json:"description"`
}

type DocumentResponse struct {
	Id          int          `json:"id"`
	DocCategory string 		 `json:"docCategory" `
	Name        string `json:"name,omitempty"`
	Directory   string `json:"directory,omitempty"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}