import { Component, input } from '@angular/core';
import { WorkerResultsDto } from '../models/dto/worker-results-dto';
import { DatePipe } from '@angular/common';

@Component({
  selector: 'app-worker-result-item',
  imports: [DatePipe],
  templateUrl: './worker-result-item.html',
  styleUrls: ['./worker-result-item.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class WorkerResultItem {
  result = input.required<WorkerResultsDto>()
}
