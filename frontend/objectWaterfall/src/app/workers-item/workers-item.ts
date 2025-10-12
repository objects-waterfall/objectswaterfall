import { Component, input, inject, signal, output } from '@angular/core';
import { WorkerItemModel } from '../models/worker/worker-item';
import { environment } from '../environments/environments';
import { HttpClient } from '@angular/common/http';

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
  private http = inject(HttpClient);

  worker = input.required<WorkerItemModel>()
  selectedItem = output<number>()
  isMininimized = signal<boolean>(true)
  isRunning = signal<boolean>(true)
  errorMessage = signal<string | null>(null)

  onSelectedHandler() {
    this.selectedItem.emit(this.worker().id)
  }

  onStop() {
    this.http.get(environment.baseAddress + 'stop?id=' + this.worker().id).subscribe({
              next: response => {
                console.log(response)
              },
              error: err => {
                this.errorMessage.set(err.error.error)
              }
            });
    this.isRunning.set(!this.isRunning())
  }
}
