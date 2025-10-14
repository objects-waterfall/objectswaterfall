package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	bbl "objectswaterfall.com/BBL"
	"objectswaterfall.com/application/dtos"
	"objectswaterfall.com/application/queries"
	"objectswaterfall.com/core/mappers"
	"objectswaterfall.com/core/models"
	"objectswaterfall.com/core/models/enums"
	"objectswaterfall.com/data/repositories"
	"objectswaterfall.com/stores"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Add(ctx *gin.Context) {
	var workerSettingsDto dtos.BackgroundWorkerSettingsDto
	if err := ctx.ShouldBindBodyWithJSON(&workerSettingsDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	repo, err := repositories.NewRepository[any]()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workerSettings := mappers.FromDtoToWorkerSettings(workerSettingsDto)
	if err := repo.AddSettings(workerSettings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func Start(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("id shouldn't be 0")})
		return
	}
	repo, err := repositories.NewRepository[any]()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name, err := repo.GetWorkerName(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store := stores.GetWorkerStore()
	if store.Exists(name) {
		ctx.JSON(http.StatusConflict, gin.H{"Error": fmt.Sprintf("The worker %s is running alredy", name)})
		return
	}

	var consumerSettings models.ConsumerSettings
	if err := ctx.ShouldBindBodyWithJSON(&consumerSettings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workerSettings, err := repo.GetWorkerSettings(name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	count, err := repo.Count(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("ther is no dummy data for %s", name)})
		return
	}
	workerSettings.ConsumerSettings = consumerSettings

	duration := time.Now().Add(time.Minute * time.Duration(workerSettings.Timer))
	context, cancel := context.WithDeadline(context.Background(), duration)
	worker := bbl.NewSendWorker(*workerSettings,
		bbl.NewTokenService(workerSettings.ConsumerSettings.AuthModel.AuthUrl,
			workerSettings.ConsumerSettings.AuthModel.Model,
			time.Duration(5*time.Minute))) //TODO: Should replace 5*time.Minute to a variable or request parameter
	worker.SetCancel(cancel)
	workerId := store.Add(id, &worker)

	go worker.DoWork(context)

	ctx.JSON(http.StatusOK, gin.H{"workerId": workerId})
}

func Stop(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("id shouldn't be 0")})
		return
	}

	store := stores.GetWorkerStore()
	err = store.CancelWork(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	worker, err := store.Get(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log := (*worker).Log()
	log.WorkerStopStatus = enums.StoppedByUser

	err = store.Remove(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "Ok"})
}

// TODO: make with id instead of name
func Seed(ctx *gin.Context) {
	var seedProc bbl.SeedProcessor
	err := ctx.ShouldBindBodyWithJSON(&seedProc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: check existance in seedProc
	if repo, err := repositories.NewRepository[any](); err == nil {
		if exists, err := repo.Exists(seedProc.WorkerName); err == nil && !exists {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Errorf("there is no worker named %s", seedProc.WorkerName).Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	errCh := make(chan error)

	go func() {
		errCh <- seedProc.ProcessJson(false, 0)
	}()

	if err = <-errCh; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "Ok"})
}

func GetWorkers(ctx *gin.Context) {
	repo, err := repositories.NewRepository[any]()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	workers, err := repo.GetAllWorkers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var workersDto []dtos.WorkerShortDto
	for _, v := range *workers {
		workersDto = append(workersDto, dtos.WorkerShortDto{
			Id:   v.Id,
			Name: v.Name,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"result": workersDto})
}

func GetRunningWorkers(ctx *gin.Context) {
	store := stores.GetWorkerStore()
	workers := store.All()
	var workersDto []dtos.WorkerShortDto
	for _, v := range *workers {
		workersDto = append(workersDto, dtos.WorkerShortDto{
			Id:   v.Id,
			Name: v.Name,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"result": workersDto})
}

func GetWorkerResults(ctx *gin.Context) {
	workerName := ctx.Query("workerName")
	repo, err := repositories.NewRepository[any]()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results, err := repo.GetWorkerResults(workerName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resultsDto []dtos.WorkerResultsDto
	for _, v := range *results {
		resultsDto = append(resultsDto, dtos.ToLogResult(v))
	}

	ctx.JSON(http.StatusOK, gin.H{"result": resultsDto})
}

var connMutex sync.Mutex
var dataCh = make(chan []byte, 100)

func WebSocketHandler(ctx *gin.Context) {
	errCh := make(chan error, 1)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()
	defer close(errCh)
	store := stores.GetWorkerStore()

	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			errCh <- err
			break
		}
		var wr queries.WorkerRequest
		err = json.Unmarshal(p, &wr)
		if err != nil {
			errCh <- err
			break
		}

		worker, err := store.Get(wr.WorkerId)
		if err != nil {
			errCh <- err
			break
		}

		(*worker).SetLogFunc(func(l models.WorkerJobLogModel) {
			data, _ := json.Marshal(l)
			dataCh <- data
			go func() {
				for msg := range dataCh {
					connMutex.Lock()
					err := conn.WriteMessage(msgType, msg)
					connMutex.Unlock()
					if err != nil {
						errCh <- err
						return
					}
					errCh <- nil
				}
			}()
		})

		err = <-errCh
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = conn.WriteMessage(msgType, p)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}
