import { asyncThunkCreator, buildCreateSlice } from "@reduxjs/toolkit";

/**
 * Utility function used in place of createSlice(),
 * which allows reducers to be defined using a creator.
 * 
 * This lets us create and handle thunk reducers directly within createSlice, instead of as separate functions
 * See: https://redux-toolkit.js.org/usage/migrating-rtk-2#createslicereducers-callback-syntax-and-thunk-support
 */
export const createAppSlice = buildCreateSlice({
    creators: {
        asyncThunk: asyncThunkCreator
    }
})