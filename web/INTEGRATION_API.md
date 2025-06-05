# Intégration Frontend Next.js avec API Go Backend

Ce document explique comment le frontend Next.js a été configuré pour fonctionner avec l'API Go backend en gérant les tokens JWT de manière sécurisée côté serveur.

## Architecture

### Gestion des Tokens
- **Stockage sécurisé** : Les tokens JWT sont stockés dans des cookies `httpOnly` côté serveur
- **Pas d'exposition côté client** : Les tokens ne sont jamais exposés au JavaScript côté client
- **Middleware d'authentification** : Vérification automatique des tokens sur toutes les routes protégées

### Structure des Fichiers

```
web/
├── lib/
│   ├── auth.ts          # Gestion de l'authentification côté serveur
│   └── api.ts           # Client API pour les appels côté client
├── app/
│   ├── api/             # Routes API proxy Next.js
│   │   ├── auth/        # Authentification (login, register, logout)
│   │   ├── profile/     # Profil utilisateur
│   │   └── users/       # Gestion des utilisateurs
│   ├── login/           # Page de connexion
│   ├── register/        # Page d'inscription
│   └── (dashboard)/     # Pages protégées du dashboard
└── middleware.ts        # Middleware d'authentification global
```

## Configuration

### Variables d'environnement

Créez un fichier `.env.local` basé sur `.env.example` :

```bash
# API Backend Configuration
API_BASE_URL=http://localhost:8080
API_VERSION=v1

# NextAuth Configuration
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key-here

# JWT Configuration for API Integration
JWT_SECRET=your-jwt-secret-key-here
```

## Fonctionnalités

### 1. Authentification Sécurisée

#### Connexion
```typescript
import { apiClient } from '@/lib/api';

const response = await apiClient.login(email, password);
if (!response.error) {
  // Redirection automatique vers le dashboard
  router.push('/dashboard');
}
```

#### Inscription
```typescript
const response = await apiClient.register({
  email,
  username,
  password,
  first_name,
  last_name
});
```

#### Déconnexion
```typescript
await apiClient.logout();
router.push('/login');
```

### 2. Appels API Authentifiés

#### Côté Serveur (Server Components)
```typescript
import { authenticatedApiCall } from '@/lib/auth';

// Dans un Server Component
const users = await authenticatedApiCall('/users');
```

#### Côté Client (Client Components)
```typescript
import { apiClient, useApi } from '@/lib/api';

// Hook pour la gestion d'état
const { loading, error, data, execute } = useApi();

// Exécuter un appel API
useEffect(() => {
  execute(() => apiClient.getUsers());
}, []);
```

### 3. Protection des Routes

Le middleware automatique protège toutes les routes sauf :
- `/login`
- `/register`
- `/` (page d'accueil)

```typescript
// middleware.ts
export { auth as middleware } from '@/lib/auth';
```

### 4. Gestion des Erreurs

Toutes les erreurs API sont gérées de manière cohérente :

```typescript
const response = await apiClient.getUsers();
if (response.error) {
  console.error('Erreur:', response.error);
} else {
  console.log('Données:', response.data);
}
```

## Avantages de cette Architecture

### Sécurité
- **Tokens httpOnly** : Impossible d'accéder aux tokens via JavaScript côté client
- **Protection CSRF** : Les cookies sont configurés avec `sameSite: 'lax'`
- **Expiration automatique** : Les tokens expirent après 7 jours

### Performance
- **Server-Side Rendering** : Les données utilisateur sont récupérées côté serveur
- **Pas de rendu client inutile** : L'authentification est vérifiée avant le rendu
- **Cache automatique** : Next.js cache les réponses des Server Components

### Développement
- **API unifiée** : Un seul point d'entrée pour tous les appels API
- **Types TypeScript** : Typage complet pour toutes les réponses API
- **Gestion d'erreurs centralisée** : Toutes les erreurs sont gérées de manière cohérente

## Utilisation

### Démarrage

1. **Démarrer l'API Go** :
   ```bash
   cd api
   make run
   ```

2. **Démarrer le frontend Next.js** :
   ```bash
   cd web
   npm install
   npm run dev
   ```

3. **Accéder à l'application** :
   - Frontend : http://localhost:3000
   - API : http://localhost:8080
   - Documentation API : http://localhost:8080/swagger/index.html

### Flux d'authentification

1. L'utilisateur visite `/login`
2. Saisie des identifiants et soumission du formulaire
3. Appel à `/api/auth/login` (route Next.js)
4. La route Next.js appelle l'API Go `/api/v1/auth/login`
5. Si succès, le token JWT est stocké dans un cookie httpOnly
6. Redirection vers `/dashboard`
7. Le middleware vérifie le token sur chaque requête

### Exemple d'utilisation complète

```typescript
'use client';

import { useState, useEffect } from 'react';
import { apiClient, useApi } from '@/lib/api';

export default function UsersPage() {
  const { loading, error, data, execute } = useApi();

  useEffect(() => {
    execute(() => apiClient.getUsers());
  }, []);

  if (loading) return <div>Chargement...</div>;
  if (error) return <div>Erreur: {error}</div>;

  return (
    <div>
      <h1>Utilisateurs</h1>
      {data?.users?.map(user => (
        <div key={user.id}>
          {user.first_name} {user.last_name} - {user.email}
        </div>
      ))}
    </div>
  );
}
```

## Endpoints Disponibles

Tous les endpoints de l'API Go sont disponibles via les routes proxy Next.js :

- `POST /api/auth/login` - Connexion
- `POST /api/auth/register` - Inscription  
- `POST /api/auth/logout` - Déconnexion
- `GET /api/profile` - Profil utilisateur
- `GET /api/users` - Liste des utilisateurs
- `GET /api/users/[id]` - Utilisateur spécifique
- `PUT /api/users/[id]` - Modifier un utilisateur
- `DELETE /api/users/[id]` - Supprimer un utilisateur

## Sécurité et Bonnes Pratiques

1. **Jamais de tokens côté client** : Les tokens JWT ne sont jamais exposés au JavaScript côté client
2. **Validation côté serveur** : Toutes les requêtes sont validées côté serveur avant d'être transmises à l'API Go
3. **Gestion des erreurs** : Les erreurs sont gérées de manière sécurisée sans exposer d'informations sensibles
4. **Expiration des sessions** : Les tokens expirent automatiquement et les utilisateurs sont redirigés vers la page de connexion
5. **HTTPS en production** : Les cookies sont configurés pour être sécurisés en production 