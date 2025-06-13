package controllers

import (
	"api/database"
	"api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserResponse représente la réponse utilisateur avec les données associées
type UserResponse struct {
	models.User
	Famille       *models.Famille       `json:"famille,omitempty"`
	Enseignant    *models.Enseignant    `json:"enseignant,omitempty"`
	Administrator *models.Administrator `json:"administrator,omitempty"`
}

// GetAllUsers récupère tous les utilisateurs (admin seulement)
// @Summary      Liste tous les utilisateurs
// @Description  Récupère la liste de tous les utilisateurs (admin seulement)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   UserResponse               "Liste des utilisateurs"
// @Failure      401  {object}  map[string]interface{}     "Non authentifié"
// @Failure      403  {object}  map[string]interface{}     "Accès refusé"
// @Router       /api/users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Preload("Addresses").Preload("Payments").Preload("Resources").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs"})
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponse := UserResponse{User: user}

		// Charger les données spécialisées selon le rôle
		switch user.Role {
		case models.RoleFamille:
			var famille models.Famille
			if err := database.DB.Where("user_id = ?", user.ID).Preload("Missions").Preload("Courses").Preload("Options").First(&famille).Error; err == nil {
				userResponse.Famille = &famille
			}
		case models.RoleEnseignant:
			var enseignant models.Enseignant
			if err := database.DB.Where("user_id = ?", user.ID).Preload("Missions").Preload("Courses").Preload("Reports").Preload("Options").Preload("Offers").First(&enseignant).Error; err == nil {
				userResponse.Enseignant = &enseignant
			}
		case models.RoleAdministrator:
			var admin models.Administrator
			if err := database.DB.Where("user_id = ?", user.ID).Preload("Reports").Preload("Offers").Preload("Resources").First(&admin).Error; err == nil {
				userResponse.Administrator = &admin
			}
		}

		userResponses = append(userResponses, userResponse)
	}

	c.JSON(http.StatusOK, userResponses)
}

// GetUserByID récupère un utilisateur spécifique par son ID
// @Summary      Récupère un utilisateur par ID
// @Description  Obtient les détails d'un utilisateur spécifique
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                        true  "ID de l'utilisateur"
// @Success      200  {object}  UserResponse               "Détails de l'utilisateur"
// @Failure      400  {object}  map[string]interface{}     "ID invalide"
// @Failure      401  {object}  map[string]interface{}     "Non authentifié"
// @Failure      404  {object}  map[string]interface{}     "Utilisateur non trouvé"
// @Router       /api/users/{id} [get]
func GetUserByID(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var user models.User
	if err := database.DB.Preload("Addresses").Preload("Payments").Preload("Resources").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	userResponse := UserResponse{User: user}

	// Charger les données spécialisées selon le rôle
	switch user.Role {
	case models.RoleFamille:
		var famille models.Famille
		if err := database.DB.Where("user_id = ?", user.ID).Preload("Missions").Preload("Courses").Preload("Options").First(&famille).Error; err == nil {
			userResponse.Famille = &famille
		}
	case models.RoleEnseignant:
		var enseignant models.Enseignant
		if err := database.DB.Where("user_id = ?", user.ID).Preload("Missions").Preload("Courses").Preload("Reports").Preload("Options").Preload("Offers").First(&enseignant).Error; err == nil {
			userResponse.Enseignant = &enseignant
		}
	case models.RoleAdministrator:
		var admin models.Administrator
		if err := database.DB.Where("user_id = ?", user.ID).Preload("Reports").Preload("Offers").Preload("Resources").First(&admin).Error; err == nil {
			userResponse.Administrator = &admin
		}
	}

	c.JSON(http.StatusOK, userResponse)
}

