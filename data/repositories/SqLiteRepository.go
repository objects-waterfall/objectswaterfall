package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"objectswaterfall.com/core/models"
	"objectswaterfall.com/data"
)

type mySqlRepositiry[T any] struct {
}

func (r mySqlRepositiry[T]) SetData(workerName string, jData T) error {
	if err := createTable(workerName); err != nil {
		return nil
	}

	stmt, err := data.DbContext.Db.Prepare(fmt.Sprintf(data.InsertData, workerName))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(jData)
	if err != nil {
		return err
	}

	return nil
}

func (r mySqlRepositiry[T]) SetChankData(workerName string, jData []T) error {
	if err := createTable(workerName); err != nil {
		return err
	}
	tx, err := data.DbContext.Db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(fmt.Sprintf(data.InsertData, workerName))
	if err != nil {
		return err
	}

	defer stmt.Close()
	for _, v := range jData {
		if _, err := stmt.Exec(v); err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r mySqlRepositiry[T]) GetData(workerName string, isRandom bool, take int, skip int64) ([]T, error) {
	rows, err := data.DbContext.Db.Query(fmt.Sprintf(data.GetJson, workerName), skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jsons []T
	for rows.Next() {
		var data T
		err := rows.Scan(&data)
		if err != nil {
			return nil, err
		}
		jsons = append(jsons, data)
	}

	return jsons, nil
}

func (r mySqlRepositiry[T]) Count(workerName string) (int64, error) {
	var count int64
	err := data.DbContext.Db.QueryRow(fmt.Sprintf(data.Count, workerName)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r mySqlRepositiry[T]) GetAllWorkers() (*[]models.WorkerShort, error) {
	rows, err := data.DbContext.Db.Query(data.Workers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var workers []models.WorkerShort
	for rows.Next() {
		var worker models.WorkerShort
		err := rows.Scan(&worker.Id, &worker.Name)
		if err != nil {
			return nil, err
		}
		workers = append(workers, worker)
	}

	return &workers, nil
}

func (r mySqlRepositiry[T]) AddSettings(settings models.BackgroundWorkerSettings) error {
	var err error
	// var existingTable string
	// err = data.DbContext.Db.QueryRow(data.Exists, settings.WorkerName).Scan(&existingTable)
	// if err != nil && err != sql.ErrNoRows {
	// 	return err
	// } else if existingTable != "" {
	// 	return fmt.Errorf("table %s already exists", settings.WorkerName)
	// }
	doExist, err := r.Exists(settings.WorkerName)
	if err != nil {
		return err
	}
	if doExist {
		return fmt.Errorf("table %s already exists", settings.WorkerName)
	}

	_, err = data.DbContext.Db.Exec(data.CreateWorkerSettingsTable)
	if err != nil {
		return err
	}

	stmt, err := data.DbContext.Db.Prepare(data.InsertWorkerSettings)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		settings.WorkerName,
		settings.Timer,
		settings.RequestDelay,
		settings.Random,
		settings.WritesNumberToSend,
		settings.TotalToSend,
		settings.StopWhenTableEnds)
	if err != nil {
		return err
	}

	err = createTable(settings.WorkerName)
	if err != nil {
		return err
	}

	return nil
}

func (r mySqlRepositiry[T]) AddWorkerResult(log models.WorkerJobLogModel) error {
	_, err := data.DbContext.Db.Exec(data.CreateWorkersResultsTable)
	if err != nil {
		return err
	}

	query := data.DbContext.Db.QueryRow(data.WorkerSettingsId, log.WorkerName)
	var workerSettingsId int64
	query.Scan(&workerSettingsId)

	stmt, err := data.DbContext.Db.Prepare(data.InsertWorkerResults)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(
		workerSettingsId,
		log.StartTime,
		log.StopTime,
		log.MedianReuestDurationTime,
		log.ItemsSended,
		log.SuccessAttemptsCount,
		log.FailedAttemptsCount,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r mySqlRepositiry[T]) GetWorkerResults(workerName string) (*[]models.WorkerJobLogModel, error) {
	rows, err := data.DbContext.Db.Query(data.GetWorkerResults, workerName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var startTimeStr, stopTimeStr string
	layout := "2006-01-02 15:04:05.999999-07:00"
	var logs []models.WorkerJobLogModel
	for rows.Next() {
		var workerName string
		var log models.WorkerJobLogModel
		err := rows.Scan(&workerName,
			&startTimeStr,
			&stopTimeStr,
			&log.MedianReuestDurationTime,
			&log.ItemsSended,
			&log.SuccessAttemptsCount,
			&log.FailedAttemptsCount)

		if err != nil {
			return nil, err
		}

		log.WorkerName = workerName
		log.StartTime, err = time.Parse(layout, startTimeStr)
		if err != nil {
			return nil, err
		}
		log.StopTime, err = time.Parse(layout, stopTimeStr)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
		if err != nil {
			return nil, err
		}
	}

	return &logs, nil
}

func (r mySqlRepositiry[T]) GetWorkerSettings(settingsWorkerName string) (*models.BackgroundWorkerSettings, error) {
	row := data.DbContext.Db.QueryRow(data.GetWorkerSettings, settingsWorkerName)

	var settings models.BackgroundWorkerSettings
	err := row.Scan(&settings.WorkerName,
		&settings.Timer,
		&settings.RequestDelay,
		&settings.Random,
		&settings.WritesNumberToSend,
		&settings.TotalToSend,
		&settings.StopWhenTableEnds)

	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (r mySqlRepositiry[T]) GetWorkerName(id int) (string, error) {
	row := data.DbContext.Db.QueryRow(data.WorkerName, id)

	var name string
	err := row.Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (r mySqlRepositiry[T]) Exists(workerName string) (bool, error) {
	var existingTable string
	err := data.DbContext.Db.QueryRow(data.Exists, workerName).Scan(&existingTable)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	} else if err == sql.ErrNoRows {
		return false, nil
	} else if existingTable != "" {
		return true, nil
	}
	return false, nil
}

func createTable(workerName string) error {
	stmt, err := data.DbContext.Db.Prepare(fmt.Sprintf(data.CreateTable, workerName))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
