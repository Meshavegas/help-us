package controllers

import (
	"api/database"
	"api/models"
	"api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest représente la structure de la requête d'inscription
type RegisterRequest struct {
	Username       string          `json:"username" binding:"required,min=3,max=50"`
	Email          string          `json:"email" binding:"required,email"`
	Password       string          `json:"password" binding:"required,min=6"`
	Role           models.UserRole `json:"role" binding:"required"`
	PhoneNumber    string          `json:"phone_number"`
	FamilyName     string          `json:"family_name,omitempty"`    // Pour les familles
	Specialization string          `json:"specialization,omitempty"` // Pour les enseignants
	Qualifications string          `json:"qualifications,omitempty"` // Pour les enseignants
}

// LoginRequest représente la structure de la requête de connexion
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse représente la réponse d'authentification
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Register godoc
// @Summary      Inscription d'un utilisateur
// @Description  Créer un nouveau compte utilisateur
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Données d'inscription"
// @Success      201      {object}  AuthResponse     "Utilisateur créé avec succès"
// @Failure      400      {object}  map[string]interface{}       "Erreur de validation"
// @Failure      409      {object}  map[string]interface{}       "Utilisateur déjà existant"
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si l'email existe déjà
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Un utilisateur avec cet email existe déjà"})
		return
	}

	// Vérifier si le nom d'utilisateur existe déjà
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Un utilisateur avec ce nom d'utilisateur existe déjà"})
		return
	}

	// Créer le nouvel utilisateur (le mot de passe sera haché automatiquement par BeforeCreate)
	user := models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password, // Le mot de passe sera haché par BeforeCreate
		Role:        req.Role,
		PhoneNumber: req.PhoneNumber,
		IsActive:    true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	// Créer les modèles spécialisés selon le rôle
	switch req.Role {
	case models.RoleFamille:
		famille := models.Famille{
			UserID:     user.ID,
			FamilyName: req.FamilyName,
		}
		if err := database.DB.Create(&famille).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du profil famille"})
			return
		}
	case models.RoleEnseignant:
		enseignant := models.Enseignant{
			UserID:         user.ID,
			Specialization: req.Specialization,
			Qualifications: req.Qualifications,
		}
		if err := database.DB.Create(&enseignant).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du profil enseignant"})
			return
		}
	case models.RoleAdministrator:
		admin := models.Administrator{
			UserID: user.ID,
		}
		if err := database.DB.Create(&admin).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du profil administrateur"})
			return
		}
	}

	// Générer le token JWT
	token, err := utils.GenerateJWT(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	// Masquer le mot de passe dans la réponse
	user.Password = ""

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login gère la connexion d'un utilisateur
// @Summary      Connexion d'un utilisateur
// @Description  Authentifier un utilisateur et retourner un token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest   true  "Identifiants de connexion"
// @Success      200      {object}  AuthResponse   "Connexion réussie"
// @Failure      400      {object}  map[string]interface{}     "Erreur de validation"
// @Failure      401      {object}  map[string]interface{}     "Identifiants incorrects"
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Rechercher l'utilisateur par email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	// Vérifier si l'utilisateur est actif
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Compte désactivé"})
		return
	}

	// Vérifier le mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	// Générer le token JWT
	token, err := utils.GenerateJWT(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	// Masquer le mot de passe dans la réponse
	user.Password = ""

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// GetProfile récupère le profil de l'utilisateur connecté
// @Summary      Récupérer le profil utilisateur
// @Description  Obtenir les informations du profil de l'utilisateur connecté
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.User                "Profil utilisateur"
// @Failure      401  {object}  map[string]interface{}     "Non authentifié"
// @Failure      404  {object}  map[string]interface{}     "Utilisateur non trouvé"
// @Router       /profile [get]
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Masquer le mot de passe
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// UpdateProfile met à jour le profil de l'utilisateur connecté
// @Summary      Mettre à jour le profil utilisateur
// @Description  Modifier les informations du profil de l'utilisateur connecté
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      map[string]interface{}  true  "Données de mise à jour"
// @Success      200      {object}  models.User             "Profil mis à jour"
// @Failure      400      {object}  map[string]interface{}  "Erreur de validation"
// @Failure      401      {object}  map[string]interface{}  "Non authentifié"
// @Router       /profile [put]
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Structure pour les données de mise à jour
	var updateData struct {
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mettre à jour les champs modifiables
	if updateData.Username != "" {
		// Vérifier si le nom d'utilisateur est déjà pris
		var existingUser models.User
		if err := database.DB.Where("username = ? AND id != ?", updateData.Username, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Ce nom d'utilisateur est déjà pris"})
			return
		}
		user.Username = updateData.Username
	}
	if updateData.PhoneNumber != "" {
		user.PhoneNumber = updateData.PhoneNumber
	}
	if updateData.Email != "" {
		// Vérifier si l'email est déjà pris
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", updateData.Email, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Cet email est déjà pris"})
			return
		}
		user.Email = updateData.Email
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du profil"})
		return
	}

	// Masquer le mot de passe
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// ChangePassword permet de changer le mot de passe
func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Vérifier le mot de passe actuel
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe actuel incorrect"})
		return
	}

	// Hacher le nouveau mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du hachage du mot de passe"})
		return
	}

	// Mettre à jour le mot de passe
	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du mot de passe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mot de passe mis à jour avec succès"})
}

// RefreshToken godoc
// @Summary      Rafraîchir le token JWT
// @Description  Générer un nouveau token JWT à partir d'un token existant
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Nouveau token généré"
// @Failure      401  {object}  map[string]interface{}  "Token invalide"
// @Router       /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token d'autorisation manquant"})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format de token invalide"})
		return
	}

	oldToken := tokenParts[1]

	// Générer un nouveau token
	newToken, err := utils.RefreshToken(oldToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Impossible de rafraîchir le token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

// Logout invalide le token (côté client)
// @Summary      Déconnexion
// @Description  Déconnecter l'utilisateur (invalide le token côté client)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Déconnexion réussie"
// @Router       /auth/logout [post]
func Logout(c *gin.Context) {
	// Dans une implémentation complète, on pourrait ajouter le token à une blacklist
	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie"})
}
