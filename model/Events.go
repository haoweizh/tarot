package model

import (
	"sync"
)

type Events struct {
	lock       sync.Mutex
	nickEvents map[string]*TarotEvent
}

func NewEvents() *Events {
	events := &Events{nickEvents: make(map[string]*TarotEvent)}
	return events
}

func (events *Events) GetUnNilAmount() (amount int) {
	events.lock.Lock()
	defer events.lock.Unlock()
	amount = 0
	if events.nickEvents == nil {
		return amount
	}
	for _, value := range events.nickEvents {
		if value != nil {
			amount++
		}
	}
	return amount
}

func (events *Events) RemoveEvent(nickName string) *TarotEvent {
	events.lock.Lock()
	defer events.lock.Unlock()
	event := events.nickEvents[nickName]
	events.nickEvents[nickName] = nil
	return event
}

func (events *Events) GetEvent(nickName string) (event *TarotEvent) {
	events.lock.Lock()
	defer events.lock.Unlock()
	return events.nickEvents[nickName]
}

func (events *Events) PutEvent(nickName string, event *TarotEvent) {
	events.lock.Lock()
	defer events.lock.Unlock()
	//if events.nickEvents[nickName] != nil {
	//	util.SocketInfo(fmt.Sprintf(`can not send %s msg, abandon %s`, event.NickName, event.SentenceType))
	//}
	events.nickEvents[nickName] = event
}
