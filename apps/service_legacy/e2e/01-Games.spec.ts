import { describe, expect, it, jest } from '@jest/globals';
import axios from 'axios';
import { createGame } from './actions';


describe('Games', () => {
    it('should create and get game via REST', async () => {
        const gameId = await createGame();
        expect(gameId).toBeTruthy();
    });

    // NOTE: this may be a good test to move into unit testing, it alone makes our e2e test suite take much longer.
    it('should delete empty games that are more than 30 seconds old', async () => {
        // Given a game
        const gameId = await createGame();

        // When I wait 40 seconds (games are up for deletion after 30 seconds, and the deletion job runs every 10 seconds)
        await new Promise((resolve) => {
            setTimeout(resolve, 40 * 1000)
        })


        // Then the game no longer exists
        const getRes = await axios.request({
            method: 'GET',
            url: 'http://localhost:8080/game'
        })

        expect(getRes.data[gameId]).toBeUndefined();
    }, 60 * 1000);
});
