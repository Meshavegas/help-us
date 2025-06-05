#!/bin/bash

# Script de test pour l'API Plateforme Éducative
# Ce script démontre l'utilisation des principales fonctionnalités de l'API

echo "=== Test de l'API Plateforme Éducative ==="
echo

# Variables
BASE_URL="http://localhost:8080/api/v1"
HEALTH_URL="http://localhost:8080/health"

# Test 1: Vérification de santé
echo "1. Test de santé de l'API..."
curl -s $HEALTH_URL | jq '.'
echo
echo

# Test 2: Inscription d'un nouvel utilisateur
echo "2. Inscription d'un nouvel utilisateur..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "demo_user",
    "email": "demo@example.com",
    "password": "password123",
    "role": "enseignant",
    "phone_number": "+33123456789",
    "specialization": "Mathématiques"
  }')

echo $REGISTER_RESPONSE | jq '.'
echo

# Extraction du token
TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token')

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo "Token reçu: ${TOKEN:0:50}..."
    echo
    
    # Test 3: Connexion
    echo "3. Test de connexion..."
    LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
      -H "Content-Type: application/json" \
      -d '{
        "email": "demo@example.com",
        "password": "password123"
      }')
    
    echo $LOGIN_RESPONSE | jq '.'
    echo
    
    # Test 4: Récupération du profil
    echo "4. Récupération du profil utilisateur..."
    curl -s -X GET $BASE_URL/profile \
      -H "Authorization: Bearer $TOKEN" | jq '.'
    echo
    echo
    
    # Test 5: Mise à jour du profil
    echo "5. Mise à jour du profil..."
    curl -s -X PUT $BASE_URL/profile \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "phone_number": "+33987654321",
        "specialization": "Physique"
      }' | jq '.'
    echo
    echo
    
    # Test 6: Déconnexion
    echo "6. Déconnexion..."
    curl -s -X POST $BASE_URL/auth/logout \
      -H "Authorization: Bearer $TOKEN" | jq '.'
    echo
    echo
    
else
    echo "Erreur: Impossible d'obtenir le token d'authentification"
    echo "Réponse: $REGISTER_RESPONSE"
fi

echo "=== Tests terminés ==="
echo
echo "Documentation Swagger disponible à: http://localhost:8080/swagger/index.html" 