package controller

import (
	"github.com/gin-gonic/gin"
	"govirt/service"
	"net/http"
)

type createvirtRequest struct {
	Disk       string `json:"disk"`
	ConfigPath string `json:"configpath"`
	Network    string `json:"network"`
}

func CreateVirtualMachine(context *gin.Context) {

	var req createvirtRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	err := service.CreateVirt(req.Disk, req.ConfigPath, req.Network)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, nil)
	return
}
