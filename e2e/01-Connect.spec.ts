import { describe, expect, it, jest } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import WebSocket from 'ws';

describe('Connect', () => {
    it('should receive a playerId', async () => {
        const { conn, messageCb } = await connect();

        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                event: 'notification.player-id',
                payload: {
                    playerId: expect.anything(),
                }
            }));
        });

        console.log(getPlayerId(messageCb));

        await disconnect(conn);
    })

    it('should receive a player list', async () => {
        const { conn, messageCb } = await connect();

        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                event: 'state.player-list',
                payload: {
                    players: expect.anything(),
                }
            }));
        });

        await disconnect(conn);
    })
})

export async function connect(): Promise<{ conn: WebSocket, messageCb: jest.Mock }> {
    let conn: WebSocket
    const openCb = jest.fn();
    const messageCb = jest.fn();

    conn = new WebSocket('ws://localhost:8080/connect')
    conn.onopen = openCb;
    conn.onmessage = (event) => {
        messageCb(JSON.parse(event.data.toString()))
    };

    await waitForExpect(() => {
        expect(openCb).toHaveBeenCalled()
    })

    return {
        conn,
        messageCb,
    }
}

export async function disconnect(conn: WebSocket): Promise<void> {
    const closeCb = jest.fn();
    conn.onclose =closeCb 

    conn.close(1000); 

    // NOTE: The websocket object seems to keep a handle open even after calling close(). Only the terminate() method actually frees the resource
    // and triggerers the close event
    conn.terminate();

    await waitForExpect(() => {
        expect(closeCb).toHaveBeenCalled()
    })
}

export function getPlayerId(messageCb: jest.Mock): string {
    for(const call of messageCb.mock.calls) {
        const body: any = call[0]
        if (body.event === 'notification.player-id') {
            return body.payload.playerId
        }
    }
    return '';
}