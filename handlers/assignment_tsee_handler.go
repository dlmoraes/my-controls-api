package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"my-controls-api/models"
)

type Handler struct {
	DB         *gorm.DB
	ApiBaseURL string
}

// POST /assignments - Cria uma nova atribuição
func (h *Handler) CreateAssignment(c *gin.Context) {
	var input models.AssignmentTsee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GET /assignments - Lista todas as atribuições
func (h *Handler) GetAssignments(c *gin.Context) {
	var assignments []models.AssignmentTsee
	if err := h.DB.Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

// PUT /assignments/:id - Atualiza uma atribuição
func (h *Handler) UpdateAssignment(c *gin.Context) {
	id := c.Param("id")
	var assignment models.AssignmentTsee
	if err := h.DB.First(&assignment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	var input models.AssignmentTsee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&assignment).Updates(input)
	c.JSON(http.StatusOK, assignment)
}

// DELETE /assignments/:id - Deleta uma atribuição
func (h *Handler) DeleteAssignment(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.AssignmentTsee{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted successfully"})
}

// POST /assignments/:id/evidence - Faz upload de um arquivo de evidência
func (h *Handler) UploadEvidence(c *gin.Context) {
	idStr := c.Param("id")
	assignmentID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var assignment models.AssignmentTsee
	if err := h.DB.First(&assignment, assignmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	filename := fmt.Sprintf("%d-%s", assignmentID, filepath.Base(file.Filename))
	savePath := filepath.Join("./uploads", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	fileURL := h.ApiBaseURL + "/uploads/" + filename
	h.DB.Model(&assignment).Update("EvidenceURL", fileURL)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "url": fileURL})
}
