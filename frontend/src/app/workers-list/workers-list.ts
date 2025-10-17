import { Component, input, output, signal } from '@angular/core';
import { WorkerItemModel } from '../models/worker/worker-item';
import { WorkersItem } from '../workers-item/workers-item';

@Component({
  selector: 'app-workers-list',
  imports: [WorkersItem],
  templateUrl: './workers-list.html',
  styleUrl: './workers-list.css'
})
export class WorkersList {
  workers = input<WorkerItemModel[]>()
  selected = output<number>()
  stoppedWorker = output<{id: number, name: string}>()

  onSelectWorker(id: number){
    this.selected.emit(id)
  }

  onWorkerStopped(worker: {id: number, name: string}){
    this.stoppedWorker.emit(worker)
  }
}
