import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchMessages } from '../wuphf/wuphfSlice';
import { AppDispatch, RootState } from '../../app/store';
import WuphfForm from '../../components/WuphfForm';
import { useParams } from 'react-router-dom';

const Wuphf: React.FC = () => {
  const { chatId } = useParams();
  const dispatch = useDispatch<AppDispatch>();
  const messages = useSelector((state: RootState) => state.wuphf.chats);
  const loading = useSelector((state: RootState) => state.wuphf.loading);
  const error = useSelector((state: RootState) => state.wuphf.error);

  useEffect(() => {
    dispatch(fetchMessages(chatId || ''));
  }, [chatId, dispatch]);

  console.log(messages)
  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold text-center mb-6">WUPHF.com</h1>
      <WuphfForm chatId={chatId || ''} />
      <ul className="mt-6 space-y-2">
        {messages?.map((message: any, index: React.Key) => (
          <li key={index} className="p-2 bg-gray-100 rounded-lg shadow">
            {message.msg}
          </li>
        ))}
      </ul>
      {loading && <p>Loading messages...</p>}
      {error && <p>{error}</p>}
    </div>
  );
};

export default Wuphf;
