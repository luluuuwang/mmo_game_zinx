package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "mmo_game_zinx/docs"
)

type ServerStatus struct {
	Version     string `json:"version"`
	StartAt     int64  `json:"startAt"`
	OnlineCount int    `json:"onlineCount"`
	SceneCount  int    `json:"sceneCount"`
}

type PlayerBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BroadcastReq struct{ Msg string `json:"msg"` }
type BroadcastAck struct{ Delivered int `json:"delivered"` }

// @title MMO Game Admin API
// @version 1.0
// @description 管理/观测 MMO Game（Zinx）服务器的 HTTP API
// @BasePath /api
func main() {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	api := r.Group("/api")
	{
		// @Summary 服务器状态
		// @Description 返回版本、启动时间、在线人数、场景/网格数量
		// @Tags status
		// @Produce json
		// @Success 200 {object} ServerStatus
		// @Failure 503 {object} map[string]string "GameOffline"
		// @Router /status [get]
		api.GET("/status", getStatus)

		// @Summary 在线玩家
		// @Tags players
		// @Produce json
		// @Param limit query int false "按数量限制" default(50)
		// @Param name  query string false "按昵称(占位)过滤"
		// @Success 200 {array} PlayerBrief
		// @Failure 503 {object} map[string]string "GameOffline"
		// @Router /players [get]
		api.GET("/players", listPlayers)

		// @Summary 世界广播（演示）
		// @Tags gm
		// @Accept json
		// @Produce json
		// @Param body body BroadcastReq true "消息体"
		// @Success 200 {object} BroadcastAck
		// @Router /gm/broadcast [post]
		api.POST("/gm/broadcast", broadcastWorld)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	_ = r.Run(":18080")
}

func getStatus(c *gin.Context) {
	resp, err := http.Get("http://127.0.0.1:19090/status")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": "GameOffline", "message": "game server not running"})
		return
	}
	defer resp.Body.Close()
	var st ServerStatus
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": "BadGateway", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}

func listPlayers(c *gin.Context) {
	url := "http://127.0.0.1:19090/players?" + c.Request.URL.RawQuery
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": "GameOffline", "message": "game server not running"})
		return
	}
	defer resp.Body.Close()
	var list []PlayerBrief
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": "BadGateway", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func broadcastWorld(c *gin.Context) {
	// 这里只做演示；你要真广播，可在游戏服侧落一个 /gm/broadcast 内部接口或走 TCP/消息队列
	var r BroadcastReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BadRequest", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, BroadcastAck{Delivered: 1})
}
