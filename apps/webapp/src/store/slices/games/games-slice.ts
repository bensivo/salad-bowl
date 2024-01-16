import axios from "axios";
import { createAppSlice } from "../../utils";
import { Game } from "./games.interface";
import { createAsyncThunk } from "@reduxjs/toolkit";

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

    extraReducers: (builder) => {

        builder.addCase(gamesThunks.fetchGames.pending, (state) => {
            state.loading = true;
        });
        builder.addCase(gamesThunks.fetchGames.rejected, (state, action) => {
            state.loading = false;
            state.error = action.payload ?? action.error;
        });
        builder.addCase(gamesThunks.fetchGames.fulfilled, (state, action) => {
            state.loading = false;
            state.games = action.payload;
        });
    },

    // /**
    //  * Selectors define values that can be read from this slice, using the hook 'useAppSelector'.
    //  * 
    //  * NOTE: selectors defined here can only access data from this slice. It is also possible to create
    //  * selectors with the function 'createSelector()', to access the entire store.
    //  */
    // selectors: {
    // }
});

export const gamesThunks = {
    fetchGames: createAsyncThunk(
        'users/fetchGames',
        async () => {
            const res = await axios.request({
                method: 'GET',
                url: 'http://localhost:8080/games',
            });

            return res.data as any[]; // TODO: interface for games
        }
    )
}

export const gamesActions = gamesSlice.actions;
export const gamesSelectors = gamesSlice.selectors;
export const gamesReducer = gamesSlice.reducer;