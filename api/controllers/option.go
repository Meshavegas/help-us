package controllers

import (
	"api/database"
	"api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OptionResponse struct {
	models.Option
}

// ListOptions godoc
// @Summary      Liste toutes les options
// @Description  Récupère la liste des options avec possibilité de filtrage
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status         query     string  false  "Statut de l'option"
// @Param        enseignant_id  query     int     false  "ID de l'enseignant"
// @Param        famille_id     query     int     false  "ID de la famille"
// @Param        offer_id       query     int     false  "ID de l'offre"
// @Success      200  {array}   models.Option
// @Failure      500  {object}  map[string]interface{}
// @Router      /options [get]
func ListOptions(c *gin.Context) {
	var options []models.Option
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
	if offerID := c.Query("offer_id"); offerID != "" {
		query = query.Where("offer_id = ?", offerID)
	}
	if err := query.Find(&options).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des options"})
		return
	}
	c.JSON(http.StatusOK, options)
}

// GetOptionByID godoc
// @Summary      Détails d'une option
// @Description  Récupère les détails d'une option spécifique
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'option"
// @Success      200  {object}  OptionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /options/{id} [get]

func GetOptionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var option models.Option
	if err := database.DB.First(&option, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option non trouvée"})
		return
	}
	c.JSON(http.StatusOK, OptionResponse{Option: option})
}

// CreateOption godoc
// @Summary      Création d'une option
// @Description  Crée une nouvelle option sur une offre
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.OptionCreateRequest  true  "Données de l'option"
// @Success      201  {object}  OptionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router      /options [post]
func CreateOption(c *gin.Context) {
	var req models.OptionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	option := models.Option{
		EnseignantID:   req.EnseignantID,
		FamilleID:      req.FamilleID,
		OfferID:        req.OfferID,
		Description:    req.Description,
		Status:         models.OptionStatusActive,
		CreationDate:   time.Now(),
		ExpirationDate: req.ExpirationDate,
	}
	if option.ExpirationDate.IsZero() {
		option.ExpirationDate = time.Now().AddDate(0, 0, 7)
	}
	if err := database.DB.Create(&option).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'option"})
		return
	}
	c.JSON(http.StatusCreated, OptionResponse{Option: option})
}

// UpdateOption godoc
// @Summary      Mise à jour d'une option
// @Description  Met à jour les informations d'une option existante
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                         true  "ID de l'option"
// @Param        request  body      models.OptionUpdateRequest  true  "Données de mise à jour"
// @Success      200  {object}  OptionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /options/{id} [put]
func UpdateOption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var req models.OptionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var option models.Option
	if err := database.DB.First(&option, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option non trouvée"})
		return
	}
	if req.Status != "" {
		option.Status = req.Status
	}
	if req.Description != "" {
		option.Description = req.Description
	}
	if req.ExpirationDate != nil {
		option.ExpirationDate = *req.ExpirationDate
	}
	database.DB.Save(&option)
	c.JSON(http.StatusOK, OptionResponse{Option: option})
}

// DeleteOption godoc
// @Summary      Suppression d'une option
// @Description  Supprime une option existante
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'option"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /options/{id} [delete]
func DeleteOption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	database.DB.Delete(&models.Option{}, id)
	c.Status(http.StatusNoContent)
}

// AcceptOption godoc
// @Summary      Acceptation d'une option
// @Description  Accepte une option en attente
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'option"
// @Success      200  {object}  OptionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /options/{id}/accept [put]

// AcceptOption - PUT/options/:id/accept
func AcceptOption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var option models.Option
	if err := database.DB.First(&option, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option non trouvée"})
		return
	}
	option.Status = models.OptionStatusAccepted
	database.DB.Save(&option)
	c.JSON(http.StatusOK, OptionResponse{Option: option})
}

// DeclineOption godoc
// @Summary      Refus d'une option
// @Description  Refuse une option en attente
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'option"
// @Success      200  {object}  OptionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /options/{id}/decline [put]
func DeclineOption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var option models.Option
	if err := database.DB.First(&option, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option non trouvée"})
		return
	}
	option.Status = models.OptionStatusExpired
	database.DB.Save(&option)
	c.JSON(http.StatusOK, OptionResponse{Option: option})
}

// CancelOption godoc
// @Summary      Annulation d'une option
// @Description  Annule une option acceptée
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'option"
// @Success      200  {object}  OptionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /options/{id}/cancel [put]
func CancelOption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var option models.Option
	if err := database.DB.First(&option, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option non trouvée"})
		return
	}
	option.Status = models.OptionStatusCancelled
	database.DB.Save(&option)
	c.JSON(http.StatusOK, OptionResponse{Option: option})
}

// ListPendingOptions godoc
// @Summary      Liste des options en attente
// @Description  Récupère la liste des options avec le statut "en attente"
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Option
// @Failure      500  {object}  map[string]interface{}
// @Router      /options/pending [get]
func ListPendingOptions(c *gin.Context) {
	var options []models.Option
	database.DB.Where("status = ?", models.OptionStatusActive).Find(&options)
	c.JSON(http.StatusOK, options)
}

// ListExpiringOptions godoc
// @Summary      Liste des options expirant bientôt
// @Description  Récupère la liste des options qui vont bientôt expirer
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Option
// @Failure      500  {object}  map[string]interface{}
// @Router      /options/expiring [get]
func ListExpiringOptions(c *gin.Context) {
	var options []models.Option
	now := time.Now()
	soon := now.Add(48 * time.Hour)
	database.DB.Where("expiration_date BETWEEN ? AND ?", now, soon).Find(&options)
	c.JSON(http.StatusOK, options)
}

// RejectOption godoc
// @Summary      Rejet d'une option
// @Description  Rejette une option en attente
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'option"
// @Success      200  {object}  OptionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /options/{id}/reject [put]
func RejectOption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var option models.Option
	if err := database.DB.First(&option, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option non trouvée"})
		return
	}
	option.Status = models.OptionStatusCancelled
	database.DB.Save(&option)
	c.JSON(http.StatusOK, OptionResponse{Option: option})
}

// ExpireOption godoc
// @Summary      Expiration d'une option
// @Description  Marque une option comme expirée
// @Tags         options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'option"
// @Success      200  {object}  OptionResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /options/{id}/expire [put]
func ExpireOption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var option models.Option
	if err := database.DB.First(&option, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option non trouvée"})
		return
	}
	option.Status = models.OptionStatusExpired
	database.DB.Save(&option)
	c.JSON(http.StatusOK, OptionResponse{Option: option})
}
