package services

import (
	"context"

	"objectswaterfall.com/core/models"
)

type WorkerStore interface {
	Add(id int, worker *Worker) int
	Get(workerId int) (*Worker, error)
	CancelWork(id int) error
	Remove(id int) error
	Exists(name string) bool
	All() *[]models.WorkerShort
}

type Worker interface {
	DoWork(ctx context.Context)
	SetCancel(context.CancelFunc)
	Cancel()
	GetWorkerName() string
	Log() *models.LogModel
	SetLogFunc(logFunc models.LogFunc)
}
