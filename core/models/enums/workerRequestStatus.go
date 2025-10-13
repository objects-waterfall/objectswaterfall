package enums

type WorkerRequestStatus int

const (
	Success WorkerRequestStatus = iota
	Failed  WorkerRequestStatus = iota
)
