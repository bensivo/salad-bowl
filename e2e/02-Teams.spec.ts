import { describe, it, expect } from '@jest/globals';
import { connect, disconnect, getPlayerId } from './01-Connect.spec';
import waitForExpect from 'wait-for-expect';

describe('Teams', () => {
    describe('join team', () => {
        it('should send a success message', async () => {
            const { conn, messageCb } = await connect();

            conn.send(JSON.stringify({
                event: "request.join-team",
                payload: {
                    "requestId": "00000000-0000-0000-0000-000000000000",
                    "team": 1,
                }
            }));

            await waitForExpect(() => {
                expect(messageCb).toHaveBeenCalledWith({
                    event: "response.join-team",
                    payload: {
                        "requestId": "00000000-0000-0000-0000-000000000000",
                        "status":"success",
                        "team": 1,
                    }
                });
            })

            await disconnect(conn);
        });


        it('should send team list', async () => {
            const { conn, messageCb } = await connect();

            const playerId = getPlayerId(messageCb);

            conn.send(JSON.stringify({
                event: "request.join-team",
                payload: {
                    "requestId": "00000000-0000-0000-0000-000000000000",
                    "team": 1,
                }
            }));

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
            const { conn, messageCb } = await connect();

            const playerId = getPlayerId(messageCb);

            conn.send(JSON.stringify({
                event: "request.join-team",
                payload: {
                    "requestId": "00000000-0000-0000-0000-000000000000",
                    "team": 1,
                }
            }));

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


            conn.send(JSON.stringify({
                event: "request.join-team",
                payload: {
                    "requestId": "00000000-0000-0000-0000-000000000000",
                    "team": 0,
                }
            }));

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
