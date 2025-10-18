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
        if (!this.socket$) {
            throw new Error('WebSocket is not connected!');
        }
        this.socket$.next(msg)
    }

    close(): void {
        this.destroy$.next();
        this.destroy$.complete();

        if (this.socket$) {
            this.socket$.complete();

            // Optional: forcibly close native socket if needed
            const nativeSocket = (this.socket$ as any)._socket;
            nativeSocket?.close?.();
        }

        this.socket$ = undefined!;
    }

    ngOnDestroy(): void {
        this.close()
    }
}