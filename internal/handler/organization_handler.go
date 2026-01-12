package handler

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	Service domain.OrganizationService
}

func NewOrganizationHandler(s domain.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{Service: s}
}

type CreateOrgRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	org, err := h.Service.CreateOrganization(req.Name, req.Description)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, org)
}

func (h *OrganizationHandler) List(c *gin.Context) {
	orgs, err := h.Service.ListOrganizations()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"list": orgs})
}

type UpdateOrgRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *OrganizationHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	org, err := h.Service.UpdateOrganization(uint(id), req.Name, req.Description)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, org)
}

func (h *OrganizationHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.Service.DeleteOrganization(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Organization deleted successfully", nil)
}
