package main

import (
	"log"
	"os"
	"strings" // 1. Garanta que o pacote 'strings' está importado

	"my-controls-api/config"
	"my-controls-api/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	r := gin.Default()

	config.ConnectDatabase()

	apiURL := os.Getenv("API_BASE_URL")
	serverPort := os.Getenv("SERVER_PORT")

	if serverPort == "" {
		serverPort = "8085"
	}

	// 2. Leia a variável com a lista de origens do .env
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		// Valor padrão para garantir que o desenvolvimento local funcione
		allowedOriginsEnv = "http://localhost:5173"
	}

	// 3. Crie o slice de origens a partir da string
	originsList := strings.Split(allowedOriginsEnv, ",")

	r.Use(cors.New(cors.Config{
		// 4. Use a lista de origens em vez do '*'
		AllowOrigins:     originsList,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./uploads")

	h := &handlers.Handler{DB: config.DB, ApiBaseURL: apiURL}

	api := r.Group("/api")
	{
		api.POST("/tsee/assignments", h.CreateAssignment)
		api.GET("/tsee/assignments", h.GetAssignments)
		api.PUT("/tsee/assignments/:id", h.UpdateAssignment)
		api.DELETE("/tsee/assignments/:id", h.DeleteAssignment)
		api.POST("/tsee/assignments/:id/evidence", h.UploadEvidence)
	}

	log.Printf("Starting server on port %s", serverPort)
	r.Run(":" + serverPort)
}
