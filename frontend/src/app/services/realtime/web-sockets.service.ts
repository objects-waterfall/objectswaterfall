import { Injectable, OnDestroy } from "@angular/core";
import { Observable, share, Subject, takeUntil } from "rxjs";
import { webSocket, WebSocketSubject } from "rxjs/webSocket";
import { ToastrService } from 'ngx-toastr';

@Injectable({
    providedIn: 'root'
})
export class WorkerRealtimeLogs implements OnDestroy {
    socket$!: WebSocketSubject<any>
    destroy$ = new Subject<void>()
    messages$!: Observable<any>

    constructor(private toastr: ToastrService) {}

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
        if (!this.socket$) {
            this.toastr.error('WebSocket connection is not established.', 'Error');
        }
        this.socket$.next(msg)
    }

    close(closeMessae: any = null): void {
        if (closeMessae !== undefined || closeMessae !== null) {
            this.send(closeMessae)
        }
        this.destroy$.next();
        this.destroy$.complete();

        if (this.socket$) {
            this.socket$.complete();
            const nativeSocket = (this.socket$ as any)._socket;
            nativeSocket?.close?.();
        }

        this.socket$ = undefined!;
    }

    ngOnDestroy(): void {
        this.close()
    }
}