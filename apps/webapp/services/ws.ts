import { WS_URL } from "@/app/constants";
import { Subject } from 'rxjs';

var conn: WebSocket | undefined;

var messages = new Subject();
export var messages$ = messages.asObservable();

export function init() {
    if (conn !== undefined){
        return;
    }

    const gameId = sessionStorage.getItem('gameId')
    const playerId = sessionStorage.getItem('playerId')

    if (!gameId) {
        console.error('Could not fetch gameId from session storage. Navigating back to home page');
        throw new Error('Could not fetch gameId from session storage. Navigating back to home page')
    }

    console.log('Connecting WS to', gameId)

    conn = new WebSocket(`${WS_URL}/game/${gameId}/connect${!!playerId ? '?playerId=' + playerId : ''}`)
    conn.onmessage = (e) => {
        const msg = JSON.parse(e.data);
        console.log('Received message', msg)

        messages.next(msg);
    }
    conn.onclose = (e) => {
        console.log('Close Event', e)
    }
}

export function disconnect() {
    if (conn == undefined) {
        return;
    }

    conn.close(1000);
    conn = undefined;
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
    disconnect,
    send,
}
export default ws;