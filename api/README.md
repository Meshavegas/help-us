# API REST Go - Gestion d'Utilisateurs et de Tâches

Une API REST complète développée en Go avec Gin, GORM et JWT pour l'authentification.

## 🚀 Fonctionnalités

- **Authentification JWT** : Inscription, connexion et protection des routes
- **Gestion des utilisateurs** : CRUD complet avec validation
- **Gestion des tâches** : Création, modification, suppression et filtrage
- **Base de données SQLite** : Stockage local avec migrations automatiques
- **Validation des données** : Validation côté serveur avec Gin
- **CORS** : Support pour les applications frontend
- **Architecture MVC** : Code organisé et maintenable

## 📋 Prérequis

- Go 1.21 ou supérieur
- Git

## 🛠️ Installation

1. Cloner le projet :
```bash
git clone <votre-repo>
cd api
```

2. Installer les dépendances :
```bash
go mod tidy
```

3. Configurer les variables d'environnement (optionnel) :
```bash
cp .env.example .env
# Modifier les valeurs dans .env si nécessaire
```

4. Lancer l'application :
```bash
go run main.go
```

L'API sera disponible sur `http://localhost:8080`

## 📚 Documentation API

### Authentification

#### Inscription
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Connexion
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Utilisateurs (Authentification requise)

#### Obtenir le profil
```http
GET /api/v1/profile
Authorization: Bearer <token>
```

#### Lister les utilisateurs
```http
GET /api/v1/users
Authorization: Bearer <token>
```

#### Obtenir un utilisateur
```http
GET /api/v1/users/{id}
Authorization: Bearer <token>
```

#### Mettre à jour un utilisateur
```http
PUT /api/v1/users/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith"
}
```

#### Supprimer un utilisateur
```http
DELETE /api/v1/users/{id}
Authorization: Bearer <token>
```

### Tâches (Authentification requise)

#### Lister les tâches
```http
GET /api/v1/tasks?status=pending&priority=high
Authorization: Bearer <token>
```

#### Créer une tâche
```http
POST /api/v1/tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Ma nouvelle tâche",
  "description": "Description de la tâche",
  "priority": "high",
  "due_date": "2024-12-31T23:59:59Z"
}
```

#### Obtenir une tâche
```http
GET /api/v1/tasks/{id}
Authorization: Bearer <token>
```

#### Mettre à jour une tâche
```http
PUT /api/v1/tasks/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "completed",
  "priority": "medium"
}
```

#### Supprimer une tâche
```http
DELETE /api/v1/tasks/{id}
Authorization: Bearer <token>
```

#### Statistiques des tâches
```http
GET /api/v1/tasks/stats
Authorization: Bearer <token>
```

## 🏗️ Structure du projet

```
api/
├── config/          # Configuration de la base de données
├── controllers/     # Contrôleurs (logique métier)
├── middleware/      # Middlewares (authentification, CORS)
├── models/          # Modèles de données
├── routes/          # Configuration des routes
├── utils/           # Utilitaires (JWT, etc.)
├── main.go          # Point d'entrée
├── go.mod           # Dépendances Go
└── .env             # Variables d'environnement
```

## 🔧 Configuration

Variables d'environnement disponibles dans `.env` :

- `PORT` : Port du serveur (défaut: 8080)
- `DB_PATH` : Chemin de la base de données SQLite
- `JWT_SECRET` : Clé secrète pour JWT
- `GIN_MODE` : Mode Gin (debug/release)

## 🧪 Test de l'API

### Avec curl

1. Inscription :
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"testuser","password":"password123","first_name":"Test","last_name":"User"}'
```

2. Connexion :
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

3. Créer une tâche (remplacer TOKEN) :
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"title":"Ma première tâche","description":"Description de test","priority":"high"}'
```

## 📝 Statuts et Priorités

### Statuts des tâches :
- `pending` : En attente
- `in_progress` : En cours
- `completed` : Terminée
- `cancelled` : Annulée

### Priorités des tâches :
- `low` : Basse
- `medium` : Moyenne
- `high` : Haute
- `urgent` : Urgente

## 🚀 Déploiement

Pour déployer en production :

1. Compiler l'application :
```bash
go build -o api main.go
```

2. Configurer les variables d'environnement :
```bash
export GIN_MODE=release
export JWT_SECRET=votre_secret_super_securise
```

3. Lancer l'application :
```bash
./api
```

## 🤝 Contribution

1. Fork le projet
2. Créer une branche feature (`git checkout -b feature/AmazingFeature`)
3. Commit les changements (`git commit -m 'Add some AmazingFeature'`)
4. Push vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails. 