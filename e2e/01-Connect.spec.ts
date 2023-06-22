import { describe, expect, it, jest, beforeEach, afterEach} from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import WebSocket from 'ws';

describe('Connect', () => {
    let conn: WebSocket
    const openCb = jest.fn();
    const closeCb = jest.fn();
    const messageCb = jest.fn();

    beforeEach(async () => {
        conn = new WebSocket('ws://localhost:8080/connect')
        conn.onopen = openCb;
        conn.onclose = closeCb;
        conn.onmessage = (event) => {
            messageCb(JSON.parse(event.data.toString()))
        };

        await waitForExpect(() => {
            expect(openCb).toHaveBeenCalled()
        })
    })

    afterEach(async () => {
        conn.terminate()

        await waitForExpect(() => {
            expect(closeCb).toHaveBeenCalled()
        })
    })

    it('should receive a playerId', async () => {
        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                "ID": expect.anything()
            }));
        });
    })

    it('should receive a player list', async () => {
        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                "Players": expect.anything()
            }));
        });
    })
})