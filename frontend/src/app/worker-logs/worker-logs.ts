import { Component, input } from '@angular/core';
import { LogModel } from '../models/worker/worker-log';
import { WorkerStopStatus } from '../models/enums/worker-stop-status';
import { WorkerRequestStatus } from '../models/enums/worker-request-status';

@Component({
  selector: 'app-worker-log',
  imports: [],
  templateUrl: './worker-logs.html',
  styleUrls:[
    './worker-logs.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class WorkerLog {
  WorkerStopStatus = WorkerStopStatus
  CurrentRequestStatus = WorkerRequestStatus
  inputLog = input.required<LogModel>()
  isWorkerChoosen = input<boolean>()
}
