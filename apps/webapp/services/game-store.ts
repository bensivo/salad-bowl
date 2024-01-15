import { createStore, select, withProps } from "@ngneat/elf";
import ws from "./ws";

export interface GameState {
    id: string;
    phase: string;
}

// TODO: This store is only updated from the game page. 
// IF the user refreshes on the wordbank page, this store doesn't get any updates. We need to move the state updates from messages into some shared component.
export class GameStore {
    initialized = false;

    private store = createStore({
        name: 'game'
    }, withProps<GameState>({
        id: '',
        phase: '',
    }));

    phase$ = this.store.pipe(
        select(s => s.phase)
    )

    init() {
        if (this.initialized) {
            return;
        }
        this.initialized = true;

        ws.messages$.subscribe((msg: any) => {
            this.store.update(s => ({
                ...s,
                phase: msg.phase
            }));
        });
    }

    reset() {
        this.store.update(s => ({
            id: '',
            phase: '',
        }))
    }
}

const gameStore = new GameStore();

export default gameStore