import { PayloadAction, createSlice } from "@reduxjs/toolkit";


export interface RouteState {
    route: string;
}

const initialState: RouteState = {
    route: '',
};

export const routeSlice = createSlice({
    name: 'route',
    initialState,
    
    /**
     * Reducers define both actions that a user can dispatch to the store
     * (name and type are inferred from the reducer function definition). And how the state
     * updates in response to those actions.
     */
    reducers: {
        setRoute: (state: RouteState, action: PayloadAction<string>) => {
            state.route = action.payload
        }
    },

    /**
     * Selectors define values that can be read from this slice, using the hook 'useAppSelector'.
     * 
     * NOTE: selectors defined here can only access data from this slice. It is also possible to create
     * selectors with the function 'createSelector()', to access the entire store.
     */
    selectors: {
        route: s => s.route,
    }
});

export const routeActions = routeSlice.actions;
export const routeSelectors = routeSlice.selectors;
export const routeReducer = routeSlice.reducer;