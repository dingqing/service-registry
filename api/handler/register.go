package handler

import (
	"log"
	"net/http"

	"github.com/dingqing/registry/configs"
	"github.com/dingqing/registry/global"
	"github.com/dingqing/registry/model"
	"github.com/dingqing/registry/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	log.Println("request api/register...")
	var req model.RequestRegister
	if e := c.ShouldBindJSON(&req); e != nil {
		log.Println("error:", e)
		err := errcode.ParamError
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}
	//bind instance
	instance := model.NewInstance(&req)
	if instance.Status != configs.StatusReceive && instance.Status != configs.StatusNotReceive {
		log.Println("register params status invalid")
		err := errcode.ParamError
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}
	//dirtytime
	if req.DirtyTimestamp > 0 {
		instance.DirtyTimestamp = req.DirtyTimestamp
	}
	global.Discovery.Registry.Register(instance, req.LatestTimestamp)
	//default do replicate. if request come from other server, req.Replication is true, ignore replicate.
	if !req.Replication {
		global.Discovery.Nodes.Load().(*model.Nodes).Replicate(configs.Register, instance)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    configs.StatusOK,
		"message": "",
		"data":    "",
	})
}
