package tower

import (
    "aoi/manager"
    "aoi/util"
    "fmt"
    "math"
)

// 地图大小25x25：x坐标范围[0,25]，y坐标范围[0,25]
// 灯塔大小5x5
// 玩家视野范围10x10
var (
    gridCountX = 5
    gridCountY = 5
    
    mapMinX = 0
    mapMinY = 0
    
    mapMaxX = 25
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
    
    visibleArea = 5
)

type AOIManager struct {
    Towers map[int]*Tower
}

func NewAOIManager(ps map[int]*manager.Player) *AOIManager {
    m := &AOIManager{
        Towers: make(map[int]*Tower),
    }
    
    for y := 0; y < gridCountY; y++ {
        for x := 0; x < gridCountX; x++ {
            tid := y*gridCountX + x
            
            m.Towers[tid] = NewTower(tid, float64(gridWidth)/2+float64(mapMinX+x*gridWidth), float64(gridHeight)/2+float64(mapMinY+y*gridHeight))
        }
    }
    
    for _, p := range ps {
        m.Enter(p)
    }
    
    return m
}

func (m *AOIManager) Enter(player *manager.Player) {
    // 获取灯塔ID
    tid := m.getTowerId(player.X, player.Y)
    // 将玩家交给灯塔管理
    m.Towers[tid].Markers[player.Id] = player
    
    // 获取玩家可见范围内的灯塔列表
    tids := m.getTowerIds(player.X, player.Y)
    for _, tid := range tids {
        tower := m.Towers[tid]
        for k, v := range tower.Watchers {
            if v.Id == player.Id {
                continue
            }
            _, exist := player.Players[k]
            if !exist {
                player.Players[k] = v
                
                // 发送玩家进入消息
                v.ReceiveEnterMessage(player)
                
                v.Players[player.Id] = player
            }
        }
        tower.Watchers[player.Id] = player
    }
}

func (m *AOIManager) Move(player *manager.Player, tx, ty int) {
    
    ftids := m.getTowerIds(player.X, player.Y)
    ttids := m.getTowerIds(tx, ty)
    
    ftid := m.getTowerId(player.X, player.Y)
    ttid := m.getTowerId(tx, ty)
    
    player.X = tx
    player.Y = ty
    
    // 修改玩家所属灯塔
    if ttid != ftid {
		delete(m.Towers[ftid].Markers, player.Id)
	
		m.Towers[ttid].Markers[player.Id] = player
	}
 
	// 玩家视野范围内灯塔没变，仅做移动操作
    if util.IntEqual(ftids, ttids) {
        for _, p := range player.Players {
            p.ReceiveMoveMessage(player)
        }
        return
    }
    
    fps := m.getTowerWatchers(ftids)
    tps := m.getTowerWatchers(ttids)
    
    for _, tid := range util.IntExcept(ftids, ttids) {
        tower := m.Towers[tid]
        // 将玩家从灯塔观察者移除
        delete(tower.Watchers, player.Id)
    }
    
    for _, p := range util.PlayerExcept(fps, tps) {
		// 从玩家列表移除
		delete(p.Players, player.Id)
    	
        _, exist := player.Players[p.Id]
        if exist {
            delete(player.Players, p.Id)
            
            p.ReceiveLeaveMessage(player)
        }
    }
    
    for _, p := range player.Players {
        p.ReceiveMoveMessage(player)
    }
    
    for _, tid := range util.IntExcept(ttids, ftids) {
        tower := m.Towers[tid]
        
        // 将玩家添加到灯塔观察者
        tower.Watchers[player.Id] = player
        
        // 将灯塔观察者交给玩家维护
        for k, v := range tower.Watchers {
            if v.Id == player.Id {
                continue
            }
            _, exist := player.Players[k]
            if !exist {
                player.Players[k] = v
                
                // 发送玩家进入消息
                v.ReceiveEnterMessage(player)
                
                v.Players[player.Id] = player
            }
        }
    }
}

func (m *AOIManager) Leave(player *manager.Player) {
    // 获取灯塔ID
    tid := m.getTowerId(player.X, player.Y)
    // 将玩家从灯塔管理移除
    delete(m.Towers[tid].Markers, player.Id)
    
    // 获取玩家可见范围内的灯塔列表
    tids := m.getTowerIds(player.X, player.Y)
    for _, tid := range tids {
        
        tower := m.Towers[tid]
        // 将玩家从灯塔观察者移除
        delete(tower.Watchers, player.Id)
        
        // 从玩家列表移除
        for _, p := range player.Players {
            // 从玩家列表移除
            delete(p.Players, player.Id)
            
            p.ReceiveLeaveMessage(player)
        }
    }
}