// UpdateUserByID met à jour un utilisateur spécifique (admin seulement)
// @Summary      Met à jour un utilisateur
// @Description  Met à jour les informations d'un utilisateur (admin seulement)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int                        true  "ID de l'utilisateur"
// @Param        request body      models.UserUpdateRequest   true  "Données de mise à jour"
// @Success      200     {object}  UserResponse               "Utilisateur mis à jour"
// @Failure      400     {object}  map[string]interface{}     "Erreur de validation"
// @Failure      401     {object}  map[string]interface{}     "Non authentifié"
// @Failure      403     {object}  map[string]interface{}     "Accès refusé"
// @Failure      404     {object}  map[string]interface{}     "Utilisateur non trouvé"
// @Router       /api/users/{id} [put]
func UpdateUserByID(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Mettre à jour les champs de base
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'utilisateur"})
		return
	}

	// Mettre à jour les données spécialisées selon le rôle
	switch user.Role {
	case models.RoleFamille:
		if req.FamilyName != "" {
			var famille models.Famille
			if err := database.DB.Where("user_id = ?", user.ID).First(&famille).Error; err == nil {
				famille.FamilyName = req.FamilyName
				database.DB.Save(&famille)
			}
		}
	case models.RoleEnseignant:
		var enseignant models.Enseignant
		if err := database.DB.Where("user_id = ?", user.ID).First(&enseignant).Error; err == nil {
			if req.Specialization != "" {
				enseignant.Specialization = req.Specialization
			}
			if req.Qualifications != "" {
				enseignant.Qualifications = req.Qualifications
			}
			database.DB.Save(&enseignant)
		}
	}

	// Recharger l'utilisateur avec toutes les relations
	if err := database.DB.Preload("Addresses").Preload("Payments").Preload("Resources").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du rechargement de l'utilisateur"})
		return
	}

	userResponse := UserResponse{User: user}
	c.JSON(http.StatusOK, userResponse)
}

// DeleteUserByID supprime un utilisateur (admin seulement)
// @Summary      Supprime un utilisateur
// @Description  Supprime un utilisateur du système (admin seulement)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                        true  "ID de l'utilisateur"
// @Success      200  {object}  map[string]interface{}     "Utilisateur supprimé"
// @Failure      400  {object}  map[string]interface{}     "ID invalide"
// @Failure      401  {object}  map[string]interface{}     "Non authentifié"
// @Failure      403  {object}  map[string]interface{}     "Accès refusé"
// @Failure      404  {object}  map[string]interface{}     "Utilisateur non trouvé"
// @Router       /api/users/{id} [delete]
func DeleteUserByID(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Supprimer les données spécialisées selon le rôle
	switch user.Role {
	case models.RoleFamille:
		database.DB.Where("user_id = ?", user.ID).Delete(&models.Famille{})
	case models.RoleEnseignant:
		database.DB.Where("user_id = ?", user.ID).Delete(&models.Enseignant{})
	case models.RoleAdministrator:
		database.DB.Where("user_id = ?", user.ID).Delete(&models.Administrator{})
	}

	// Supprimer l'utilisateur (soft delete)
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de l'utilisateur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur supprimé avec succès"})
}

// GetUserAddresses récupère toutes les adresses d'un utilisateur
// @Summary      Récupère les adresses d'un utilisateur
// @Description  Obtient toutes les adresses associées à un utilisateur
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                        true  "ID de l'utilisateur"
// @Success      200  {array}   models.Address             "Liste des adresses"
// @Failure      400  {object}  map[string]interface{}     "ID invalide"
// @Failure      401  {object}  map[string]interface{}     "Non authentifié"
// @Failure      404  {object}  map[string]interface{}     "Utilisateur non trouvé"
// @Router       /api/users/{id}/addresses [get]
func GetUserAddresses(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	// Vérifier que l'utilisateur existe
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var addresses []models.Address
	if err := database.DB.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des adresses"})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// GetUserPayments récupère tous les paiements d'un utilisateur
// @Summary      Récupère les paiements d'un utilisateur
// @Description  Obtient tous les paiements associés à un utilisateur
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                        true  "ID de l'utilisateur"
// @Success      200  {array}   models.Payment             "Liste des paiements"
// @Failure      400  {object}  map[string]interface{}     "ID invalide"
// @Failure      401  {object}  map[string]interface{}     "Non authentifié"
// @Failure      404  {object}  map[string]interface{}     "Utilisateur non trouvé"
// @Router       /api/users/{id}/payments [get]
func GetUserPayments(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	// Vérifier que l'utilisateur existe
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var payments []models.Payment
	if err := database.DB.Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des paiements"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetUserResources récupère toutes les ressources accessibles à un utilisateur
// @Summary      Récupère les ressources d'un utilisateur
// @Description  Obtient toutes les ressources accessibles à un utilisateur
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                        true  "ID de l'utilisateur"
// @Success      200  {array}   models.Resource            "Liste des ressources"
// @Failure      400  {object}  map[string]interface{}     "ID invalide"
// @Failure      401  {object}  map[string]interface{}     "Non authentifié"
// @Failure      404  {object}  map[string]interface{}     "Utilisateur non trouvé"
// @Router       /api/users/{id}/resources [get]
func GetUserResources(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	// Vérifier que l'utilisateur existe
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var resources []models.Resource
	if err := database.DB.Model(&user).Association("Resources").Find(&resources); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des ressources"})
		return
	}

	c.JSON(http.StatusOK, resources)
}
