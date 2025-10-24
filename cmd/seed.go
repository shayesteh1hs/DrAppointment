package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/database"
	"github.com/shayesteh1hs/DrAppointment/internal/utils"
)

// Sample specialties data
var sampleSpecialties = []struct {
	Name string
}{
	{"Cardiology"},
	{"Dermatology"},
	{"Neurology"},
	{"Orthopedics"},
	{"Pediatrics"},
	{"Psychiatry"},
	{"Radiology"},
	{"Surgery"},
	{"Internal Medicine"},
	{"Gynecology"},
}

// Sample doctors data
var sampleDoctors = []struct {
	Name        string
	Specialty   string // Will be matched to specialty name
	PhoneNumber string
	AvatarURL   string
	Description string
}{
	{
		Name:        "Dr. Sarah Johnson",
		Specialty:   "Cardiology",
		PhoneNumber: "+1-555-0101",
		AvatarURL:   "https://example.com/avatars/sarah-johnson.jpg",
		Description: "Experienced cardiologist with 15 years of practice in heart disease treatment and prevention.",
	},
	{
		Name:        "Dr. Michael Chen",
		Specialty:   "Dermatology",
		PhoneNumber: "+1-555-0102",
		AvatarURL:   "https://example.com/avatars/michael-chen.jpg",
		Description: "Board-certified dermatologist specializing in skin cancer detection and cosmetic dermatology.",
	},
	{
		Name:        "Dr. Emily Rodriguez",
		Specialty:   "Neurology",
		PhoneNumber: "+1-555-0103",
		AvatarURL:   "https://example.com/avatars/emily-rodriguez.jpg",
		Description: "Neurologist with expertise in treating epilepsy, multiple sclerosis, and movement disorders.",
	},
	{
		Name:        "Dr. James Wilson",
		Specialty:   "Orthopedics",
		PhoneNumber: "+1-555-0104",
		AvatarURL:   "https://example.com/avatars/james-wilson.jpg",
		Description: "Orthopedic surgeon specializing in joint replacement and sports medicine.",
	},
	{
		Name:        "Dr. Lisa Thompson",
		Specialty:   "Pediatrics",
		PhoneNumber: "+1-555-0105",
		AvatarURL:   "https://example.com/avatars/lisa-thompson.jpg",
		Description: "Pediatrician with 12 years of experience caring for children from infancy to adolescence.",
	},
	{
		Name:        "Dr. Robert Davis",
		Specialty:   "Psychiatry",
		PhoneNumber: "+1-555-0106",
		AvatarURL:   "https://example.com/avatars/robert-davis.jpg",
		Description: "Psychiatrist specializing in anxiety disorders, depression, and cognitive behavioral therapy.",
	},
	{
		Name:        "Dr. Jennifer Brown",
		Specialty:   "Radiology",
		PhoneNumber: "+1-555-0107",
		AvatarURL:   "https://example.com/avatars/jennifer-brown.jpg",
		Description: "Radiologist with expertise in MRI, CT scans, and diagnostic imaging.",
	},
	{
		Name:        "Dr. David Miller",
		Specialty:   "Surgery",
		PhoneNumber: "+1-555-0108",
		AvatarURL:   "https://example.com/avatars/david-miller.jpg",
		Description: "General surgeon with specialization in minimally invasive procedures and laparoscopic surgery.",
	},
	{
		Name:        "Dr. Amanda Garcia",
		Specialty:   "Internal Medicine",
		PhoneNumber: "+1-555-0109",
		AvatarURL:   "https://example.com/avatars/amanda-garcia.jpg",
		Description: "Internal medicine physician focusing on preventive care and chronic disease management.",
	},
	{
		Name:        "Dr. Christopher Lee",
		Specialty:   "Gynecology",
		PhoneNumber: "+1-555-0110",
		AvatarURL:   "https://example.com/avatars/christopher-lee.jpg",
		Description: "Gynecologist providing comprehensive women's health care including obstetrics and gynecology.",
	},
}

func main() {
	// Database configuration
	dbConfig := database.Config{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnvInt("DB_PORT", 5432),
		User:     utils.GetEnv("DB_USER", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", "postgres"),
		DBName:   utils.GetEnv("DB_NAME", "drgo"),
		SSLMode:  utils.GetEnv("DB_SSL_MODE", "disable"),
	}

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, &dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run seeder
	if err := seedDatabase(ctx, db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully!")
}

func seedDatabase(ctx context.Context, db *sql.DB) error {
	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Seed specialties
	specialtyMap, err := seedSpecialties(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to seed specialties: %w", err)
	}

	// Seed doctors
	if err := seedDoctors(ctx, tx, specialtyMap); err != nil {
		return fmt.Errorf("failed to seed doctors: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func seedSpecialties(ctx context.Context, tx *sql.Tx) (map[string]uuid.UUID, error) {
	specialtyMap := make(map[string]uuid.UUID)

	query := `
		INSERT INTO specialties (id, name, created_at, updated_at) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO NOTHING
		RETURNING id, name
	`

	for _, specialty := range sampleSpecialties {
		// Check if specialty already exists
		var existingID uuid.UUID
		checkQuery := `SELECT id FROM specialties WHERE name = $1`
		err := tx.QueryRowContext(ctx, checkQuery, specialty.Name).Scan(&existingID)
		if err == nil {
			// Specialty exists, add to map
			specialtyMap[specialty.Name] = existingID
			log.Printf("Specialty '%s' already exists, skipping", specialty.Name)
			continue
		} else if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to check existing specialty: %w", err)
		}

		// Create new specialty
		id := uuid.New()
		now := time.Now()

		var insertedID uuid.UUID
		var insertedName string
		err = tx.QueryRowContext(ctx, query, id, specialty.Name, now, now).Scan(&insertedID, &insertedName)
		if err != nil {
			return nil, fmt.Errorf("failed to insert specialty '%s': %w", specialty.Name, err)
		}

		specialtyMap[insertedName] = insertedID
		log.Printf("Inserted specialty: %s (ID: %s)", insertedName, insertedID)
	}

	return specialtyMap, nil
}

func seedDoctors(ctx context.Context, tx *sql.Tx, specialtyMap map[string]uuid.UUID) error {
	query := `
		INSERT INTO doctors (id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	for _, doctor := range sampleDoctors {
		// Get specialty ID
		specialtyID, exists := specialtyMap[doctor.Specialty]
		if !exists {
			log.Printf("Warning: Specialty '%s' not found for doctor '%s', skipping", doctor.Specialty, doctor.Name)
			continue
		}

		// Check if doctor already exists
		var existingID uuid.UUID
		checkQuery := `SELECT id FROM doctors WHERE phone_number = $1`
		err := tx.QueryRowContext(ctx, checkQuery, doctor.PhoneNumber).Scan(&existingID)
		if err == nil {
			log.Printf("Doctor with phone number '%s' already exists, skipping", doctor.PhoneNumber)
			continue
		} else if err != sql.ErrNoRows {
			return fmt.Errorf("failed to check existing doctor: %w", err)
		}

		// Create new doctor
		id := uuid.New()
		now := time.Now()

		_, err = tx.ExecContext(ctx, query, id, doctor.Name, specialtyID, doctor.PhoneNumber, doctor.AvatarURL, doctor.Description, now, now)
		if err != nil {
			return fmt.Errorf("failed to insert doctor '%s': %w", doctor.Name, err)
		}

		log.Printf("Inserted doctor: %s (Specialty: %s)", doctor.Name, doctor.Specialty)
	}

	return nil
}
