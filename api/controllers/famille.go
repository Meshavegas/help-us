package controllers

import (
	"api/database"
	"api/middleware"
	"api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FamilleResponse regroupe l'utilisateur et son profil famille + relations principales
type FamilleResponse struct {
	models.User
	Famille  models.Famille   `json:"famille"`
	Missions []models.Mission `json:"missions,omitempty"`
	Courses  []models.Course  `json:"courses,omitempty"`
	Options  []models.Option  `json:"options,omitempty"`
}

// ListFamilles godoc
// @Summary      Liste toutes les familles
// @Description  Récupère la liste de toutes les familles
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   FamilleResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /familles [get]
func ListFamilles(c *gin.Context) {
	var users []models.User
	if err := database.DB.Where("role = ?", models.RoleFamille).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des familles"})
		return
	}

	var resp []FamilleResponse
	for _, u := range users {
		var fam models.Famille
		database.DB.Where("user_id = ?", u.ID).First(&fam)
		resp = append(resp, FamilleResponse{User: u, Famille: fam})
	}
	c.JSON(http.StatusOK, resp)
}

// GetFamilleByID godoc
// @Summary      Détails d'une famille
// @Description  Récupère les détails d'une famille spécifique
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la famille"
// @Success      200  {object}  FamilleResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /familles/{id} [get]
func GetFamilleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil || user.Role != models.RoleFamille {
		c.JSON(http.StatusNotFound, gin.H{"error": "Famille non trouvée"})
		return
	}

	var fam models.Famille
	database.DB.Where("user_id = ?", user.ID).First(&fam)

	// Relations rapides
	database.DB.Where("famille_id = ?", fam.UserID).Find(&fam.Missions)
	database.DB.Where("famille_id = ?", fam.UserID).Find(&fam.Courses)
	database.DB.Where("famille_id = ?", fam.UserID).Find(&fam.Options)

	c.JSON(http.StatusOK, FamilleResponse{User: user, Famille: fam, Missions: fam.Missions, Courses: fam.Courses, Options: fam.Options})
}

// UpdateFamille godoc
// @Summary      Mise à jour d'une famille
// @Description  Met à jour les informations d'une famille (admin ou propriétaire)
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                        true  "ID de la famille"
// @Param        request  body      models.UserUpdateRequest   true  "Données de mise à jour"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /familles/{id} [put]
func UpdateFamille(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Autorisation: admin ou propriétaire
	if !middleware.IsAdmin(c) {
		currentID, _ := middleware.GetUserID(c)
		if currentID != uint(id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé"})
			return
		}
	}

	var req struct {
		FamilyName string `json:"family_name"`
		Phone      string `json:"phone_number"`
		Email      string `json:"email"`
		Username   string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil || user.Role != models.RoleFamille {
		c.JSON(http.StatusNotFound, gin.H{"error": "Famille non trouvée"})
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.PhoneNumber = req.Phone
	}
	database.DB.Save(&user)

	var fam models.Famille
	database.DB.Where("user_id = ?", user.ID).First(&fam)
	if req.FamilyName != "" {
		fam.FamilyName = req.FamilyName
		database.DB.Save(&fam)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Famille mise à jour"})
}

// DeleteFamille godoc
// @Summary      Suppression d'une famille
// @Description  Supprime un compte famille (admin seulement)
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la famille"
// @Success      204  {object}  nil
// @Failure      403  {object}  map[string]interface{}
// @Router       /familles/{id} [delete]
func DeleteFamille(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	database.DB.Delete(&models.User{}, id)
	database.DB.Where("user_id = ?", id).Delete(&models.Famille{})
	c.Status(http.StatusNoContent)
}

// GetFamilleTeachers godoc
// @Summary      Liste des enseignants d'une famille
// @Description  Récupère la liste des enseignants associés à une famille
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la famille"
// @Success      200  {array}   models.User
// @Router       /familles/{id}/teachers [get]
func GetFamilleTeachers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	// Trouver enseignants via missions ou cours
	var teacherIDs []uint
	database.DB.Model(&models.Mission{}).Where("famille_id = ?", id).Distinct().Pluck("enseignant_id", &teacherIDs)
	database.DB.Model(&models.Course{}).Where("famille_id = ?", id).Distinct().Pluck("enseignant_id", &teacherIDs)

	var teachers []models.User
	if len(teacherIDs) > 0 {
		database.DB.Where("id IN ?", teacherIDs).Find(&teachers)
	}
	c.JSON(http.StatusOK, teachers)
}

// GetFamilleMissions godoc
// @Summary      Liste des missions d'une famille
// @Description  Récupère la liste des missions associées à une famille
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la famille"
// @Success      200  {array}   models.Mission
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /familles/{id}/missions [get]
func GetFamilleMissions(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var missions []models.Mission
	database.DB.Where("famille_id = ?", id).Find(&missions)
	c.JSON(http.StatusOK, missions)
}

// GetFamilleCourses godoc
// @Summary      Liste des cours d'une famille
// @Description  Récupère la liste des cours associés à une famille
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la famille"
// @Success      200  {array}   models.Course
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /familles/{id}/courses [get]
func GetFamilleCourses(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var courses []models.Course
	database.DB.Where("famille_id = ?", id).Find(&courses)
	c.JSON(http.StatusOK, courses)
}

// GetFamillePayments godoc
// @Summary      Liste des paiements d'une famille
// @Description  Récupère la liste des paiements associés à une famille
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la famille"
// @Success      200  {array}   models.Payment
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /familles/{id}/payments [get]
func GetFamillePayments(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var payments []models.Payment
	database.DB.Where("user_id = ?", id).Find(&payments)
	c.JSON(http.StatusOK, payments)
}

// PostFamilleReview godoc
// @Summary      Ajouter un avis sur une famille
// @Description  Ajoute un avis sur une famille (non implémenté)
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int  true  "ID de la famille"
// @Param        review  body      object  true  "Données de l'avis"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /familles/{id}/reviews [post]
func PostFamilleReview(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Fonctionnalité review non implémentée"})
}

// GetFamilleOptions godoc
// @Summary      Liste des options d'une famille
// @Description  Récupère la liste des options associées à une famille
// @Tags         familles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la famille"
// @Success      200  {array}   models.Option
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /familles/{id}/options [get]
func GetFamilleOptions(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var options []models.Option
	database.DB.Where("famille_id = ?", id).Find(&options)
	c.JSON(http.StatusOK, options)
}
