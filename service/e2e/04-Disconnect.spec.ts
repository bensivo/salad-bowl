import { beforeAll, describe, expect, it } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import { connect, createLobby, disconnect, getPlayerId } from './actions';

describe('Disconnect', () => {
    let lobbyId: string;
    beforeAll(async () => {
        lobbyId = await createLobby();
    });

    it('should update the player list, with the given player offline', async () => {
        const res1 = await connect(lobbyId)
        const [conn1, messageCb1] = [res1.conn, res1.messageCb];

        const res2 = await connect(lobbyId)
        const [conn2, messageCb2] = [res2.conn, res2.messageCb];
        const playerId2 = getPlayerId(messageCb2)


        // Given player 2 has joined
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith(expect.objectContaining({
                payload: {
                    players: expect.arrayContaining([{
                        id: playerId2,
                        status: 'online',
                    }])
                }
            }))
        });

        messageCb1.mockClear();

        // When player 2 disconnects
        await disconnect(conn2);

        // Then player 1 receives an updated player list without player 2
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith(expect.objectContaining({
                payload: {
                    players: expect.arrayContaining([{
                        id: playerId2,
                        status: 'offline',
                    }])
                }
            }))
        });

        await disconnect(conn1);
    })
})
