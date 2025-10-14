package repositories

import (
	"fmt"

	"objectswaterfall.com/core/models"
	"objectswaterfall.com/data"
)

const SQ_LITE = "sqlite3"

type Repository[T any] interface {
	SetData(workerName string, data T) error
	SetChankData(workerName string, jData []T) error
	GetData(workerName string, isRandom bool, take int, skip int64) ([]T, error)
	Count(workerName string) (int64, error)
	Exists(workerName string) (bool, error)
}

type SqLiteRepository[T any] interface {
	Repository[T]
	GetAllWorkers() (*[]models.WorkerShort, error)
	AddSettings(settings models.BackgroundWorkerSettings) error
	GetWorkerSettings(settingsWorkerName string) (*models.BackgroundWorkerSettings, error)
	GetWorkerName(id int) (string, error)
	AddWorkerResult(log models.WorkerJobLogModel) error
	GetWorkerResults(workerName string) (*[]models.WorkerJobLogModel, error)
}

func NewRepository[T any]() (SqLiteRepository[T], error) {
	switch data.DbContext.Driver {
	case SQ_LITE:
		return mySqlRepositiry[T]{}, nil
	default:
		return nil, fmt.Errorf("there is no repository for %s driver", data.DbContext.Driver)
	}
}
