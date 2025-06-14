package database

import (
	"api/models"
	"fmt"
	"log"
	"os"

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

// SeedDatabase ajoute des données de test
func SeedDatabase() error {
	// Vérifier si des données existent déjà
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		fmt.Println("La base de données contient déjà des données")
		return nil
	}

	fmt.Println("Ajout de données de test...")

	// Créer un administrateur par défaut
	admin := models.User{
		Username:    "admin",
		Email:       "mesha@mm.com",
		Password:    "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Role:        models.RoleAdministrator,
		PhoneNumber: "+33123456789",
		IsActive:    true,
	}

	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	// Créer le profil administrateur
	adminProfile := models.Administrator{
		UserID: admin.ID,
	}

	if err := DB.Create(&adminProfile).Error; err != nil {
		return err
	}

	// Créer un enseignant de test
	teacher := models.User{
		Username:    "teacher1",
		Email:       "teacher@example.com",
		Password:    "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Role:        models.RoleEnseignant,
		PhoneNumber: "+33123456790",
		IsActive:    true,
	}

	if err := DB.Create(&teacher).Error; err != nil {
		return err
	}

	// Créer le profil enseignant
	teacherProfile := models.Enseignant{
		UserID:         teacher.ID,
		Specialization: "Mathématiques",
		Qualifications: "Master en Mathématiques, 5 ans d'expérience",
	}

	if err := DB.Create(&teacherProfile).Error; err != nil {
		return err
	}

	// Créer une famille de test
	family := models.User{
		Username:    "famille1",
		Email:       "famille@example.com",
		Password:    "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Role:        models.RoleFamille,
		PhoneNumber: "+33123456791",
		IsActive:    true,
	}

	if err := DB.Create(&family).Error; err != nil {
		return err
	}

	// Créer le profil famille
	familyProfile := models.Famille{
		UserID:     family.ID,
		FamilyName: "Famille Dupont",
	}

	if err := DB.Create(&familyProfile).Error; err != nil {
		return err
	}

	fmt.Println("Données de test ajoutées avec succès")
	fmt.Println("Comptes créés:")
	fmt.Println("- Admin: admin@example.com / password")
	fmt.Println("- Enseignant: teacher@example.com / password")
	fmt.Println("- Famille: famille@example.com / password")

	return nil
}
