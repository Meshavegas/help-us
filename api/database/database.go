package database

import (
	"api/models"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bxcodec/faker/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initialise la connexion à la base de données
func InitDatabase() {
	var err error

	// Récupérer les informations de connexion depuis les variables d'environnement
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	// Construire la chaîne de connexion PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)

	// Configuration du logger GORM
	logLevel := logger.Silent
	if os.Getenv("DB_DEBUG") == "true" {
		logLevel = logger.Info
	}

	// Connexion à la base de données PostgreSQL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatal("Erreur lors de la connexion à la base de données:", err)
	}

	fmt.Println("Connexion à la base de données établie")

	// Migration automatique des modèles
	err = AutoMigrate()
	if err != nil {
		log.Fatal("Erreur lors de la migration:", err)
	}

	fmt.Println("Migration de la base de données terminée")
}

// AutoMigrate effectue la migration automatique de tous les modèles
func AutoMigrate() error {
	return DB.AutoMigrate(
		// Modèles de base
		&models.User{},
		&models.Address{},

		// Modèles spécialisés d'utilisateurs
		&models.Famille{},
		&models.Enseignant{},
		&models.Administrator{},

		// Modèles de cours et missions
		&models.Course{},
		&models.Mission{},
		&models.Report{},

		// Modèles de paiement et offres
		&models.Payment{},
		&models.Offer{},
		&models.Option{},

		// Modèles de ressources
		&models.Resource{},
	)
}

// CloseDatabase ferme la connexion à la base de données
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB retourne l'instance de la base de données
func GetDB() *gorm.DB {
	return DB
}

// ResetDatabase supprime et recrée toutes les tables (utile pour les tests)
func ResetDatabase() error {
	// Supprimer toutes les tables
	err := DB.Migrator().DropTable(
		&models.User{},
		&models.Address{},
		&models.Famille{},
		&models.Enseignant{},
		&models.Administrator{},
		&models.Course{},
		&models.Mission{},
		&models.Report{},
		&models.Payment{},
		&models.Offer{},
		&models.Option{},
		&models.Resource{},
		"user_resources",    // Table de liaison many2many
		"enseignant_offers", // Table de liaison many2many
	)
	if err != nil {
		return err
	}

	// Recréer les tables
	return AutoMigrate()
}

// --- SEEDING CONSTANTS ---
const (
	NB_ADMINS                   = 2
	NB_ENSEIGNANTS              = 50
	NB_FAMILLES                 = 200
	NB_OFFERS                   = 30
	NB_MISSIONS_PER_FAMILLE_MAX = 3
	NB_COURS_PER_MISSION_MAX    = 15
	NB_RESOURCES                = 100
)

// SeedDatabase remplace les données existantes par un jeu de données de test complet et cohérent.
func SeedDatabase() error {
	if os.Getenv("SEED_DATABASE") != "true" {
		fmt.Println("Variable SEED_DATABASE non activée, seeding ignoré.")
		return nil
	}

	fmt.Println("--- Début du Seeding de la base de données ---")
	fmt.Println("Réinitialisation de la base de données...")
	if err := ResetDatabase(); err != nil {
		return fmt.Errorf("erreur lors de la réinitialisation de la DB: %w", err)
	}

	// Hashing du mot de passe commun
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erreur lors du hachage du mot de passe: %w", err)
	}
	password := string(hashedPassword)

	// --- Seeding ---
	fmt.Println("Création des administrateurs...")
	admins, err := seedUsers(password, models.RoleAdministrator, NB_ADMINS)
	if err != nil {
		return err
	}

	fmt.Println("Création des enseignants...")
	enseignants, err := seedUsers(password, models.RoleEnseignant, NB_ENSEIGNANTS)
	if err != nil {
		return err
	}

	fmt.Println("Création des familles...")
	familles, err := seedUsers(password, models.RoleFamille, NB_FAMILLES)
	if err != nil {
		return err
	}

	allUsers := append(append(admins, enseignants...), familles...)

	fmt.Println("Création des adresses...")
	userAddresses, err := seedAddresses(allUsers)
	if err != nil {
		return err
	}

	fmt.Println("Création des offres...")
	offers, err := seedOffers(admins)
	if err != nil {
		return err
	}

	fmt.Println("Création cohérente des missions, cours, paiements et rapports...")
	if err := seedCoherentMissionsAndSubData(familles, enseignants, userAddresses); err != nil {
		return err
	}

	fmt.Println("Création des options...")
	if err := seedOptions(offers, familles, enseignants); err != nil {
		return err
	}

	fmt.Println("Création des ressources...")
	if err := seedResources(admins, allUsers); err != nil {
		return err
	}

	fmt.Println("--- Seeding terminé avec succès ---")
	fmt.Printf("Comptes créés : %d Admins, %d Enseignants, %d Familles. Le mot de passe pour tous est 'password'.\n", len(admins), len(enseignants), len(familles))

	return nil
}

