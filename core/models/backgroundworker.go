package models

import "context"

type Work func(ctx context.Context) int

type BackgroundWorkerSettings struct {
	WorkerName         string           `json:"workerName"`
	Timer              float64          `json:"timer"`
	RequestDelay       float64          `json:"requestDellay"`
	Random             bool             `json:"random"`
	WritesNumberToSend int              `json:"writesNumberToSend"`
	TotalToSend        int64            `json:"totalToSend"`
	StopWhenTableEnds  bool             `json:"stopWhenTableEnds"`
	ConsumerSettings   ConsumerSettings `json:"consumerSettings"`
}

// TODO: make dto for this
type ConsumerSettings struct {
	Host      string    `json:"host"`
	AuthModel AuthModel `json:"authModel"`
}

type AuthModel struct {
	AuthUrl string `json:"authUrl"`
	Model   string `json:"model"`
}
