import { Component, input, output, signal } from '@angular/core';
import { WorkerItemModel } from '../models/worker/worker-item';

@Component({
  selector: 'app-workers-item',
  imports: [],
  templateUrl: './workers-item.html',
  styleUrls: [
    './workers-item.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class WorkersItem {
  worker = input.required<WorkerItemModel>()
  selectedItem = output<number>()
  isMininimized = signal<boolean>(true)
  isRunning = signal<boolean>(true)

  onSelectedHandler() {
    this.selectedItem.emit(this.worker().id)
  }

  onStop() {
    this.isRunning.set(!this.isRunning())
  }
}
