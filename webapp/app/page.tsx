'use client';

import { useState } from 'react';
import { CharInput } from '../components/char-input';
import { useRouter } from 'next/navigation';
import axios from 'axios';
import './page.css';
import { HTTP_URL } from './constants';

export default function HomePage() {
    const router  = useRouter();
    const [joinCode, setJoinCode] = useState('');

    const onClickNewGame = async () => {
        const res = await axios.request({
            method: 'post',
            url: `${HTTP_URL}/lobbies`,
        })

        // TODO: error notification on failure

        const lobbyId = res.data.lobbyId;
        sessionStorage.setItem('lobbyId', lobbyId);
        router.push('/lobby')
    }

    const onClickJoinGame = async () => {
        console.log('Joining Game with code', joinCode);
        if (joinCode.length != 7) {
            console.error('Fill in the entire join code before joining')
            return;
        }

        const res = await axios.request({
            method: 'get',
            url: `${HTTP_URL}/lobbies`,
        });
        if (!res.data[joinCode]) {
            console.error(`Lobby ${joinCode} not found`)
            // TODO: error notification
            return;
        }

        sessionStorage.setItem('lobbyId', joinCode);
        router.push('/lobby')
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
