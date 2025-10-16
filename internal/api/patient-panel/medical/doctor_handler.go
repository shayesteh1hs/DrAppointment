package medical

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"drgo/internal/domain/medical"
	medicalFilter "drgo/internal/filter/medical"
	"drgo/internal/pagination"
	medicalRepo "drgo/internal/repository/medical"
)

type Handler struct {
	repo medicalRepo.DoctorRepository
}

func NewHandler(repo medicalRepo.DoctorRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	var paginationParams pagination.LimitOffsetParams
	if err := c.ShouldBindQuery(&paginationParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}
	if err := paginationParams.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var filterPrams medicalFilter.DoctorQueryParam
	if err := c.ShouldBindQuery(&filterPrams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters"})
		return
	}
	if err := filterPrams.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paginator := pagination.NewLimitOffsetPaginator[medical.Doctor](paginationParams)
	result, err := h.repo.GetAllOffset(c.Request.Context(), filterPrams, paginator)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	doctorRoutes := router.Group("/doctors")
	{
		doctorRoutes.GET("", h.GetAll)
	}
}
