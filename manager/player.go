package manager

import "fmt"

type Player struct {
	Id int `json:"id"`
	X  int `json:"x"`
	Y  int `json:"y"`

	Players map[int]*Player
}

func NewPlayer(id int, x, y int) *Player {
	return &Player{
		Id:      id,
		X:       x,
		Y:       y,
		Players: make(map[int]*Player),
	}
}

func (p *Player) ReceiveEnterMessage(player *Player) {
	fmt.Printf("[玩家%d](%d, %d) 收到消息：[玩家%d](%d, %d) 进入视野.\n",
		p.Id, p.X, p.Y, player.Id, player.X, player.Y)
}

func (p *Player) ReceiveMoveMessage(player *Player) {
	fmt.Printf("[玩家%d](%d, %d) 收到消息: [玩家%d]移动到(%d, %d).\n",
		p.Id, p.X, p.Y, player.Id, player.X, player.Y)
}

func (p *Player) ReceiveLeaveMessage(player *Player) {
	fmt.Printf("[玩家%d](%d, %d) 收到消息：[玩家%d](%d, %d) 离开视野.\n",
		p.Id, p.X, p.Y, player.Id, player.X, player.Y)
}
