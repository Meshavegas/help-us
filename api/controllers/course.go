package controllers

import (
	"api/database"
	"api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseResponse struct {
	models.Course
	Payments []models.Payment `json:"payments,omitempty"`
}

// ListCourses godoc
// @Summary      Liste tous les cours
// @Description  Récupère la liste des cours avec possibilité de filtrage par statut, enseignant, famille ou mission
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status         query     string  false  "Statut du cours"
// @Param        enseignant_id  query     int     false  "ID de l'enseignant"
// @Param        famille_id     query     int     false  "ID de la famille"
// @Param        mission_id     query     int     false  "ID de la mission"
// @Success      200  {array}   CourseResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/courses [get]
func ListCourses(c *gin.Context) {
	var courses []models.Course
	query := database.DB

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if enseignantID := c.Query("enseignant_id"); enseignantID != "" {
		query = query.Where("enseignant_id = ?", enseignantID)
	}
	if familleID := c.Query("famille_id"); familleID != "" {
		query = query.Where("famille_id = ?", familleID)
	}
	if missionID := c.Query("mission_id"); missionID != "" {
		query = query.Where("mission_id = ?", missionID)
	}
	if err := query.Preload("Payments").Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des cours"})
		return
	}
	var resp []CourseResponse
	for _, course := range courses {
		resp = append(resp, CourseResponse{Course: course, Payments: course.Payments})
	}
	c.JSON(http.StatusOK, resp)
}

// GetCourseByID godoc
// @Summary      Détails d'un cours
// @Description  Récupère les détails d'un cours spécifique
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID du cours"
// @Success      200  {object}  CourseResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/courses/{id} [get]
func GetCourseByID(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}
	var course models.Course
	if err := database.DB.Preload("Payments").First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cours non trouvé"})
		return
	}
	c.JSON(http.StatusOK, CourseResponse{Course: course, Payments: course.Payments})
}

// CreateCourse godoc
// @Summary      Création d'un cours
// @Description  Crée un nouveau cours
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.CourseCreateRequest  true  "Données du cours"
// @Success      201  {object}  CourseResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/courses [post]
func CreateCourse(c *gin.Context) {
	var req models.CourseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	course := models.Course{
		ScheduledTime: req.ScheduledTime,
		Duration:      req.Duration,
		Location:      req.Location,
		EnseignantID:  req.EnseignantID,
		AddressID:     req.AddressID,
	}
	if missionID := c.Query("mission_id"); missionID != "" {
		if id, err := strconv.ParseUint(missionID, 10, 32); err == nil {
			course.MissionID = uint(id)
		}
	}
	if familleID := c.Query("famille_id"); familleID != "" {
		if id, err := strconv.ParseUint(familleID, 10, 32); err == nil {
			course.FamilleID = uint(id)
		}
	}
	if err := database.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du cours"})
		return
	}
	c.JSON(http.StatusCreated, CourseResponse{Course: course})
}

// UpdateCourse godoc
// @Summary      Mise à jour d'un cours
// @Description  Met à jour les informations d'un cours existant
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                       true  "ID du cours"
// @Param        request  body      models.CourseUpdateRequest  true  "Données de mise à jour"
// @Success      200  {object}  CourseResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/courses/{id} [put]
func UpdateCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}
	var req models.CourseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cours non trouvé"})
		return
	}
	if req.ScheduledTime != nil {
		course.ScheduledTime = *req.ScheduledTime
	}
	if req.Duration != nil {
		course.Duration = *req.Duration
	}
	if req.Location != "" {
		course.Location = req.Location
	}
	if req.Status != "" {
		course.Status = req.Status
	}
	if err := database.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du cours"})
		return
	}
	c.JSON(http.StatusOK, CourseResponse{Course: course})
}

// DeleteCourse godoc
// @Summary      Suppression d'un cours
// @Description  Supprime un cours existant
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID du cours"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/courses/{id} [delete]
func DeleteCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}
	if err := database.DB.Delete(&models.Course{}, courseID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du cours"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ScheduleCourse godoc
// @Summary      Planification d'un cours
// @Description  Planifie un cours à une date spécifique
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                          true  "ID du cours"
// @Param        request  body      models.CourseScheduleRequest  true  "Données de planification"
// @Success      200  {object}  CourseResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/courses/{id}/schedule [put]
func ScheduleCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cours non trouvé"})
		return
	}
	course.Status = models.CourseStatusScheduled
	if err := database.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la planification du cours"})
		return
	}
	c.JSON(http.StatusOK, CourseResponse{Course: course})
}

// CancelCourse godoc
// @Summary      Annulation d'un cours
// @Description  Annule un cours planifié
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID du cours"
// @Success      200  {object}  CourseResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/courses/{id}/cancel [put]
func CancelCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cours non trouvé"})
		return
	}
	course.Status = models.CourseStatusCancelled
	if err := database.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'annulation du cours"})
		return
	}
	c.JSON(http.StatusOK, CourseResponse{Course: course})
}

// CompleteCourse godoc
// @Summary      Validation d'un cours
// @Description  Marque un cours comme terminé
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID du cours"
// @Success      200  {object}  CourseResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/courses/{id}/complete [put]
func CompleteCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cours non trouvé"})
		return
	}
	course.Status = models.CourseStatusCompleted
	if err := database.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la validation du cours"})
		return
	}
	c.JSON(http.StatusOK, CourseResponse{Course: course})
}

// DeclareCourse godoc
// @Summary      Déclaration des heures d'un cours
// @Description  Déclare les heures effectuées pour un cours
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                          true  "ID du cours"
// @Param        request  body      map[string]float64           true  "Heures effectuées"
// @Success      200  {object}  CourseResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/courses/{id}/declare [post]
func DeclareCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}
	var payload struct {
		Hours float64 `json:"hours" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cours non trouvé"})
		return
	}
	// Ici, on pourrait stocker les heures déclarées dans un champ ou un modèle associé
	// Pour l'instant, on change juste le statut
	course.Status = models.CourseStatusInProgress
	if err := database.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la déclaration du cours"})
		return
	}
	c.JSON(http.StatusOK, CourseResponse{Course: course})
}

// GetCoursePayments godoc
// @Summary      Liste des paiements d'un cours
// @Description  Récupère tous les paiements associés à un cours
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID du cours"
// @Success      200  {array}   models.Payment
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/v1/courses/{id}/payments [get]
func GetCoursePayments(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}
	var payments []models.Payment
	if err := database.DB.Where("course_id = ?", courseID).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des paiements"})
		return
	}
	c.JSON(http.StatusOK, payments)
}