func (m *AOIManager) Info() {
    for _, t := range m.Towers {
        fmt.Printf("灯塔%d(%f,%f):\n", t.Id, t.X, t.Y)
        fmt.Print("管理对象:\n")
        for _, p := range t.Markers {
            fmt.Printf(" [玩家%d](%d,%d):", p.Id, p.X, p.Y)
            for _, p1 := range p.Players {
                fmt.Printf("->[玩家%d](%d,%d)", p1.Id, p1.X, p1.Y)
            }
            fmt.Println()
        }
        fmt.Printf("\n观察者对象:\n")
        for _, p := range t.Watchers {
            fmt.Printf(" [玩家%d](%d,%d)\n", p.Id, p.X, p.Y)
        }
        fmt.Println()
    }
}

// 根据坐标获取格子(灯塔)ID
// 格子边界上的坐标只会在一个格子内
// 计算方式：(y-地图最小Y坐标)/格子高度*X方向格子数量+(x-地图最小X坐标)/格子宽度
func (m *AOIManager) getTowerId(x, y int) int {
    // 如果坐标处于地图边界 特殊处理
    if x == mapMaxX {
        x = x - 1
    }
    if y == mapMaxY {
        y = y - 1
    }
    return (y-mapMinY)/gridHeight*gridCountX + (x-mapMinY)/gridWidth
}

// 判断坐标是否在灯塔范围内
func (m *AOIManager) inTower(tid, x, y int) bool {
    minX := math.Max(float64(x-visibleArea), float64(mapMinX))
    maxX := math.Min(float64(x+visibleArea), float64(mapMaxX))
    minY := math.Max(float64(y-visibleArea), float64(mapMinY))
    maxY := math.Min(float64(y+visibleArea), float64(mapMaxY))
    
    return minX <= m.Towers[tid].X && m.Towers[tid].X <= maxX && minY <= m.Towers[tid].Y && m.Towers[tid].Y <= maxY
}

// 根据灯塔坐标获取灯塔上的所有观察者
func (m *AOIManager) getTowerWatchers(towerIds []int) (players []*manager.Player) {
    pmap := make(map[int]*manager.Player)
    for _, tid := range towerIds {
        tower := m.Towers[tid]
        for k, v := range tower.Watchers {
            pmap[k] = v
        }
    }
    
    for _, v := range pmap {
        players = append(players, v)
    }
    return
}

// 根据玩家坐标获取视野范围内的所有灯塔
func (m *AOIManager) getTowerIds(x, y int) (towerIds []int) {
    
    tid := m.getTowerId(x, y)
    
    towerIds = append(towerIds, tid)
    
    idx := tid % gridCountX
    idy := tid / gridCountX
    
    if idx-1 >= 0 {
        if m.inTower(tid-1, x, y) {
            towerIds = append(towerIds, tid-1)
        }
        
        if idy-1 >= 0 && m.inTower(tid-1-gridCountX, x, y) {
            towerIds = append(towerIds, tid-1-gridCountX)
        }
        
        if idy+1 < gridCountY && m.inTower(tid-1+gridCountX, x, y) {
            towerIds = append(towerIds, tid-1+gridCountX)
        }
    }
    
    if idx+1 < gridCountX {
        if m.inTower(tid+1, x, y) {
            towerIds = append(towerIds, tid+1)
        }
        
        if idy-1 >= 0 && m.inTower(tid+1-gridCountX, x, y) {
            towerIds = append(towerIds, tid+1-gridCountX)
        }
        
        if idy+1 < gridCountY && m.inTower(tid+1+gridCountX, x, y) {
            towerIds = append(towerIds, tid+1+gridCountX)
        }
    }
    
    if idy-1 >= 0 && m.inTower(tid-gridCountX, x, y) {
        towerIds = append(towerIds, tid-gridCountX)
    }
    
    if idy+1 < gridCountY && m.inTower(tid+gridCountX, x, y) {
        towerIds = append(towerIds, tid+gridCountX)
    }
    
    return
}

type Tower struct {
    Id int
    
    X float64
    
    Y float64
    
    Watchers map[int]*manager.Player
    
    Markers map[int]*manager.Player
}

func NewTower(id int, x float64, y float64) *Tower {
    return &Tower{
        Id:       id,
        X:        x,
        Y:        y,
        Watchers: make(map[int]*manager.Player),
        Markers:  make(map[int]*manager.Player),
    }
}
