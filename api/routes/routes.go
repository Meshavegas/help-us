package routes

import (
	"net/http"

	"api/controllers"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

// HealthCheck vérifie l'état de l'API
// @Summary      Vérification de santé
// @Description  Vérifie l'état de l'API
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "API fonctionnelle"
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "API Plateforme Éducative fonctionne correctement",
	})
}

func SetupRoutes(router *gin.Engine) {
	router.GET("/health", HealthCheck)

	// Groupe API v1
	v1 := router.Group("/api/v1")
	{
		// Routes d'authentification (publiques)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/refresh", controllers.RefreshToken)
		}

		// Routes protégées (nécessitent une authentification)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", controllers.GetProfile)
			protected.PUT("/profile", controllers.UpdateProfile)
			protected.POST("/auth/logout", controllers.Logout)

			// Routes administrateur
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireAdmin())
			{
				// TODO: Ajouter les routes d'administration
				admin.GET("/users", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Route admin - liste des utilisateurs"})
				})
			}

			// Routes enseignant
			teacher := protected.Group("/teacher")
			teacher.Use(middleware.RequireTeacherOrAdmin())
			{
				// TODO: Ajouter les routes enseignant
				teacher.GET("/courses", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Route enseignant - liste des cours"})
				})
			}

			// Routes famille
			family := protected.Group("/family")
			family.Use(middleware.RequireParentOrAdmin())
			{
				// TODO: Ajouter les routes famille
				family.GET("/missions", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Route famille - liste des missions"})
				})
			}
		}
	}

	// Route 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route non trouvée",
		})
	})
}
