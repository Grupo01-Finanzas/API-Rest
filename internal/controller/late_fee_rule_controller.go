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

// LateFeeRuleController handles API requests related to late fee rules.
type LateFeeRuleController struct {
	lateFeeRuleService service.LateFeeRuleService
}

// NewLateFeeRuleController creates a new LateFeeRuleController instance.
func NewLateFeeRuleController(lateFeeRuleService service.LateFeeRuleService) *LateFeeRuleController {
	return &LateFeeRuleController{lateFeeRuleService: lateFeeRuleService}
}

// CreateLateFeeRule godoc
// @Summary Create a late fee rule
// @Description Creates a new late fee rule for an establishment (or globally).
// @Tags LateFeeRules
// @Accept json
// @Produce json
// @Param lateFeeRule body request.CreateLateFeeRuleRequest true "Late fee rule details"
// @Success 201 {object} response.LateFeeRuleResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /late-fee-rules [post]
func (c *LateFeeRuleController) CreateLateFeeRule(ctx *gin.Context) {
	var req request.CreateLateFeeRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	lateFeeRule, err := c.lateFeeRuleService.CreateLateFeeRule(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, lateFeeRule)
}

// GetLateFeeRuleByID godoc
// @Summary Get late fee rule by ID
// @Description Retrieves a late fee rule by its ID.
// @Tags LateFeeRules
// @Produce json
// @Param id path int true "Late Fee Rule ID"
// @Success 200 {object} response.LateFeeRuleResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /late-fee-rules/{id} [get]
func (c *LateFeeRuleController) GetLateFeeRuleByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid late fee rule ID"})
		return
	}

	lateFeeRule, err := c.lateFeeRuleService.GetLateFeeRuleByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Late fee rule not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lateFeeRule)
}

// UpdateLateFeeRule godoc
// @Summary Update a late fee rule
// @Description Updates an existing late fee rule by its ID.
// @Tags LateFeeRules
// @Accept json
// @Produce json
// @Param id path int true "Late Fee Rule ID"
// @Param lateFeeRule body request.UpdateLateFeeRuleRequest true "Updated late fee rule details"
// @Success 200 {object} response.LateFeeRuleResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /late-fee-rules/{id} [put]
func (c *LateFeeRuleController) UpdateLateFeeRule(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid late fee rule ID"})
		return
	}

	var req request.UpdateLateFeeRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	lateFeeRule, err := c.lateFeeRuleService.UpdateLateFeeRule(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Late fee rule not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lateFeeRule)
}

// DeleteLateFeeRule godoc
// @Summary Delete a late fee rule
// @Description Deletes a late fee rule by its ID.
// @Tags LateFeeRules
// @Produce json
// @Param id path int true "Late Fee Rule ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /late-fee-rules/{id} [delete]
func (c *LateFeeRuleController) DeleteLateFeeRule(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid late fee rule ID"})
		return
	}

	if err := c.lateFeeRuleService.DeleteLateFeeRule(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Late fee rule not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetAllLateFeeRules godoc
// @Summary Get all late fee rules
// @Description Retrieves all late fee rules.
// @Tags LateFeeRules
// @Produce json
// @Success 200 {array} response.LateFeeRuleResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /late-fee-rules [get]
func (c *LateFeeRuleController) GetAllLateFeeRules(ctx *gin.Context) {
	lateFeeRules, err := c.lateFeeRuleService.GetAllLateFeeRules()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lateFeeRules)
}

// GetLateFeeRulesByEstablishmentID godoc
// @Summary Get late fee rules by establishment ID
// @Description Retrieves all late fee rules for a specific establishment.
// @Tags LateFeeRules
// @Produce json
// @Param establishmentID path int true "Establishment ID"
// @Success 200 {array} response.LateFeeRuleResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /establishments/{establishmentID}/late-fee-rules [get]
func (c *LateFeeRuleController) GetLateFeeRulesByEstablishmentID(ctx *gin.Context) {
	establishmentID, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	lateFeeRules, err := c.lateFeeRuleService.GetLateFeeRulesByEstablishmentID(uint(establishmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lateFeeRules)
}
