import { Component, input, signal, inject, OnInit } from '@angular/core';
import { WorkerResultsDto } from '../models/dto/worker-results-dto';
import { WorkerResultItem } from '../worker-result-item/worker-result-item';
import { WorkerItemModel } from '../models/worker/worker-item';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environments';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-worker-result-list',
  imports: [WorkerResultItem, 
    FormsModule
  ],
  templateUrl: './worker-result-list.html',
  styleUrls: ['./worker-result-list.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class WorkerResultList implements OnInit {
  private http = inject(HttpClient);
  workers = input<WorkerItemModel[]>()
  results = signal<WorkerResultsDto[]>([])
  selected = signal(-1)
  workerName = signal<string>("")

  constructor(private toastr: ToastrService) {}

  ngOnInit(): void {
    this.http.get<{result: WorkerResultsDto[]}>(environment.baseAddress + `getWorkerResults?take=9`).subscribe({
      next: res => {
        if (res.result === null) {
          this.results.set([]);
          return;
        }
        this.results.set([...res.result]);
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

  onSelect(event: Event){
    const selectedWorker = (event.target as HTMLSelectElement).value;
    const worker = this.workers()?.find(w => w.id === Number(selectedWorker));
    this.workerName.set(worker?.name || "")

    this.getWorkerResults(`getWorkerResults?workerName=${worker?.name}&take=9`);
  }

  getWorkerResults(path: string) {
    this.http.get<{result: WorkerResultsDto[]}>(environment.baseAddress + path).subscribe({
      next: res => {
        if (res.result === null) {
          console.log(res);
          this.results.set([]);
          return;
        }
        this.results.set([...res.result]);
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
}
