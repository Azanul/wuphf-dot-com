import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface WuphfState {
  messages: string[];
}

const initialState: WuphfState = {
  messages: [],
};

const wuphfSlice = createSlice({
  name: 'wuphf',
  initialState,
  reducers: {
    sendWuphf: (state, action: PayloadAction<string>) => {
      state.messages.push(action.payload);
    },
  },
});

export const { sendWuphf } = wuphfSlice.actions;
export default wuphfSlice.reducer;
