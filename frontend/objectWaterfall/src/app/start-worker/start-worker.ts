import { Component, signal, inject, input } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { WorkerItemModel } from '../models/worker/worker-item';
import { FormsModule } from '@angular/forms';
import { StartWorkerData } from '../models/worker/start-worker';
import { environment } from '../environments/environments';

@Component({
  selector: 'app-start-worker',
  imports: [FormsModule],
  templateUrl: './start-worker.html',
  styleUrls: [
    './start-worker.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class StartWorker {
  placeholdertext = '{\n "UserName": "Name", \n "UserPassword": "Password" \n}'
  errorMessage = signal<string | null>(null)
  isMinimized = signal(false)
  startWorkerData = signal(new StartWorkerData())
  workers = input<WorkerItemModel[]>()
  private http = inject(HttpClient);
  selected = signal(0)

  onSelect(event: Event){
    const selectedWorker = (event.target as HTMLSelectElement).value;
    this.selected.set(Number(selectedWorker))
  }

  onStart(){
        this.http.post(environment.baseAddress + 'start?id=' + this.selected(), this.startWorkerData()).subscribe({
          next: response => {
            console.log(response)
          },
          error: err => {
            this.errorMessage.set(err.error.error)
          }
        });
  }

  resize() {
    this.isMinimized.set(!this.isMinimized())
  }
}
