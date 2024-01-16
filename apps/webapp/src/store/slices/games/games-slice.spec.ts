import { EnhancedStore, configureStore } from '@reduxjs/toolkit';
import axios from 'axios';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { gamesActions, gamesReducer, gamesThunks } from './games-slice';

vi.mock('axios');

describe('games slice', () => {
    let store: EnhancedStore<any, any, any>;
    let stateListener: (a: any) => void;

    beforeEach(() => {
        store = configureStore({
            reducer: {
                games: gamesReducer,
            },
        });

        // Mock function that gets called with every state update
        stateListener = vi.fn();
        store.subscribe(() => {
            stateListener(store.getState())
        });
    })
    // Mock store that only has a games reducer

    afterEach(() => {
        vi.resetAllMocks();
    });


    describe('fetchGames', () => {
        it('should toggle loading appropriately', async () => {
            vi.mocked(axios.request).mockResolvedValue({
                data: []
            });

            await store.dispatch(gamesThunks.fetchGames());
            expect(stateListener).toHaveBeenCalledWith({
                games: {
                    loading: true,
                    games: [],
                    error: null,
                }
            });
            
            expect(stateListener).toHaveBeenCalledWith({
                games: {
                    loading: false,
                    games: [],
                    error: null,
                }
            });
        });
        it('should set HTTP response in state', async () => {
            const games = [
                {
                    id: 'asdf'
                }
            ];
            vi.mocked(axios.request).mockResolvedValue({
                data: games
            });

            await store.dispatch(gamesActions.fetchGames(null));
            expect(stateListener).toHaveBeenCalledWith({
                games: {
                    loading: true,
                    games: [],
                    error: null,
                }
            })

            expect(store.getState().games.loading).toEqual(false);
            expect(store.getState().games.games).toEqual(games);
        });

        it('should set error response if needed', async () => {
            const error = new Error('asdf');
            vi.mocked(axios.request).mockRejectedValue(error);

            await store.dispatch(gamesActions.fetchGames(null));

            // Errors get serialized during handling, so the saved 'error' 
            // instance isn't literally the same one we threw
            expect(store.getState().games.error.name).toEqual(error.name);
            expect(store.getState().games.error.message).toEqual(error.message);
        });
    });
});