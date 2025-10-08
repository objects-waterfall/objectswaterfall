package bbl

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"objectswaterfall.com/core/errors"
	"objectswaterfall.com/core/models"
	"objectswaterfall.com/core/services"
	"objectswaterfall.com/data/repositories"
	"objectswaterfall.com/utils/stopwatch"
)

const totalMessage = "|| Total amount of records have been sent %d of %d"
const successMessage = "Request %d of %s was success " + totalMessage
const failedMessage = "Request %d of %s was failed " + totalMessage

type SendWorker struct {
	settings     models.BackgroundWorkerSettings
	cancelFunc   context.CancelFunc
	group        sync.WaitGroup
	repo         repositories.Repository[string]
	totalSended  int64
	tokenService TokenService
	medianValue  models.MedianValue
	log          *models.LogModel
	models.LogFunc
}

type dataResult struct {
	data []string
	err  error
}

type requestResult struct {
	requestRes models.ResponseResult
	err        error
}

func NewSendWorker(settings models.BackgroundWorkerSettings, tokenService TokenService) services.Worker {
	repo, err := repositories.NewRepository[string]()
	if err != nil {
		panic(err)
	}
	return &SendWorker{
		settings:     settings,
		repo:         repo,
		tokenService: tokenService,
		log:          &models.LogModel{},
		medianValue:  models.NewMedianValue(),
	}
}

func (w *SendWorker) DoWork(ctx context.Context) {
	log.Printf("Worker was started at %v", time.Now())
	var counter int64 = 0
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker was stoped at %v, because of: %s ", time.Now(), ctx.Err().Error())
			return
		default:
			w.work(&counter)
		}
	}
}

func (w *SendWorker) SetCancel(cancel context.CancelFunc) {
	w.cancelFunc = cancel
}

func (w *SendWorker) Cancel() {
	w.cancelFunc()
	w.group.Wait()
}

func (w *SendWorker) GetWorkerName() string {
	return w.settings.WorkerName
}

func (w *SendWorker) Log() *models.LogModel {
	return w.log
}

func (w *SendWorker) SetLogFunc(logFunc models.LogFunc) {
	w.LogFunc = logFunc
}

func (w *SendWorker) work(counter *int64) {
	w.group.Add(1)

	go w.actualWork()
	*counter += 1

	if w.totalSended >= w.settings.TotalToSend {
		w.group.Wait()
		log.Printf("Worker is done. Sent %d of %d", w.totalSended, w.settings.TotalToSend)
		w.log.Log = fmt.Sprintf("Worker is done. Sent %d of %d", w.totalSended, w.settings.TotalToSend)
		w.log.MedianReuestDurationTime = w.medianValue.FindMedian()
		if w.LogFunc != nil {
			w.LogFunc(*w.log)
		}
		w.Cancel()
		return
	}

	time.Sleep(time.Duration(w.settings.RequestDelay) * time.Second)

	tableTotalRecordsNumber, _ := w.repo.Count(w.settings.WorkerName) // Move somewhere for do not making calculation each time. A good place is Worker.Settings. Set during init of worker.

	if tableTotalRecordsNumber <= w.totalSended {
		switch {
		case !w.settings.StopWhenTableEnds && !w.settings.Random:
			w.totalSended = 0
		case w.settings.StopWhenTableEnds && !w.settings.Random:
			w.Cancel()
		}
	}
}

func (w *SendWorker) actualWork() {
	defer w.group.Done()
	dataCh := make(chan dataResult)
	go w.getData(dataCh)

	dataResult := <-dataCh
	if dataResult.err != nil {
		log.Println(dataResult.err)
		return
	}

	respCh := make(chan requestResult)
	sw := stopwatch.NewStopWatch()
	sw.Start()
	go w.sendRequest(dataResult, respCh)
	respRes := <-respCh
	requstDuration := sw.Elapsed(time.Second)
	if respRes.err != nil {
		w.medianValue.AddNum(requstDuration)
		w.setLog(fmt.Sprintf(failedMessage+" error %s", w.log.SuccessAttemptsCount+w.log.FailedAttemptsCount, w.settings.WorkerName, w.totalSended, w.settings.TotalToSend, respRes.err), requstDuration, false)
		if w.LogFunc != nil {
			w.LogFunc(*w.log)
		}
		if _, ok := respRes.err.(errors.TokenRecievingError); ok {
			log.Println(respRes.err)
			w.Cancel()
			return
		}
		log.Println(respRes.err)
		return
	}
	w.medianValue.AddNum(requstDuration)
	w.setLog(fmt.Sprintf(successMessage, w.log.SuccessAttemptsCount+w.log.FailedAttemptsCount, w.settings.WorkerName, w.totalSended, w.settings.TotalToSend), requstDuration, true)
	if w.LogFunc != nil {
		w.LogFunc(*w.log)
	}

	log.Println(w.log)
}

func (w *SendWorker) getData(dataCh chan dataResult) {
	defer close(dataCh)
	var skip int64
	if w.settings.Random {
		count, err := w.repo.Count(w.settings.WorkerName)
		if err != nil {
			dataCh <- dataResult{
				data: nil,
				err:  err,
			}
			return
		}
		skip = rand.Int63n(count)
	} else {
		skip = w.totalSended
	}

	data, err := w.repo.GetData(w.settings.WorkerName, w.settings.Random, w.settings.WritesNumberToSend, skip)
	dataCh <- dataResult{
		data: data,
		err:  err,
	}
}

func (w *SendWorker) sendRequest(data dataResult, respCh chan requestResult) {
	defer close(respCh)
	sending := NewSendingService()
	var (
		token   string
		headers = make(map[string]string)
	)
	if w.settings.ConsumerSettings.AuthModel.AuthUrl != "" && w.settings.ConsumerSettings.AuthModel.Model != "" {
		var err error
		token, err = w.tokenService.Token()
		if err != nil {
			respCh <- requestResult{
				requestRes: models.ResponseResult{},
				err:        fmt.Errorf("auth error. %w", err),
			}
			return
		}
		headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	}
	resp, err := sending.SendRequest(w.settings.ConsumerSettings.Host, data.data, headers)
	w.totalSended += int64(len(data.data))
	respCh <- requestResult{
		requestRes: resp,
		err:        err,
	}
}

func (w *SendWorker) setLog(msg string, duration float64, isSuccess bool) {
	w.log.Log = msg
	w.log.RequestDirationTime = duration
	w.log.MedianReuestDurationTime = w.medianValue.FindMedian()
	if isSuccess {
		w.log.SuccessAttemptsCount++
	} else {
		w.log.FailedAttemptsCount++
	}
}
