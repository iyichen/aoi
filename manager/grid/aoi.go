package grid

import (
    "aoi/manager"
    "aoi/util"
    "fmt"
)

// 地图大小25x25：x坐标范围[0,25]，y坐标范围[0,25]
// 格子大小5x5，x方向5个格子 y方向5个格子 格子高度5 宽度5
// 格子ID为[0,24]
var (
    gridCountX = 5
    gridCountY = 5
    
    mapMinX = 0
    mapMaxX = 25
    
    mapMinY = 0
    mapMaxY = 25
    
    //gridCountX = 10
    //gridCountY = 7
    //
    //mapMinX = 0
    //mapMaxX = 50
    //
    //mapMinY = 0
    //mapMaxY = 70
    
    gridWidth  = (mapMaxX - mapMinX) / gridCountX
    gridHeight = (mapMaxY - mapMinY) / gridCountY
)

type AOIManager struct {
    Grids map[int]*Grid
}

func NewAOIManager(ps map[int]*manager.Player) *AOIManager {
    m := &AOIManager{
        Grids: make(map[int]*Grid),
    }
    
    // 生成格子ID
    for y := 0; y < gridCountY; y++ {
        for x := 0; x < gridCountX; x++ {
            gid := y*gridCountX + x
            
            m.Grids[gid] = NewGrid(gid, mapMinX+x*gridWidth, mapMinX+(x+1)*gridWidth, mapMinY+y*gridHeight, mapMinY+(y+1)*gridHeight)
        }
    }
    
    for _, p := range ps {
        gid := m.getGid(p.X, p.Y)
        m.Grids[gid].Players[p.Id] = p
    }
    
    return m
}

func (m *AOIManager) Enter(player *manager.Player) {
    m.Grids[m.getGid(player.X, player.Y)].Players[player.Id] = player
    
    for _, p := range m.getPlayers(player) {
        p.ReceiveEnterMessage(player)
    }
    
}

func (m *AOIManager) Move(player *manager.Player, tx, ty int) {
    // 获取移动前的所属格子
    fgid := m.getGid(player.X, player.Y)
    // 获取移动后的所属格子
    tgid := m.getGid(tx, ty)
    
    // 获取移动前的格子列表
    fps := m.getPlayers(player)
    
    player.X = tx
    player.Y = ty
    // 获取移动后的格子列表
    tps := m.getPlayers(player)
    
    // 前后所属格子不同，删除后添加
    if fgid != tgid {
        delete(m.Grids[fgid].Players, player.Id)
        
        m.Grids[tgid].Players[player.Id] = player
    }
    
    // 此处获取差集交集可以优化，此处是对玩家列表处理，可以调整成对格子列表处理
    // 差集，发送离开视野消息
    for _, p := range util.PlayerExcept(fps, tps) {
        p.ReceiveLeaveMessage(player)
    }
    
    // 交集，发送移动消息
    for _, p := range util.PlayerIntersect(fps, tps) {
        p.ReceiveMoveMessage(player)
    }
    
    // 差集，发送进入游戏消息
    for _, p := range util.PlayerExcept(tps, fps) {
        p.ReceiveEnterMessage(player)
    }
}

func (m *AOIManager) Leave(player *manager.Player) {
    delete(m.Grids[m.getGid(player.X, player.Y)].Players, player.Id)
    
    for _, p := range m.getPlayers(player) {
        p.ReceiveLeaveMessage(player)
    }
}

func (m *AOIManager) Info() {
    for _, g := range m.Grids {
        fmt.Printf("格子%d:", g.Id)
        for _, p := range g.Players {
            fmt.Printf("[玩家%d](%d,%d),", p.Id, p.X, p.Y)
        }
        fmt.Println()
    }
}

// 获取九宫格内的所有玩家
// 获取顺序为玩家所属格子(8)-7-2-12-9-4-14-3-13
func (m *AOIManager) getPlayers(player *manager.Player) (players []*manager.Player) {
    
    gid := m.getGid(player.X, player.Y)
    
    var grids []*Grid
    grids = append(grids, m.Grids[gid])
    
    // 获取x方向格子索引
    idx := gid % gridCountX
    // 获取y方向格子索引
    idy := gid / gridCountX
    
    // 左边界判断
    // 如果格子左边有格子，判断该格子上方和下方是否有格子
    if idx-1 >= 0 {
        grids = append(grids, m.Grids[gid-1])
        
        if idy-1 >= 0 {
            grids = append(grids, m.Grids[gid-1-gridCountX])
        }
        
        if idy+1 < gridCountY {
            grids = append(grids, m.Grids[gid-1+gridCountX])
        }
    }
    
    // 右边界判断
    // 如果格子右边有格子，判断该格子上方和下方是否有格子
    if idx+1 < gridCountX {
        grids = append(grids, m.Grids[gid+1])
        
        if idy-1 >= 0 {
            grids = append(grids, m.Grids[gid+1-gridCountX])
        }
        
        if idy+1 < gridCountY {
            grids = append(grids, m.Grids[gid+1+gridCountX])
        }
    }
    
    // 下边界判断
    if idy-1 >= 0 {
        grids = append(grids, m.Grids[gid-gridCountX])
    }
    
    // 上边界判断
    if idy+1 < gridCountY {
        grids = append(grids, m.Grids[gid+gridCountX])
    }
    
    // 获取格子内的所有玩家
    for _, g := range grids {
        for _, v := range g.Players {
            if v.Id != player.Id {
                players = append(players, v)
            }
        }
    }
    return
}

// 根据坐标获取格子ID
// 格子边界上的坐标只会在一个格子内
// 计算方式：(y-地图最小Y坐标)/格子高度*X方向格子数量+(x-地图最小X坐标)/格子宽度
func (m *AOIManager) getGid(x, y int) int {
    // 如果坐标处于地图边界 特殊处理
    // 可以定义玩家可移动区域，在这种情况下，玩家就无法到达边界
    if x == mapMaxX {
        x = x - 1
    }
    if y == mapMaxY {
        y = y - 1
    }
    return (y-mapMinY)/gridHeight*gridCountX + (x-mapMinY)/gridWidth
}

// 格子信息
type Grid struct {
    Id      int
    MinX    int
    MaxX    int
    MinY    int
    MaxY    int
    Players map[int]*manager.Player
}

func NewGrid(id, minX, maxX, minY, maxY int) *Grid {
    return &Grid{
        Id:      id,
        MinX:    minX,
        MaxX:    maxX,
        MinY:    minY,
        MaxY:    maxY,
        Players: make(map[int]*manager.Player),
    }
}
