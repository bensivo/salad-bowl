'use client'
import { useRouter } from 'next/navigation';
import { useObservableState } from "observable-hooks";
import { useEffect, useState } from "react";
import * as uuid from 'uuid';
import playerStore from "../../services/player-store";
import ws from '../../services/ws';

import gameStore from '@/services/game-store';
import wordStore from '@/services/wordbank-store';
import axios from 'axios';
import { HTTP_URL } from '../constants';
import './lobby.css';

export default function LobbyPage() {
    const router = useRouter();
    const [gameId, setLobbyId] = useState<string | null>('');
    const teams = useObservableState(playerStore.teams$, []);
    const myPlayerId = useObservableState(playerStore.myPlayerId$, '')

    useEffect(() => {
        init();
    }, []) // passing an empty array in the second arg makes this effect only run once

    function init() {
        console.log('Initializing lobbby page')
        const gameId = sessionStorage.getItem('gameId')
        setLobbyId(gameId);

        ws.init();
        gameStore.init();
        playerStore.init();
        wordStore.init();

        ws.messages$.subscribe((msg: any) => {
            if (msg.phase !== 'lobby') {
                router.push('/game') // TODO: put all the instance-specific pages on the same sub-page. Prevent routing
            }
        });
    }

    async function joinTeam(i: number) {
        await axios.request({
            method: 'POST',
            url: `${HTTP_URL}/games/${gameId}/event`,
            data: {
                name: 'player-joined',
                timestamp: new Date(Date.now()).toISOString(),
                payload: {
                    playerId: '1111',
                    playerName: 'alice'
                }
            },
        })

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
        <div id='lobby'>
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