import React, { useState, useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { sendWuphf } from '../features/wuphf/wuphfSlice';

type Props = { chatId: string };

const WuphfForm: React.FC<Props> = (props) => {
  const dispatch = useDispatch();
  const [message, setMessage] = useState('');
  const [status, setStatus] = useState('');

  useEffect(() => {
    if (status) {
      const timer = setTimeout(() => {
        setStatus('');
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [status]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await fetch(`${process.env.REACT_APP_BASE_URL}/notification`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': localStorage.getItem('token') || '',
        },
        body: JSON.stringify({ sender: localStorage.getItem('user_id'), chat_id: props.chatId, msg: message }),
      });
      if (response.ok) {
        setStatus('Message sent successfully');
        setMessage('');
        dispatch(sendWuphf({ chatId: props.chatId, message: message }));
      } else {
        setStatus('Failed to send message');
      }
    } catch (error) {
      setStatus('Error sending message');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col items-center space-y-4">
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        placeholder="Send a WUPHF"
        className="w-full p-2 border border-gray-300 rounded-lg"
      />
      <button type="submit" className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600">
        Send
      </button>
      {status && <p>{status}</p>}
    </form>
  );
};

export default WuphfForm;
export type { Props };