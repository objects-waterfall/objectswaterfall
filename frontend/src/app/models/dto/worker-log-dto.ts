import { WorkerRequestStatus } from "../enums/worker-request-status"
import { WorkerStopStatus } from "../enums/worker-stop-status"

// TODO: remove
export class WorkerLogDto {
    WorkerName = ""
    WorkerStopStatus = WorkerStopStatus.runnning
    TotalItemsToSend = 0
    ItemsSended = 0
    CurrentRequestStatus = WorkerRequestStatus.Success
    RequestErrorMessage = ""
    RequestDirationTime = 0
    MedianReuestDurationTime = 0
    StartTime = new Date()
    StopTime = new Date()
    RequestNumber = 0
    SuccessAttemptsCount = 0
    FailedAttemptsCount = 0
}