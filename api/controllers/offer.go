package controllers

import (
	"api/database"
	"api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OfferResponse struct {
	models.Offer
	Options []models.Option `json:"options,omitempty"`
}

// ListOffers godoc
// @Summary      Liste toutes les offres
// @Description  Récupère la liste des offres avec possibilité de filtrage par statut, sujet et niveau
// @Tags         offers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status   query     string  false  "Statut de l'offre"
// @Param        subject  query     string  false  "Matière enseignée"
// @Param        level    query     string  false  "Niveau d'étude"
// @Success      200  {array}   OfferResponse
// @Failure      500  {object}  map[string]interface{}
// @Router      /offers [get]
func ListOffers(c *gin.Context) {
	var offers []models.Offer
	query := database.DB
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if subject := c.Query("subject"); subject != "" {
		query = query.Where("subject = ?", subject)
	}
	if level := c.Query("level"); level != "" {
		query = query.Where("level = ?", level)
	}
	if err := query.Preload("Options").Find(&offers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des offres"})
		return
	}
	var resp []OfferResponse
	for _, o := range offers {
		resp = append(resp, OfferResponse{Offer: o, Options: o.Options})
	}
	c.JSON(http.StatusOK, resp)
}

// GetOfferByID godoc
// @Summary      Détails d'une offre
// @Description  Récupère les détails d'une offre spécifique
// @Tags         offers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'offre"
// @Success      200  {object}  OfferResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /offers/{id} [get]
func GetOfferByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var offer models.Offer
	if err := database.DB.Preload("Options").First(&offer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offre non trouvée"})
		return
	}
	c.JSON(http.StatusOK, OfferResponse{Offer: offer, Options: offer.Options})
}

// CreateOffer godoc
// @Summary      Création d'une offre
// @Description  Crée une nouvelle offre
// @Tags         offers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.OfferCreateRequest  true  "Données de l'offre"
// @Success      201  {object}  OfferResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router      /offers [post]
func CreateOffer(c *gin.Context) {
	var req models.OfferCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offer := models.Offer{
		Title:           req.Title,
		Description:     req.Description,
		HourlyRate:      req.HourlyRate,
		Requirements:    req.Requirements,
		Subject:         req.Subject,
		Level:           req.Level,
		Status:          models.OfferStatusOpen,
		PublicationDate: time.Now(),
	}
	if err := database.DB.Create(&offer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'offre"})
		return
	}
	c.JSON(http.StatusCreated, OfferResponse{Offer: offer})
}

// UpdateOffer godoc
// @Summary      Mise à jour d'une offre
// @Description  Met à jour les informations d'une offre existante
// @Tags         offers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                       true  "ID de l'offre"
// @Param        request  body      models.OfferUpdateRequest  true  "Données de mise à jour"
// @Success      200  {object}  OfferResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /offers/{id} [put]
func UpdateOffer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var req models.OfferUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var offer models.Offer
	if err := database.DB.First(&offer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offre non trouvée"})
		return
	}
	if req.Title != "" {
		offer.Title = req.Title
	}
	if req.Description != "" {
		offer.Description = req.Description
	}
	if req.HourlyRate != 0 {
		offer.HourlyRate = req.HourlyRate
	}
	if req.Status != "" {
		offer.Status = req.Status
	}
	if req.Requirements != "" {
		offer.Requirements = req.Requirements
	}
	if req.Subject != "" {
		offer.Subject = req.Subject
	}
	if req.Level != "" {
		offer.Level = req.Level
	}
	database.DB.Save(&offer)
	c.JSON(http.StatusOK, OfferResponse{Offer: offer})
}

// DeleteOffer godoc
// @Summary      Suppression d'une offre
// @Description  Supprime une offre existante
// @Tags         offers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'offre"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /offers/{id} [delete]
func DeleteOffer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	database.DB.Delete(&models.Offer{}, id)
	c.Status(http.StatusNoContent)
}

// GetOfferOptions godoc
// @Summary      Liste des options d'une offre
// @Description  Récupère toutes les options associées à une offre
// @Tags         offers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'offre"
// @Success      200  {array}   models.Option
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /offers/{id}/options [get]
func GetOfferOptions(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var options []models.Option
	database.DB.Where("offer_id = ?", id).Find(&options)
	c.JSON(http.StatusOK, options)
}

// CloseOffer godoc
// @Summary      Fermeture d'une offre
// @Description  Marque une offre comme fermée
// @Tags         offers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'offre"
// @Success      200  {object}  OfferResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router      /offers/{id}/close [put]
func CloseOffer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var offer models.Offer
	if err := database.DB.First(&offer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offre non trouvée"})
		return
	}
	offer.Status = models.OfferStatusClosed
	database.DB.Save(&offer)
	c.JSON(http.StatusOK, OfferResponse{Offer: offer})
}

// ListActiveOffers godoc
// @Summary      Liste des offres actives
// @Description  Récupère la liste des offres avec le statut 'open'
// @Tags         offers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   OfferResponse
// @Failure      500  {object}  map[string]interface{}
// @Router      /offers/active [get]
func ListActiveOffers(c *gin.Context) {
	var offers []models.Offer
	database.DB.Where("status = ?", models.OfferStatusOpen).Find(&offers)
	c.JSON(http.StatusOK, offers)
}

// SearchOffers godoc
// @Summary      Recherche d'offres
// @Description  Recherche des offres selon différents critères
// @Tags         offers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        query  query     string  false  "Terme de recherche"
// @Success      200  {array}   OfferResponse
// @Failure      500  {object}  map[string]interface{}
// @Router      /offers/search [get]
func SearchOffers(c *gin.Context) {
	// Pour l'instant, même logique que ListOffers avec plus de filtres
	ListOffers(c)
}
