import { Component, signal, inject, input } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { SeedingData } from '../models/data/seeding-data';
import { WorkerItemModel } from '../models/worker/worker-item';
import { environment } from '../environments/environments';



@Component({
  selector: 'app-seed-data',
  imports: [FormsModule],
  templateUrl: './seed-data.html',
  styleUrls: [
    './seed-data.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class SeedData {
  data = signal<SeedingData>(new SeedingData())
  errorMessage = signal<string | null>(null)
  isLoading = signal<boolean>(false)
  isMinimized = signal(true)
  workers = input<WorkerItemModel[]>()
  private http = inject(HttpClient);
  selected = signal("")


  onSubmit() {
    this.errorMessage.set(null)
    this.isLoading.set(true)
    this.sendSettings()
  }

  private sendSettings() {
    this.data().workerName = this.selected()
    const payload = {
      ...this.data(),
      jStr: this.data().jStr
    };
    this.http.post(environment.baseAddress + 'seed', payload).subscribe({
      next: response => {
        this.isLoading.set(false)
      },
      error: err => {
        this.errorMessage.set(err.error.error)
        this.isLoading.set(false)
      }
    });
  }

  resize() {
    this.isMinimized.set(!this.isMinimized())
  }

  onSelect(event: Event){
    const selectedWorker = (event.target as HTMLSelectElement).value;
    this.selected.set(selectedWorker)
  }
}
