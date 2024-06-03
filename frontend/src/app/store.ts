import { configureStore } from '@reduxjs/toolkit';
import wuphfReducer from '../features/wuphf/wuphfSlice';

export const store = configureStore({
  reducer: {
    wuphf: wuphfReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
