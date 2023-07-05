import { WS_URL } from "@/app/constants";
import { Subject } from 'rxjs';

var conn: WebSocket;

var messages = new Subject();
export var messages$ = messages.asObservable();

export function connect() {
    if (conn !== undefined) {
        return;
    }

    const lobbyId = sessionStorage.getItem('lobbyId')
    const playerId = sessionStorage.getItem('playerId')

    if (!lobbyId) {
        console.error('Could not fetch lobbyId from session storage. Navigating back to home page');
        throw new Error('Could not fetch lobbyId from session storage. Navigating back to home page')
    }

    conn = new WebSocket(`${WS_URL}/lobbies/${lobbyId}/connect${!!playerId ? '?playerId=' + playerId : ''}`)
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
    connect,
    send,
}
export default ws;