import { describe, expect, it, jest } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import WebSocket from 'ws';
import { connect, disconnect, getPlayerId } from './01-Connect.spec';
import { joinTeam } from './02-Teams.spec';

describe('Disconnect', () => {
    it('should update the player list, with the given player removed', async () => {
        const res1 = await connect()
        const [conn1, messageCb1] = [res1.conn, res1.messageCb];

        const res2 = await connect()
        const [conn2, messageCb2] = [res2.conn, res2.messageCb];
        const playerId2 = getPlayerId(messageCb2)


        // Given player 2 has joined
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith(expect.objectContaining({
                payload: {
                    players: expect.arrayContaining([playerId2])
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
                    players: expect.not.arrayContaining([playerId2])
                }
            }))
        });

        await disconnect(conn1);
    })

    it('should update the team list, with the given player removed', async () => {
        const res1 = await connect()
        const [conn1, messageCb1] = [res1.conn, res1.messageCb];

        const res2 = await connect()
        const [conn2, messageCb2] = [res2.conn, res2.messageCb];
        const playerId2 = getPlayerId(messageCb2)


        // Given player2 has joined team 0
        await joinTeam(conn2, messageCb2, 0);
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith({
                event: "state.teams",
                payload: {
                    "teams": [
                        expect.arrayContaining([playerId2]),
                        expect.anything(),
                    ],
                }
            });
        });
        messageCb1.mockClear();

        // When player 2 disconnects
        await disconnect(conn2);

        // Then player 1 receives an updated team list without player 2
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith({
                event: "state.teams",
                payload: {
                    "teams": [
                        expect.not.arrayContaining([playerId2]),
                        expect.anything(),
                    ],
                }
            });
        });

        await disconnect(conn1);
    })
})
