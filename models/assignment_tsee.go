package models

import (
	"time"

	"gorm.io/gorm"
)

type AssignmentTsee struct {
	gorm.Model               // Inclui ID, CreatedAt, UpdatedAt, DeletedAt
	Agent          string    `json:"agent"`
	Name           string    `json:"name"`
	Local          string    `json:"local"`
	AssignmentDate time.Time `json:"assignmentDate"`
	Quantity       int       `json:"quantity"`
	EvidenceURL    string    `json:"evidenceUrl,omitempty"` // Link para o arquivo salvo
}
