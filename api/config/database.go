package config

import (
	"log"
	"os"

	"api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase initialise la connexion à la base de données
func ConnectDatabase() {
	var err error

	// Configuration du logger GORM
	gormLogger := logger.Default
	if os.Getenv("GIN_MODE") == "release" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Connexion à la base de données SQLite
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Fatal("Erreur lors de la connexion à la base de données:", err)
	}

	// Migration automatique des modèles
	err = DB.AutoMigrate(
		// Base models
		&models.User{},
		&models.Address{},

		// User role models
		&models.Famille{},
		&models.Enseignant{},
		&models.Administrator{},

		// Core business models
		&models.Mission{},
		&models.Course{},
		&models.Payment{},
		&models.Report{},
		&models.Offer{},
		&models.Option{},
		&models.Resource{},
	)

	if err != nil {
		log.Fatal("Erreur lors de la migration de la base de données:", err)
	}

	log.Println("Base de données connectée et migrée avec succès")
}

// GetDB retourne l'instance de la base de données
func GetDB() *gorm.DB {
	return DB
}
