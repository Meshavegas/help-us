package controllers

import (
	"api/database"
	"api/middleware"
	"api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EnseignantResponse regroupe le user + profile + relations
type EnseignantResponse struct {
	models.User
	Enseignant models.Enseignant `json:"enseignant"`
	Missions   []models.Mission  `json:"missions,omitempty"`
	Courses    []models.Course   `json:"courses,omitempty"`
	Reports    []models.Report   `json:"reports,omitempty"`
	Options    []models.Option   `json:"options,omitempty"`
}

// ListEnseignants godoc
// @Summary      Liste tous les enseignants
// @Description  Récupère la liste de tous les enseignants
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   EnseignantResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /enseignants [get]
func ListEnseignants(c *gin.Context) {
	var users []models.User
	if err := database.DB.Where("role = ?", models.RoleEnseignant).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}
	var resp []EnseignantResponse
	for _, u := range users {
		var prof models.Enseignant
		database.DB.Where("user_id = ?", u.ID).First(&prof)
		resp = append(resp, EnseignantResponse{User: u, Enseignant: prof})
	}
	c.JSON(http.StatusOK, resp)
}

// GetEnseignantByID godoc
// @Summary      Détails d'un enseignant
// @Description  Récupère les détails d'un enseignant spécifique
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'enseignant"
// @Success      200  {object}  EnseignantResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /enseignants/{id} [get]
func GetEnseignantByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil || user.Role != models.RoleEnseignant {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enseignant non trouvé"})
		return
	}

	var prof models.Enseignant
	database.DB.Where("user_id = ?", user.ID).First(&prof)

	// relations
	database.DB.Where("enseignant_id = ?", user.ID).Find(&prof.Missions)
	database.DB.Where("enseignant_id = ?", user.ID).Find(&prof.Courses)
	database.DB.Where("enseignant_id = ?", user.ID).Find(&prof.Reports)
	database.DB.Where("enseignant_id = ?", user.ID).Find(&prof.Options)

	c.JSON(http.StatusOK, EnseignantResponse{User: user, Enseignant: prof, Missions: prof.Missions, Courses: prof.Courses, Reports: prof.Reports, Options: prof.Options})
}

