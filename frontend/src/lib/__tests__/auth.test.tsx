import React from 'react';
import { render, screen, waitFor, act } from '@testing-library/react';
import { AuthProvider, useAuth } from '../auth';
import { api } from '../api';

jest.mock('../api', () => ({
    api: {
        setToken: jest.fn(),
        getToken: jest.fn(),
        getCurrentUser: jest.fn(),
        login: jest.fn(),
        register: jest.fn(),
    },
}));

// Test component that uses the auth hook
function TestComponent() {
    const { user, loading, isAdmin, login, register, logout } = useAuth();

    return (
        <div>
            <div data-testid="loading">{loading ? 'true' : 'false'}</div>
            <div data-testid="user">{user ? user.name : 'null'}</div>
            <div data-testid="isAdmin">{isAdmin ? 'true' : 'false'}</div>
            <button data-testid="login-btn" onClick={() => login('test@test.com', 'pass')}>Login</button>
            <button data-testid="register-btn" onClick={() => register('test@test.com', 'pass', 'Test')}>Register</button>
            <button data-testid="logout-btn" onClick={logout}>Logout</button>
        </div>
    );
}

describe('AuthProvider', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        (api.getToken as jest.Mock).mockReturnValue(null);
        (api.getCurrentUser as jest.Mock).mockResolvedValue({ data: null });
    });

    it('provides initial loading state', async () => {
        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        // Initially loading should be true or false depending on token
        await waitFor(() => {
            expect(screen.getByTestId('loading')).toBeInTheDocument();
        });
    });

    it('provides null user when not authenticated', async () => {
        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('user')).toHaveTextContent('null');
        });
    });

    it('fetches current user when token exists', async () => {
        (api.getToken as jest.Mock).mockReturnValue('valid-token');
        (api.getCurrentUser as jest.Mock).mockResolvedValue({
            data: { id: '1', name: 'Test User', email: 'test@test.com', role: 'user' },
        });

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('user')).toHaveTextContent('Test User');
        });
    });

    it('sets isAdmin to true for admin users', async () => {
        (api.getToken as jest.Mock).mockReturnValue('admin-token');
        (api.getCurrentUser as jest.Mock).mockResolvedValue({
            data: { id: '1', name: 'Admin', email: 'admin@test.com', role: 'admin' },
        });

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('isAdmin')).toHaveTextContent('true');
        });
    });

    it('sets isAdmin to false for regular users', async () => {
        (api.getToken as jest.Mock).mockReturnValue('user-token');
        (api.getCurrentUser as jest.Mock).mockResolvedValue({
            data: { id: '1', name: 'User', email: 'user@test.com', role: 'user' },
        });

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('isAdmin')).toHaveTextContent('false');
        });
    });

    it('handles login success', async () => {
        (api.login as jest.Mock).mockResolvedValue({
            data: { token: 'new-token', user: { id: '1', name: 'Logged In', role: 'user' } },
        });

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('loading')).toHaveTextContent('false');
        });

        await act(async () => {
            screen.getByTestId('login-btn').click();
        });

        await waitFor(() => {
            expect(api.login).toHaveBeenCalledWith('test@test.com', 'pass');
        });
    });

    it('handles login failure', async () => {
        (api.login as jest.Mock).mockResolvedValue({
            data: null,
            error: 'Invalid credentials',
        });

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('loading')).toHaveTextContent('false');
        });

        await act(async () => {
            screen.getByTestId('login-btn').click();
        });

        await waitFor(() => {
            expect(screen.getByTestId('user')).toHaveTextContent('null');
        });
    });

    it('handles register success', async () => {
        (api.register as jest.Mock).mockResolvedValue({
            data: { token: 'new-token', user: { id: '1', name: 'New User', role: 'user' } },
        });

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('loading')).toHaveTextContent('false');
        });

        await act(async () => {
            screen.getByTestId('register-btn').click();
        });

        await waitFor(() => {
            expect(api.register).toHaveBeenCalledWith('test@test.com', 'pass', 'Test');
        });
    });

    it('handles logout', async () => {
        (api.getToken as jest.Mock).mockReturnValue('token');
        (api.getCurrentUser as jest.Mock).mockResolvedValue({
            data: { id: '1', name: 'User', role: 'user' },
        });

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('user')).toHaveTextContent('User');
        });

        await act(async () => {
            screen.getByTestId('logout-btn').click();
        });

        await waitFor(() => {
            expect(api.setToken).toHaveBeenCalledWith(null);
            expect(screen.getByTestId('user')).toHaveTextContent('null');
        });
    });

    it('clears token when getCurrentUser fails', async () => {
        (api.getToken as jest.Mock).mockReturnValue('invalid-token');
        (api.getCurrentUser as jest.Mock).mockResolvedValue({
            data: null,
            error: 'Unauthorized',
        });

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(api.setToken).toHaveBeenCalledWith(null);
        });
    });
});

describe('useAuth hook', () => {
    it('throws error when used outside AuthProvider', () => {
        // Suppress console.error for this test
        const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => { });

        expect(() => {
            render(<TestComponent />);
        }).toThrow('useAuth must be used within an AuthProvider');

        consoleSpy.mockRestore();
    });
});
