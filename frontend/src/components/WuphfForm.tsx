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
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        placeholder="Send a WUPHF"
      />
      <button type="submit">Send</button>
    </form>
  );
};

export default WuphfForm;
