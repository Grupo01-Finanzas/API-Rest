package controller

import (
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

// GetAllClients godoc
// @Summary      Get all clients
// @Description  Gets a list of all clients.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200     {array}   response.ClientResponse
// @Failure      500     {object}  response.ErrorResponse
// @Router       /clients [get]
func (c *ClientController) GetAllClients(ctx *gin.Context) {
	clients, err := c.clientService.GetAllClients()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []response.ClientResponse
	for _, client := range clients {
		userResponse := c.clientService.NewUserResponse(client.User)
		resp = append(resp, response.ClientResponse{
			ID:        client.ID,
			User:      userResponse,
			IsActive:  client.IsActive,  // Access IsActive from entities.Client
			CreatedAt: client.CreatedAt, // Access CreatedAt from entities.Client
			UpdatedAt: client.UpdatedAt, // Access UpdatedAt from entities.Client
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
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  response.ClientResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
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
	userResponse := c.clientService.NewUserResponse(client.User)

	resp := response.ClientResponse{
		ID:        client.ID,
		User:      userResponse,
		IsActive:  client.IsActive,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateClient godoc
// @Summary      Update client
// @Description  Updates a client's data.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id     path      int                      true  "Client ID"
// @Param        client  body      request.UpdateClientRequest  true  "Updated client data"
// @Success      200     {object}  map[string]string  
// @Failure      400     {object}  response.ErrorResponse
// @Failure      404     {object}  response.ErrorResponse
// @Failure      500     {object}  response.ErrorResponse
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
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
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
