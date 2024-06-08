import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

type chat = { chatId: string, messages: any[] };

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

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

export const fetchChats = createAsyncThunk('wuphf/fetchChats', async () => {
  let retries = 3;
  let retryDelay = 2000;

  for (let attempt = 0; attempt < retries; attempt++) {
    const response = await fetch(`${process.env.REACT_APP_BASE_URL}/history?userId=${localStorage.getItem('user_id')}`,
      {
        headers: {
          'Authorization': localStorage.getItem('token') || '',
        }
      });

    if (response.ok) {
      return await response.json();
    }

    if (response.status === 404) {
      if (attempt < retries - 1) {
        await delay(retryDelay);
      } else {
        throw new Error('Failed to fetch chats after multiple attempts');
      }
    } else {
      throw new Error('Failed to fetch chats');
    }
  }
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
  const messages = await response.json();
  return { chatId, messages };
});

const wuphfSlice = createSlice({
  name: 'wuphf',
  initialState,
  reducers: {
    sendWuphf: (state, action: PayloadAction<{ chatId: string; message: string }>) => {
      const { chatId, message } = action.payload;
      state.chats.find((chat) => chat.chatId === chatId)?.messages.push({ sender: localStorage.getItem('user_id'), msg: message });
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
        const updatedChat = action.payload;
        const chatIndex = state.chats.findIndex(chat => chat.chatId === updatedChat.chatId);
        if (chatIndex >= 0) {
          state.chats[chatIndex] = updatedChat;
        } else {
          state.chats.push(updatedChat);
        }
      })
      .addCase(fetchMessages.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch messages';
      });
  },
});

export const { sendWuphf } = wuphfSlice.actions;
export default wuphfSlice.reducer;
