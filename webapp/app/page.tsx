'use client';

import { useState } from 'react';
import { CharInput } from '../components/char-input';
import axios from 'axios';
import './page.css';

export default function HomePage() {

    const [joinCode, setJoinCode] = useState('');

    const onClickNewGame = async () => {
        const res = await axios.request({
            method: 'post',
            url: 'http://localhost:8080/lobbies',
        })

        // TODO: error notification on failure

        const lobbyId = res.data.lobbyId;
        sessionStorage.setItem('lobbyId', lobbyId);
        window.location.href = '/lobby'
    }

    const onClickJoinGame = async () => {
        console.log('Joining Game with code', joinCode);
        if (joinCode.length != 7) {
            console.error('Fill in the entire join code before joining')
            return;
        }

        const res = await axios.request({
            method: 'get',
            url: 'http://localhost:8080/lobbies',
        });
        if (!res.data[joinCode]) {
            console.error(`Lobby ${joinCode} not found`)
            // TODO: error notification
            return;
        }

        sessionStorage.setItem('lobbyId', joinCode);
        window.location.href = '/lobby'
    }

    return (
        <div id="homepage">
            <div id='content-main' className='content-main card'>
                <section id='title'>
                    <h1>Salad Bowl</h1>
                    <h3>by bensivo</h3>
                </section>

                <section id='menu'>
                    <button onClick={onClickNewGame}>New Game</button>

                    <div id='char-input-container'>
                        <CharInput template="xxx-xxx" onChange={setJoinCode}></CharInput>
                    </div>
                    <button onClick={onClickJoinGame}>Join Game</button>
                </section>
            </div>
        </div>
    )
}
