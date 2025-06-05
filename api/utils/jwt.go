package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims représente les claims du JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT génère un token JWT pour un utilisateur
func GenerateJWT(userID uint, role string) (string, error) {
	// Récupérer la clé secrète depuis les variables d'environnement
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default-secret-key-change-in-production"
	}

	// Récupérer la durée d'expiration depuis les variables d'environnement
	expirationHours := 24 // Par défaut 24 heures
	if envExpiration := os.Getenv("JWT_EXPIRATION_HOURS"); envExpiration != "" {
		if hours, err := strconv.Atoi(envExpiration); err == nil {
			expirationHours = hours
		}
	}

	// Créer les claims
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "educational-platform-api",
			Subject:   strconv.Itoa(int(userID)),
		},
	}

	// Créer le token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signer le token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT valide un token JWT et retourne les claims
func ValidateJWT(tokenString string) (*Claims, error) {
	// Récupérer la clé secrète
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default-secret-key-change-in-production"
	}

	// Parser le token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Vérifier la méthode de signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("méthode de signature inattendue")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Vérifier si le token est valide
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token invalide")
}

// ExtractUserIDFromToken extrait l'ID utilisateur d'un token JWT
func ExtractUserIDFromToken(tokenString string) (uint, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ExtractRoleFromToken extrait le rôle d'un token JWT
func ExtractRoleFromToken(tokenString string) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}

// IsTokenExpired vérifie si un token est expiré
func IsTokenExpired(tokenString string) bool {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return true
	}

	return time.Now().After(claims.ExpiresAt.Time)
}

// RefreshToken génère un nouveau token avec une nouvelle expiration
func RefreshToken(tokenString string) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}

	// Générer un nouveau token avec les mêmes informations utilisateur
	return GenerateJWT(claims.UserID, claims.Role)
}
