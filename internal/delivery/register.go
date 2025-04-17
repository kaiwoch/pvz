package delivery

import (
	"net/http"
	"pvz/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	registerUsecase usecase.UserUsecase
}

func NewRegisterHandler(registerUsecase usecase.UserUsecase) *RegisterHandler {
	return &RegisterHandler{registerUsecase: registerUsecase}
}

func (h *RegisterHandler) Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id, err := h.registerUsecase.Register(input.Email, input.Password, input.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
