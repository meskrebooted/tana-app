package config

import (
	"fmt"
	"log"
	"os"

	"coworking-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect opens the Postgres connection using env vars and runs migrations.
func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "coworking"),
		getEnv("DB_PASSWORD", "coworking"),
		getEnv("DB_NAME", "coworking"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Space{}, &models.Booking{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	DB = db
	log.Println("database connected and migrated")

	seedSpaces(db)
}

// seedSpaces inserts a starter set of coworking spaces if the table is empty.
func seedSpaces(db *gorm.DB) {
	var count int64
	db.Model(&models.Space{}).Count(&count)
	if count > 0 {
		return
	}

	spaces := []models.Space{
		{Name: "La Tana del Procione", Type: "desk", Capacity: 1, Location: "Open space - vicino alla finestra",
			Vibe: "Per chi lavora meglio con la luce naturale e un caffè sempre a portata di mano."},
		{Name: "L'Angolo del Silenzio", Type: "desk", Capacity: 1, Location: "Open space - zona fondo",
			Vibe: "Cuffie, focus e nessuna interruzione: la postazione per le sessioni di deep work."},
		{Name: "Il Faro", Type: "desk", Capacity: 1, Location: "Open space - centro sala",
			Vibe: "Postazione rialzata, ottima se ti piace tenere d'occhio tutta la stanza."},
		{Name: "Sala Vulcano", Type: "room", Capacity: 6, Location: "Piano 2",
			Vibe: "Whiteboard enorme e sedie comode: qui nascono i brainstorming che contano."},
		{Name: "Sala Bussola", Type: "room", Capacity: 10, Location: "Piano 2",
			Vibe: "La sala per le riunioni serie, con schermo grande e acustica pulita."},
		{Name: "Cabina Sussurro", Type: "booth", Capacity: 1, Location: "Open space - corridoio",
			Vibe: "Phone booth insonorizzata: perfetta per la call che non vuoi far sentire a tutti."},
	}

	if err := db.Create(&spaces).Error; err != nil {
		log.Printf("warning: could not seed spaces: %v", err)
	} else {
		log.Println("seeded starter coworking spaces")
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
