import React from 'react';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import WuphfForm from './components/WuphfForm';
import Login from './components/Login';
import ProtectedRoute from './components/ProtectedRoute';

const router = createBrowserRouter([
  {
    path: "/",
    element: <ProtectedRoute exact path="/" component={WuphfForm} />,
  },
  {
    path: "/login",
    element: <Login />,
  }
]);

const App: React.FC = () => {
  return (
    <RouterProvider router={router} />
  );
};

export default App;
