package delivery

import (
	"net/http"
	"pvz/internal/usecase"

	"github.com/gin-gonic/gin"
)

type DummyLoginHandler struct {
	dummyLoginUsecase usecase.UserUsecase
}

func NewDummyLoginHandler(loginUsecase usecase.UserUsecase) *DummyLoginHandler {
	return &DummyLoginHandler{dummyLoginUsecase: loginUsecase}
}

func (h *DummyLoginHandler) DummyLogin(c *gin.Context) {

	var input struct {
		Role string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	switch input.Role {
	case "moderator":
		token, err := h.dummyLoginUsecase.Login("dummy_moderator@test.com", "supersecretpassword")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})

	case "employee":
		token, err := h.dummyLoginUsecase.Login("dummy_employee@test.com", "supersecretpassword")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

}
