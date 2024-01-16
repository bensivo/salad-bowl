import axios from "axios";
import { createAppSlice } from "../../utils";
import { Game } from "./games.interface";

export interface GamesState {
    loading: boolean;
    games: Game[];
    error: any;
}

const initialState: GamesState = {
    loading: false,
    games: [],
    error: null,
};

export const gamesSlice = createAppSlice({
    name: 'games',
    initialState,

    /**
     * Reducers define both actions that a user can dispatch to the store
     * (name and type are inferred from the reducer function definition). And how the state
     * updates in response to those actions.
     */
    reducers: (create) => ({
        fetchGames: create.asyncThunk(
            async () => { // TODO: fix typing once this MR is released: https://github.com/reduxjs/redux-toolkit/issues/4060
                
                await new Promise((resolve) => setTimeout(resolve, Math.random() * 500));
                const res = await axios.request({
                    method: 'GET',
                    url: 'http://localhost:8080/games',
                });

                return res.data as Game[];
            },
            {
                pending: (state) => {
                    state.loading = true;
                },
                rejected: (state, action) => {
                    state.error = action.payload ?? action.error;
                },
                fulfilled: (state, action) => {
                    state.games = action.payload;
                },
                settled: (state) => {
                    state.loading = false;
                },
            }
        ),
    }),

    // /**
    //  * Selectors define values that can be read from this slice, using the hook 'useAppSelector'.
    //  * 
    //  * NOTE: selectors defined here can only access data from this slice. It is also possible to create
    //  * selectors with the function 'createSelector()', to access the entire store.
    //  */
    selectors: {
        games: s => s.games,
    }
});

export const gamesActions = gamesSlice.actions;
export const gamesSelectors = gamesSlice.selectors;
export const gamesReducer = gamesSlice.reducer;