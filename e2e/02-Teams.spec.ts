import { describe, it } from '@jest/globals';
import { connect, disconnect } from './01-Connect.spec';

describe('Teams', () => {
    it('should return the team list when I join a team', async () => {
        const { conn, messageCb } = await connect();

        conn.send(JSON.stringify({
            event: "request.join-team",
            payload: {
                "team": "1",
            }
        }));

        await disconnect(conn);
    })
})
