package tower

import (
	"aoi/manager"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewAOIManager(t *testing.T) {
	players := make(map[int]*manager.Player)

	rand.Seed(time.Now().Unix())

	aoi := NewAOIManager(players)

	for i := 1; i < 20; i++ {
		px := rand.Intn(mapMaxX)
		py := rand.Intn(mapMaxY)
		fmt.Printf("%d:(%d, %d)\n", aoi.getTowerIds(px, py), px, py)
	}
}
