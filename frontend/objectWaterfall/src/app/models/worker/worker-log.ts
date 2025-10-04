export class LogModel {
	Log = "";
  requestDirationTime = 0;
  medianReuestDurationTime = 0;
  SuccessAttemptsCount = 0;
  FailedAttemptsCount = 0;

  constructor(data?: Partial<LogModel>) {
    if (data) {
      this.Log = data.Log ?? "";
      this.requestDirationTime = this.formatDuration(data.requestDirationTime);
      this.medianReuestDurationTime = this.formatDuration(data.medianReuestDurationTime);
      this.SuccessAttemptsCount = data.SuccessAttemptsCount ?? 0;
      this.FailedAttemptsCount = data.FailedAttemptsCount ?? 0;
    }
  }

  private formatDuration(value?: number): number {
    return value !== undefined ? parseFloat(value.toFixed(4)) : 0;
  }
}