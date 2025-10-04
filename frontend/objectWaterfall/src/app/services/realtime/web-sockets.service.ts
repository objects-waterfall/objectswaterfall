import { Injectable, OnDestroy } from "@angular/core";
import { Observable, share, Subject, takeUntil } from "rxjs";
import { webSocket, WebSocketSubject } from "rxjs/webSocket";

@Injectable({
    providedIn: 'root'
})
export class WorkerRealtimeLogs implements OnDestroy {
    socket$!: WebSocketSubject<any>
    destroy$ = new Subject<void>()
    messages$!: Observable<any>

    startConnection(url: string) {
        if (this.socket$){
            return
        }

        this.socket$ = webSocket(url)

        this.messages$ = this.socket$.pipe(
            takeUntil(this.destroy$),
            share()
        )
    }

    send(msg: any) {
        if (!this.socket$){
            new Error('WebSocket is not connected!');
        }
        this.socket$.next(msg)
    }

    close(): void {
        this.destroy$.next();
        this.socket$?.complete()
    }

    ngOnDestroy(): void {
        this.close()
    }
}