// CreateEnseignant godoc
// @Summary      Création d'un enseignant
// @Description  Crée un nouveau compte enseignant (admin seulement)
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.UserCreateRequest  true  "Données de l'enseignant"
// @Success      201  {object}  EnseignantResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /enseignants [post]
func CreateEnseignant(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé"})
		return
	}
	var req struct {
		Username       string `json:"username" binding:"required"`
		Email          string `json:"email" binding:"required"`
		Password       string `json:"password" binding:"required"`
		Phone          string `json:"phone_number"`
		Specialization string `json:"specialization"`
		Qualifications string `json:"qualifications"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password, // hash handled in hooks
		PhoneNumber: req.Phone,
		Role:        models.RoleEnseignant,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création utilisateur"})
		return
	}
	prof := models.Enseignant{
		UserID:         user.ID,
		Specialization: req.Specialization,
		Qualifications: req.Qualifications,
	}
	database.DB.Create(&prof)

	c.JSON(http.StatusCreated, EnseignantResponse{User: user, Enseignant: prof})
}

// UpdateEnseignant godoc
// @Summary      Mise à jour d'un enseignant
// @Description  Met à jour les informations d'un enseignant (admin ou propriétaire)
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                        true  "ID de l'enseignant"
// @Param        request  body      models.UserUpdateRequest   true  "Données de mise à jour"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /enseignants/{id} [put]
func UpdateEnseignant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	// admin or owner
	if !middleware.IsAdmin(c) {
		currentID, _ := middleware.GetUserID(c)
		if currentID != uint(id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé"})
			return
		}
	}
	var req struct {
		Username       string `json:"username"`
		Email          string `json:"email"`
		Phone          string `json:"phone_number"`
		Specialization string `json:"specialization"`
		Qualifications string `json:"qualifications"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil || user.Role != models.RoleEnseignant {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enseignant non trouvé"})
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

	var prof models.Enseignant
	database.DB.Where("user_id = ?", user.ID).First(&prof)
	if req.Specialization != "" {
		prof.Specialization = req.Specialization
	}
	if req.Qualifications != "" {
		prof.Qualifications = req.Qualifications
	}
	database.DB.Save(&prof)

	c.JSON(http.StatusOK, gin.H{"message": "Enseignant mis à jour"})
}

// DeleteEnseignant godoc
// @Summary      Suppression d'un enseignant
// @Description  Supprime un compte enseignant (admin seulement)
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'enseignant"
// @Success      204  {object}  nil
// @Failure      403  {object}  map[string]interface{}
// @Router       /enseignants/{id} [delete]
func DeleteEnseignant(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	database.DB.Delete(&models.User{}, id)
	database.DB.Where("user_id = ?", id).Delete(&models.Enseignant{})
	c.Status(http.StatusNoContent)
}

// GetEnseignantStudents godoc
// @Summary      Liste des élèves d'un enseignant
// @Description  Récupère la liste des élèves associés à un enseignant
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'enseignant"
// @Success      200  {array}   models.User
// @Router       /enseignants/{id}/students [get]
func GetEnseignantStudents(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var familyIDs []uint
	database.DB.Model(&models.Mission{}).Where("enseignant_id = ?", id).Distinct().Pluck("famille_id", &familyIDs)
	database.DB.Model(&models.Course{}).Where("enseignant_id = ?", id).Distinct().Pluck("famille_id", &familyIDs)
	var families []models.User
	if len(familyIDs) > 0 {
		database.DB.Where("id IN ?", familyIDs).Find(&families)
	}
	c.JSON(http.StatusOK, families)
}

// GetEnseignantMissions godoc
// @Summary      Liste des missions d'un enseignant
// @Description  Récupère la liste des missions associées à un enseignant
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'enseignant"
// @Success      200  {array}   models.Mission
// @Router       /enseignants/{id}/missions [get]
func GetEnseignantMissions(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var missions []models.Mission
	database.DB.Where("enseignant_id = ?", id).Find(&missions)
	c.JSON(http.StatusOK, missions)
}

// GetEnseignantCourses godoc
// @Summary      Liste des cours d'un enseignant
// @Description  Récupère la liste des cours associés à un enseignant
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'enseignant"
// @Success      200  {array}   models.Course
// @Router       /enseignants/{id}/courses [get]
func GetEnseignantCourses(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var courses []models.Course
	database.DB.Where("enseignant_id = ?", id).Find(&courses)
	c.JSON(http.StatusOK, courses)
}

// GetEnseignantPayments godoc
// @Summary      Liste des paiements d'un enseignant
// @Description  Récupère la liste des paiements associés à un enseignant
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'enseignant"
// @Success      200  {array}   models.Payment
// @Router       /enseignants/{id}/payments [get]
func GetEnseignantPayments(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var payments []models.Payment
	database.DB.Where("user_id = ?", id).Find(&payments)
	c.JSON(http.StatusOK, payments)
}

// GetEnseignantReports godoc
// @Summary      Liste des rapports d'un enseignant
// @Description  Récupère la liste des rapports soumis par un enseignant
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'enseignant"
// @Success      200  {array}   models.Report
// @Router       /enseignants/{id}/reports [get]
func GetEnseignantReports(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var reports []models.Report
	database.DB.Where("enseignant_id = ?", id).Find(&reports)
	c.JSON(http.StatusOK, reports)
}

// GetEnseignantOptions godoc
// @Summary      Liste des options d'un enseignant
// @Description  Récupère la liste des options associées à un enseignant
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'enseignant"
// @Success      200  {array}   models.Option
// @Router       /enseignants/{id}/options [get]
func GetEnseignantOptions(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var options []models.Option
	database.DB.Where("enseignant_id = ?", id).Find(&options)
	c.JSON(http.StatusOK, options)
}

// GetEnseignantsNearby godoc
// @Summary      Liste des enseignants à proximité
// @Description  Récupère la liste des enseignants proches d'une localisation donnée
// @Tags         enseignants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        lat     query     number  true   "Latitude"
// @Param        lng     query     number  true   "Longitude"
// @Param        radius  query     number  false  "Rayon de recherche en km"
// @Success      200  {array}   EnseignantResponse
// @Router       /enseignants/nearby [get]
func GetEnseignantsNearby(c *gin.Context) {
	// Stub: renvoie tous les enseignants pour l'instant
	ListEnseignants(c)
}
