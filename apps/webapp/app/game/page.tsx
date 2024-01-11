'use client';

import gameStore from "@/services/game-store";
import ws from "@/services/ws";
import { useObservableState } from "observable-hooks";
import { useEffect } from "react";
import WordbankPage from "./wordbank";
import LobbyPage from "./lobby";
import playerStore from "@/services/player-store";
import wordStore from "@/services/wordbank-store";

export default function GamePage() {
    const phase = useObservableState(gameStore.phase$);


    useEffect(() => {
        init()
    }, []);

    function init() {
        console.log('Initializing game page')
        ws.init();
        gameStore.init();
        playerStore.init();
        wordStore.init();
    }


    switch(phase) {
        case 'lobby':
            return (
                <LobbyPage></LobbyPage>
            )
        case 'word-bank':
            return (
                <WordbankPage></WordbankPage>
            )
        case 'lobby':
            return (
                <div>TODO</div>
            )
    }
}