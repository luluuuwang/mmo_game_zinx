package core

import "fmt"


// 定义一些AOI边界值的宏
const (
	AOI_MIN_X int = 85
	AOI_MAX_X int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y int = 75
	AOI_MAX_Y int = 400
	AOI_CNTS_Y int = 20
)


// AOI区域管理模块

type AOIManager struct {
	MinX int // 区域的左边界坐标

	MaxX int // 区域的右边界坐标

	CntsX int // X方向格子的数量

	MinY int // 区域的上边界坐标

	MaxY int // 区域的下边界坐标

	CntsY int // Y方向格子的数量

	grids map[int]*Grid // 当前区域内有哪些格子 格子ID|格子对象
}

// 初始化一个AOI区域管理模块

func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX: minX,
		MaxX: maxX,
		CntsX: cntsX,
		MinY: minY,
		MaxY: maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	// 给AOI初始化区域内的所有格子进行编号和初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 格子编号
			gid := y*cntsX + x

			aoiMgr.grids[gid] = NewGrid(gid, 
				aoiMgr.MinX + x * aoiMgr.gridWidth(),
				aoiMgr.MinX + (x+1) * aoiMgr.gridWidth(),
				aoiMgr.MinY + y * aoiMgr.gridLength(),
				aoiMgr.MinY + (y+1) * aoiMgr.gridLength(),	
			)
		}
	}


	return aoiMgr
}


// 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在Y轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}


// 打印格子信息
func (m *AOIManager) String() string {
	// 打印AOIManager信息
	s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, CntsX:%d, MinY:%d, MaxY:%d, CntsY:%d\n Grids in AOIManager:\n", 
	 m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)

	// 打印全部格子信息
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}


// 根据格子GID得到周边九宫格格子的ID集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 判断gID是否在AOIManager中
	if _, ok := m.grids[gID]; !ok {
		return 
	}

	// 初始化grids返回值切片, 将当前gID放进来
	grids = append(grids, m.grids[gID])

	// 判断gID左右是否有格子 idx = id % nx,  idy = id /nx， 如果有，放在gidsX集合中
	idx := gID % m.CntsX
	
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}
	
	gidsX := make([]int, 0, len(grids))

	for  _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	// 遍历gidsX中每个gid，判断上下是否有格子
	for _, id := range gidsX {
		idy := id / m.CntsX
		if idy > 0 {
			grids = append(grids, m.grids[id - m.CntsX])
		}
		if idy < m.CntsY - 1{
			grids = append(grids, m.grids[id + m.CntsX])
		}
	
	}
	
	return grids
}


// 根据玩家坐标得到GID格子编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntsY + idx
}


// 通过玩家横纵坐标的到周边九宫格内全部的PlayerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家的GID格子id
	gID := m.GetGidByPos(x, y)

	// 通过GID得到周边九宫格信息
	grids := m.GetSurroundGridsByGid(gID)

	// 将九宫格内的全部player的id累加到playerIDs中
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
		// fmt.Printf("===> grid ID : %d, pids : %v ====\n", grid.GID, grid.GetPlayerIDs())
	}

	return 

}



// 添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// 移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// 通过GID获取全部的PlayerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return playerIDs
}

// 通过坐标将Player添加到一个格子中
func (m *AOIManager) AddPidToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

// 通过坐标把一个Player从一个格子中删除
func (m *AOIManager) RemovePidFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Remove(pID)

}