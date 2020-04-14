package whole

import (
	"aoi/manager"
	"aoi/util"
	"fmt"
	"math"
)

// 地图大小25x25：x坐标范围[0,25]，y坐标范围[0,25]
// 玩家视野范围10x10
var visibleArea = 5

type AOIManager struct {
	Players map[int]*manager.Player
}

func NewAOIManager(ps map[int]*manager.Player) *AOIManager {
	return &AOIManager{Players: ps}
}

func (m *AOIManager) Enter(player *manager.Player) {
	// 玩家进入地图
	m.Players[player.Id] = player

	// 发送消息
	for _, p := range m.getPlayers(player) {
		p.ReceiveEnterMessage(player)
	}
}

func (m *AOIManager) Move(player *manager.Player, tx, ty int) {

	// 移动前视野范围内玩家列表
	fps := m.getPlayers(player)

	player.X = tx
	player.Y = ty
	// 移动后视野范围内玩家列表
	tps := m.getPlayers(player)

	// 获取差集，玩家离开视野
	for _, p := range util.PlayerExcept(fps, tps) {
		p.ReceiveLeaveMessage(player)
	}

	// 获取交集，玩家移动
	for _, p := range util.PlayerIntersect(fps, tps) {
		p.ReceiveMoveMessage(player)
	}

	// 获取差集，玩家进入视野
	for _, p := range util.PlayerExcept(tps, fps) {
		p.ReceiveEnterMessage(player)
	}
}

func (m *AOIManager) Leave(player *manager.Player) {
	// 从地图上删除玩家
	delete(m.Players, player.Id)

	// 发送离开消息
	for _, p := range m.getPlayers(player) {
		p.ReceiveLeaveMessage(player)
	}
}

func (m *AOIManager) Info() {
	fmt.Println("全局玩家:")
	for _, p := range m.Players {
		fmt.Printf("[玩家%d](%d, %d). ", p.Id, p.X, p.Y)
	}
	fmt.Println()
}

// 获取玩家视野范围内的所有玩家
func (m *AOIManager) getPlayers(player *manager.Player) (players []*manager.Player) {
	for _, p := range m.Players {
		if p.Id != player.Id && math.Abs(float64(p.X-player.X)) <= float64(visibleArea) && math.Abs(float64(p.Y-player.Y)) <= float64(visibleArea) {
			players = append(players, p)
		}
	}
	return
}
