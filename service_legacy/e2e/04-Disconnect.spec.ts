import { beforeAll, describe, expect, it } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import { connect, createGame, disconnect, getPlayerId } from './actions';

describe('Disconnect', () => {
    let gameId: string;
    beforeAll(async () => {
        gameId = await createGame();
    });

    it('should update the player list, with the given player offline', async () => {
        const res1 = await connect(gameId)
        const [conn1, messageCb1] = [res1.conn, res1.messageCb];

        const res2 = await connect(gameId)
        const [conn2, messageCb2] = [res2.conn, res2.messageCb];
        const playerId2 = await getPlayerId(messageCb2)


        // Given player 2 has joined
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith(expect.objectContaining({
                payload: {
                    players: expect.arrayContaining([{
                        id: playerId2,
                        status: 'online',
                        team: 0,
                    }])
                }
            }))
        });

        messageCb1.mockClear();

        // When player 2 disconnects
        await disconnect(conn2);

        // Then player 1 receives an updated player list with player 2 set offline
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith(expect.objectContaining({
                payload: {
                    players: expect.arrayContaining([{
                        id: playerId2,
                        status: 'offline',
                        team: 0,
                    }])
                }
            }))
        });

        await disconnect(conn1);
    })
})
