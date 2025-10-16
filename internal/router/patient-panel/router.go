package patient_panel

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	medical_api "drgo/internal/api/patient-panel/medical"
	"drgo/internal/repository/medical"
)

func SetupPatientPanelRoutes(rg *gin.RouterGroup, db *sql.DB) {
	doctorRepo := medical.NewDoctorRepository(db)
	doctorHandler := medical_api.NewHandler(doctorRepo)
	doctorHandler.RegisterRoutes(rg)

}
