package controllers

import (
	"api/database"
	"api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// MissionResponse représente la mission avec ses relations principales
// pour éviter des boucles d'objets trop profondes dans le JSON.
type MissionResponse struct {
	models.Mission
	Courses  []models.Course  `json:"courses,omitempty"`
	Reports  []models.Report  `json:"reports,omitempty"`
	Payments []models.Payment `json:"payments,omitempty"`
}

// ListMissions godoc
// @Summary      Liste toutes les missions
// @Description  Récupère la liste des missions avec possibilité de filtrage par statut, enseignant ou famille
// @Tags         missions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status         query     string  false  "Statut de la mission"
// @Param        enseignant_id  query     int     false  "ID de l'enseignant"
// @Param        famille_id     query     int     false  "ID de la famille"
// @Success      200  {array}   MissionResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /missions [get]
func ListMissions(c *gin.Context) {
	var missions []models.Mission

	query := database.DB

	// Filtres éventuels
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if enseignantID := c.Query("enseignant_id"); enseignantID != "" {
		query = query.Where("enseignant_id = ?", enseignantID)
	}
	if familleID := c.Query("famille_id"); familleID != "" {
		query = query.Where("famille_id = ?", familleID)
	}

	if err := query.Preload("Courses").Preload("Reports").Find(&missions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des missions"})
		return
	}

	var resp []MissionResponse
	for _, m := range missions {
		resp = append(resp, MissionResponse{
			Mission: m,
			Courses: m.Courses,
			Reports: m.Reports,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// GetMissionByID godoc
// @Summary      Détails d'une mission
// @Description  Récupère les détails d'une mission spécifique
// @Tags         missions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la mission"
// @Success      200  {object}  MissionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /missions/{id} [get]
func GetMissionByID(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mission invalide"})
		return
	}

	var mission models.Mission
	if err := database.DB.Preload("Courses").Preload("Reports").First(&mission, missionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission non trouvée"})
		return
	}

	c.JSON(http.StatusOK, MissionResponse{Mission: mission, Courses: mission.Courses, Reports: mission.Reports})
}

// CreateMission godoc
// @Summary      Création d'une mission
// @Description  Crée une nouvelle mission
// @Tags         missions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.MissionCreateRequest  true  "Données de la mission"
// @Success      201  {object}  MissionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /missions [post]
func CreateMission(c *gin.Context) {
	var req models.MissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mission := models.Mission{
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Description:  req.Description,
		EnseignantID: req.EnseignantID,
	}
	// FamilleID: on peut récupérer depuis contexte utilisateur si rôle famille
	if familleIDStr := c.Query("famille_id"); familleIDStr != "" {
		if id, err := strconv.ParseUint(familleIDStr, 10, 32); err == nil {
			mission.FamilleID = uint(id)
		}
	}

	if err := database.DB.Create(&mission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la mission"})
		return
	}

	c.JSON(http.StatusCreated, MissionResponse{Mission: mission})
}

// UpdateMission godoc
// @Summary Mise à jour d'une mission
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID mission"
// @Param request body models.MissionUpdateRequest true "Champs à mettre à jour"
// @Success 200 {object} MissionResponse
// @Failure 404 {object} map[string]string
// @Router /missions/{id} [put]
func UpdateMission(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mission invalide"})
		return
	}

	var req models.MissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var mission models.Mission
	if err := database.DB.First(&mission, missionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission non trouvée"})
		return
	}

	// Appliquer les modifications
	if req.EndDate != nil {
		mission.EndDate = req.EndDate
	}
	if req.Description != "" {
		mission.Description = req.Description
	}
	if req.Status != "" {
		mission.Status = req.Status
	}

	if err := database.DB.Save(&mission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, MissionResponse{Mission: mission})
}

// DeleteMission godoc
// @Summary      Suppression d'une mission
// @Description  Supprime une mission existante
// @Tags         missions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la mission"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /missions/{id} [delete]
func DeleteMission(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mission invalide"})
		return
	}

	if err := database.DB.Delete(&models.Mission{}, missionID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.Status(http.StatusNoContent)
}

// StopMission met le statut à stopped
func StopMission(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var mission models.Mission
	if err := database.DB.First(&mission, missionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission non trouvée"})
		return
	}

	mission.Status = models.MissionStatusStopped
	end := time.Now()
	mission.EndDate = &end

	database.DB.Save(&mission)
	c.JSON(http.StatusOK, MissionResponse{Mission: mission})
}

// ExtendMission change la date de fin
func ExtendMission(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var payload struct {
		EndDate time.Time `json:"end_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var mission models.Mission
	if err := database.DB.First(&mission, missionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission non trouvée"})
		return
	}

	mission.EndDate = &payload.EndDate
	mission.Status = models.MissionStatusActive
	database.DB.Save(&mission)

	c.JSON(http.StatusOK, MissionResponse{Mission: mission})
}

// GetMissionCourses godoc
// @Summary      Liste des cours d'une mission
// @Description  Récupère la liste des cours associés à une mission
// @Tags         missions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la mission"
// @Success      200  {array}   models.Course
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /missions/{id}/courses [get]
func GetMissionCourses(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mission invalide"})
		return
	}

	var courses []models.Course
	if err := database.DB.Where("mission_id = ?", missionID).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des cours"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// GetMissionReports godoc
// @Summary      Liste des rapports d'une mission
// @Description  Récupère la liste des rapports associés à une mission
// @Tags         missions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la mission"
// @Success      200  {array}   models.Report
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /missions/{id}/reports [get]
func GetMissionReports(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mission invalide"})
		return
	}

	var reports []models.Report
	if err := database.DB.Where("mission_id = ?", missionID).Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des rapports"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// GetMissionPayments godoc
// @Summary      Liste des paiements d'une mission
// @Description  Récupère la liste des paiements associés à une mission
// @Tags         missions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de la mission"
// @Success      200  {array}   models.Payment
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /missions/{id}/payments [get]
func GetMissionPayments(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mission invalide"})
		return
	}

	// Récupérer la mission avec ses cours
	var mission models.Mission
	if err := database.DB.Preload("Courses").First(&mission, missionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission non trouvée"})
		return
	}

	if len(mission.Courses) == 0 {
		c.JSON(http.StatusOK, []models.Payment{})
		return
	}

	var courseIDs []uint
	for _, course := range mission.Courses {
		courseIDs = append(courseIDs, course.ID)
	}

	var payments []models.Payment
	if err := database.DB.Where("course_id IN ?", courseIDs).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des paiements"})
		return
	}
	c.JSON(http.StatusOK, payments)
}
