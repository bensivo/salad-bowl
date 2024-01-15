'use client';
import playerStore, { Player } from '@/services/player-store';
import { useObservableState } from 'observable-hooks';
import { useEffect, useState } from 'react';
import { Observable } from 'rxjs';
import { combineLatestWith, map, tap } from 'rxjs/operators';
import * as uuid from 'uuid';
import wordStore from '../../services/wordbank-store';
import ws from '../../services/ws';
import gameStore from '@/services/game-store';

import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './wordbank.css';

/**
 * Emits an object mapping playerIds to the words they have submitted.
 * 
 * Example:
 * If there are 2 players - 000-000 and 111-111, and 000-000 has submitted the word 'asdf'
 * {
 *  '000-000': ['asdf'],
 *  '111-111': [],
 * }
 */
const playerWords$: Observable<Record<string, string[]>> = wordStore.submittedWords$.pipe(
    combineLatestWith(playerStore.players$),
    map(([submittedWords, players]) => {
        const playerIdToPlayer: Record<string, Player> = {}
        for (const player of players) {
            playerIdToPlayer[player.id] = player;
        }

        const playerIdToWords: Record<string, string[]> = {};
        for (const playerId of Object.keys(playerIdToPlayer)) {
            playerIdToWords[playerId] = [];
        }

        for (const submittedWord of submittedWords) {
            const player = playerIdToPlayer[submittedWord.playerId];
            if (!player) {
                throw new Error(`Word ${submittedWord.word} submitted by unknown player ${submittedWord.playerId}`)
            }

            playerIdToWords[submittedWord.playerId].push(submittedWord.word);
        }

        return playerIdToWords
    }),
    tap(pw => console.log(pw))
)

/**
 * Emits the words this player has submitted
 */
const myPlayerWords$: Observable<string[]> = playerWords$.pipe(
    combineLatestWith(playerStore.myPlayerId$),
    map(([playerWords, myPlayerId]) => {
        return playerWords[myPlayerId]
    })
);

export default function WordbankPage() {
    const playerWords = useObservableState(playerWords$);
    const myPlayerWords = useObservableState(myPlayerWords$);

    const [wordInput, setWordInput] = useState('');
    const [initialized, setInitialized] = useState(false);

    useEffect(() => {
        console.log('Initializing wordbank page')
        ws.init();
        gameStore.init();
        playerStore.init();
        wordStore.init();


        const subscription = ws.messages$.subscribe((msg: any) => {
            switch (msg.event) {
                case 'response.add-word':
                    if (msg.payload.status !== 'success') {
                        toast(`Error: ${msg.payload.description}`, {
                            position: 'bottom-center',
                            autoClose: 60000,
                            hideProgressBar: true,
                            closeOnClick: true,
                            pauseOnHover: true,
                            draggable: false,
                            theme: 'light',
                        })
                    }
                    break;
            }
        });

        return () => {
            subscription.unsubscribe();
        }
    }, []);

    function submitWord() {
        if (wordInput.length == 0) {
            return;
        }
        ws.send(JSON.stringify({
            event: 'request.add-word',
            payload: {
                requestId: uuid.v4(),
                word: wordInput,
            }
        }))
        setWordInput('');
    }

    return (
        <div id="game">
            <div id='title' className='content-main card'>
                <h1>Word Bank</h1>
                <span>Add 3 words to the word bank</span>
            </div>

            <div id="words" className='content-main card'>
                <h3>Submitted Words:</h3>
                {Object.entries(playerWords ?? {}).map(entry => {
                    const playerId = entry[0];
                    const words = entry[1];
                    return (
                        <div key={playerId} className="player-words">
                            <label>{playerId}: </label>
                            {words.map((w, index) => (
                                <span className="submitted-word" key={index}>✅</span>
                            ))}
                        </div>
                    )
                })}
            </div>

            {
                myPlayerWords && myPlayerWords.length < 3 ?
                    (
                        <div className='content-main card'>
                            <label>Word: </label>
                            <input type="text" value={wordInput} onChange={e => setWordInput(e.target.value)}></input>
                            <button className="btn-primary" onClick={submitWord}>Submit</button>
                        </div>
                    )
                    :
                    (
                        <div className='content-main card'>
                            <span>All words submitted. Wait for everyone to finish.</span>
                        </div>
                    )
            }

            <ToastContainer className="toast" />
        </div>
    )
}