package apis

import (
	"fmt"
	"mmo_game_zinx/core"
	"mmo_game_zinx/pb"
	"zinx/ziface"
	"zinx/znet"

	"google.golang.org/protobuf/proto"
)

// 世界聊天 路由业务
type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	//1 解析客户端传递进来的proto协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Talk Unmarshal error: ", err)
		return
	}

	//2 当前的聊天数据是属于哪个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")

	//3 根据pid得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//4 将聊天数据广播给其他全部在线玩家
	player.Talk(proto_msg.Content)
}
