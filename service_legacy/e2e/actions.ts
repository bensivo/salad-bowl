import { expect, jest } from '@jest/globals';
import axios from 'axios';
import waitForExpect from 'wait-for-expect';
import WebSocket from 'ws';

export async function createGame(): Promise<string> {
    let createRes;
    try {
        createRes = await axios.request({
            method: 'POST',
            url: 'http://localhost:8080/game',
            data: {}
        });
    } catch (e) {
        throw new Error('Create game failed' + e);
    }

    const gameId = createRes.data.gameId;

    const getRes = await axios.request({
        method: 'GET',
        url: 'http://localhost:8080/game'
    })

    expect(getRes.data[gameId]).toBeTruthy();

    return gameId;
}

export async function connect(gameId: string, playerId?: string): Promise<{ conn: WebSocket, messageCb: jest.Mock }> {
    let conn: WebSocket
    const openCb = jest.fn();
    const messageCb = jest.fn();

    conn = !!playerId ?
        new WebSocket(`ws://localhost:8080/game/${gameId}/connect?playerId=${playerId}`)
        :
        new WebSocket(`ws://localhost:8080/game/${gameId}/connect`)
        ;
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
    conn.onclose = closeCb

    conn.close(1000);

    // NOTE: The websocket object seems to keep a handle open even after calling close(). Only the terminate() method actually frees the resource
    // and triggerers the close event
    conn.terminate();

    await waitForExpect(() => {
        expect(closeCb).toHaveBeenCalled()
    })
}

export async function getPlayerId(messageCb: jest.Mock): Promise<string> {
    await waitForExpect(() => {
        expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
            event: 'notification.player-id',
        }))
    })
    for (const call of messageCb.mock.calls) {
        const body: any = call[0]
        if (body.event === 'notification.player-id') {
            return body.payload.playerId
        }
    }

    throw new Error('Could not find playerId');
}