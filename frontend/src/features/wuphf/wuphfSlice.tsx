import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

type chat = { chatId: string, messages: string[] };

interface WuphfState {
  chats: chat[];
  loading: boolean;
  error: string | null;
}

const initialState: WuphfState = {
  chats: [],
  loading: false,
  error: null,
};

export const fetchChats = createAsyncThunk('wuphf/fetchChats', async () => {
  const response = await fetch(`${process.env.REACT_APP_BASE_URL}/history?userId=${localStorage.getItem('user_id')}`, {
    headers: {
      'Authorization': localStorage.getItem('token') || '',
    },
  });
  if (!response.ok) {
    throw new Error('Failed to fetch chats');
  }
  return await response.json();
});

export const fetchMessages = createAsyncThunk('wuphf/fetchMessages', async (chatId: string) => {
  const response = await fetch(`${process.env.REACT_APP_BASE_URL}/history?chatId=${chatId}`, {
    headers: {
      'Authorization': localStorage.getItem('token') || '',
    },
  });
  if (!response.ok) {
    throw new Error('Failed to fetch chats');
  }
  return await response.json();
});

const wuphfSlice = createSlice({
  name: 'wuphf',
  initialState,
  reducers: {
    sendWuphf: (state, action: PayloadAction<{ chatId: string; message: string }>) => {
      const { chatId, message } = action.payload;
      const chat = state.chats.find((chat) => chat.chatId === chatId);
      if (chat) {
        chat.messages.push(message);
      }
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchChats.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchChats.fulfilled, (state, action) => {
        state.loading = false;
        state.chats = action.payload;
      })
      .addCase(fetchChats.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch chats';
      })
      .addCase(fetchMessages.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchMessages.fulfilled, (state, action) => {
        state.loading = false;
        state.chats = action.payload;
      })
      .addCase(fetchMessages.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch messages';
      });
  },
});

export const { sendWuphf } = wuphfSlice.actions;
export default wuphfSlice.reducer;
