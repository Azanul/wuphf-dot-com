import React from 'react';

interface ChatItemProps {
  itemKey: React.Key;
  chat: {
    id: string;
    name: string;
  };
  onSelectChat: (chatId: string) => void;
}

const ChatItem: React.FC<ChatItemProps> = ({ itemKey, chat, onSelectChat }) => {
  return (
    <li key={itemKey} onClick={() => onSelectChat(chat.id)} className='bg-green-100 p-4 m-1'>
      {chat.name}
    </li>
  );
};

export default ChatItem;
