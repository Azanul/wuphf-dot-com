import React from 'react';

interface ChatItemProps {
  chat: {
    id: string;
    name: string;
  };
  onSelectChat: (chatId: string) => void;
}

const ChatItem: React.FC<ChatItemProps> = ({ chat, onSelectChat }) => {
  return (
    <li onClick={() => onSelectChat(chat.id)}>
      {chat.name}
    </li>
  );
};

export default ChatItem;
