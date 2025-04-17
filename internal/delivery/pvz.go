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
	pvzUsecase usecase.PVZUsecase
}

func NewPVZHandler(pvzUsecase usecase.PVZUsecase) *PVZHandler {
	return &PVZHandler{pvzUsecase: pvzUsecase}
}

func (h *PVZHandler) PostPVZ(c *gin.Context) {
	role, _ := c.Get("role")
	id, _ := c.Get("userID")

	if role.(string) != "moderator" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permision denied"})
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

	if input.Id.IsNil() || input.City == "" {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pvz)
}

func (h *PVZHandler) GetPVZs(c *gin.Context) {
	var filter entity.Filter
	role, _ := c.Get("role")

	if role.(string) != "moderator" && role.(string) != "employee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permision denied"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}
	filter.Page = page

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	filter.Limit = limit

	if startDateStr := c.Query("startDate"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate format"})
			return
		}
		filter.StartDate = &startDate
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate format"})
			return
		}
		filter.EndDate = &endDate
	}

	// Вызов usecase
	response, err := h.pvzUsecase.GetPVZsWithFilter(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
