package controllers

import (
	"api/database"
	"api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddressResponse struct {
	models.Address
}

// ListAddresses godoc
// @Summary      Liste toutes les adresses
// @Description  Récupère la liste de toutes les adresses (admin seulement)
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Address
// @Failure      500  {object}  map[string]interface{}
// @Router       /addresses [get]
func ListAddresses(c *gin.Context) {
	var addresses []models.Address
	if err := database.DB.Find(&addresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des adresses"})
		return
	}
	c.JSON(http.StatusOK, addresses)
}

// GetAddressByID godoc
// @Summary      Détails d'une adresse
// @Description  Récupère les détails d'une adresse spécifique
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'adresse"
// @Success      200  {object}  AddressResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /addresses/{id} [get]
func GetAddressByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var address models.Address
	if err := database.DB.Preload("User").First(&address, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Adresse non trouvée"})
		return
	}
	c.JSON(http.StatusOK, AddressResponse{Address: address})
}

// CreateAddress godoc
// @Summary      Création d'une adresse
// @Description  Crée une nouvelle adresse
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.AddressCreateRequest  true  "Données de l'adresse"
// @Success      201  {object}  AddressResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /addresses [post]
func CreateAddress(c *gin.Context) {
	var req models.AddressCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	address := models.Address{
		Street:     req.Street,
		City:       req.City,
		PostalCode: req.PostalCode,
		Country:    req.Country,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
	}
	if err := database.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'adresse"})
		return
	}
	c.JSON(http.StatusCreated, AddressResponse{Address: address})
}

// UpdateAddress godoc
// @Summary      Mise à jour d'une adresse
// @Description  Met à jour les informations d'une adresse existante
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                          true  "ID de l'adresse"
// @Param        request  body      models.AddressUpdateRequest  true  "Données de mise à jour"
// @Success      200  {object}  AddressResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /addresses/{id} [put]
func UpdateAddress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	var req models.AddressUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var address models.Address
	if err := database.DB.First(&address, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Adresse non trouvée"})
		return
	}
	if req.Street != "" {
		address.Street = req.Street
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.PostalCode != "" {
		address.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		address.Country = req.Country
	}
	if req.Latitude != 0 {
		address.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		address.Longitude = req.Longitude
	}
	database.DB.Save(&address)
	c.JSON(http.StatusOK, AddressResponse{Address: address})
}

// DeleteAddress godoc
// @Summary      Suppression d'une adresse
// @Description  Supprime une adresse existante
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID de l'adresse"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Router       /addresses/{id} [delete]
func DeleteAddress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	database.DB.Delete(&models.Address{}, id)
	c.Status(http.StatusNoContent)
}

// GeocodeAddress godoc
// @Summary      Géocodage d'une adresse
// @Description  Convertit une adresse en coordonnées géographiques (latitude/longitude)
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        street       query     string  true   "Rue"
// @Param        city        query     string  true   "Ville"
// @Param        postal_code query     string  true   "Code postal"
// @Param        country     query     string  true   "Pays"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /addresses/geocode [get]
func GeocodeAddress(c *gin.Context) {
	// Stub: retourne des coordonnées factices
	c.JSON(http.StatusOK, gin.H{
		"latitude":  48.8566,
		"longitude": 2.3522,
		"message":   "Géocodage simulé (Paris)",
	})
}

// CalculateRoute godoc
// @Summary      Calcul d'itinéraire
// @Description  Calcule l'itinéraire entre deux adresses
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        origin_id      query     int     true  "ID de l'adresse de départ"
// @Param        destination_id query     int     true  "ID de l'adresse d'arrivée"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /addresses/route [get]
func CalculateRoute(c *gin.Context) {
	originID, _ := strconv.ParseUint(c.Query("origin_id"), 10, 32)
	destID, _ := strconv.ParseUint(c.Query("destination_id"), 10, 32)
	var origin, dest models.Address
	database.DB.First(&origin, originID)
	database.DB.First(&dest, destID)
	// Stub: retourne une distance/durée factice
	c.JSON(http.StatusOK, gin.H{
		"distance": "5 km",
		"duration": "15 min",
		"route":    []string{"Point A", "Point B"},
	})
}
