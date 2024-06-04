import React, { useState } from 'react';

const Login: React.FC = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const response = await fetch(`${process.env.REACT_APP_BASE_URL}/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password }),
            });
            if (response.ok) {
                const token = response.headers.get('AUTHORIZATION');
                if (token) {
                    localStorage.setItem('token', token);
                    setError('');
                } else {
                    setError('Login failed');
                }
            } else {
                setError('Login failed');
            }
        } catch (error) {
            setError('Error logging in');
        }
    };

    return (
        <div className="w-screen h-screen flex justify-center items-center">
            <form onSubmit={handleLogin} className="flex flex-col items-center space-y-4 p-4">
                <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Email"
                    className="w-full p-2 border border-gray-300 rounded-lg"
                />
                <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Password"
                    className="w-full p-2 border border-gray-300 rounded-lg"
                />
                <button type="submit" className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600">
                    Login
                </button>
                {error && <p className="text-red-500">{error}</p>}
            </form>
        </div>
    );
}

export default Login;