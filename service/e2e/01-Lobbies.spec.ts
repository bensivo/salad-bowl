import { describe, expect, it, jest } from '@jest/globals';
import axios from 'axios';
import { createLobby } from './actions';


describe('Lobbies', () => {
    it('should create and get lobby via REST', async () => {
        const lobbyId = await createLobby();
        expect(lobbyId).toBeTruthy();
    });

    // NOTE: this may be a good test to move into unit testing, it alone makes our e2e test suite take much longer.
    it('should delete empty lobbies that are more than 30 seconds old', async () => {
        // Given a lobby
        const lobbyId = await createLobby();

        // When I wait 40 seconds (lobbies are up for deletion after 30 seconds, and the deletion job runs every 10 seconds)
        await new Promise((resolve) => {
            setTimeout(resolve, 40 * 1000)
        })


        // Then the lobby no longer exists
        const getRes = await axios.request({
            method: 'GET',
            url: 'http://localhost:8080/lobbies'
        })

        expect(getRes.data[lobbyId]).toBeUndefined();
    }, 60 * 1000);
});
