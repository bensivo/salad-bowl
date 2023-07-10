import { describe, expect, it } from '@jest/globals';
import waitForExpect from 'wait-for-expect';
import { joinTeam } from './03-Teams.spec';
import { connect, createLobby, disconnect, getPlayerId } from './actions';

describe('Wordbank', () => {
    it('start game - should notify all players', async () => {
        // Given a lobby with 2 players
        const lobbyId = await createLobby();
        const res1 = await connect(lobbyId);
        const conn1 = res1.conn;
        const messageCb1 = res1.messageCb
        await joinTeam(conn1, messageCb1, 0);

        const res2 = await connect(lobbyId);
        const conn2 = res2.conn;
        const messageCb2 = res2.messageCb
        await joinTeam(conn2, messageCb2, 1);

        // When someone sends the game-start request 
        conn1.send(JSON.stringify({
            event: 'request.start-game',
            payload: {
                "requestId": '00000000-0000-0000-0000-000000000000',
            },
        }))

        // Then they get a response
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith({
                event: "response.start-game",
                payload: {
                    "requestId": "00000000-0000-0000-0000-000000000000",
                    "status": "success",
                }
            });
        }, 5000, 1000);

        // Then everyone gets a notification
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith({
                event: "notification.game-started",
                payload: {}
            });
            expect(messageCb2).toHaveBeenCalledWith({
                event: "notification.game-started",
                payload: {}
            });
        }, 5000, 1000);


        await disconnect(conn1);
        await disconnect(conn2);
    });

    it('add-word - should update all players', async () => {
        // Given a lobby with 2 players, which has started
        const lobbyId = await createLobby();
        const res1 = await connect(lobbyId);
        const conn1 = res1.conn;
        const messageCb1 = res1.messageCb
        await joinTeam(conn1, messageCb1, 0);
        const playerId1 = getPlayerId(messageCb1)

        const res2 = await connect(lobbyId);
        const conn2 = res2.conn;
        const messageCb2 = res2.messageCb
        await joinTeam(conn2, messageCb2, 1);
        const playerId2 = getPlayerId(messageCb2)

        conn1.send(JSON.stringify({
            event: 'request.start-game',
            payload: {
                "requestId": '00000000-0000-0000-0000-000000000000',
            },
        }))

        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith({
                event: "notification.game-started",
                payload: {}
            });
        }, 5000, 1000);

        // When someone sends request to add word
        conn1.send(JSON.stringify({
            event: 'request.add-word',
            payload: {
                requestId: "00000000-0000-0000-0000-000000000000",
                word: 'asdf',
            },
        }))


        // Then they get a response
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith({
                event: "response.add-word",
                payload: {
                    requestId: "00000000-0000-0000-0000-000000000000",
                    status: "success",
                }
            });
        }, 5000, 1000);


        // Then everyone gets a state update
        await waitForExpect(() => {
            expect(messageCb1).toHaveBeenCalledWith({
                event: "state.word-bank",
                payload: {
                    submittedWords: [
                        {
                            word: 'asdf',
                            playerId: playerId1
                        }
                    ]
                }
            });
            expect(messageCb2).toHaveBeenCalledWith({
                event: "state.word-bank",
                payload: {
                    submittedWords: [
                        {
                            word: 'asdf',
                            playerId: playerId1
                        }
                    ]
                }
            });
        }, 5000, 1000);

        await disconnect(conn1);
        await disconnect(conn2);
    });
});
