'use client';

import { useState } from 'react';
import { CharInput } from '../components/char-input';
import './page.css';

export default function HomePage() {

    const [joinCode, setJoinCode] = useState('');

    const onClickNewGame = () => {
        // TODO: make a request to the server to make sure the lobby can be created
        window.location.href = '/lobby' // TODO: add the lobby's url here
    }

    const onClickJoinGame = () => {
        console.log('Joining Game with code', joinCode);
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
