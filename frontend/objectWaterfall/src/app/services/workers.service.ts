import { inject, Injectable } from "@angular/core";
import { WorkerItemModel } from "../models/worker/worker-item";
import { environment } from "../environments/environments";
import { HttpClient } from '@angular/common/http';
import { ResultMap } from "../models/dto/result-map";
import { Result } from "../models/result";

@Injectable({
    providedIn: 'root'
})
export class WorkersService {
    private http = inject(HttpClient);

    getWorkers(path: string): Result<WorkerItemModel[]>{
        let workers: WorkerItemModel[] = []
        this.http.get<ResultMap>(environment.baseAddress + path).subscribe({
        next: response => {
                if (response.result !== null){
                for (let i = 0; i < response.result.length; i++) {
                            workers[i] = new WorkerItemModel(response.result[i].id, 
                                response.result[i].name);
                        }
                }
                
            },
        error: err => {
                return new Result(null, err.error.error)
            }
        });
        return new Result(workers, "")
    }
}