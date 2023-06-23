import { describe, expect, it, jest } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import WebSocket from 'ws';

describe('Connect', () => {
    it('should receive a playerId', async () => {
        const { conn, messageCb } = await connect();

        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                "ID": expect.anything()
            }));
        });

        await disconnect(conn);
    })

    it('should receive a player list', async () => {
        const { conn, messageCb } = await connect();

        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                "Players": expect.anything()
            }));
        });

        await disconnect(conn);
    })
})

async function connect(): Promise<{ conn: WebSocket, messageCb: jest.Mock }> {
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

async function disconnect(conn: WebSocket): Promise<void> {
    const closeCb = jest.fn();
    conn.onclose = closeCb;

    conn.terminate();

    await waitForExpect(() => {
        expect(closeCb).toHaveBeenCalled()
    })
}