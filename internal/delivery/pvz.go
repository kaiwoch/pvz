package delivery

import (
	"net/http"
	"pvz/internal/storage/migrations/entity"
	"pvz/internal/usecase"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type PVZHandler struct {
	pvzUsecase *usecase.PVZUsecase
}

func NewPVZHandler(pvzUsecase *usecase.PVZUsecase) *PVZHandler {
	return &PVZHandler{pvzUsecase: pvzUsecase}
}

func (h *PVZHandler) PostPVZ(c *gin.Context) {
	role, _ := c.Get("role")
	id, _ := c.Get("userID")

	if role.(string) != "moderator" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permision denied"})
		return
	}

	var input struct {
		Id               uuid.UUID `json:"id"`
		RegistrationDate time.Time `json:"registrationDate"`
		City             string    `json:"city"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	idx, err := uuid.FromString(id.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pvz, err := h.pvzUsecase.CreatePVZ(input.Id, idx, input.City, input.RegistrationDate)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pvz": pvz})
}

func (h *PVZHandler) GetPVZs(c *gin.Context) {
	role, _ := c.Get("role")

	if role.(string) != "moderator" && role.(string) != "employee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permision denied"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var startDate, endDate *time.Time
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &t
		}
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &t
		}
	}

	filter := entity.Filter{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	}

	// Вызов usecase
	response, err := h.pvzUsecase.GetPVZsWithFilter(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
