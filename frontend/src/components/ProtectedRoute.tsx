import React, { useEffect, useState } from 'react';
import { Navigate } from 'react-router-dom';

interface ProtectedRouteProps {
    component: React.ComponentType<any>;
    exact?: boolean;
    path: string;
    [key: string]: any;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ component: Component, ...rest }) => {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const token = localStorage.getItem('token') || '';
                const user_id = localStorage.getItem('user_id') || '';
                const response = await fetch(`${process.env.REACT_APP_BASE_URL}/user?id=${user_id}`, {
                    headers: {
                        'Authorization': token
                    }
                });
                if (response.ok) {
                    setIsAuthenticated(true);
                } else {
                    setIsAuthenticated(false);
                }
            } catch (error) {
                setIsAuthenticated(false);
            }
        };
        fetchUser();
    }, []);

    if (isAuthenticated === null) {
        return <div>Loading...</div>;
    }

    return (
        isAuthenticated ? (
            <Component {...rest} />
        ) : (
            <Navigate to="/login" replace={true} />
        )
    );
};

export default ProtectedRoute;
