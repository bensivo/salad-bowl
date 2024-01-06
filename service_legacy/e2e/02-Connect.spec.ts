import { beforeAll, describe, expect, it } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import { connect, createGame, disconnect, getPlayerId } from './actions';

describe('Connect', () => {
    let gameId: string;

    beforeAll(async () => {
        gameId = await createGame();
    });

    it('should receive a playerId', async () => {
        const { conn, messageCb } = await connect(gameId);

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
        const { conn, messageCb } = await connect(gameId);
        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                event: 'notification.player-id',
                payload: {
                    playerId: expect.anything(),
                }
            }));
        });

        const playerId = await getPlayerId(messageCb);

        await disconnect(conn);

        // When I connect again, giving a playerId
        const res = await connect(gameId, playerId);
        const conn2 = res.conn;
        const messageCb2 = res.messageCb;


        // Then I get the same playerId back
        const playerId2 = await getPlayerId(messageCb2);
        expect(playerId2).toEqual(playerId);

        await disconnect(conn2);
    })

    it('should receive a player list', async () => {
        const gameId = await createGame();
        const { conn, messageCb } = await connect(gameId);

        const playerId = await getPlayerId(messageCb);

        await waitForExpect(() => {
            expect(messageCb).toHaveBeenCalledWith(expect.objectContaining({
                event: 'state.player-list',
                payload: {
                    players: expect.arrayContaining([{
                        id: playerId,
                        status: 'online',
                        team: 0,
                    }])
                }
            }));
        });

        await disconnect(conn);
    })
})
