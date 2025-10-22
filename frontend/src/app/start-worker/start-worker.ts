import { Component, signal, inject, input, output } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { WorkerItemModel } from '../models/worker/worker-item';
import { FormsModule } from '@angular/forms';
import { StartWorkerData } from '../models/worker/start-worker';
import { environment } from '../environments/environments';
import { ToastrService } from 'ngx-toastr';
import e from 'express';

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
  private http = inject(HttpClient);

  placeholdertext = '{\n "UserName": "Name", \n "UserPassword": "Password" \n}'
  isMinimized = signal(false)
  // TODO: separate models (StartWorkerData and AuthModel) in here and make a checkbox like "use auth" or something
  startWorkerData = signal(new StartWorkerData())
  workers = input<WorkerItemModel[]>()
  selected = signal(0)
  newWorkerStarted = output<number>()

  constructor(private toastr: ToastrService) {}

  onSelect(event: Event){
    const selectedWorker = (event.target as HTMLSelectElement).value;
    this.selected.set(Number(selectedWorker))
  }

  onStart(){
        this.http.post<{workerId: number}>(environment.baseAddress + 'start?id=' + this.selected(), this.startWorkerData()).subscribe({
          next: res => {
            this.newWorkerStarted.emit(res.workerId)
          },
          error: err => {
            if (err.name === 'HttpErrorResponse'){
              this.toastr.error(`Connection error. ${err.message}`, 'Error');
            } else {
              this.toastr.error(err.message, 'Error');
            }
          }
        });
  }

  resize() {
    this.isMinimized.set(!this.isMinimized())
  }
}
