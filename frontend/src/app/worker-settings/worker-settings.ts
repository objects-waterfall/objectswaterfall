import { Component, inject } from '@angular/core';
import { signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { WorkerSettingsModel } from '../models/worker/worker-settings';
import { environment } from '../environments/environments';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-worker-settings',
  imports: [FormsModule],
  templateUrl: './worker-settings.html',
  styleUrls: [
    './worker-settings.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class WorkerSettings {
  newSettings = signal<WorkerSettingsModel>(new WorkerSettingsModel())
  isLoading = signal<boolean>(false)
  isMinimized = signal(true)
  private http = inject(HttpClient);

  constructor(private toastr: ToastrService) {}

  onAdd(){
    this.isLoading.set(true)
    this.sendSettings()
  }

  updateRandom(event: Event) {
    const isChecked = (event.target as HTMLInputElement).checked;
    this.newSettings.update(settings => ({
      ...settings,
      random: isChecked
    }));
  }

  updateEndData(event: Event) {
    const isChecked = (event.target as HTMLInputElement).checked;
    this.newSettings.update(settings => ({
      ...settings,
      StopWhenTableEnds: isChecked
    }));
  }

  sendSettings() {
    const payload = this.newSettings();
    console.log(payload)
    this.http.post(environment.baseAddress + 'add', payload).subscribe({
      next: _ => {
        this.isLoading.set(false)
      },
      error: err => {
        if (err.name === 'HttpErrorResponse'){
              this.toastr.error(`Connection error. ${err.message}`, 'Error');
            } else {
              this.toastr.error(err.message, 'Error');
            }
        this.isLoading.set(false)
      }
    });
  }

  resize() {
    this.isMinimized.set(!this.isMinimized())
  }
}
