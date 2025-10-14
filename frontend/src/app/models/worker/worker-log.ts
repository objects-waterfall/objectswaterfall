import { WorkerRequestStatus } from "../enums/worker-request-status"
import { WorkerStopStatus } from "../enums/worker-stop-status"

export class LogModel {
	  WorkerName = ""
    WorkerStopStatus: WorkerStopStatus | undefined
    TotalItemsToSend = 0
    ItemsSended = 0
    CurrentRequestStatus = WorkerRequestStatus.Failed
    RequestErrorMessage: string | undefined
    RequestDirationTime = 0
    MedianReuestDurationTime = 0
    StartTime = new Date()
    StopTime: Date | null = null
    RequestNumber = 0
    SuccessAttemptsCount = 0
    FailedAttemptsCount = 0

  constructor(data?: Partial<LogModel>) {
    if (data) {
      this.WorkerName = data.WorkerName ?? ""
      this.WorkerStopStatus = data.WorkerStopStatus ?? WorkerStopStatus.runnning
      this.TotalItemsToSend = data.TotalItemsToSend ?? 0
      this.ItemsSended = data.ItemsSended ?? 0
      this.CurrentRequestStatus = data.CurrentRequestStatus ?? WorkerRequestStatus.Failed
      this.RequestErrorMessage = data.RequestErrorMessage ?? undefined
      this.RequestDirationTime = this.formatDuration(data.RequestDirationTime);
      this.MedianReuestDurationTime = this.formatDuration(data.MedianReuestDurationTime);
      this.StartTime = data.StartTime ?? new Date()
      this.StopTime = data.StopTime ?? null
      this.RequestNumber = data.RequestNumber ?? 0
      this.SuccessAttemptsCount = data.SuccessAttemptsCount ?? 0;
      this.FailedAttemptsCount = data.FailedAttemptsCount ?? 0;
    }
  }

  private formatDuration(value?: number): number {
    return value !== undefined ? parseFloat(value.toFixed(4)) : 0;
  }
}