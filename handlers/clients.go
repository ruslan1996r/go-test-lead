package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"leads/storage"
)

type ClientsHandlers struct {
	*BasicHandler
	basePath string
	storage  *storage.Storage
}

func NewClientsHandlers(storage *storage.Storage) *ClientsHandlers {
	return &ClientsHandlers{
		storage: storage,
	}
}

func (h *ClientsHandlers) InstallRoutes(r gin.IRouter) {
	c := r.Group("/clients")

	c.POST("/", h.CreateClient)
	c.GET("/", h.GetClients)
	c.GET("/:id", h.GetClient)
	c.POST("/assign", h.AssignLead)
}

// CreateClient creates a new client
//
// @Summary Creates a new client
// @Param _ body storage.ClientRequest true "New client payload"
// @Tags client
// @Produce json
// @Failure	500	{object} ErrorResponse
// @Success 200 {bool} true
// @Router /clients [post]
func (h *ClientsHandlers) CreateClient(c *gin.Context) {
	var body storage.ClientRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.sendInternalServerError(c, err)
		return
	}

	err := h.storage.CreateClient(c, body)
	if err != nil {
		h.sendInternalServerError(c, err)
		return
	}

	h.sendOk(c, true)
}

// GetClients receives a list of clients
//
// @Summary Receives a list of clients
// @Tags client
// @Produce json
// @Failure	500	{object} ErrorResponse
// @Success 200 {object} []storage.Client
// @Router /clients [get]
func (h *ClientsHandlers) GetClients(c *gin.Context) {
	clients, err := h.storage.GetClients(c, nil)
	if err != nil {
		h.sendInternalServerError(c, err)
		return
	}

	h.sendOk(c, clients)
}

// GetClient get client by clientID
//
// @Summary Get client by clientID
// @Description Returns a single client array with the found user or a string with an error in case user is not found
// @Param id path string true "Client ID"
// @Tags client
// @Produce json
// @Failure	500	{object} ErrorResponse
// @Failure	404	{object} string
// @Success 200 {object} []storage.Client
// @Router /clients/{id} [get]
func (h *ClientsHandlers) GetClient(c *gin.Context) {
	clientID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.sendInternalServerError(c, err)
		return
	}

	client, err := h.storage.GetClients(c, &clientID)
	if err != nil {
		h.sendInternalServerError(c, err)
		return
	}

	if client == nil {
		h.notFound(c, fmt.Sprintf("client with ID '%s' was not found", c.Param("id")))
		return
	}

	h.sendOk(c, client)
}

// AssignLead assigns a Lead to a suitable client
//
// @Summary Assigns a Lead to a suitable client
// @Param _ body storage.AssignLeadRequest true "Assign lead payload"
// @Description Selects a suitable client for assignment. Assigns a Lead to him and returns ID of this client.
// @Description Initially sort users by their availability and suitable time frames.
// @Description Then sort users by their priority and percentage of free capacity. Select the user with the highest indicator.
// @Tags client
// @Produce json
// @Failure	500	{object} ErrorResponse
// @Success 200 {object} storage.Lead
// @Router /clients/assign [post]
func (h *ClientsHandlers) AssignLead(c *gin.Context) {
	var lead storage.AssignLeadRequest
	if err := c.ShouldBindJSON(&lead); err != nil {
		h.sendInternalServerError(c, err)
		return
	}

	createdLead, err := h.storage.AssignLead(c, lead)
	if err != nil {
		h.sendInternalServerError(c, err)
		return
	}

	h.sendOk(c, createdLead)
}
