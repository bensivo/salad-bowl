'use client'
import * as uuid from 'uuid';
import { useObservableState } from "observable-hooks"
import { useEffect, useState } from "react"
import playerStore from "../../services/player-store"
import { useRouter } from 'next/navigation';
import ws from '../../services/ws';

import './page.css';

export default function LobbyPage() {
    const router = useRouter();
    const [gameId, setLobbyId] = useState<string | null>('');
    const teams = useObservableState(playerStore.teams$, []);
    const myPlayerId = useObservableState(playerStore.myPlayerId$, '')

    useEffect(() => {
        init();
    }, []) // passing an empty array in the second arg makes this effect only run once

    function init() {
        const gameId = sessionStorage.getItem('gameId')
        setLobbyId(gameId);

        ws.init();
        playerStore.init();

        ws.messages$.subscribe((msg: any) => {
            switch (msg.event) {
                case 'state.game-phase':
                    router.push('/game') // TODO: put all the instance-specific pages on the same sub-page. Prevent routing
            }
        });
    }

    function joinTeam(i: number) {
        ws.send(JSON.stringify({
            event: 'request.join-team',
            payload: {
                "requestId": uuid.v4(),
                "team": i
            }
        }))
    }

    function startGame() {
        ws.send(JSON.stringify({
            event: 'request.start-game',
            payload: {
                "requestId": uuid.v4(),
            }
        }))
    }

    return (
        <div id='game'>

            <div id='title' className='content-main card'>
                <h1>Lobby</h1>
                <h3>Join Code: {gameId}</h3>

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
                                {team.map(player => (
                                    <li className={`${player.status == 'offline' ? 'player-offline' : ''}`} key={player.id}>{player.id === myPlayerId ? `${player.id} (me)` : player.id}</li>
                                ))}
                            </ul>
                        </div>
                    ))
                }
            </div>

            <div className='content-main card'>
                <button onClick={startGame}> Start Game </button>
            </div>
        </div>
    )
}