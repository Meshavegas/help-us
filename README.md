### Description du projet

Ce projet modélise une plateforme de gestion et de mise en relation entre familles, enseignants et administrateurs autour de missions éducatives, de cours particuliers et de ressources pédagogiques.

## Principaux acteurs

- **User** : Classe de base représentant un utilisateur de la plateforme (identifiant, nom d'utilisateur, mot de passe, email, téléphone, date de création, gestion de connexion et de profil).
- **Famille** : Hérite de User. Peut consulter et sélectionner des enseignants, planifier des sessions, gérer les paiements, accéder à des ressources, donner des avis, etc.
- Enseignant : Hérite de User. Peut compléter son profil, consulter et postuler à des offres, planifier et déclarer des cours, soumettre des rapports, accéder à des documents/vidéos, etc.
- Administrateur : Hérite de User. Gère les comptes, les étudiants, les enseignants, les offres, valide les rapports, fournit du support, etc. Domaines fonctionnels
- Course : Représente un cours programmé (horaire, durée, lieu, statut, gestion de la planification et validation).
- Payment : Gère les paiements (montant, date, statut, type, génération de factures, historique).
- Mission : Encadre un ensemble de cours (dates, statut, description, création, arrêt, prolongation).
- Report : Rapports soumis par les enseignants et validés par l’administration (contenu, statut, soumission, validation).
- Offer : Offres de mission ou de cours (titre, description, taux horaire, publication, statut, candidatures).
- Resource : Ressources pédagogiques (titre, type, URL, date de dépôt, gestion de l’accès et du téléchargement).
- Option : Options de réservation ou de candidature (dates, statut, création, annulation, acceptation).
- Address : Adresses physiques (rue, ville, code postal, pays, coordonnées GPS, gestion et calcul d’itinéraires).

### Relations principales

- Les familles embauchent des enseignants et planifient des cours.
- Les enseignants enseignent pour des familles, conduisent des cours, soumettent des rapports, créent des options et postulent à des offres.
- Les administrateurs gèrent les utilisateurs, valident les rapports, créent des offres et gèrent les ressources.
- Les missions regroupent plusieurs cours et rapports.
- Les paiements sont associés aux utilisateurs et aux cours.
- Les ressources sont accessibles par les utilisateurs et gérées par les administrateurs.
- Les adresses sont liées aux utilisateurs et aux lieux de cours.
  ![Class Diagram](https://www.mermaidchart.com/raw/5c8cf39f-38b7-4d81-a99e-7e0e244be859?theme=light&version=v0.1&format=svg)

### Comprehensive API Endpoints for Educational Platform

Based on the class diagram, here's a complete list of API endpoints that would cover all functionalities of the educational platform. The endpoints are organized by entity and include all necessary operations for the system.

## Authentication & User Management

| Method | Endpoint                   | Description                                                 |
| ------ | -------------------------- | ----------------------------------------------------------- |
| POST   | `/api/auth/register`       | Register a new user (famille, enseignant, or administrator) |
| POST   | `/api/auth/login`          | Authenticate a user and receive a token                     |
| POST   | `/api/auth/logout`         | Log out the current user                                    |
| GET    | `/api/auth/me`             | Get the current authenticated user's profile                |
| PUT    | `/api/auth/me`             | Update the current user's profile                           |
| POST   | `/api/auth/password/reset` | Request a password reset                                    |
| PUT    | `/api/auth/password/reset` | Reset password with token                                   |

## User Profiles

| Method | Endpoint                   | Description                            |
| ------ | -------------------------- | -------------------------------------- |
| GET    | `/api/users`               | List all users (admin only)            |
| GET    | `/api/users/:id`           | Get a specific user's details          |
| PUT    | `/api/users/:id`           | Update a user's details (admin only)   |
| DELETE | `/api/users/:id`           | Delete a user (admin only)             |
| GET    | `/api/users/:id/addresses` | Get all addresses for a user           |
| GET    | `/api/users/:id/payments`  | Get all payments for a user            |
| GET    | `/api/users/:id/resources` | Get all resources accessible to a user |

## Famille Endpoints

| Method | Endpoint                     | Description                               |
| ------ | ---------------------------- | ----------------------------------------- |
| GET    | `/api/familles`              | List all families (admin only)            |
| GET    | `/api/familles/:id`          | Get a specific family's details           |
| POST   | `/api/familles`              | Create a new family profile               |
| PUT    | `/api/familles/:id`          | Update a family's details                 |
| DELETE | `/api/familles/:id`          | Delete a family profile                   |
| GET    | `/api/familles/:id/teachers` | Get all teachers working with this family |
| GET    | `/api/familles/:id/missions` | Get all missions for this family          |
| GET    | `/api/familles/:id/courses`  | Get all courses for this family           |
| GET    | `/api/familles/:id/payments` | Get all payments made by this family      |
| POST   | `/api/familles/:id/reviews`  | Submit a review for a teacher             |
| GET    | `/api/familles/:id/options`  | Get all options received by this family   |

## Enseignant Endpoints

| Method | Endpoint                        | Description                                       |
| ------ | ------------------------------- | ------------------------------------------------- |
| GET    | `/api/enseignants`              | List all teachers                                 |
| GET    | `/api/enseignants/:id`          | Get a specific teacher's details                  |
| POST   | `/api/enseignants`              | Create a new teacher profile                      |
| PUT    | `/api/enseignants/:id`          | Update a teacher's details                        |
| DELETE | `/api/enseignants/:id`          | Delete a teacher profile                          |
| GET    | `/api/enseignants/:id/students` | Get all students (families) for this teacher      |
| GET    | `/api/enseignants/:id/missions` | Get all missions for this teacher                 |
| GET    | `/api/enseignants/:id/courses`  | Get all courses for this teacher                  |
| GET    | `/api/enseignants/:id/payments` | Get all payments to this teacher                  |
| GET    | `/api/enseignants/:id/reports`  | Get all reports submitted by this teacher         |
| GET    | `/api/enseignants/:id/options`  | Get all options created by this teacher           |
| GET    | `/api/enseignants/nearby`       | Find teachers near a location (with query params) |

## Administrator Endpoints

| Method | Endpoint                        | Description                                           |
| ------ | ------------------------------- | ----------------------------------------------------- |
| GET    | `/api/administrators`           | List all administrators (super admin only)            |
| GET    | `/api/administrators/:id`       | Get a specific administrator's details                |
| POST   | `/api/administrators`           | Create a new administrator profile (super admin only) |
| PUT    | `/api/administrators/:id`       | Update an administrator's details                     |
| DELETE | `/api/administrators/:id`       | Delete an administrator profile (super admin only)    |
| GET    | `/api/administrators/dashboard` | Get dashboard statistics and metrics                  |

## Mission Endpoints

| Method | Endpoint                     | Description                                |
| ------ | ---------------------------- | ------------------------------------------ |
| GET    | `/api/missions`              | List all missions (with filtering options) |
| GET    | `/api/missions/:id`          | Get a specific mission's details           |
| POST   | `/api/missions`              | Create a new mission                       |
| PUT    | `/api/missions/:id`          | Update a mission's details                 |
| DELETE | `/api/missions/:id`          | Delete a mission                           |
| GET    | `/api/missions/:id/courses`  | Get all courses for this mission           |
| GET    | `/api/missions/:id/reports`  | Get all reports for this mission           |
| GET    | `/api/missions/:id/payments` | Get all payments for this mission          |
| PUT    | `/api/missions/:id/stop`     | Stop/end a mission                         |
| PUT    | `/api/missions/:id/extend`   | Extend a mission's end date                |

## Course Endpoints

| Method | Endpoint                    | Description                               |
| ------ | --------------------------- | ----------------------------------------- |
| GET    | `/api/courses`              | List all courses (with filtering options) |
| GET    | `/api/courses/:id`          | Get a specific course's details           |
| POST   | `/api/courses`              | Create a new course                       |
| PUT    | `/api/courses/:id`          | Update a course's details                 |
| DELETE | `/api/courses/:id`          | Delete a course                           |
| PUT    | `/api/courses/:id/schedule` | Schedule or reschedule a course           |
| PUT    | `/api/courses/:id/cancel`   | Cancel a course                           |
| PUT    | `/api/courses/:id/complete` | Mark a course as completed                |
| POST   | `/api/courses/:id/declare`  | Declare hours for a completed course      |
| GET    | `/api/courses/:id/payments` | Get all payments for this course          |

## Payment Endpoints

| Method | Endpoint                    | Description                                |
| ------ | --------------------------- | ------------------------------------------ |
| GET    | `/api/payments`             | List all payments (with filtering options) |
| GET    | `/api/payments/:id`         | Get a specific payment's details           |
| POST   | `/api/payments`             | Create a new payment                       |
| PUT    | `/api/payments/:id`         | Update a payment's details                 |
| DELETE | `/api/payments/:id`         | Delete a payment (admin only)              |
| GET    | `/api/payments/statistics`  | Get payment statistics (admin only)        |
| POST   | `/api/payments/process`     | Process a pending payment                  |
| GET    | `/api/payments/advances`    | Get all advance payments                   |
| POST   | `/api/payments/invoice/:id` | Generate an invoice for a payment          |

## Report Endpoints

| Method | Endpoint                    | Description                               |
| ------ | --------------------------- | ----------------------------------------- |
| GET    | `/api/reports`              | List all reports (with filtering options) |
| GET    | `/api/reports/:id`          | Get a specific report's details           |
| POST   | `/api/reports`              | Create a new report                       |
| PUT    | `/api/reports/:id`          | Update a report's details                 |
| DELETE | `/api/reports/:id`          | Delete a report                           |
| PUT    | `/api/reports/:id/validate` | Validate a submitted report               |
| PUT    | `/api/reports/:id/reject`   | Reject a submitted report                 |
| GET    | `/api/reports/pending`      | Get all pending reports (admin only)      |

## Offer Endpoints

| Method | Endpoint                  | Description                              |
| ------ | ------------------------- | ---------------------------------------- |
| GET    | `/api/offers`             | List all offers (with filtering options) |
| GET    | `/api/offers/:id`         | Get a specific offer's details           |
| POST   | `/api/offers`             | Create a new offer                       |
| PUT    | `/api/offers/:id`         | Update an offer's details                |
| DELETE | `/api/offers/:id`         | Delete an offer                          |
| GET    | `/api/offers/:id/options` | Get all options for this offer           |
| PUT    | `/api/offers/:id/close`   | Close an offer                           |
| GET    | `/api/offers/active`      | Get all active offers                    |
| GET    | `/api/offers/search`      | Search offers with specific criteria     |

## Resource Endpoints

| Method | Endpoint                            | Description                                 |
| ------ | ----------------------------------- | ------------------------------------------- |
| GET    | `/api/resources`                    | List all resources (with filtering options) |
| GET    | `/api/resources/:id`                | Get a specific resource's details           |
| POST   | `/api/resources`                    | Create a new resource                       |
| PUT    | `/api/resources/:id`                | Update a resource's details                 |
| DELETE | `/api/resources/:id`                | Delete a resource                           |
| POST   | `/api/resources/:id/access`         | Grant a user access to a resource           |
| DELETE | `/api/resources/:id/access/:userId` | Remove a user's access to a resource        |
| GET    | `/api/resources/documents`          | Get all document resources                  |
| GET    | `/api/resources/videos`             | Get all video resources                     |
| GET    | `/api/resources/search`             | Search resources with specific criteria     |

## Option Endpoints

| Method | Endpoint                   | Description                               |
| ------ | -------------------------- | ----------------------------------------- |
| GET    | `/api/options`             | List all options (with filtering options) |
| GET    | `/api/options/:id`         | Get a specific option's details           |
| POST   | `/api/options`             | Create a new option                       |
| PUT    | `/api/options/:id`         | Update an option's details                |
| DELETE | `/api/options/:id`         | Delete an option                          |
| PUT    | `/api/options/:id/accept`  | Accept an option                          |
| PUT    | `/api/options/:id/decline` | Decline an option                         |
| PUT    | `/api/options/:id/cancel`  | Cancel an option                          |
| GET    | `/api/options/pending`     | Get all pending options                   |
| GET    | `/api/options/expiring`    | Get all options expiring soon             |

## Address Endpoints

| Method | Endpoint                 | Description                                 |
| ------ | ------------------------ | ------------------------------------------- |
| GET    | `/api/addresses`         | List all addresses (admin only)             |
| GET    | `/api/addresses/:id`     | Get a specific address's details            |
| POST   | `/api/addresses`         | Create a new address                        |
| PUT    | `/api/addresses/:id`     | Update an address's details                 |
| DELETE | `/api/addresses/:id`     | Delete an address                           |
| GET    | `/api/addresses/geocode` | Geocode an address (convert to coordinates) |
| GET    | `/api/addresses/route`   | Calculate route between two addresses       |

## Search and Filtering Endpoints

| Method | Endpoint               | Description                                             |
| ------ | ---------------------- | ------------------------------------------------------- |
| GET    | `/api/search/teachers` | Search for teachers with specific criteria              |
| GET    | `/api/search/families` | Search for families with specific criteria (admin only) |
| GET    | `/api/search/courses`  | Search for courses with specific criteria               |
| GET    | `/api/search/missions` | Search for missions with specific criteria              |
| GET    | `/api/search/global`   | Global search across multiple entities                  |

## Statistics and Reporting Endpoints

| Method | Endpoint                   | Description                               |
| ------ | -------------------------- | ----------------------------------------- |
| GET    | `/api/statistics/users`    | Get user statistics                       |
| GET    | `/api/statistics/courses`  | Get course statistics                     |
| GET    | `/api/statistics/payments` | Get payment statistics                    |
| GET    | `/api/statistics/missions` | Get mission statistics                    |
| GET    | `/api/reports/financial`   | Generate financial reports                |
| GET    | `/api/reports/activity`    | Generate activity reports                 |
| GET    | `/api/reports/performance` | Generate performance reports for teachers |

## Notification Endpoints

| Method | Endpoint                      | Description                                |
| ------ | ----------------------------- | ------------------------------------------ |
| GET    | `/api/notifications`          | Get all notifications for the current user |
| GET    | `/api/notifications/:id`      | Get a specific notification's details      |
| PUT    | `/api/notifications/:id/read` | Mark a notification as read                |
| PUT    | `/api/notifications/read-all` | Mark all notifications as read             |
| DELETE | `/api/notifications/:id`      | Delete a notification                      |
| POST   | `/api/notifications/settings` | Update notification preferences            |

## System Configuration Endpoints (Admin Only)

| Method | Endpoint              | Description                       |
| ------ | --------------------- | --------------------------------- |
| GET    | `/api/config/system`  | Get system configuration          |
| PUT    | `/api/config/system`  | Update system configuration       |
| GET    | `/api/config/payment` | Get payment configuration         |
| PUT    | `/api/config/payment` | Update payment configuration      |
| GET    | `/api/config/email`   | Get email templates configuration |
| PUT    | `/api/config/email`   | Update email templates            |
| GET    | `/api/system/logs`    | Get system logs                   |
| POST   | `/api/system/backup`  | Create a system backup            |

## Implementation Considerations

When implementing these endpoints, consider the following:

1. **Authentication and Authorization**: Implement proper authentication for all endpoints and ensure authorization checks (e.g., a famille should only access their own data).
2. **Pagination**: For list endpoints, implement pagination to handle large datasets efficiently.
3. **Filtering and Sorting**: Allow filtering and sorting options for list endpoints to improve usability.
4. **Validation**: Implement thorough request validation to ensure data integrity.
5. **Error Handling**: Provide clear error messages and appropriate HTTP status codes.
6. **Documentation**: Create comprehensive API documentation with examples for each endpoint.
7. **Rate Limiting**: Implement rate limiting to prevent abuse.
8. **Versioning**: Consider implementing API versioning to allow for future changes without breaking existing clients.
9. **CORS**: Configure CORS appropriately to allow access from authorized client applications.
10. **Logging**: Implement request logging for debugging and audit purposes.

This comprehensive list of endpoints covers all the functionality represented in the class diagram and provides a solid foundation for building the educational platform's API.
