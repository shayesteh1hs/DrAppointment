package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"

	"drgo/internal/database"
	"drgo/internal/domain"
	medicalRepo "drgo/internal/repository/medical"
)

func main() {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "drgo"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

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
		doc := &domain.Doctor{
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

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
