import 'jest';
import axios from 'axios';

describe('games rest api', () => {
    describe('POST /games', () => {
        it('creates a new game', async () => {
            const res = await axios.request({
                method: 'POST',
                url: 'http://localhost:8080/games',
                data: {},
            });

            expect(res.data.id).toBeTruthy();
        });
    });
    describe('GET /games', () => {
        it('returns previously created games', async () => {
            // Given there we post a game
            const postRes = await axios.request({
                method: 'POST',
                url: 'http://localhost:8080/games',
                data: {},
            });

            // When we call get
            const getRes = await axios.request({
                method: 'GET',
                url: 'http://localhost:8080/games'
            });

            // Then the response includes our posted game from earlier
            const found = (getRes.data as any[]).some(g => g.id == postRes.data.id);
            expect(found).toEqual(true)
        });
    });
    describe('GET /games/{id}', () => {
        it('returns the game', async () => {
            // Given there we post a game
            const postRes = await axios.request({
                method: 'POST',
                url: 'http://localhost:8080/games',
                data: {},
            });

            // When we call get
            const getRes = await axios.request({
                method: 'GET',
                url: `http://localhost:8080/games/${postRes.data.id}`
            });

            expect(getRes.data.id).toEqual(postRes.data.id);
        })

        it('returns 404 on not found', async () => {
            // When we call get, without a previous post
            const getRes = await axios.request({
                method: 'GET',
                url: `http://localhost:8080/games/1111`,
                validateStatus: () => true,
            });

            expect(getRes.status).toEqual(404);
        })
    });

    describe('DELETE /games/{id}', () => {

        it('removes the game from the GET /games endpoint', async () => {
            // Given there we post a game
            const postRes = await axios.request({
                method: 'POST',
                url: 'http://localhost:8080/games',
                data: {},
            });

            // When we call get 
            let getRes = await axios.request({
                method: 'GET',
                url: `http://localhost:8080/games`,
            });

            // Then the game is there
            let found = (getRes.data as any[]).some(g => g.id == postRes.data.id);
            expect(found).toEqual(true);

            // Given we call delete
            await axios.request({
                method: 'DELETE',
                url: `http://localhost:8080/games/${postRes.data.id}`,
            });

            // When we call get 
            getRes = await axios.request({
                method: 'GET',
                url: `http://localhost:8080/games`,
            });

            // Then the game is not there
            found = (getRes.data as any[]).some(g => g.id == postRes.data.id);
            expect(found).toEqual(false);
        })
    });
});