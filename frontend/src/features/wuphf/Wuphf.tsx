import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../../app/store';
import WuphfForm from '../../components/WuphfForm';

const Wuphf: React.FC = () => {
  const messages = useSelector((state: RootState) => state.wuphf.messages);

  return (
    <div>
      <h1>WUPHF.com</h1>
      <WuphfForm />
      <ul>
        {messages.map((message, index) => (
          <li key={index}>{message}</li>
        ))}
      </ul>
    </div>
  );
};

export default Wuphf;
