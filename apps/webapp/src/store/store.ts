import { configureStore } from '@reduxjs/toolkit';
import { gamesReducer } from './slices/games/games-slice';
import { routeReducer } from './slices/route/route-slice';


export const store = configureStore({
    reducer: { 
        route: routeReducer,
        games: gamesReducer,
    },
    devTools: true,
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch