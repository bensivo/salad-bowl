import { Store, createStore, select, setProp, withProps } from "@ngneat/elf";
import { map } from 'rxjs';

export interface Player {
    id: string;
    status: 'online' | 'offline';
}
export interface PlayerState {
    players: Player[];
    teams: Player[][];
}

export class PlayerStore {
    private store = createStore({
            name: 'message'
        }, withProps<PlayerState>({
           players: [],
           teams: [],
        }));

    players$ = this.store.pipe(
        map(s => s.players)
    );

    teams$ = this.store.pipe(
        select(s => s.teams)
    );

    setPlayers(players: Player[]) {
        console.log('Setting players', players)
        this.store.update(s => ({
            ...s,
            players,
        }))

    }

    setTeams(teams: Player[][]) {
        this.store.update(setProp('teams', teams));
    }
}

const playerStore = new PlayerStore();

export default playerStore