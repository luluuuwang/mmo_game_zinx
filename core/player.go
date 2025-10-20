package core

import (
	"mmo_game_zinx/pb"
	"fmt"
	"math/rand"
	"sync"
	"zinx/ziface"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	Pid  int32              // 玩家ID
	Conn ziface.IConnection // 当前玩家与 客户端 的链接
	X    float32            // 平面的x坐标
	Y    float32            // 高度
	Z    float32            // 平面y坐标
	V    float32            // 旋转的0-360角度
}

var PidGen int32 = 1  // 用来生产玩家ID的计数器
var IdLock sync.Mutex // 保护PidGen的Mutex

// 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	// 创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}

	return p
}

// 提供一个发送给客户端消息的方法
// 将pb的protobuf数据序列化之后，在调用zinx的SendMsg方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	// 将proto Message结构体序列化，转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}

	// 将二进制文件通过zinx框架的SendMsg发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error!")
		return
	}

	return
}

// 告知客户端玩家Pid，同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	// 组建MsgID:1 的proto数据
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	// 将消息发送给客户端
	p.SendMsg(1, proto_msg)
}

// 广播玩家的出生地点
func (p *Player) BroadCastStartPosition() {
	// 组建 MsgID:200 的 proto 数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, // Tp2代表广播位置坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 将消息发送给客户端
	p.SendMsg(200, proto_msg)
}

// 玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	// 1.组建 MsgID:200,Tp:1 的proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, // Tp-1代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	// 2.得到当前所有的在线玩家
	players := WorldMgrObj.GetAllPlayer()

	// 3.向所有的在线玩家发送信息
	for _, player := range players {
		// player分别给对应的客户端发送消息
		player.SendMsg(200, proto_msg)
	}
}

// 向周围玩家广播当前玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	// 1.获取当前玩家周围的玩家有哪些(九宫格)
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	// 2.将当前玩家的位置信息通过MsgID:200 发送给周围的玩家(让其他玩家看到自己)
	// 组建MsgID=200的数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 全部周围的玩家都向自己的客户端发送200消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	// 3.将周围玩家的位置信息通过MsgID:202 发送给当前的玩家(让自己看到其他玩家)
	// 组建MsgID=200的数据
	// 制作pb.Player.slice
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		// 制作一个message Player
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}
	// 封装为新建的SyncPlayer protobuf数据
	SyncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}

	// 将组建好的数据发送给当前玩家的客户端
	p.SendMsg(202, SyncPlayers_proto_msg)

}

// 获取当前玩家附近九宫格之内的玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	// 得到当前玩家九宫格内的所有玩家PID
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	// 将pid对应的player放到Players切片中
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}

// 广播当前玩家的位置移动信息
func (p *Player) UpdatePos(x, y, z, v float32) {
	// 更新当前玩家的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	// 组建MsgID=200，Tp-4的数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 获取当前玩家的周边玩家AOI九宫格之内的玩家
	players := p.GetSurroundingPlayers()

	// 依次给每个玩家对应的客户端发送当前玩家位置更新的消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

}

// 玩家下线业务
func (p *Player) Offline() {
	// 得到当前玩家周边的
	players := p.GetSurroundingPlayers()

	// 给周围玩家广播MsgID:201的消息
	// 组建消息
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	// 广播
	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
