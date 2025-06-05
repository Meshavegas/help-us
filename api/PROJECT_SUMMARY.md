# 🚀 API REST Go - Gestion d'Utilisateurs et de Tâches

## 📋 Résumé du Projet

Cette API REST complète a été développée en Go avec le framework Gin pour la gestion d'utilisateurs et de tâches. Elle inclut un système d'authentification JWT, une documentation Swagger interactive, et une architecture modulaire bien structurée.

## ✨ Fonctionnalités Principales

### 🔐 Authentification & Autorisation
- **Inscription d'utilisateurs** avec validation des données
- **Connexion sécurisée** avec génération de tokens JWT
- **Middleware d'authentification** pour protéger les routes
- **Gestion des profils utilisateurs**

### 👥 Gestion des Utilisateurs
- **CRUD complet** pour les utilisateurs
- **Validation des données** (email unique, username unique)
- **Gestion du statut actif/inactif**
- **Suppression en cascade** des tâches associées

### 📝 Gestion des Tâches
- **CRUD complet** pour les tâches
- **Système de statuts** : pending, in_progress, completed, cancelled
- **Système de priorités** : low, medium, high, urgent
- **Filtrage avancé** par statut et priorité
- **Statistiques des tâches** par utilisateur
- **Dates d'échéance** optionnelles

## 🏗️ Architecture

```
api/
├── config/          # Configuration de la base de données
├── controllers/     # Logique métier et handlers HTTP
├── middleware/      # Middleware d'authentification
├── models/          # Modèles de données et structures
├── routes/          # Configuration des routes
├── utils/           # Utilitaires (JWT, etc.)
├── docs/            # Documentation Swagger générée
├── main.go          # Point d'entrée de l'application
├── Makefile         # Commandes de développement
├── examples.http    # Exemples de requêtes HTTP
└── go.mod           # Dépendances Go
```

## 🛠️ Technologies Utilisées

- **Go 1.21+** - Langage de programmation
- **Gin** - Framework web HTTP
- **GORM** - ORM pour Go
- **SQLite** - Base de données (configurable)
- **JWT** - Authentification par tokens
- **Swagger** - Documentation API interactive
- **bcrypt** - Hachage des mots de passe

## 🚀 Démarrage Rapide

### Prérequis
- Go 1.21 ou supérieur
- Git

### Installation et Lancement

```bash
# Cloner le projet
git clone <repository-url>
cd api

# Installer les dépendances
make install

# Lancer en mode développement
make dev

# Ou directement avec Go
go run main.go
```

L'API sera accessible sur `http://localhost:8080`

## 📚 Documentation

### Swagger UI
- **URL** : http://localhost:8080/swagger/index.html
- **Documentation interactive** avec possibilité de tester les endpoints
- **Schémas de données** détaillés
- **Exemples de requêtes/réponses**

### Endpoints Principaux

#### Authentification
- `POST /api/v1/auth/register` - Inscription
- `POST /api/v1/auth/login` - Connexion
- `GET /api/v1/profile` - Profil utilisateur

#### Utilisateurs
- `GET /api/v1/users` - Liste des utilisateurs
- `GET /api/v1/users/{id}` - Détails d'un utilisateur
- `PUT /api/v1/users/{id}` - Mise à jour
- `DELETE /api/v1/users/{id}` - Suppression

#### Tâches
- `GET /api/v1/tasks` - Liste des tâches (avec filtres)
- `POST /api/v1/tasks` - Création d'une tâche
- `GET /api/v1/tasks/{id}` - Détails d'une tâche
- `PUT /api/v1/tasks/{id}` - Mise à jour
- `DELETE /api/v1/tasks/{id}` - Suppression
- `GET /api/v1/tasks/stats` - Statistiques

## 🧪 Tests et Exemples

### Fichier examples.http
Le fichier `examples.http` contient des exemples prêts à utiliser pour tous les endpoints de l'API.

### Tests avec curl
```bash
# Vérification de santé
curl http://localhost:8080/health

# Inscription
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"testuser","password":"password123","first_name":"Test","last_name":"User"}'

# Connexion
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## 🔧 Commandes Make Disponibles

```bash
make install      # Installer les dépendances
make dev          # Lancer en mode développement
make build        # Compiler l'application
make test         # Lancer les tests
make swagger      # Générer la documentation Swagger
make clean        # Nettoyer les fichiers générés
make help         # Afficher l'aide
```

## 🔒 Sécurité

- **Hachage bcrypt** pour les mots de passe
- **Tokens JWT** pour l'authentification
- **Validation des données** d'entrée
- **Middleware CORS** configuré
- **Protection des routes** sensibles

## 📊 Fonctionnalités Avancées

### Filtrage des Tâches
```bash
# Filtrer par statut
GET /api/v1/tasks?status=pending

# Filtrer par priorité
GET /api/v1/tasks?priority=high

# Combinaison de filtres
GET /api/v1/tasks?status=pending&priority=high
```

### Statistiques
L'endpoint `/api/v1/tasks/stats` fournit :
- Nombre total de tâches
- Répartition par statut
- Nombre de tâches haute priorité

## 🌟 Points Forts

1. **Architecture modulaire** et maintenable
2. **Documentation complète** avec Swagger
3. **Sécurité robuste** avec JWT et bcrypt
4. **Validation des données** complète
5. **Gestion d'erreurs** appropriée
6. **Code bien structuré** et commenté
7. **Exemples pratiques** fournis
8. **Makefile** pour faciliter le développement

## 🔄 Prochaines Étapes Possibles

- [ ] Tests unitaires et d'intégration
- [ ] Pagination pour les listes
- [ ] Système de notifications
- [ ] Upload de fichiers
- [ ] Cache Redis
- [ ] Déploiement Docker
- [ ] CI/CD Pipeline
- [ ] Monitoring et logs

## 📞 Support

Pour toute question ou problème :
1. Consultez la documentation Swagger
2. Vérifiez les exemples dans `examples.http`
3. Consultez les logs de l'application
4. Utilisez `make help` pour les commandes disponibles

---

**Statut** : ✅ Projet complet et fonctionnel  
**Version** : 1.0  
**Dernière mise à jour** : Juin 2025 