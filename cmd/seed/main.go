package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"

	"drgo/internal/database"
	medicalDomain "drgo/internal/domain/medical"
	medicalRepo "drgo/internal/repository/medical"
	"drgo/internal/utils"
)

func main() {
	dbConfig := database.Config{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnvInt("DB_PORT", 5432),
		User:     utils.GetEnv("DB_USER", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", "postgres"),
		DBName:   utils.GetEnv("DB_NAME", "drgo"),
		SSLMode:  utils.GetEnv("DB_SSL_MODE", "disable"),
	}

	db, err := database.Connect(&dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to close database")
		}
	}(db)

	doctorRepo := medicalRepo.NewDoctorRepository(db)

	log.Println("Starting to seed doctors...")

	specialties := []string{
		"Cardiologist",
		"Dermatologist",
		"Neurologist",
		"Orthopedist",
		"Pediatrician",
		"Psychiatrist",
		"General Practitioner",
		"Ophthalmologist",
		"ENT Specialist",
		"Gynecologist",
	}

	doctorCount := 50
	for i := 0; i < doctorCount; i++ {
		doc := &medicalDomain.Doctor{
			ID:          uuid.New(),
			Name:        gofakeit.Name(),
			Specialty:   specialties[gofakeit.Number(0, len(specialties)-1)],
			PhoneNumber: gofakeit.Phone(),
			AvatarURL:   gofakeit.ImageURL(200, 200),
			Description: gofakeit.Sentence(15),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := doctorRepo.Create(context.Background(), doc)
		if err != nil {
			log.Printf("Failed to create doctor %s: %v", doc.Name, err)
			continue
		}

		log.Printf("Created doctor: %s (%s)", doc.Name, doc.Specialty)
	}

	log.Printf("Successfully seeded %d doctors!", doctorCount)
}
