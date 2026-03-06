import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import Header from '../Header';
import { AuthContext } from '@/lib/auth';
import { User } from '@/lib/api';

const mockPush = jest.fn();
jest.mock('next/navigation', () => ({
    useRouter: () => ({ push: mockPush }),
}));

const mockAuthContext: {
    user: User | null;
    loading: boolean;
    isAdmin: boolean;
    login: jest.Mock;
    register: jest.Mock;
    logout: jest.Mock;
} = {
    user: null,
    loading: false,
    isAdmin: false,
    login: jest.fn(),
    register: jest.fn(),
    logout: jest.fn(),
};

const renderWithAuth = (authValue = mockAuthContext) => {
    return render(
        <AuthContext.Provider value={authValue}>
            <Header />
        </AuthContext.Provider>
    );
};

describe('Header', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('renders the header element', () => {
        renderWithAuth();
        expect(screen.getByTestId('header')).toBeInTheDocument();
    });

    it('displays the logo with correct text', () => {
        renderWithAuth();
        expect(screen.getByText("Matthew's Galaxy")).toBeInTheDocument();
    });

    it('displays the logo emoji', () => {
        renderWithAuth();
        expect(screen.getByText('🌌')).toBeInTheDocument();
    });

    it('has navigation links', () => {
        renderWithAuth();
        expect(screen.getByText('Home')).toBeInTheDocument();
        expect(screen.getByText('Blog')).toBeInTheDocument();
        expect(screen.getByText('About')).toBeInTheDocument();
    });

    it('shows Sign In button when not logged in', () => {
        renderWithAuth();
        expect(screen.getByText('Sign In')).toBeInTheDocument();
    });

    it('shows user name and logout when logged in', () => {
        renderWithAuth({
            ...mockAuthContext,
            user: { id: '1', email: 'test@test.com', name: 'Test User', role: 'user', created_at: '' },
        });

        expect(screen.getByText('👋 Test User')).toBeInTheDocument();
        expect(screen.getByTestId('logout-btn')).toBeInTheDocument();
    });

    it('shows Admin link when user is admin', () => {
        renderWithAuth({
            ...mockAuthContext,
            user: { id: '1', email: 'admin@test.com', name: 'Admin', role: 'admin', created_at: '' },
            isAdmin: true,
        });

        expect(screen.getByText('Admin')).toBeInTheDocument();
    });

    it('does not show Admin link for regular users', () => {
        renderWithAuth({
            ...mockAuthContext,
            user: { id: '1', email: 'user@test.com', name: 'User', role: 'user', created_at: '' },
            isAdmin: false,
        });

        expect(screen.queryByText('Admin')).not.toBeInTheDocument();
    });

    it('calls logout when logout button is clicked', () => {
        const logoutMock = jest.fn();
        renderWithAuth({
            ...mockAuthContext,
            user: { id: '1', email: 'test@test.com', name: 'Test', role: 'user', created_at: '' },
            logout: logoutMock,
        });

        fireEvent.click(screen.getByTestId('logout-btn'));
        expect(logoutMock).toHaveBeenCalled();
    });

    it('toggles mobile menu when menu button is clicked', () => {
        renderWithAuth();
        const menuButton = screen.getByTestId('menu-toggle');
        const nav = screen.getByTestId('nav');

        expect(nav).not.toHaveClass('open');

        fireEvent.click(menuButton);
        expect(nav).toHaveClass('open');

        fireEvent.click(menuButton);
        expect(nav).not.toHaveClass('open');
    });

    it('closes mobile menu when a link is clicked', () => {
        renderWithAuth();
        const menuButton = screen.getByTestId('menu-toggle');
        const nav = screen.getByTestId('nav');

        fireEvent.click(menuButton);
        expect(nav).toHaveClass('open');

        fireEvent.click(screen.getByText('Blog'));
        expect(nav).not.toHaveClass('open');
    });

    it('has correct link hrefs', () => {
        renderWithAuth();

        expect(screen.getByText('Home').closest('a')).toHaveAttribute('href', '/');
        expect(screen.getByText('Blog').closest('a')).toHaveAttribute('href', '/blog');
        expect(screen.getByText('About').closest('a')).toHaveAttribute('href', '/about');
    });
});
