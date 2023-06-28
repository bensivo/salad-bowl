'use client'
import * as uuid from 'uuid';
import { useObservableState } from "observable-hooks"
import { useEffect, useState } from "react"
import playerStore from "../store/player-store"

import './page.css';

var conn: WebSocket;

export default function LobbyPage() {
    const lobbyId = sessionStorage.getItem('lobbyId')
    if (!lobbyId) {
        console.error('Could not fetch lobbyId from session storage. Navigating back to home page');
        // TODO: show an error message to the user once they get back to the homepage.
        window.location.href = '/';
        return;
    }

    const [myPlayerId, setMyPlayerId] = useState<string>('');

    useEffect(() => {
        connect();
    }, []) // passing an empty array in the second arg makes this effect only run once

    const players = useObservableState(playerStore.players$, []);
    const teams = useObservableState(playerStore.teams$, []);

    function connect() {
        conn = new WebSocket(`ws://localhost:8080/lobbies/${lobbyId}/connect`)
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
            <h1>Lobby: {lobbyId}</h1>
            <h3>Player ID: {myPlayerId}</h3>
            <div>
                <h3>Players</h3>
                {players.map((playerId, i) => (
                    <div key={i}>ID: {(playerId == myPlayerId) ? `${playerId} (me)` : playerId}</div>
                ))}
            </div>

            <div>
                <h3>Teams</h3>
                {
                    teams.map((team, i) => (
                        <div key={i}>
                            <h4>Team {i}</h4>
                            <button onClick={() => joinTeam(i)}> Join Team </button>
                            {team.map(playerId => (
                                <span key={playerId}>{playerId},</span>
                            ))
                            }
                        </div>
                    ))
                }
            </div>
        </div>
    )
}