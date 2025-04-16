package delivery

import (
	"net/http"
	"pvz/internal/storage/usecase"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginUsecase *usecase.UserUsecase
}

func NewLoginHandler(loginUsecase *usecase.UserUsecase) *LoginHandler {
	return &LoginHandler{loginUsecase: loginUsecase}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.loginUsecase.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
