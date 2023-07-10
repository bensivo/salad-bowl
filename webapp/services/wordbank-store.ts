import { createStore, select, withProps } from "@ngneat/elf";
import ws from "./ws";

export interface SubmittedWord {
    word: string;
    playerId: string;
}

export interface WordState {
    submittedWords: SubmittedWord[];
}

export class WordBankStore {
    initialized = false;

    private store = createStore({
            name: 'word'
        }, withProps<WordState>({
            submittedWords: [],
        }));
    
    
    submittedWords$ = this.store.pipe(
        select(s => s.submittedWords)
    );

    init() {
        if(this.initialized) {
            return;
        }
        this.initialized = true;

        ws.messages$.subscribe((msg: any) => {
            switch(msg.event) {
                case 'state.word-bank':
                    this.store.update(s => ({
                        ...s,
                        submittedWords: msg.payload.submittedWords,
                    }))
                    break;
            }
        })
    }

    setWords(words: SubmittedWord[]) {
        this.store.update(s => ({
            submittedWords: words,
        }));
    }

}

const wordStore = new WordBankStore();

export default wordStore