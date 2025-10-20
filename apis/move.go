package apis

import (
	"mmo_game_zinx/core"
	"mmo_game_zinx/pb"
	"fmt"
	"zinx/ziface"
	"zinx/znet"

	"google.golang.org/protobuf/proto"
)

// 玩家移动
type MoveApi struct {
	znet.BaseRouter
}

// 向其他玩家广播当前玩家位置移动信息
func (m *MoveApi) Handle(request ziface.IRequest) {
	// 解析客户端传来的proto协议
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move : Position Unmarshal error", err)
		return 
	}

	// 得到当前发送位置的是哪个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error, ", err)
		return
	}

	fmt.Printf("Player pid = %d, move(%f,%f,%f,%f)\n", 
				pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	// 给其他玩家进行广播
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	// 广播并更新当前玩家的坐标
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}