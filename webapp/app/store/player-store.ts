import { Store, createStore, select, setProp, withProps } from "@ngneat/elf";
import { map } from 'rxjs';

export interface PlayerState {
    players: string[];
    teams: string[][];
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

    setPlayers(players: string[]) {
        console.log('Setting players', players)
        this.store.update(s => ({
            ...s,
            players,
        }))

    }

    setTeams(teams: string[][]) {
        this.store.update(setProp('teams', teams));
    }
}

const playerStore = new PlayerStore();

export default playerStore