import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../../app/store';
import WuphfForm from '../../components/WuphfForm';

const Wuphf: React.FC = () => {
  const messages = useSelector((state: RootState) => state.wuphf.messages);

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold text-center mb-6">WUPHF.com</h1>
      <WuphfForm />
      <ul className="mt-6 space-y-2">
        {messages.map((message, index) => (
          <li key={index} className="p-2 bg-gray-100 rounded-lg shadow">
            {message}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Wuphf;
