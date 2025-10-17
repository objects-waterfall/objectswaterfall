import { Component, inject, OnDestroy, OnInit, signal } from '@angular/core';
import { WorkerSettings } from "./worker-settings/worker-settings";
import { SeedData } from "./seed-data/seed-data";
import { StartWorker } from "./start-worker/start-worker";
import { WorkerLog } from "./worker-logs/worker-logs";
import { WorkersList } from "./workers-list/workers-list";
import { WorkerItemModel } from './models/worker/worker-item';

import { WorkerRealtimeLogs } from './services/realtime/web-sockets.service';
import { single, Subscription } from 'rxjs';
import { environment } from './environments/environments';
import { LogModel } from './models/worker/worker-log';
import { WorkersService } from './services/workers.service';
import { firstValueFrom } from 'rxjs';
import { WarningPopup } from "./warning-popup/warning-popup";
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-root',
  imports: [WorkerSettings, SeedData, StartWorker, WorkerLog, WorkersList, WorkerLog, WarningPopup],
  templateUrl: './app.html',
  styleUrls: ['./app.css',
    '../assets/styles/settings-controls.css'
  ]
})
export class App implements OnInit, OnDestroy {
  private http = inject(HttpClient);
  private websocketService = inject(WorkerRealtimeLogs)
  private workersService = inject(WorkersService)
  private subscription!: Subscription
  private workerForStop: {id: number, name: string} | undefined

  runningWorkers = signal<WorkerItemModel[]>([])
  existingWorkers = signal<WorkerItemModel[]>([])
  isLoading = signal<boolean>(false)
  errorMessage = signal<string | null>(null)
  workerLog = signal<LogModel>(new LogModel())
  selectedForStartWorker = signal<WorkerItemModel>(new WorkerItemModel(0, ""))
  selectedRunningWorker = signal<WorkerItemModel>(new WorkerItemModel(0, ""))
  isRunningWorkerSet = signal<boolean>(false)
  showWarningPopup = signal<boolean>(false)
  warningTitle = signal<string>("Stopping the shit worker")
  warningMessage = signal<string>(`You are about to ptop the worker "shit". Are you sure?`)

  ngOnDestroy(): void {
    this.subscription.unsubscribe()
    this.websocketService.close()
  }

  async ngOnInit() {
    await this.getRunningWorkers()
    this.getExistingWorkers()
    this.websocketService.startConnection(environment.baseAddress + 'logsWs')
    this.subscription = this.websocketService.messages$.subscribe({
      next: (msg: LogModel) => {
        const selected = this.selectedRunningWorker().name
        if (selected === msg.WorkerName){
          this.workerLog.set(new LogModel(msg))
        }
      },
      error: err => this.errorMessage = err,
      complete: () => console.log('Socket close')
    })
  }

  sendPing(): void {
    this.websocketService.send({ type: 'PING', timestamp: Date.now() });
  }

  async getRunningWorkers() {
    this.isLoading.set(true);
    const res = await firstValueFrom(this.workersService.getWorkers('getRunningWorkers'));
    if (res.Err !== "") {
      this.errorMessage.set(res.Err);
      this.isLoading.set(false);
      return;
    }
    this.runningWorkers.set([...res.Data!]); // ensure new array reference
    this.isLoading.set(false);
  }

  async getExistingWorkers() {
    this.isLoading.set(true);
    this.workersService.getWorkers('getWorkers').subscribe(res => {
      if (res.Err !== "") {
        this.errorMessage.set(res.Err);
        this.isLoading.set(false);
        return;
      }
      this.existingWorkers.set([...res.Data!]);
      this.isLoading.set(false);
    });
  }

  onSelectedRunningWorker(id: number) {
    this.websocketService.send({"workerId" : id})
    this.selectedRunningWorker.set(this.runningWorkers().find(x => x.id == id)!)
    this.isRunningWorkerSet.set(true)
  }

  async onNewWorkerStarted(id: number) {
    await this.getRunningWorkers();
    const newWorker = this.runningWorkers().find(x => x.id === id);
    if (newWorker) {
      this.onSelectedRunningWorker(id)
    }
  }

  onStoppedWorker(w: {id: number, name: string}){
    this.warningTitle.set(`Stopping the "${w.name}" worker`)
    this.warningMessage.set(`You are about to ptop the worker "${w.name}". Are you sure?`)
    this.workerForStop = w
    this.showWarningPopup.set(true)
  }

  onConfirmStopWorker(confirmed: boolean) {
    if (!confirmed) {
      this.showWarningPopup.set(false)
      return
    }

    this.http.get(environment.baseAddress + 'stop?id=' + this.workerForStop!.id).subscribe({
              next: _ => {
                this.removeStoppedWorkerFromList()
              },
              error: err => {
                this.errorMessage.set(err.error.error)
              }
            });
    this.showWarningPopup.set(false)
  }

  private removeStoppedWorkerFromList(){
    const worker = this.runningWorkers().find(x => x.id ===  this.workerForStop!.id)
    if (!worker)
    {
      return
    }
    const index = this.runningWorkers().indexOf(worker)
    if (index > -1){
      this.runningWorkers.set(this.runningWorkers().slice(index, this.runningWorkers().length - 1))
    }
  }
}
