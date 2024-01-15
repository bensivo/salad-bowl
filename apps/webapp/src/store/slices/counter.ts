import { PayloadAction, createSlice } from "@reduxjs/toolkit";


export interface CounterState {
    count: number;
}

const initialState: CounterState = {
    count: 0,
};

export const counterSlice = createSlice({
    name: 'counter',
    initialState,
    
    /**
     * Reducers define both actions that a user can dispatch to the store
     * (name and type are inferred from the reducer function definition). And how the state
     * updates in response to those actions.
     */
    reducers: {
        incrementBy: (state: CounterState, action: PayloadAction<number>) => {
            state.count += action.payload
        }
    },

    /**
     * Selectors define values that can be read from this slice, using the hook 'useAppSelector'.
     * 
     * NOTE: selectors defined here can only access data from this slice. It is also possible to create
     * selectors with the function 'createSelector()', to access the entire store.
     */
    selectors: {
        count: (state) => state.count,
    }
});

export const counterActions = counterSlice.actions;
export const counterSelectors = counterSlice.selectors;
export const counterReducer = counterSlice.reducer;