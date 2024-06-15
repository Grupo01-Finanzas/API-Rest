package controller

import (
	"errors"
	"net/http"
	"strconv"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Retrieve a list of all products.
// @Tags Products
// @Produce json
// @Success 200 {array} response.ProductResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /products [get]
func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.productService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Retrieve a product by its ID.
// @Tags Products
// @Produce json
// @Param id path uint true "Product ID"
// @Success 200 {object} response.ProductResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /products/{id} [get]
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	product, err := c.productService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// GetProductsByEstablishmentID godoc
// @Summary Get products by establishment ID
// @Description Retrieve products associated with a specific establishment ID.
// @Tags Products
// @Produce json
// @Param establishment_id path uint true "Establishment ID"
// @Success 200 {array} response.ProductResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /establishments/{establishment_id}/products [get]
func (c *ProductController) GetProductsByEstablishmentID(ctx *gin.Context) {
	establishmentID, err := strconv.ParseUint(ctx.Param("establishment_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	products, err := c.productService.GetByEstablishmentID(uint(establishmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product associated with an establishment.
// @Tags Products
// @Accept json
// @Produce json
// @Param product body request.CreateProductRequest true "Product details"
// @Success 201 {object} response.ProductResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req request.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	product, err := c.productService.Create(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update an existing product by ID.
// @Tags Products
// @Accept json
// @Produce json
// @Param id path uint true "Product ID"
// @Param product body request.UpdateProductRequest true "Product details"
// @Success 200 {object} response.ProductResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	var req request.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	product, err := c.productService.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID.
// @Tags Products
// @Produce json
// @Param id path uint true "Product ID"
// @Success 204
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	err = c.productService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
