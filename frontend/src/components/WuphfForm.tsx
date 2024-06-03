import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { sendWuphf } from '../features/wuphf/wuphfSlice';

const WuphfForm: React.FC = () => {
  const [message, setMessage] = useState('');
  const dispatch = useDispatch();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    dispatch(sendWuphf(message));
    setMessage('');
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
    </form>
  );
};

export default WuphfForm;
