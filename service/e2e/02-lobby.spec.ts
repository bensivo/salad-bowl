import 'jest';
import axios from 'axios';

describe('lobby', () => {
    it('player join and leave', async () => {
        // Create game
        let res = await axios.request({
            method: 'POST',
            url: 'http://localhost:8080/games',
            data: {},
        });
        const gameId = res.data.id;

        // Send player-joined
        await sendEvent(gameId, {
            name: 'player-joined',
            timestamp: new Date().toISOString(),
            payload: {
                playerId: '1111',
                playerName: 'alice',
            }
        });

        // Verfy player is in game
        res = await axios.request({
            method: 'GET',
            url: `http://localhost:8080/games/${gameId}`,
        })
        expect(res.data.players[0].playerId).toEqual('1111')
        expect(res.data.players[0].playerName).toEqual('alice')


        // Send player-left
        await sendEvent(gameId, {
            name: 'player-left',
            timestamp: new Date().toISOString(),
            payload: {
                playerId: '1111',
            }
        })

        // Verify player is removed
        res = await axios.request({
            method: 'GET',
            url: `http://localhost:8080/games/${gameId}`,
        })
        expect(res.data.players).toEqual([])

        // Cleanup
        await axios.request({
            method: 'DELETE',
            url: `http://localhost:8080/games/${gameId}`,
        })
    });

    it('player join team', async () => {
        // Create game
        let res = await axios.request({
            method: 'POST',
            url: 'http://localhost:8080/games',
            data: {},
        });
        const gameId = res.data.id;

        // Send player-joined
        await sendEvent(gameId, {
            name: 'player-joined',
            timestamp: new Date().toISOString(),
            payload: {
                playerId: '1111',
                playerName: 'alice',
            }
        });

        // Verfy player is in game
        res = await axios.request({
            method: 'GET',
            url: `http://localhost:8080/games/${gameId}`,
        })
        expect(res.data.players[0].playerId).toEqual('1111')
        expect(res.data.players[0].playerName).toEqual('alice')


        // Send team-joined
        await sendEvent(gameId, {
            name: 'team-joined',
            timestamp: new Date().toISOString(),
            payload: {
                playerId: '1111',
                teamName: 'Red'
            }
        })

        // Verify player is in the team
        res = await axios.request({
            method: 'GET',
            url: `http://localhost:8080/games/${gameId}`,
        })
        expect(res.data.teams[0].playerIds[0]).toEqual('1111')

        // Send player-left
        await sendEvent(gameId, {
            name: 'player-left',
            timestamp: new Date().toISOString(),
            payload: {
                playerId: '1111',
            }
        })

        // Verify player has been removed from the team
        res = await axios.request({
            method: 'GET',
            url: `http://localhost:8080/games/${gameId}`,
        })
        expect(res.data.teams[0].playerIds).toEqual([])

        // Cleanup
        await axios.request({
            method: 'DELETE',
            url: `http://localhost:8080/games/${gameId}`,
        })
    });
});

async function sendEvent(gameId: string, event: any) {
    await axios.request({
        method: "POST",
        url: `http://localhost:8080/games/${gameId}/event`,
        data: event
    });
}