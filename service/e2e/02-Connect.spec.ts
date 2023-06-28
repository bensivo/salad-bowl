import { beforeAll, describe, expect, it } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import { connect, createLobby, disconnect } from './actions';

describe('Connect', () => {
    let lobbyId: string;

    beforeAll(async () => {
        lobbyId = await createLobby();
    });


    it('should receive a playerId', async () => {
        const { conn, messageCb } = await connect(lobbyId);

        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                event: 'notification.player-id',
                payload: {
                    playerId: expect.anything(),
                }
            }));
        });

        await disconnect(conn);
    })

    it('should receive a player list', async () => {
        const lobbyId = await createLobby();
        const { conn, messageCb } = await connect(lobbyId);

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
