import { Component, computed, input, signal } from '@angular/core';
import { LogModel } from '../models/worker/worker-log';
import { WorkerLogItem } from '../worker-log-item/worker-log-item';
import { sign } from 'node:crypto';

@Component({
  selector: 'app-worker-logs',
  imports: [WorkerLogItem],
  templateUrl: './worker-logs.html',
  styleUrls:[
    './worker-logs.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class WorkerLogs {
  workerName = input<string>()
  inputLogs = input<LogModel[]>()
  isWorkerChoosen = input<boolean>()
}
