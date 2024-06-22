package controller

import (
	"net/http"
	"strconv"

	"ApiRestFinance/internal/middleware"
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/service"

	"github.com/gin-gonic/gin"
)

// ProductController handles product-related endpoints.
type ProductController struct {
	productService service.ProductService
}

// NewProductController creates a new instance of ProductController.
func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

// CreateProduct godoc
// @Summary      Create Product
// @Description  Creates a new product for the authenticated admin\'s establishment.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string                  true  "Bearer {token}"
// @Param        product        body      request.CreateProductRequest  true  "Product data"
// @Success      201  {object}  response.ProductResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req request.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Only admins can create products
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can create products"})
		return
	}

	establishmentID := middleware.GetEstablishmentIDFromContext(ctx)
	req.EstablishmentID = establishmentID

	product, err := c.productService.CreateProduct(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// GetProductByID godoc
// @Summary      Get Product by ID
// @Description  Gets a product by its ID.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id             path      int  true  "Product ID"
// @Success      200  {object}  response.ProductResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products/{id} [get]
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	product, err := c.productService.GetProductByID(uint(productID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// GetAllProductsByEstablishmentID godoc
// @Summary      Get Products by Establishment ID
// @Description  Gets all products associated with an establishment.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        establishmentID   path      int  true  "Establishment ID"
// @Success      200  {array}   response.ProductResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /establishments/{establishmentID}/products [get]
func (c *ProductController) GetAllProductsByEstablishmentID(ctx *gin.Context) {
	establishmentID, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	// You may want to add authorization logic here (admin only or public?)

	products, err := c.productService.GetAllProductsByEstablishmentID(uint(establishmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// UpdateProduct godoc
// @Summary      Update Product
// @Description  Updates an existing product. Only admins can update products.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id             path      int                      true  "Product ID"
// @Param        product        body      request.UpdateProductRequest  true  "Updated product data"
// @Success      200  {object}  response.ProductResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	var req request.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Check user role - Only Admins can update products
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can update products"})
		return
	}

	updatedProduct, err := c.productService.UpdateProduct(uint(productID), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct godoc
// @Summary      Delete Product
// @Description  Deletes a product by its ID. Only Admins can delete products.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id             path      int  true  "Product ID"
// @Success      204  "No Content"
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	// Only Admins can delete products
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can delete products"})
		return
	}

	if err := c.productService.DeleteProduct(uint(productID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent) // 204 No Content on successful deletion
}
