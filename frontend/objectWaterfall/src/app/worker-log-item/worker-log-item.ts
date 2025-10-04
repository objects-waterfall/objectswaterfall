import { Component, input } from '@angular/core';
import { LogModel } from '../models/worker/worker-log';

@Component({
  selector: 'app-worker-log-item',
  imports: [],
  templateUrl: './worker-log-item.html',
  styleUrls:[ './worker-log-item.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class WorkerLogItem {
  workerLog = input.required<LogModel>()
}
