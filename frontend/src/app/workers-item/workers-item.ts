import { Component, input, inject, signal, output } from '@angular/core';
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
  stoppedWorker = output<{id: number, name: string}>()
  isMininimized = signal<boolean>(true)
  isRunning = signal<boolean>(true)
  errorMessage = signal<string | null>(null)

  onSelectedHandler() {
    this.selectedItem.emit(this.worker().id)
  }

  onStop() {
    this.stoppedWorker.emit({
                  id: this.worker().id,
                  name: this.worker().name
                })
    // this.http.get(environment.baseAddress + 'stop?id=' + this.worker().id).subscribe({
    //           next: _ => {
    //             this.stoppedWorker.emit({
    //               id: this.worker().id,
    //               name: this.worker().name
    //             })
    //             // this.isRunning.set(!this.isRunning())
    //             // console.log(response)
    //           },
    //           error: err => {
    //             this.errorMessage.set(err.error.error)
    //           }
    //         });
  }
}
