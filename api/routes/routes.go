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

			// Routes utilisateurs (accessibles à tous les utilisateurs authentifiés pour leur propre profil)
			users := protected.Group("/users")
			{
				// Routes spécifiques à un utilisateur
				users.GET("/:id", controllers.GetUserByID)
				users.GET("/:id/addresses", controllers.GetUserAddresses)
				users.GET("/:id/payments", controllers.GetUserPayments)
				users.GET("/:id/resources", controllers.GetUserResources)
			}

			// Endpoints utilisateurs administrateur (liste, update, delete) au chemin /users...
			protected.GET("/users", middleware.RequireAdmin(), controllers.GetAllUsers)
			protected.PUT("/users/:id", middleware.RequireAdmin(), controllers.UpdateUserByID)
			protected.DELETE("/users/:id", middleware.RequireAdmin(), controllers.DeleteUserByID)

			// Routes administrateur
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireAdmin())
			{
				// Gestion des utilisateurs (admin seulement)
				admin.GET("/users", controllers.GetAllUsers)
				admin.PUT("/users/:id", controllers.UpdateUserByID)
				admin.DELETE("/users/:id", controllers.DeleteUserByID)
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

			// Familles routes
			familles := protected.Group("/familles")
			{
				// list (admin only)
				familles.GET("", middleware.RequireAdmin(), controllers.ListFamilles)

				familles.GET("/:id", controllers.GetFamilleByID)
				familles.PUT("/:id", controllers.UpdateFamille)
				familles.DELETE("/:id", middleware.RequireAdmin(), controllers.DeleteFamille)

				familles.GET("/:id/teachers", controllers.GetFamilleTeachers)
				familles.GET("/:id/missions", controllers.GetFamilleMissions)
				familles.GET("/:id/courses", controllers.GetFamilleCourses)
				familles.GET("/:id/payments", controllers.GetFamillePayments)
				familles.POST("/:id/reviews", controllers.PostFamilleReview)
				familles.GET("/:id/options", controllers.GetFamilleOptions)
			}

			// Missions routes
			missions := protected.Group("/missions")
			{
				missions.GET("", controllers.ListMissions)
				missions.POST("", controllers.CreateMission)
				missions.GET("/:id", controllers.GetMissionByID)
				missions.PUT("/:id", controllers.UpdateMission)
				missions.DELETE("/:id", controllers.DeleteMission)

				missions.GET("/:id/courses", controllers.GetMissionCourses)
				missions.GET("/:id/reports", controllers.GetMissionReports)
				missions.GET("/:id/payments", controllers.GetMissionPayments)

				missions.PUT("/:id/stop", controllers.StopMission)
				missions.PUT("/:id/extend", controllers.ExtendMission)
			}

			// Courses routes
			courses := protected.Group("/courses")
			{
				courses.GET("", controllers.ListCourses)
				courses.POST("", controllers.CreateCourse)
				courses.GET("/:id", controllers.GetCourseByID)
				courses.PUT("/:id", controllers.UpdateCourse)
				courses.DELETE("/:id", controllers.DeleteCourse)

				courses.PUT("/:id/schedule", controllers.ScheduleCourse)
				courses.PUT("/:id/cancel", controllers.CancelCourse)
				courses.PUT("/:id/complete", controllers.CompleteCourse)
				courses.POST("/:id/declare", controllers.DeclareCourse)
				courses.GET("/:id/payments", controllers.GetCoursePayments)
			}

			// Enseignants routes
			enseignants := protected.Group("/enseignants")
			{
				enseignants.GET("", controllers.ListEnseignants)
				enseignants.POST("", middleware.RequireAdmin(), controllers.CreateEnseignant)

				enseignants.GET("/:id", controllers.GetEnseignantByID)
				enseignants.PUT("/:id", controllers.UpdateEnseignant)
				enseignants.DELETE("/:id", middleware.RequireAdmin(), controllers.DeleteEnseignant)

				enseignants.GET("/:id/students", controllers.GetEnseignantStudents)
				enseignants.GET("/:id/missions", controllers.GetEnseignantMissions)
				enseignants.GET("/:id/courses", controllers.GetEnseignantCourses)
				enseignants.GET("/:id/payments", controllers.GetEnseignantPayments)
				enseignants.GET("/:id/reports", controllers.GetEnseignantReports)
				enseignants.GET("/:id/options", controllers.GetEnseignantOptions)

				enseignants.GET("/nearby", controllers.GetEnseignantsNearby)
			}

			// Offers routes
			offers := protected.Group("/offers")
			{
				offers.GET("", controllers.ListOffers)
				offers.POST("", controllers.CreateOffer)
				offers.GET("/:id", controllers.GetOfferByID)
				offers.PUT("/:id", controllers.UpdateOffer)
				offers.DELETE("/:id", controllers.DeleteOffer)

				offers.GET("/:id/options", controllers.GetOfferOptions)
				offers.PUT("/:id/close", controllers.CloseOffer)
				offers.GET("/active", controllers.ListActiveOffers)
				offers.GET("/search", controllers.SearchOffers)
			}

			// Options routes
			options := protected.Group("/options")
			{
				options.GET("", controllers.ListOptions)
				options.POST("", controllers.CreateOption)
				options.GET("/:id", controllers.GetOptionByID)
				options.PUT("/:id", controllers.UpdateOption)
				options.DELETE("/:id", controllers.DeleteOption)

				options.PUT("/:id/accept", controllers.AcceptOption)
				options.PUT("/:id/decline", controllers.DeclineOption)
				options.PUT("/:id/cancel", controllers.CancelOption)
				options.GET("/pending", controllers.ListPendingOptions)
				options.GET("/expiring", controllers.ListExpiringOptions)
			}

			// Addresses routes
			addresses := protected.Group("/addresses")
			{
				addresses.GET("", middleware.RequireAdmin(), controllers.ListAddresses)
				addresses.POST("", controllers.CreateAddress)
				addresses.GET("/:id", controllers.GetAddressByID)
				addresses.PUT("/:id", controllers.UpdateAddress)
				addresses.DELETE("/:id", controllers.DeleteAddress)
				addresses.GET("/geocode", controllers.GeocodeAddress)
				addresses.GET("/route", controllers.CalculateRoute)
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
