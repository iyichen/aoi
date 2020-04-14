package manager

type AOIStrategyManager interface {
	// 进入游戏
	Enter(player *Player)

	// 游戏内移动
	Move(player *Player, tx, ty int)

	// 离开游戏
	Leave(player *Player)

	// 当前AOI信息
	Info()
}
