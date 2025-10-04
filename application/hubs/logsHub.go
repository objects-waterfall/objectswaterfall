package hubs

import (
	"sync"

	"github.com/philippseith/signalr"
	"golang.org/x/sync/semaphore"
	"objectswaterfall.com/core/models"
	"objectswaterfall.com/core/services"
	"objectswaterfall.com/stores"
)

type LogsHub struct {
	signalr.Hub
	listeners   map[string]int
	workerStore services.WorkerStore
	wg          sync.WaitGroup
	sem         *semaphore.Weighted
}

func NewLogHub(gorCount int64) LogsHub {
	return LogsHub{
		listeners:   make(map[string]int),
		workerStore: stores.GetWorkerStore(),
		sem:         semaphore.NewWeighted(gorCount),
		wg:          sync.WaitGroup{},
	}
}

func (h *LogsHub) OnConnected(connectionID string) {
	h.listeners[connectionID] = -1
}

func (h *LogsHub) OnDisconnected(connectionID string) {
	delete(h.listeners, connectionID)
}

func (h *LogsHub) PushLogs(workerId int) {
	worker, err := h.workerStore.Get(workerId)
	if err != nil {
		go h.Clients().Caller().Send("/receiveError", err.Error())
	}
	(*worker).SetLogFunc(func(l models.LogModel) {
		go h.Clients().Caller().Send("/recieveLog", l)
	})
}
