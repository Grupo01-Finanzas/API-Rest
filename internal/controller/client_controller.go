package controller

import (
	"ApiRestFinance/internal/model/entities"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/service"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientController struct {
	clientService service.ClientService
}

func NewClientController(clientService service.ClientService) *ClientController {
	return &ClientController{clientService: clientService}
}

// CreateClient godoc
// @Summary      Create a new client
// @Description  Creates a new client with the provided data.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        client  body      request.CreateClientRequest  true  "Client data"
// @Success      201     {object}  response.ClientResponse
// @Failure      400     {object}  map[string]string  "Invalid request"
// @Failure      500     {object}  map[string]string  "Internal server error"
// @Router       /clients [post]
func (c *ClientController) CreateClient(ctx *gin.Context) {
	var req request.CreateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &entities.Client{
		UserID:      req.UserID,
		Phone:       req.Phone,
		Email:       req.Email,
		CreditLimit: req.CreditLimit,
		IsActive:    req.IsActive,
	}

	if err := c.clientService.CreateClient(client); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// You might want to return the created client in the response
	// Get the created client by ID or UserID
	// createdClient, err := c.clientService.GetClientByUserID(req.UserID) // Or GetClientByID
	// ... handle error ...

	ctx.JSON(http.StatusCreated, gin.H{"message": "Client created successfully"})
}

// GetAllClients godoc
// @Summary      Get all clients
// @Description  Gets a list of all clients.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Success      200     {array}   response.ClientResponse
// @Failure      500     {object}  map[string]string  "Internal server error"
// @Router       /clients [get]
func (c *ClientController) GetAllClients(ctx *gin.Context) {
	clients, err := c.clientService.GetAllClients()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []response.ClientResponse
	for _, client := range clients {
		resp = append(resp, response.ClientResponse{
			ID:          client.ID,          // Access ID from entities.Client
			UserID:      client.UserID,      // Access UserID from entities.Client
			Phone:       client.Phone,       // Access Phone from entities.Client
			Email:       client.Email,       // Access Email from entities.Client
			CreditLimit: client.CreditLimit, // Access CreditLimit from entities.Client
			IsActive:    client.IsActive,    // Access IsActive from entities.Client
			CreatedAt:   client.CreatedAt,   // Access CreatedAt from entities.Client
			UpdatedAt:   client.UpdatedAt,   // Access UpdatedAt from entities.Client
		})
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetClientByID godoc
// @Summary      Get client by ID
// @Description  Gets a client by its ID.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  response.ClientResponse
// @Failure      400  {object}  map[string]string  "Invalid client ID"
// @Failure      404  {object}  map[string]string  "Client not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /clients/{id} [get]
func (c *ClientController) GetClientByID(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	client, err := c.clientService.GetClientByID(uint(clientID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	resp := response.ClientResponse{
		ID:          client.ID,
		UserID:      client.UserID,
		Phone:       client.Phone,
		Email:       client.Email,
		CreditLimit: client.CreditLimit,
		IsActive:    client.IsActive,
		CreatedAt:   client.CreatedAt,
		UpdatedAt:   client.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateClient godoc
// @Summary      Update client
// @Description  Updates a client's data.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        id     path      int                      true  "Client ID"
// @Param        client  body      request.UpdateClientRequest  true  "Updated client data"
// @Success      200     {object}  map[string]string  "Client updated successfully"
// @Failure      400     {object}  map[string]string  "Invalid client ID or request body"
// @Failure      404     {object}  map[string]string  "Client not found"
// @Failure      500     {object}  map[string]string  "Internal server error"
// @Router       /clients/{id} [put]
func (c *ClientController) UpdateClient(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var req request.UpdateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the existing client
	client, err := c.clientService.GetClientByID(uint(clientID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	// Update the client fields
	client.Phone = req.Phone
	client.Email = req.Email
	client.CreditLimit = req.CreditLimit
	client.IsActive = req.IsActive

	// Save the updated client
	if err := c.clientService.UpdateClient(client); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Client updated successfully"})
}

// DeleteClient godoc
// @Summary      Delete client
// @Description  Deletes a client by its ID.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  map[string]string  "Client deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid client ID"
// @Failure      404  {object}  map[string]string  "Client not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /clients/{id} [delete]
func (c *ClientController) DeleteClient(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	if err := c.clientService.DeleteClient(uint(clientID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}
