package bbl

import (
	"encoding/json"
	"errors"
	"sync"

	"objectswaterfall.com/core/models"
	"objectswaterfall.com/data/repositories"
	"objectswaterfall.com/utils"
)

type SeedProcessor struct {
	WorkerName string `json:"workerName"`
	Jstr       string `json:"jStr"`
	Count      int    `json:"count"`
}

// TODO: Add context.Context for canceling
func (s SeedProcessor) ProcessJson(areChanks bool, inChank int) error {
	if s.Jstr == "" {
		return errors.New("received json string is empty")
	}

	var jArray models.JArray
	err := json.Unmarshal([]byte(s.Jstr), &jArray)
	if err != nil {
		return err
	}

	if areChanks {
		return insertChanksToDatabase(&jArray, s.WorkerName, s.Count, inChank)
	}
	return insertToDatabase(&jArray, s.WorkerName, s.Count)
}

func insertToDatabase(jArray *models.JArray, workerName string, howMuch int) error {
	errCh := make(chan error)
	repo, err := repositories.NewRepository[string]()
	if err != nil {
		return err
	}
	for howMuch != 0 {
		howMuch--
		for k, v := range *jArray {
			ptr := v
			utils.FillWithDummyData(&ptr)
			(*jArray)[k] = ptr
		}

		dummyData, err := json.Marshal(jArray)
		if err != nil {
			return err
		}
		jStr := string(dummyData)
		go func() {
			errCh <- repo.SetData(workerName, jStr)
		}()

		if err = <-errCh; err != nil {
			return err
		}
	}

	return nil
}

func insertChanksToDatabase(jArray *models.JArray, workerName string, howMuch, inChank int) error {
	errCh := make(chan error, inChank)
	var chank []string
	var wg sync.WaitGroup
	repo, err := repositories.NewRepository[string]()
	if err != nil {
		return err
	}
	for i := 0; i < howMuch; i++ {
		for k, v := range *jArray {
			ptr := v
			utils.FillWithDummyData(&ptr)
			(*jArray)[k] = ptr
		}

		dummyData, err := json.Marshal(jArray)
		if err != nil {
			return err
		}

		chank = append(chank, string(dummyData))
		if len(chank) >= inChank {
			wg.Add(1)
			chankCopy := append([]string(nil), chank...)
			go func(chankData []string) {
				errCh <- repo.SetChankData(workerName, chankData)
				wg.Done()
			}(chankCopy)
			chank = chank[:0]
		}
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()
	for e := range errCh {
		if e != nil {
			return e
		}
	}

	return nil
}
