package medical

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shayesteh1hs/DrAppointment/internal/domain/medical"
	medicalFilter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	medicalRepo "github.com/shayesteh1hs/DrAppointment/internal/repository/medical"
)

type Handler struct {
	repo medicalRepo.DoctorRepository
}

func NewHandler(repo medicalRepo.DoctorRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) GetAllPaginated(c *gin.Context) {
	var paginationParams pagination.LimitOffsetParams
	if err := c.ShouldBindQuery(&paginationParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}
	if err := paginationParams.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var filterParams medicalFilter.DoctorQueryParam
	if err := c.ShouldBindQuery(&filterParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters"})
		return
	}
	if err := filterParams.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalCount, err := h.repo.Count(c.Request.Context(), filterParams)
	if err != nil {
		log.Printf("failed to fetch doctors count: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors count"})
	}

	paginator := pagination.NewLimitOffsetPaginator[medical.Doctor](paginationParams)
	doctors, err := h.repo.GetAllPaginated(c.Request.Context(), filterParams, paginator)
	if err != nil {
		log.Printf("failed to fetch doctors: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}

	result := paginator.CreatePaginationResult(doctors, totalCount)
	c.JSON(http.StatusOK, result)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	doctorRoutes := router.Group("/doctors")

	doctorRoutes.GET("", h.GetAllPaginated)
}
