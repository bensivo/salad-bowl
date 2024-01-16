import { configureStore } from '@reduxjs/toolkit';
import axios from 'axios';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { gamesActions, gamesReducer } from './games-slice';

vi.mock('axios');

describe('games slice', () => {

    const createMockStore = () => {
        const store = configureStore({
            reducer: {
                games: gamesReducer,
            },
        });

        const listener = vi.fn(); // Mock function that gets called with every state update
        store.subscribe(() => {
            listener(store.getState())
        });

        return {
            store,
            listener,
        }
    }

    afterEach(() => {
        vi.resetAllMocks();
    });

    describe('fetchGames', () => {
        it('should toggle loading appropriately', async () => {
            const { store, listener } = createMockStore();
            vi.mocked(axios.request).mockResolvedValue({
                data: []
            });

            await store.dispatch(gamesActions.fetchGames(null));
            expect(listener).toHaveBeenCalledWith({
                games: {
                    loading: true,
                    games: [],
                    error: null,
                }
            });

            expect(listener).toHaveBeenCalledWith({
                games: {
                    loading: false,
                    games: [],
                    error: null,
                }
            });
        });
        it('should set HTTP response in state', async () => {
            const { store, listener } = createMockStore();
            const games = [
                {
                    id: 'asdf'
                }
            ];
            vi.mocked(axios.request).mockResolvedValue({
                data: games
            });

            await store.dispatch(gamesActions.fetchGames(null));
            expect(listener).toHaveBeenCalledWith({
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
            const { store } = createMockStore();

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
