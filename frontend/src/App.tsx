import React from 'react';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import WuphfForm from './components/WuphfForm';
import AuthForm from './components/AuthForm';
import ProtectedRoute from './components/ProtectedRoute';

const router = createBrowserRouter([
  {
    path: "/",
    element: <ProtectedRoute exact path="/" component={WuphfForm} />,
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
