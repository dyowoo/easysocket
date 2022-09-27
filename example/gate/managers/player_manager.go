/**
 * @Author: Jason
 * @Description:
 * @File: player_manager.go
 * @Date: 2022/9/27 9:43
 **/

package managers

import (
	"github.com/dyowoo/easysocket"
	"sync"
)

type Player struct {
	Session easysocket.ISession
}

type PlayerManager struct {
	sync.RWMutex

	players map[uint32]*Player
}

func (pm *PlayerManager) Add(player *Player) {
	pm.Lock()
	defer pm.Unlock()

	if pm.players == nil {
		pm.players = make(map[uint32]*Player)
	}

	pm.players[player.Session.GetConnId()] = player
}

func (pm *PlayerManager) Remove(session easysocket.ISession) {
	pm.Lock()
	defer pm.Unlock()

	delete(pm.players, session.GetConnId())
}

func (pm *PlayerManager) Get(connId uint32) *Player {
	pm.RLock()
	defer pm.RUnlock()

	return pm.players[connId]
}

var PlayerMgr = PlayerManager{}
