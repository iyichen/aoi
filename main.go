package main

import (
	"aoi/manager"
	"aoi/manager/grid"
	"aoi/manager/list"
	"aoi/manager/tower"
	"aoi/manager/whole"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

var players map[int]*manager.Player

func init() {
	data, err := ioutil.ReadFile("data/players.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &players); err != nil {
		panic(err)
	}

	for _, p := range players {
		p.Players = make(map[int]*manager.Player)
	}

	fmt.Println("玩家列表初始化完成...")

	rand.Seed(time.Now().Unix())
}

func main() {

	var strategy manager.AOIStrategyManager
	for ; ; {
		fmt.Printf("1-全局匹配\n2-九宫格\n3-灯塔\n4-十字链表\n")
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			switch scanner.Text() {
			case "1":
				fmt.Println("策略初始化. 当前策略：全局匹配")
				strategy = whole.NewAOIManager(players)
				strategy.Info()
				fmt.Println("策略初始化完成. 当前策略：全局匹配")
				break
			case "2":
				fmt.Println("策略初始化. 当前策略：九宫格")
				strategy = grid.NewAOIManager(players)
				strategy.Info()
				fmt.Println("策略初始化完成. 当前策略：九宫格")
				break
			case "3":
				fmt.Println("策略初始化. 当前策略：灯塔")
				strategy = tower.NewAOIManager(players)
				strategy.Info()
				fmt.Println("策略初始化完成. 当前策略：灯塔")
				break
			case "4":
				fmt.Println("策略初始化. 当前策略：十字链表")
				strategy = list.NewAOIManager(players)
				strategy.Info()
				fmt.Println("策略初始化完成. 当前策略：十字链表")
				break
			default:
				fmt.Println("策略不支持")
				continue
			}
		}

		pid := 100
		px := rand.Intn(25)
		py := rand.Intn(25)

		px = 13
		py = 9
		player := manager.NewPlayer(pid, px, py)
		fmt.Printf("[玩家%d]进入游戏，当前坐标(%d, %d).\n", pid, px, py)
		strategy.Enter(player)

		for scanner.Scan() {
			exit := false
			switch scanner.Text() {
			case "r":
				if px >= 25 {
					fmt.Printf("[玩家%d]无法移动，已到达边界，当前坐标(%d, %d).\n", pid, px, py)
					break
				}
				px = player.X + 1
				fmt.Printf("[玩家%d]向右移动，从坐标(%d, %d)到坐标(%d, %d).\n", pid, player.X, player.Y, px, py)
				strategy.Move(player, px, py)
				break
			case "l":
				if px <= 0 {
					fmt.Printf("[玩家%d]无法移动，已到达边界，当前坐标(%d, %d).\n", pid, px, py)
					break
				}
				px = player.X - 1
				fmt.Printf("[玩家%d]向左移动，从坐标(%d, %d)到坐标(%d, %d).\n", pid, player.X, player.Y, px, py)
				strategy.Move(player, px, py)
				break
			case "d":
				if py <= 0 {
					fmt.Printf("[玩家%d]无法移动，已到达边界，当前坐标(%d, %d).\n", pid, px, py)
					break
				}
				py = player.Y - 1
				fmt.Printf("[玩家%d]向下移动，从坐标(%d, %d)到坐标(%d, %d).\n", pid, player.X, player.Y, px, py)
				strategy.Move(player, px, py)
				break
			case "u":
				if py >= 25 {
					fmt.Printf("[玩家%d]无法移动，已到达边界，当前坐标(%d, %d).\n", pid, px, py)
					break
				}
				py = player.Y + 1
				fmt.Printf("[玩家%d]向上移动，从坐标(%d, %d)到坐标(%d, %d).\n", pid, player.X, player.Y, px, py)
				strategy.Move(player, px, py)
				break
			case "info":
				strategy.Info()
				break
			case "exit":
				exit = true
				fmt.Printf("[玩家%d]离开游戏，当前坐标(%d, %d).\n", pid, px, py)
				strategy.Leave(player)
				strategy.Info()
				break
			default:
				fmt.Printf("请输入命令: u | d | l | r | g | exit.\n")
			}
			if exit {
				break
			}
		}

	}
}
