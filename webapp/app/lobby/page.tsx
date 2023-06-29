'use client'
import * as uuid from 'uuid';
import { useObservableState } from "observable-hooks"
import { useEffect, useState } from "react"
import playerStore from "../store/player-store"
import { useRouter } from 'next/navigation';

import './page.css';

var conn: WebSocket;

export default function LobbyPage() {
    const router  = useRouter();
    const [lobbyId, setLobbyId] = useState<string | null>('');
    const [myPlayerId, setMyPlayerId] = useState<string>('');

    useEffect(() => {
        connect();
    }, []) // passing an empty array in the second arg makes this effect only run once

    const players = useObservableState(playerStore.players$, []);
    const teams = useObservableState(playerStore.teams$, []);

    function connect() {
        const lobbyId = sessionStorage.getItem('lobbyId') // find out how to keep this from failing in the next server
        if (!lobbyId) {
            console.error('Could not fetch lobbyId from session storage. Navigating back to home page');
            // TODO: show an error message to the user once they get back to the homepage.
            router.push('/');
        }
        setLobbyId(lobbyId)

        conn = new WebSocket(`ws://api.saladbowl.bensivo.com/lobbies/${lobbyId}/connect`)
        conn.onmessage = (e) => {
            const msg = JSON.parse(e.data);
            console.log('Received message', msg)

            switch (msg.event) {
                case 'notification.player-id':
                    setMyPlayerId(msg.payload.playerId)
                    break;
                case 'state.player-list':
                    playerStore.setPlayers(msg.payload.players)
                    break;
                case 'state.teams':
                    playerStore.setTeams(msg.payload.teams)
                    break;
            }
        }
    }

    function joinTeam(i: number) {
        conn.send(JSON.stringify({
            event: 'request.join-team',
            payload: {
                "requestId": uuid.v4(),
                "team": i
            }
        }))
    }

    return (
        <div id='lobby'>

            <div id='title' className='content-main card'>
                <h1>Lobby</h1>
                <h3>Join Code: {lobbyId}</h3>

                <label>Player ID: {myPlayerId}</label>
            </div>

            <div>
                {
                    teams.map((team, i) => (
                        <div key={i} className='content-main card'>
                            <div className="card-header">
                                <h4>Team {i}</h4>
                                <button onClick={() => joinTeam(i)}> Join Team </button>
                            </div>
                            <ul>
                                {team.map(playerId => (
                                    <li key={playerId}>{playerId === myPlayerId ? `${playerId} (me)` : playerId}</li>
                                ))}
                            </ul>
                        </div>
                    ))
                }
            </div>
        </div>
    )
}