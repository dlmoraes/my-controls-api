package main

import (
	"my-controls-api/config"
	"my-controls-api/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	r := gin.Default()

	// Conecta ao banco de dados
	config.ConnectDatabase()

	clientURL := os.Getenv("CLIENT_ORIGIN_URL")
	apiURL := os.Getenv("API_BASE_URL")
	serverPort := os.Getenv("SERVER_PORT")

	if serverPort == "" {
		serverPort = "8085"
	}

	// Configura o CORS (Cross-Origin Resource Sharing) para permitir requisições do seu frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{clientURL}, // Endereço do seu frontend Next.js
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Rota para servir os arquivos estáticos da pasta /uploads
	r.Static("/uploads", "./uploads")

	// Instancia o handler com a conexão do banco
	h := &handlers.Handler{DB: config.DB, ApiBaseURL: apiURL}

	// Agrupa as rotas da API
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
