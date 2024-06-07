import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchChats } from '../features/wuphf/wuphfSlice';
import ChatItem from './ChatItem';
import { RootState, AppDispatch } from '../app/store';
import { useNavigate } from 'react-router-dom';

const ChatList: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>();
  const chats = useSelector((state: RootState) => state.wuphf.chats);
  const chatStatus = useSelector((state: RootState) => state.wuphf.error);

  const handleSelectChat = (chatId: string) => {
    navigate(`/chat/${chatId}`);
  };

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
      {chats.map((chat, idx) => (
        <ChatItem key={idx} itemKey={idx} chat={{id: chat.chatId, name: chat.chatId}} onSelectChat={handleSelectChat} />
      ))}
    </ul>
  );
};

export default ChatList;
