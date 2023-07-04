import { beforeAll, describe, expect, it } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import { connect, createLobby, disconnect, getPlayerId } from './actions';

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

    it('should use an existing playerId if given', async () => {
        // Given a past connection that was then disconnected
        const { conn, messageCb } = await connect(lobbyId);
        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                event: 'notification.player-id',
                payload: {
                    playerId: expect.anything(),
                }
            }));
        });

        const playerId = getPlayerId(messageCb);

        await disconnect(conn);

        // When I connect again, giving a playerId
        const res = await connect(lobbyId, playerId);
        const conn2 = res.conn;
        const messageCb2 = res.messageCb;


        // Then I get the same playerId back
        const playerId2 = getPlayerId(messageCb2);
        expect(playerId2).toEqual(playerId);

        await disconnect(conn2);
    })

    it('should receive a player list', async () => {
        const lobbyId = await createLobby();
        const { conn, messageCb } = await connect(lobbyId);

        const playerId = getPlayerId(messageCb);

        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                event: 'state.player-list',
                payload: {
                    players: expect.arrayContaining([{
                        id: playerId,
                        status: 'online'
                    }])
                }
            }));
        });

        await disconnect(conn);
    })
})
