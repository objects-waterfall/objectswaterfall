import { Component, inject, OnDestroy, OnInit, signal } from '@angular/core';
import { WorkerSettings } from "./worker-settings/worker-settings";
import { SeedData } from "./seed-data/seed-data";
import { StartWorker } from "./start-worker/start-worker";
import { WorkerLog } from "./worker-logs/worker-logs";
import { WorkersList } from "./workers-list/workers-list";
import { WorkerItemModel } from './models/worker/worker-item';

import { WorkerRealtimeLogs } from './services/realtime/web-sockets.service';
import { Subscription } from 'rxjs';
import { environment } from './environments/environments';
import { LogModel } from './models/worker/worker-log';
import { WorkersService } from './services/workers.service';
import { firstValueFrom } from 'rxjs';
import { registerLocaleData } from '@angular/common';

@Component({
  selector: 'app-root',
  imports: [WorkerSettings, SeedData, StartWorker, WorkerLog, WorkersList, WorkerLog],
  templateUrl: './app.html',
  styleUrls: ['./app.css',
    '../assets/styles/settings-controls.css'
  ]
})
export class App implements OnInit, OnDestroy {
  private websocketService = inject(WorkerRealtimeLogs)
  private workersService = inject(WorkersService)
  private subscription!: Subscription

  runningWorkers = signal<WorkerItemModel[]>([])
  existingWorkers = signal<WorkerItemModel[]>([])
  isLoading = signal<boolean>(false)
  errorMessage = signal<string | null>(null)
  workerLog = signal<LogModel>(new LogModel())
  selectedForStartWorker = signal<WorkerItemModel>(new WorkerItemModel(0, ""))
  selectedRunningWorker = signal<WorkerItemModel>(new WorkerItemModel(0, ""))
  isRunningWorkerSet = signal<boolean>(false);

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
}
