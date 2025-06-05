package main

import (
	"log"
	"os"

	"api/database"
	"api/routes"

	_ "api/docs" // This line is necessary for go-swagger to find your docs!

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API REST Go - Plateforme Éducative
// @version 1.0
// @description API REST complète pour la gestion d'une plateforme éducative avec familles, enseignants et administrateurs
// @termsOfService http://swagger.io/terms/

// @contact.name Support API
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Charger les variables d'environnement
	if err := godotenv.Load(); err != nil {
		log.Println("Aucun fichier .env trouvé")
	}

	// Initialiser la base de données
	database.InitDatabase()

	// Ajouter des données de test si nécessaire
	if os.Getenv("SEED_DATABASE") == "true" {
		if err := database.SeedDatabase(); err != nil {
			log.Printf("Erreur lors de l'ajout des données de test: %v", err)
		}
	}

	// Configurer Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Créer le routeur
	router := gin.Default()

	// Middleware CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Documentation Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configurer les routes de l'API
	routes.SetupRoutes(router)

	// Démarrer le serveur
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Serveur démarré sur le port %s", port)
	log.Printf("📚 Documentation Swagger disponible sur: http://localhost:%s/swagger/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur:", err)
	}
}
