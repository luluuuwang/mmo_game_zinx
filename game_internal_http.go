// file: game_internal_http.go
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mmo_game_zinx/core"
)

type ServerStatus struct {
	Version     string `json:"version"`
	StartAt     int64  `json:"startAt"`
	OnlineCount int    `json:"onlineCount"`
	SceneCount  int    `json:"sceneCount"`
}

type PlayerBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"` // Player 没有名字字段，先留空
}

func startInternalHTTP() {
	mux := http.NewServeMux()

	// 实时状态：127.0.0.1:19090/status
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		st := ServerStatus{
			Version:     "1.0.0",
			StartAt:     time.Now().Add(-1 * time.Hour).Unix(), // 如有真实启动时间可替换
			OnlineCount: worldOnlineCount(),
			SceneCount:  worldSceneCount(),
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(st)
	})

	// 在线玩家列表：127.0.0.1:19090/players?limit=50&name=xxx
	mux.HandleFunc("/players", func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		if s := r.URL.Query().Get("limit"); s != "" {
			if n, err := strconv.Atoi(s); err == nil && n > 0 {
				limit = n
			}
		}

		all := listWorldPlayers()
		out := make([]PlayerBrief, 0, limit)
		for _, p := range all {
			// 若未来你在 Player 里增加名字，这里可做 nameLike 模糊匹配
			out = append(out, *p)
			if len(out) >= limit {
				break
			}
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(out)
	})

	go http.ListenAndServe("127.0.0.1:19090", mux)
}

// ====== 下面 3 个函数用真实 API 实现 ======

// 在线人数（基于 WorldMgrObj.Players）
func worldOnlineCount() int {
	if core.WorldMgrObj == nil || core.WorldMgrObj.Players == nil {
		return 0
	}
	return len(core.WorldMgrObj.Players)
}

// 场景数量（你的项目没有直接概念，先给一个固定值，后面你再改）
func worldSceneCount() int {
	return 1
}

// 列出在线玩家（用 GetAllPlayer()）
func listWorldPlayers() []*PlayerBrief {
	out := make([]*PlayerBrief, 0, 64)
	if core.WorldMgrObj == nil {
		return out
	}
	for _, p := range core.WorldMgrObj.GetAllPlayer() {
		out = append(out, &PlayerBrief{
			ID:   int(p.Pid),
			Name: "", // 你的 Player 没有名字字段，先置空
		})
	}
	return out
}
