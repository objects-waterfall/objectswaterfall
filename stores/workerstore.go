package stores

import (
	"errors"
	"fmt"
	"sort"

	"objectswaterfall.com/core/models"
	"objectswaterfall.com/core/services"
)

type workerStore struct {
	workers map[int]*services.Worker
}

var store services.WorkerStore

func GetWorkerStore() services.WorkerStore {
	if store != nil {
		return store
	}
	store = &workerStore{
		workers: map[int]*services.Worker{},
	}
	return store
}

func (w *workerStore) Add(id int, worker *services.Worker) int {
	if id != 0 {
		w.workers[id] = worker
		return id
	}

	lenght := len(w.workers)
	keys := make([]int, 0, lenght)
	for k := range w.workers {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	last := keys[lenght-1]
	workerId := last + 1
	w.workers[workerId] = worker

	return workerId
}

func (w *workerStore) Get(workerId int) (*services.Worker, error) {
	if worker, ok := w.workers[workerId]; ok {
		return worker, nil
	}

	return nil, errors.New("wrong worker identifire")
}

func (w *workerStore) All() *[]models.WorkerShort {
	var runningWorkers []models.WorkerShort
	for k, v := range w.workers {
		runningWorkers = append(runningWorkers, models.WorkerShort{
			Id:   k,
			Name: (*v).GetWorkerName(),
		})
	}
	return &runningWorkers
}

func (w *workerStore) Exists(name string) bool {
	for _, v := range w.workers {
		if (*v).GetWorkerName() == name {
			return true
		}
	}

	return false
}

func (w *workerStore) CancelWork(workerId int) error {
	if _, ok := (*w).workers[workerId]; !ok {
		return fmt.Errorf("there is no worker with id %d", workerId)
	}
	(*w.workers[workerId]).Cancel()
	return nil
}

func (w *workerStore) Remove(workerId int) error {
	if _, ok := w.workers[workerId]; !ok {
		return errors.New("wrong worker identifire")
	}
	(*w.workers[workerId]).Cancel()
	delete(w.workers, workerId)
	return nil
}
