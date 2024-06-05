import React from 'react';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import Wuphf from './features/wuphf/Wuphf';
import AuthForm from './components/AuthForm';
import ProtectedRoute from './components/ProtectedRoute';

const router = createBrowserRouter([
  {
    path: "/",
    element: <ProtectedRoute exact path="/" component={Wuphf} />,
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