func seedUsers(password string, role models.UserRole, count int) ([]models.User, error) {
	var users []models.User
	for i := 0; i < count; i++ {
		user := models.User{
			Username:    faker.Username(),
			Email:       faker.Email(),
			Password:    password,
			Role:        role,
			PhoneNumber: faker.Phonenumber(),
			IsActive:    true,
		}
		if err := DB.Create(&user).Error; err != nil {
			return nil, err
		}

		switch role {
		case models.RoleAdministrator:
			profile := models.Administrator{UserID: user.ID}
			DB.Create(&profile)
		case models.RoleEnseignant:
			profile := models.Enseignant{
				UserID:         user.ID,
				Specialization: faker.Word(),
				Qualifications: faker.Sentence(),
			}
			DB.Create(&profile)
		case models.RoleFamille:
			profile := models.Famille{
				UserID:     user.ID,
				FamilyName: faker.LastName(),
			}
			DB.Create(&profile)
		}
		users = append(users, user)
	}
	return users, nil
}

func seedAddresses(users []models.User) (map[uint]models.Address, error) {
	addresses := make(map[uint]models.Address)
	for _, user := range users {
		address := models.Address{
			Street:     faker.Sentence(),
			City:       faker.Word(),
			PostalCode: faker.Word(),
			Country:    "France",
			Latitude:   rand.Float64()*2 + 48, // Lat around Paris
			Longitude:  rand.Float64()*2 + 2,  // Lon around Paris
			UserID:     user.ID,
		}
		if err := DB.Create(&address).Error; err != nil {
			return nil, err
		}
		addresses[user.ID] = address
	}
	return addresses, nil
}

