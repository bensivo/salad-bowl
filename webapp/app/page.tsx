'use client';

import gameStore from '@/services/game-store';
import playerStore from '@/services/player-store';
import wordStore from '@/services/wordbank-store';
import axios from 'axios';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { CharInput } from '../components/char-input';
import { HTTP_URL } from './constants';
import './page.css';
import ws from '@/services/ws';

export default function HomePage() {
    const router  = useRouter();
    const [joinCode, setJoinCode] = useState('');

    useEffect(() => {
        init();
    }, [])

    function init() {
        ws.disconnect(); // Handles case where player navigated here using the back button. We've found desktop browsers don't actually end the websocket in that case.
    }

    const onClickNewGame = async () => {
        // Clear all stores, in case there is leftover state from a previous game
        gameStore.reset();
        playerStore.reset();
        wordStore.reset();

        const res = await axios.request({
            method: 'post',
            url: `${HTTP_URL}/game`,
        })

        // TODO: error notification on failure

        const gameId = res.data.gameId;
        sessionStorage.setItem('gameId', gameId);
        router.push('/game')
    }

    const onClickJoinGame = async () => {
        console.log('Joining Game with code', joinCode);
        if (joinCode.length != 7) {
            console.error('Fill in the entire join code before joining')
            return;
        }

        const res = await axios.request({
            method: 'get',
            url: `${HTTP_URL}/game`,
        });
        if (!res.data[joinCode]) {
            console.error(`Lobby ${joinCode} not found`)
            // TODO: error notification
            return;
        }

        sessionStorage.setItem('gameId', joinCode);
        router.push('/game')
    }

    return (
        <div id="homepage">
            <div id='content-main' className='content-main card'>
                <section id='title'>
                    <h1>Salad Bowl</h1>
                    <h3>by bensivo & dayboobian</h3>
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
