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

func CancelHandler(c *gin.Context) {
	log.Println("request api/cancel...")
	var req model.RequestCancel
	if e := c.ShouldBindJSON(&req); e != nil {
		err := errcode.ParamError
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}
	instance, err := global.Discovery.Registry.Cancel(req.Env, req.AppId, req.Hostname, req.LatestTimestamp)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}
	//replication to other server
	if !req.Replication {
		global.Discovery.Nodes.Load().(*model.Nodes).Replicate(configs.Cancel, instance)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    configs.StatusOK,
		"message": "",
	})
}
