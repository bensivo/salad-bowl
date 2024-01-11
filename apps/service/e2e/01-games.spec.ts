import 'jest';
import axios from 'axios';

describe('games rest api', () => {
    it('create, get, delete, get', async () => {
        // Create a game
        let res = await axios.request({
            method: 'POST',
            url: 'http://localhost:8080/games',
            data: {},
        });
        const gameId = res.data.id;

        // Get all games, find our gameId in the list
        res = await axios.request({
            method: 'GET',
            url: 'http://localhost:8080/games'
        });
        let found = (res.data as any[]).some(g => g.id == gameId);
        expect(found).toEqual(true)

        // Get our game specifically
        res = await axios.request({
            method: 'GET',
            url: `http://localhost:8080/games/${gameId}`
        });
        expect(res.data.id).toEqual(gameId);


        // Delete the game
        await axios.request({
            method: 'DELETE',
            url: `http://localhost:8080/games/${gameId}`,
        });

        // Call get all, game is not in list
        res = await axios.request({
            method: 'GET',
            url: `http://localhost:8080/games`,
        });
        found = (res.data as any[]).some(g => g.id == gameId);
        expect(found).toEqual(false);

        // Call get one, returns 404
        res = await axios.request({
            method: 'GET',
            url: `http://localhost:8080/games/${gameId}`,
            validateStatus: () => true,
        });
        expect(res.status).toEqual(404);
    });
});