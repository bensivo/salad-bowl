import { createStore, select, withProps } from "@ngneat/elf";

export interface Player {
    id: string;
    status: 'online' | 'offline';
    team: number;
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
           teams: [[], []],
        }));

    teams$ = this.store.pipe(
        select(s => s.teams)
    );

    setPlayers(players: Player[]) {
        const teams: Player[][] = [[],[]];

        for(const player of players) {
            if (player.team == 0 || player.team == 1) {
                teams[player.team].push(player)
            } else {
                console.error(`Cannot put player ${player.id} in team ${player.team}. Only teams 0 and 1 are allowed`)
            }
        }

        this.store.update(s => ({
            players,
            teams,
        }))

    }
}

const playerStore = new PlayerStore();

export default playerStore