import React from 'react';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import AuthForm from './components/AuthForm';
import ProtectedRoute from './components/ProtectedRoute';
import ChatList from './components/ChatList';
import Wuphf from './features/wuphf/Wuphf';

const router = createBrowserRouter([
  {
    path: "/",
    element: <ProtectedRoute exact path="/" component={ChatList} />,
  },
  {
    path: "/chat/:chatId",
    element: <ProtectedRoute path="/chat" component={Wuphf} />,
  },
  {
    path: "/login",
    element: <AuthForm mode='login' />,
  },
  {
    path: "/register",
    element: <AuthForm mode='register' />,
  }
]);

const App: React.FC = () => {
  return (
    <RouterProvider router={router} />
  );
};

export default App;
