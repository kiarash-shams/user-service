package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"user-service/api/dto"
	"user-service/api/helper"
	"user-service/config"
	"user-service/services"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



type DocumentHandler struct {
	service *services.DocumentService
}

func NewDocumentHandler(cfg *config.Config) *DocumentHandler {
	return &DocumentHandler{
		service: services.NewDocumentService(cfg),
	}
}


// CreateDocument godoc
// @Summary Create a Document
// @Description Create a Document
// @Tags Documents
// @Accept x-www-form-urlencoded
// @produces json
// @Param file formData dto.UploadFileRequest true "Create a file"
// @Param file formData file true "Create a file"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.DocumentResponse} "Document response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/documents/ [post]
// @Security AuthBearer
func (h *DocumentHandler) Create(c *gin.Context) {
	upload := dto.UploadFileRequest{}
	err := c.ShouldBind(&upload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	req := dto.CreateDocumentRequest{}
	req.Description = upload.Description
	req.MimeType = upload.File.Header.Get("Content-Type")
	req.Directory = "uploads"
	req.Name, err = saveUploadedFile(upload.File, req.Directory)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	res, err := h.service.Create(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, helper.Success))

}


// GetDocument godoc
// @Summary Get a Document
// @Description Get a Document
// @Tags Documents
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.DocumentResponse} "Document response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v2/auth/resource/documents/{id} [get]
// @Security AuthBearer
func (h *DocumentHandler) GetById(c *gin.Context) {
	GetById(c, h.service.GetById)
}


func saveUploadedFile(file *multipart.FileHeader, directory string) (string, error) {
	// test.txt -> 95239855629856.txt
	randFileName := uuid.New()
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return "", err
	}
	fileName := file.Filename
	fileNameArr := strings.Split(fileName, ".")
	fileExt := fileNameArr[len(fileNameArr)-1]
	fileName = fmt.Sprintf("%s.%s", randFileName, fileExt)
	dst := fmt.Sprintf("%s/%s", directory, fileName)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}
	return fileName, nil
}