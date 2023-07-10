import { WS_URL } from "@/app/constants";
import { Subject } from 'rxjs';

var conn: WebSocket;

var messages = new Subject();
export var messages$ = messages.asObservable();

export function init() {
    if (conn !== undefined) {
        return;
    }

    const gameId = sessionStorage.getItem('gameId')
    const playerId = sessionStorage.getItem('playerId')

    if (!gameId) {
        console.error('Could not fetch gameId from session storage. Navigating back to home page');
        throw new Error('Could not fetch gameId from session storage. Navigating back to home page')
    }

    conn = new WebSocket(`${WS_URL}/game/${gameId}/connect${!!playerId ? '?playerId=' + playerId : ''}`)
    conn.onmessage = (e) => {
        const msg = JSON.parse(e.data);
        console.log('Received message', msg)

        messages.next(msg);
    }
}

export function send(msg: string) {
    if (conn === undefined) {
        throw new Error('Cannot send message. Websocket is undefined');
    }

    conn.send(msg);
}

const ws = {
    messages$,
    init,
    send,
}
export default ws;