package delivery

import (
	"net/http"
	"pvz/internal/storage/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
}

func NewProductHandler(productUsecase *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
}

func (h *ProductHandler) Reception(c *gin.Context) {
	role, _ := c.Get("role")

	if role.(string) != "employee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permision denied"})
		return
	}

	var input struct {
		ProductType string    `json:"type"`
		ID          uuid.UUID `json:"pvzId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if input.ID.String() == "" || input.ProductType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	product, err := h.productUsecase.CreateProduct(input.ID, input.ProductType)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteLastProduct(c *gin.Context) {
	role, _ := c.Get("role")

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

	err = h.productUsecase.DeleteLastProduct(pvz_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, "")
}
