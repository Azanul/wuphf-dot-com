import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchMessages } from '../wuphf/wuphfSlice';
import { AppDispatch, RootState } from '../../app/store';
import WuphfForm, { Props } from '../../components/WuphfForm';

const Wuphf: React.FC<Props> = (props) => {
  const dispatch = useDispatch<AppDispatch>();
  const messages = useSelector((state: RootState) => state.wuphf.chats.find((chat) => chat.chatId === props.chatId)?.messages);
  const loading = useSelector((state: RootState) => state.wuphf.loading);
  const error = useSelector((state: RootState) => state.wuphf.error);

  useEffect(() => {
    dispatch(fetchMessages());
  }, [dispatch]);

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold text-center mb-6">WUPHF.com</h1>
      <WuphfForm receiver_id={props.receiver_id} chatId={props.chatId} />
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