func seedOffers(admins []models.User) ([]models.Offer, error) {
	var offers []models.Offer
	for i := 0; i < NB_OFFERS; i++ {
		offer := models.Offer{
			Title:           faker.Sentence(),
			Description:     faker.Paragraph(),
			HourlyRate:      float64(rand.Intn(25) + 15), // Rate between 15 and 40
			PublicationDate: time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
			Status:          models.OfferStatusOpen,
			Requirements:    faker.Paragraph(),
			Subject:         faker.Word(),
			Level:           []string{"Collège", "Lycée", "Supérieur"}[rand.Intn(3)],
			CreatedByID:     admins[rand.Intn(len(admins))].ID,
		}
		if err := DB.Create(&offer).Error; err != nil {
			return nil, err
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func seedCoherentMissionsAndSubData(familles, enseignants []models.User, addresses map[uint]models.Address) error {
	for _, famille := range familles {
		// Assign 1 to 3 teachers to this family
		assignedTeacherCount := rand.Intn(NB_MISSIONS_PER_FAMILLE_MAX) + 1
		for i := 0; i < assignedTeacherCount; i++ {
			enseignant := enseignants[rand.Intn(len(enseignants))]

			// Create a mission for this pair
			mission := models.Mission{
				StartDate:    time.Now().Add(-time.Duration(rand.Intn(60)) * 24 * time.Hour),
				Status:       models.MissionStatusActive,
				Description:  fmt.Sprintf("Mission de %s pour la famille %s", enseignant.Username, famille.Username),
				FamilleID:    famille.ID,
				EnseignantID: enseignant.ID,
			}
			if err := DB.Create(&mission).Error; err != nil {
				return err
			}

			// For this mission, create some courses
			courseCount := rand.Intn(NB_COURS_PER_MISSION_MAX) + 1
			for j := 0; j < courseCount; j++ {
				familleAddress := addresses[famille.ID]
				course := models.Course{
					ScheduledTime: time.Now().Add(time.Duration(rand.Intn(30)-15) * 24 * time.Hour),
					Duration:      []int{60, 90, 120}[rand.Intn(3)],
					Location:      "À domicile",
					Status:        models.CourseStatusScheduled,
					FamilleID:     mission.FamilleID,    // Coherent
					EnseignantID:  mission.EnseignantID, // Coherent
					MissionID:     mission.ID,           // Coherent
					AddressID:     familleAddress.ID,    // Coherent
				}
				if err := DB.Create(&course).Error; err != nil {
					return err
				}

				// Create a payment for this course (70% chance)
				if rand.Intn(10) < 7 {
					payment := models.Payment{
						Amount:      float64(course.Duration) / 60.0 * float64(rand.Intn(25)+15), // Random hourly rate
						PaymentDate: course.ScheduledTime.Add(time.Hour),
						Status:      []models.PaymentStatus{models.PaymentStatusCompleted, models.PaymentStatusPending}[rand.Intn(2)],
						Type:        models.PaymentTypeCourse,
						Description: fmt.Sprintf("Paiement pour le cours #%d", course.ID),
						UserID:      course.FamilleID, // Coherent
						CourseID:    course.ID,        // Coherent
					}
					if err := DB.Create(&payment).Error; err != nil {
						return err
					}
				}
			}

			// Create a report for this mission (50% chance)
			if rand.Intn(2) == 0 {
				report := models.Report{
					SubmissionDate: time.Now(),
					Content:        faker.Paragraph(),
					Status:         models.ReportStatusSubmitted,
					EnseignantID:   mission.EnseignantID, // Coherent
					MissionID:      mission.ID,           // Coherent
				}
				if err := DB.Create(&report).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func seedOptions(offers []models.Offer, familles, enseignants []models.User) error {
	for _, offer := range offers {
		optionCount := rand.Intn(5) + 1
		for i := 0; i < optionCount; i++ {
			teacherUser := enseignants[rand.Intn(len(enseignants))]
			familyUser := familles[rand.Intn(len(familles))]

			// 1) On s'assure que le profil enseignant existe
			var teacherProfile models.Enseignant
			if err := DB.First(&teacherProfile, "user_id = ?", teacherUser.ID).Error; err != nil {
				continue // passe au suivant, le profil n'a pas été trouvé (ne devrait pas arriver)
			}

			option := models.Option{
				CreationDate:   time.Now(),
				ExpirationDate: time.Now().Add(7 * 24 * time.Hour),
				Status:         models.OptionStatusActive,
				Description:    "Option pour l'offre " + offer.Title,
				EnseignantID:   teacherProfile.UserID, // cohérent
				FamilleID:      familyUser.ID,
				OfferID:        offer.ID,
			}
			if err := DB.Create(&option).Error; err != nil {
				return err
			}

			// 2) on renseigne la table many-to-many enseignant_offers
			DB.Model(&teacherProfile).Association("Offers").Append(&offer)
		}
	}
	return nil
}

func seedResources(admins, users []models.User) error {
	for i := 0; i < NB_RESOURCES; i++ {
		resource := models.Resource{
			Title:       faker.Sentence(),
			Type:        []models.ResourceType{"document", "video", "link"}[rand.Intn(3)],
			URL:         faker.URL(),
			Description: faker.Paragraph(),
			IsPublic:    rand.Intn(2) == 0,
			ManagedByID: admins[rand.Intn(len(admins))].ID,
		}
		if err := DB.Create(&resource).Error; err != nil {
			return err
		}

		// Associate with some users if not public
		if !resource.IsPublic {
			numUsers := rand.Intn(5) + 1
			for j := 0; j < numUsers; j++ {
				user := users[rand.Intn(len(users))]
				DB.Model(&resource).Association("Users").Append(&user)
			}
		}
	}
	return nil
}
