import { createStore, select, withProps } from "@ngneat/elf";
import { tap } from 'rxjs/operators';
import ws from "./ws";

export interface Player {
    id: string;
    status: 'online' | 'offline';
    team: number;
}
export interface PlayerState {
    myPlayerId: string;
    players: Player[];
    teams: Player[][];
}

// TODO: This store is only updated from the game page. 
// IF the user refreshes on the wordbank page, this store doesn't get any updates. We need to move the state updates from messages into some shared component.
export class PlayerStore {
    initialized = false;

    private store = createStore({
        name: 'message'
    }, withProps<PlayerState>({
        myPlayerId: '',
        players: [],
        teams: [[], []],
    }));


    myPlayerId$ = this.store.pipe(
        select(s => s.myPlayerId),
        tap(s => console.log(s))
    );

    teams$ = this.store.pipe(
        select(s => s.teams)
    );

    players$ = this.store.pipe(
        select(s => s.players)
    );

    init() {
        if (this.initialized) {
            return;
        }
        this.initialized = true;

        ws.messages$.subscribe((msg: any) => {
            switch (msg.event) {
                case 'notification.player-id':
                    sessionStorage.setItem('playerId', msg.payload.playerId);
                    this.store.update(s => ({
                        ...s,
                        myPlayerId: msg.payload.playerId,
                    }));
                    break;
                case 'state.player-list':
                    playerStore.setPlayers(msg.payload.players)
                    break;
            }
        });

    }

    reset() {
        this.store.update(s => ({
            myPlayerId: '',
            players: [],
            teams: [[],[]],
        }))
    }

    setPlayers(players: Player[]) {
        const teams: Player[][] = [[], []];

        for (const player of players) {
            if (player.team == 0 || player.team == 1) {
                teams[player.team].push(player)
            } else {
                console.error(`Cannot put player ${player.id} in team ${player.team}. Only teams 0 and 1 are allowed`)
            }
        }

        this.store.update(s => ({
            ...s,
            players,
            teams,
        }));
    }
}

const playerStore = new PlayerStore();

export default playerStore