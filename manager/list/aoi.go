package list

import (
	"aoi/manager"
	"aoi/util"
	"container/list"
	"fmt"
	"math"
)

// 地图大小25x25：x坐标范围[0,25]，y坐标范围[0,25]
// 玩家视野范围10x10
var visibleArea = 5

type AOIManager struct {
	XList *list.List
	YList *list.List
}

func NewAOIManager(ps map[int]*manager.Player) *AOIManager {
	m := &AOIManager{
		XList: list.New(),
		YList: list.New(),
	}
	for _, p := range ps {
		m.insertList(p)
	}

	return m
}

func (m *AOIManager) Enter(player *manager.Player) {
	// 插入XY链表
	m.insertList(player)

	// 发送进入游戏消息
	for _, p := range m.getPlayers(player) {
		p.ReceiveEnterMessage(player)
	}
}

func (m *AOIManager) Move(player *manager.Player, tx, ty int) {

	fps := m.getPlayers(player)

	m.removeList(player)

	player.X = tx
	player.Y = ty
	tps := m.getPlayers(player)

	m.insertList(player)

	for _, p := range util.PlayerExcept(fps, tps) {
		p.ReceiveLeaveMessage(player)
	}

	for _, p := range util.PlayerIntersect(fps, tps) {
		p.ReceiveMoveMessage(player)
	}

	for _, p := range util.PlayerExcept(tps, fps) {
		p.ReceiveEnterMessage(player)
	}

}

func (m *AOIManager) Leave(player *manager.Player) {
	m.removeList(player)

	for _, p := range m.getPlayers(player) {
		p.ReceiveLeaveMessage(player)
	}
}

func (m *AOIManager) Info() {
	fmt.Print("X链表信息:")
	for e := m.XList.Front(); e != nil; e = e.Next() {
		ep := e.Value.(*manager.Player)
		fmt.Printf(" -> [玩家%d](%d, %d)", ep.Id, ep.X, ep.Y)
	}
	fmt.Println()
	fmt.Print("Y链表信息:")
	for e := m.YList.Front(); e != nil; e = e.Next() {
		ep := e.Value.(*manager.Player)
		fmt.Printf(" -> [玩家%d](%d, %d)", ep.Id, ep.X, ep.Y)
	}
	fmt.Println()
}

func (m *AOIManager) Get(pid int, list2 *list.List) *manager.Player {
	var player *manager.Player
	for e := list2.Front(); e != nil; e = e.Next() {
		ep := e.Value.(*manager.Player)
		if ep.Id == pid {
			player = ep
			break
		}
	}
	return player
}

// 查找玩家视野内的玩家列表
// 分别查询X，Y链表玩家，然后对两个玩家列表取交集
func (m *AOIManager) getPlayers(player *manager.Player) (players []*manager.Player) {
	// 获取X链表视野内的玩家，判断条件|x1-x|<visibleArea
	var xps []*manager.Player
	for e := m.XList.Front(); e != nil; e = e.Next() {
		ep := e.Value.(*manager.Player)
		if ep.Id != player.Id && math.Abs(float64(ep.X-player.X)) <= float64(visibleArea) {
			xps = append(xps, ep)
		}
	}

	// 获取Y链表视野内的玩家，判断条件|y1-y|<visibleArea
	var yps []*manager.Player
	for e := m.YList.Front(); e != nil; e = e.Next() {
		ep := e.Value.(*manager.Player)
		if ep.Id != player.Id && math.Abs(float64(ep.Y-player.Y)) <= float64(visibleArea) {
			yps = append(yps, ep)
		}
	}

	// 获取2个玩家列表交集
	return util.PlayerIntersect(xps, yps)
}

// 将玩家从XY链表中移除
func (m *AOIManager) removeList(p *manager.Player) {
	// 从X链表移除
	for e := m.XList.Front(); e != nil; e = e.Next() {
		ep := e.Value.(*manager.Player)
		if ep.Id == p.Id {
			m.XList.Remove(e)
			break
		}
	}

	// 从Y链表移除
	for e := m.YList.Front(); e != nil; e = e.Next() {
		ep := e.Value.(*manager.Player)
		if ep.Id == p.Id {
			m.YList.Remove(e)
			break
		}
	}
}

// 将玩家插入到XY链表
// XY链表按照从小到大排列
func (m *AOIManager) insertList(p *manager.Player) {
	// 插入X链表
	xDone := false
	for e := m.XList.Front(); e != nil; e = e.Next() {
		ep := e.Value.(*manager.Player)
		if ep.X > p.X {
			m.XList.InsertBefore(p, e)
			xDone = true
			break
		}
	}
	// 没有找到，插入X链表末端
	if !xDone {
		m.XList.PushBack(p)
	}

	// 插入Y链表
	yDone := false
	for e := m.YList.Front(); e != nil; e = e.Next() {
		ep := e.Value.(*manager.Player)
		if ep.Y > p.Y {
			m.YList.InsertBefore(p, e)
			yDone = true
			break
		}
	}
	// 没有找到，插入Y链表末端
	if !yDone {
		m.YList.PushBack(p)
	}
}
