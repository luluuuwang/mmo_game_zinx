package core

import (
	"fmt"
	"sync"
)

// AOI地图中的一个格子类型

type Grid struct {
	GID int // 格子ID

	MinX int // 格子的左边边界坐标

	MaxX int // 格子的右边边界坐标

	MinY int // 格子的上边边界坐标

	MaxY int // 格子的下边边界坐标

	playerIDs map[int]bool // 当前格子内玩家或者物体成员的ID集合

	pIDLock sync.RWMutex // 保护当前集合的锁

}


// 初始化当前的格子的方法
func NewGrid(gID,minX,maxX,minY,maxY int) *Grid {
	return &Grid{
		GID: gID,
		MinX: minX,
		MaxX: maxX,
		MinY: minY,
		MaxY: maxY,
		playerIDs: make(map[int]bool),
	}
}

// 给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// 从格子中删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// 得到当前格子中所有的玩家ID
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

// 调试使用-打印出 格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id : %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
						g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}

 