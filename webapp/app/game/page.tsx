'use client';
import { useEffect } from 'react';
import ws from '../../services/ws';
import './page.css'

export default function GamePage() {

    useEffect(() => {
        init()
    }, []);

    function init() {
        ws.connect();
    }

    function submitWord() {
        ws.send(JSON.stringify({
            event: 'request.add-word',
            payload: {
                word: 'asdf'
            }
        }))
    }

    return (
        <div id="game">
            <div id='title' className='content-main card'>
                <h1>Word Bank</h1>
                <span>Add 3 words to the word bank</span>
            </div>

            <div className='content-main card'>
                <label>Word: </label>
                <input type="text"></input>
                <button onClick={submitWord}>Submit</button>
            </div>
        </div>
    )
}