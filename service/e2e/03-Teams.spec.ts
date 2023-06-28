import { describe, it, jest, expect, beforeAll } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import { WebSocket } from 'ws';
import { connect, createLobby, disconnect, getPlayerId } from './actions';

describe('Teams', () => {
    let lobbyId: string;

    beforeAll(async () => {
        lobbyId = await createLobby();
    });

    describe('join team', () => {
        it('should send a success message', async () => {
            const { conn, messageCb } = await connect(lobbyId);

            await joinTeam(conn, messageCb, 1);

            await disconnect(conn);
        });


        it('should send team list', async () => {
            const { conn, messageCb } = await connect(lobbyId);

            const playerId = getPlayerId(messageCb);

            await joinTeam(conn, messageCb, 1);

            await waitForExpect(() => {
                expect(messageCb).toHaveBeenCalledWith({
                    event: "state.teams",
                    payload: {
                        "teams": [
                            expect.anything(),
                            expect.arrayContaining([playerId]),
                        ],
                    }
                });
            })

            await disconnect(conn);
        });
    });

    describe('switch teams', () => {
        it('should remove me from the old team and put me in the new one', async () => {
            const { conn, messageCb } = await connect(lobbyId);

            const playerId = getPlayerId(messageCb);

            await joinTeam(conn, messageCb, 1);

            await waitForExpect(() => {
                expect(messageCb).toHaveBeenCalledWith({
                    event: "state.teams",
                    payload: {
                        "teams": [
                            expect.anything(),
                            expect.arrayContaining([playerId]),
                        ],
                    }
                });
            })

            await joinTeam(conn, messageCb, 0);

            await waitForExpect(() => {
                expect(messageCb).toHaveBeenCalledWith({
                    event: "state.teams",
                    payload: {
                        "teams": [
                            expect.arrayContaining([playerId]),
                            expect.not.arrayContaining([playerId]),
                        ],
                    }
                });
            })

            await disconnect(conn);
        })
    })
});

export async function joinTeam(conn: WebSocket, messageCb: jest.Mock, num: number) {
    conn.send(JSON.stringify({
        event: "request.join-team",
        payload: {
            "requestId": "00000000-0000-0000-0000-000000000000",
            "team": num,
        }
    }));

    await waitForExpect(() => {
        expect(messageCb).toHaveBeenCalledWith({
            event: "response.join-team",
            payload: {
                "requestId": "00000000-0000-0000-0000-000000000000",
                "status":"success",
                "team": num,
            }
        });
    })
}