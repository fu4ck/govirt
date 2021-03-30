package server

import (
	"github.com/gin-gonic/gin"
	"govirt/controller"
	"govirt/runtime"
	"net/http"
)

func StartRouter(v *runtime.GOVirt) {
	engine := gin.New()
	engine.POST("/createVirualMachine", controller.CreateVirtualMachine)

	http.ListenAndServe(v.Cfg.Address, engine)
}
