# Documentation des Endpoints API

## Base URL
```
http://localhost:8080
```

## Authentification
L'API utilise JWT (JSON Web Tokens) pour l'authentification. Après connexion, incluez le token dans l'en-tête Authorization :
```
Authorization: Bearer <votre_token_jwt>
```

---

## 🏥 Santé de l'API

### GET /health
Vérifie l'état de l'API.

**Réponse :**
```json
{
  "status": "OK",
  "message": "API REST Go fonctionne correctement"
}
```

---

## 🔐 Authentification

### POST /api/v1/auth/register
Inscription d'un nouvel utilisateur.

**Corps de la requête :**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Réponse (201) :**
```json
{
  "message": "Utilisateur créé avec succès",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### POST /api/v1/auth/login
Connexion d'un utilisateur existant.

**Corps de la requête :**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Réponse (200) :**
```json
{
  "message": "Connexion réussie",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## 👤 Profil Utilisateur (Authentification requise)

### GET /api/v1/profile
Récupère le profil de l'utilisateur connecté.

**En-têtes :**
```
Authorization: Bearer <token>
```

**Réponse (200) :**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 👥 Gestion des Utilisateurs (Authentification requise)

### GET /api/v1/users
Liste tous les utilisateurs.

**En-têtes :**
```
Authorization: Bearer <token>
```

**Réponse (200) :**
```json
{
  "users": [
    {
      "id": 1,
      "email": "user@example.com",
      "username": "username",
      "first_name": "John",
      "last_name": "Doe",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "count": 1
}
```

### GET /api/v1/users/{id}
Récupère un utilisateur spécifique.

**En-têtes :**
```
Authorization: Bearer <token>
```

**Réponse (200) :**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### PUT /api/v1/users/{id}
Met à jour un utilisateur.

**En-têtes :**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Corps de la requête :**
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane@example.com",
  "username": "janesmith",
  "is_active": false
}
```

**Réponse (200) :**
```json
{
  "message": "Utilisateur mis à jour avec succès",
  "user": {
    "id": 1,
    "email": "jane@example.com",
    "username": "janesmith",
    "first_name": "Jane",
    "last_name": "Smith",
    "is_active": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### DELETE /api/v1/users/{id}
Supprime un utilisateur.

**En-têtes :**
```
Authorization: Bearer <token>
```

**Réponse (200) :**
```json
{
  "message": "Utilisateur supprimé avec succès"
}
```

---

## ✅ Gestion des Tâches (Authentification requise)

### GET /api/v1/tasks
Liste les tâches de l'utilisateur connecté.

**En-têtes :**
```
Authorization: Bearer <token>
```

**Paramètres de requête optionnels :**
- `status` : Filtrer par statut (pending, in_progress, completed, cancelled)
- `priority` : Filtrer par priorité (low, medium, high, urgent)

**Exemple :**
```
GET /api/v1/tasks?status=pending&priority=high
```

**Réponse (200) :**
```json
{
  "tasks": [
    {
      "id": 1,
      "title": "Ma première tâche",
      "description": "Description de la tâche",
      "status": "pending",
      "priority": "high",
      "due_date": "2024-12-31T23:59:59Z",
      "user_id": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "count": 1
}
```

### POST /api/v1/tasks
Crée une nouvelle tâche.

**En-têtes :**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Corps de la requête :**
```json
{
  "title": "Nouvelle tâche",
  "description": "Description de la tâche",
  "priority": "high",
  "due_date": "2024-12-31T23:59:59Z"
}
```

**Réponse (201) :**
```json
{
  "message": "Tâche créée avec succès",
  "task": {
    "id": 1,
    "title": "Nouvelle tâche",
    "description": "Description de la tâche",
    "status": "pending",
    "priority": "high",
    "due_date": "2024-12-31T23:59:59Z",
    "user_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### GET /api/v1/tasks/{id}
Récupère une tâche spécifique.

**En-têtes :**
```
Authorization: Bearer <token>
```

**Réponse (200) :**
```json
{
  "task": {
    "id": 1,
    "title": "Ma tâche",
    "description": "Description",
    "status": "pending",
    "priority": "medium",
    "due_date": null,
    "user_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### PUT /api/v1/tasks/{id}
Met à jour une tâche.

**En-têtes :**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Corps de la requête :**
```json
{
  "title": "Titre modifié",
  "description": "Nouvelle description",
  "status": "in_progress",
  "priority": "urgent",
  "due_date": "2024-12-31T23:59:59Z"
}
```

**Réponse (200) :**
```json
{
  "message": "Tâche mise à jour avec succès",
  "task": {
    "id": 1,
    "title": "Titre modifié",
    "description": "Nouvelle description",
    "status": "in_progress",
    "priority": "urgent",
    "due_date": "2024-12-31T23:59:59Z",
    "user_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### DELETE /api/v1/tasks/{id}
Supprime une tâche.

**En-têtes :**
```
Authorization: Bearer <token>
```

**Réponse (200) :**
```json
{
  "message": "Tâche supprimée avec succès"
}
```

### GET /api/v1/tasks/stats
Récupère les statistiques des tâches de l'utilisateur.

**En-têtes :**
```
Authorization: Bearer <token>
```

**Réponse (200) :**
```json
{
  "stats": {
    "total": 10,
    "pending": 3,
    "in_progress": 2,
    "completed": 4,
    "cancelled": 1,
    "high_priority": 5
  }
}
```

---

## 📊 Codes de Statut HTTP

- **200 OK** : Requête réussie
- **201 Created** : Ressource créée avec succès
- **400 Bad Request** : Données de requête invalides
- **401 Unauthorized** : Authentification requise ou token invalide
- **403 Forbidden** : Accès interdit
- **404 Not Found** : Ressource non trouvée
- **409 Conflict** : Conflit (ex: email déjà utilisé)
- **500 Internal Server Error** : Erreur serveur

---

## 🔧 Valeurs Autorisées

### Statuts des Tâches
- `pending` : En attente
- `in_progress` : En cours
- `completed` : Terminée
- `cancelled` : Annulée

### Priorités des Tâches
- `low` : Basse
- `medium` : Moyenne (par défaut)
- `high` : Haute
- `urgent` : Urgente

---

## 🚨 Gestion des Erreurs

Toutes les erreurs retournent un objet JSON avec un message d'erreur :

```json
{
  "error": "Message d'erreur descriptif"
}
```

### Erreurs d'Authentification
```json
{
  "error": "Token d'autorisation requis"
}
```

### Erreurs de Validation
```json
{
  "error": "Key: 'UserCreateRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
``` 