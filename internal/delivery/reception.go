package delivery

import (
	"net/http"
	"pvz/internal/storage/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type ReceptionHandler struct {
	receptionUsecase *usecase.ReceptionUsecase
}

func NewReceptionHandler(receptionUsecase *usecase.ReceptionUsecase) *ReceptionHandler {
	return &ReceptionHandler{receptionUsecase: receptionUsecase}
}

func (h *ReceptionHandler) Reception(c *gin.Context) {
	role, _ := c.Get("role")
	//id, _ := c.Get("userID")

	if role.(string) != "employee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permision denied"})
		return
	}

	var input struct {
		ID uuid.UUID `json:"pvzId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if input.ID.String() == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	reception, err := h.receptionUsecase.CreateReception(input.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reception)
}

func (h *ReceptionHandler) UpdateReceptionStatus(c *gin.Context) {
	role, _ := c.Get("role")
	//id, _ := c.Get("userID")

	if role.(string) != "employee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permision denied"})
		return
	}

	pvz_id_string := c.Param("pvzId")

	pvz_id, err := uuid.FromString(pvz_id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong query"})
		return
	}

	reception, err := h.receptionUsecase.UpdateReceptionStatus(pvz_id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reception)
}
