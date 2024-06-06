import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchChats } from '../features/wuphf/wuphfSlice';
import ChatItem from './ChatItem';
import { RootState, AppDispatch } from '../app/store';

interface ChatListProps {
  onSelectChat: (chatId: string) => void;
}

const ChatList: React.FC<ChatListProps> = ({ onSelectChat }) => {
  const dispatch = useDispatch<AppDispatch>();
  const chats = useSelector((state: RootState) => state.wuphf.chats);
  const chatStatus = useSelector((state: RootState) => state.wuphf.error);

  useEffect(() => {
    if (chatStatus === null) {
      dispatch(fetchChats());
    }
  }, [chatStatus, dispatch]);

  if (chatStatus === 'loading') {
    return <div>Loading...</div>;
  }

  if (chatStatus === 'failed') {
    return <div>Error fetching chats.</div>;
  }

  return (
    <ul>
      {chats.map((chat) => (
        <ChatItem key={chat.chatId} chat={{id: chat.chatId, name: chat.chatId}} onSelectChat={onSelectChat} />
      ))}
    </ul>
  );
};

export default ChatList;
