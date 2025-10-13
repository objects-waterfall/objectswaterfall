import { inject, Injectable } from "@angular/core";
import { WorkerItemModel } from "../models/worker/worker-item";
import { environment } from "../environments/environments";
import { HttpClient } from '@angular/common/http';
import { ResultMap } from "../models/dto/result-map";
import { Result } from "../models/result";
import { catchError, map, Observable, of } from "rxjs";

@Injectable({
    providedIn: 'root'
})
export class WorkersService {
    private http = inject(HttpClient);

    // getWorkers(path: string): Observable<Result<WorkerItemModel[]>>{
    //     let workers: WorkerItemModel[] = []
    //     this.http.get<ResultMap>(environment.baseAddress + path).subscribe({
    //     next: response => {
    //             if (response.result !== null){
    //             for (let i = 0; i < response.result.length; i++) {
    //                         workers[i] = new WorkerItemModel(response.result[i].id, 
    //                             response.result[i].name);
    //                     }
    //             }
    //         },
    //     error: err => {
    //             return new Result(null, err.error.error)
    //         }
    //     });
    //     return new Result(workers, "")
    // }

    getWorkers(path: string): Observable<Result<WorkerItemModel[] | null>>{
        return this.http.get<ResultMap>(environment.baseAddress + path).pipe(
            map(resp => {
                let workers: WorkerItemModel[] = []
                if (resp.result !== null){
                for (let i = 0; i < resp.result.length; i++) {
                            workers[i] = new WorkerItemModel(resp.result[i].id, 
                                resp.result[i].name);
                        }
                }
                return new Result([...workers], "")
            }),
            catchError(err => of(new Result(null, err.error.error)))
        )
    }
}