package middleware

import (
	"api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware vérifie l'authentification JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token d'autorisation manquant"})
			c.Abort()
			return
		}

		// Vérifier le format "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format de token invalide"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Valider le token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		// Ajouter les informations utilisateur au contexte
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RequireRole vérifie que l'utilisateur a le rôle requis
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Rôle utilisateur non trouvé"})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de type de rôle"})
			c.Abort()
			return
		}

		// Vérifier si le rôle est autorisé
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé - rôle insuffisant"})
		c.Abort()
	}
}

// RequireAdmin vérifie que l'utilisateur est un administrateur
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireTeacher vérifie que l'utilisateur est un enseignant
func RequireTeacher() gin.HandlerFunc {
	return RequireRole("teacher")
}

// RequireParent vérifie que l'utilisateur est un parent
func RequireParent() gin.HandlerFunc {
	return RequireRole("parent")
}

// RequireChild vérifie que l'utilisateur est un enfant
func RequireChild() gin.HandlerFunc {
	return RequireRole("child")
}

// RequireTeacherOrAdmin vérifie que l'utilisateur est enseignant ou admin
func RequireTeacherOrAdmin() gin.HandlerFunc {
	return RequireRole("teacher", "admin")
}

// RequireParentOrAdmin vérifie que l'utilisateur est parent ou admin
func RequireParentOrAdmin() gin.HandlerFunc {
	return RequireRole("parent", "admin")
}

// RequireAnyRole vérifie que l'utilisateur a au moins un des rôles spécifiés
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return RequireRole(roles...)
}

// GetUserID récupère l'ID utilisateur depuis le contexte
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

// GetUserRole récupère le rôle utilisateur depuis le contexte
func GetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("userRole")
	if !exists {
		return "", false
	}

	role, ok := userRole.(string)
	return role, ok
}

// IsAdmin vérifie si l'utilisateur actuel est un administrateur
func IsAdmin(c *gin.Context) bool {
	role, exists := GetUserRole(c)
	return exists && role == "admin"
}

// IsTeacher vérifie si l'utilisateur actuel est un enseignant
func IsTeacher(c *gin.Context) bool {
	role, exists := GetUserRole(c)
	return exists && role == "teacher"
}

// IsParent vérifie si l'utilisateur actuel est un parent
func IsParent(c *gin.Context) bool {
	role, exists := GetUserRole(c)
	return exists && role == "parent"
}

// IsChild vérifie si l'utilisateur actuel est un enfant
func IsChild(c *gin.Context) bool {
	role, exists := GetUserRole(c)
	return exists && role == "child"
}

// CanAccessUser vérifie si l'utilisateur peut accéder aux données d'un autre utilisateur
func CanAccessUser(c *gin.Context, targetUserID uint) bool {
	currentUserID, exists := GetUserID(c)
	if !exists {
		return false
	}

	// L'utilisateur peut toujours accéder à ses propres données
	if currentUserID == targetUserID {
		return true
	}

	// Les admins peuvent accéder aux données de tous les utilisateurs
	if IsAdmin(c) {
		return true
	}

	// TODO: Ajouter la logique pour les parents qui peuvent accéder aux données de leurs enfants
	// et les enseignants qui peuvent accéder aux données de leurs élèves

	return false
}
