import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

const AuthForm: React.FC<{ mode: 'login' | 'register' }> = ({ mode }) => {
    const navigate = useNavigate();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');

    const handleFormSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError('');

        const baseUrl = process.env.REACT_APP_BASE_URL;
        const url = mode === 'login' ? `${baseUrl}/auth/login` : `${baseUrl}/auth/register`;
        const method = 'POST';
        const body = mode === 'login'
            ? JSON.stringify({ email, password })
            : JSON.stringify({ email, password, confirmPassword });

        try {
            const response = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body,
            });

            if (response.ok) {
                const token = response.headers.get('Authorization');

                if (token) {
                    localStorage.setItem('token', token);
                    navigate('/');
                } else {
                    setError('Login/Registration failed');
                }
            } else {
                setError(await response.text());
            }
        } catch (error) {
            setError('Error logging in/registering');
        }
    };

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = event.target;
        switch (name) {
            case 'email':
                setEmail(value);
                break;
            case 'password':
                setPassword(value);
                break;
            case 'confirmPassword':
                setConfirmPassword(value);
                break;
            default:
                break;
        }
    };

    return (
        <div className="w-screen h-screen flex justify-center items-center">
            <form onSubmit={handleFormSubmit} className="flex flex-col items-center space-y-4 p-4">
                <h2 className="text-xl mb-4">{mode === 'login' ? 'Login' : 'Sign Up'}</h2>
                {mode === 'login' ? (
                    <>
                        <p>Not registered yet? <Link to="/register">Signup</Link></p>
                        <input
                            name="email"
                            type="email"
                            value={email}
                            onChange={handleInputChange}
                            placeholder="Email"
                            className="w-full p-2 border border-gray-300 rounded-lg"
                        />
                        <input
                            name="password"
                            type="password"
                            value={password}
                            onChange={handleInputChange}
                            placeholder="Password"
                            className="w-full p-2 border border-gray-300 rounded-lg"
                        />
                    </>
                ) : (
                    <>
                        <p>Already have an account? <Link to="/login">Login</Link></p>
                        <input
                            name="email"
                            type="email"
                            value={email}
                            onChange={handleInputChange}
                            placeholder="Email"
                            className="w-full p-2 border border-gray-300 rounded-lg"
                        />
                        <input
                            name="password"
                            type="password"
                            value={password}
                            onChange={handleInputChange}
                            placeholder="Password"
                            className="w-full p-2 border border-gray-300 rounded-lg"
                        />
                        <input
                            name="confirmPassword"
                            type="password"
                            value={confirmPassword}
                            onChange={handleInputChange}
                            placeholder="Confirm Password"
                            className="w-full p-2 border border-gray-300 rounded-lg"
                        />
                    </>
                )}
                <button type="submit" className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-green-500"> Submit </button>
                {error && <p className="text-red-500">{error}</p>}
            </form>
        </div>
    );
}

export default AuthForm;