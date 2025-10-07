package main

import (
	"demo/apis"
	"demo/core"
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// 当前客户端建立连接之后的Hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个Player对象
	player := core.NewPlayer(conn)

	// 给客户端发送MsgID:1的消息，同步当前Player的ID给客户端
	player.SyncPid()

	// 给客户端发送MsgID:200的消息，同步当前Player的初始位置给客户端
	player.BroadCastStartPosition()

	// 将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	// 将该连接绑定一个Pid，记录当前连接属于哪个玩家
	conn.SetProperty("pid", player.Pid)

	// 告知周边玩家当前玩家已上线，广播当前玩家位置信息
	player.SyncSurrounding()

	fmt.Println("====> Player pid = ", player.Pid, " is arrived <====")

}

// 当前客户端断开连接之后的Hook函数
func OnConnectionLost(conn ziface.IConnection) {
	// 通过conn拿到pid
	pid, _ := conn.GetProperty("pid")
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 玩家下线业务
	player.Offline()

	fmt.Println("======> Player pid = ", pid, " offline ... <======")

}


func main() {
	// 创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	// 连接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// 注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})  // 聊天的接口
	s.AddRouter(3, &apis.MoveApi{})  // 向周围广播当前玩家位置变化


	// 启动服务
	s.Serve()
}