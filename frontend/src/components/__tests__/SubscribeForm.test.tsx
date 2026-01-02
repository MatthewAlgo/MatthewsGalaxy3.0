import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SubscribeForm from '../SubscribeForm';
import { api } from '@/lib/api';

jest.mock('@/lib/api', () => ({
    api: {
        subscribe: jest.fn(),
    },
}));

describe('SubscribeForm', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('renders the subscribe form', () => {
        render(<SubscribeForm />);
        expect(screen.getByTestId('subscribe-form')).toBeInTheDocument();
    });

    it('displays the title', () => {
        render(<SubscribeForm />);
        expect(screen.getByText('Stay in Orbit')).toBeInTheDocument();
    });

    it('displays the description', () => {
        render(<SubscribeForm />);
        expect(screen.getByText(/Get notified when new posts launch/)).toBeInTheDocument();
    });

    it('displays the email icon', () => {
        render(<SubscribeForm />);
        expect(screen.getByText('✉️')).toBeInTheDocument();
    });

    it('has an email input field', () => {
        render(<SubscribeForm />);
        const input = screen.getByTestId('subscribe-input');
        expect(input).toBeInTheDocument();
        expect(input).toHaveAttribute('type', 'email');
        expect(input).toHaveAttribute('required');
    });

    it('has a subscribe button', () => {
        render(<SubscribeForm />);
        const button = screen.getByTestId('subscribe-button');
        expect(button).toBeInTheDocument();
        expect(button).toHaveTextContent('Subscribe 🚀');
    });

    it('updates email input value when typing', async () => {
        render(<SubscribeForm />);
        const input = screen.getByTestId('subscribe-input');

        await userEvent.type(input, 'test@example.com');
        expect(input).toHaveValue('test@example.com');
    });

    it('shows loading state when submitting', async () => {
        (api.subscribe as jest.Mock).mockImplementation(
            () => new Promise((resolve) => setTimeout(() => resolve({ data: { message: 'Success' } }), 100))
        );

        render(<SubscribeForm />);
        const input = screen.getByTestId('subscribe-input');
        const button = screen.getByTestId('subscribe-button');

        await userEvent.type(input, 'test@example.com');
        fireEvent.click(button);

        expect(button).toHaveTextContent('...');
        expect(button).toBeDisabled();
        expect(input).toBeDisabled();
    });

    it('shows success message on successful subscription', async () => {
        (api.subscribe as jest.Mock).mockResolvedValue({
            data: { message: 'Successfully subscribed!' },
            error: null,
        });

        render(<SubscribeForm />);
        const input = screen.getByTestId('subscribe-input');
        const button = screen.getByTestId('subscribe-button');

        await userEvent.type(input, 'test@example.com');
        fireEvent.click(button);

        await waitFor(() => {
            expect(screen.getByTestId('subscribe-message')).toHaveTextContent('Successfully subscribed!');
        });
    });

    it('clears input on successful subscription', async () => {
        (api.subscribe as jest.Mock).mockResolvedValue({
            data: { message: 'Success' },
            error: null,
        });

        render(<SubscribeForm />);
        const input = screen.getByTestId('subscribe-input');
        const button = screen.getByTestId('subscribe-button');

        await userEvent.type(input, 'test@example.com');
        fireEvent.click(button);

        await waitFor(() => {
            expect(input).toHaveValue('');
        });
    });

    it('shows error message on failed subscription', async () => {
        (api.subscribe as jest.Mock).mockResolvedValue({
            data: null,
            error: 'Email already subscribed',
        });

        render(<SubscribeForm />);
        const input = screen.getByTestId('subscribe-input');
        const button = screen.getByTestId('subscribe-button');

        await userEvent.type(input, 'test@example.com');
        fireEvent.click(button);

        await waitFor(() => {
            expect(screen.getByTestId('subscribe-message')).toHaveTextContent('Email already subscribed');
        });
    });

    it('calls api.subscribe with correct email', async () => {
        (api.subscribe as jest.Mock).mockResolvedValue({ data: { message: 'Success' } });

        render(<SubscribeForm />);
        const input = screen.getByTestId('subscribe-input');
        const button = screen.getByTestId('subscribe-button');

        await userEvent.type(input, 'test@example.com');
        fireEvent.click(button);

        await waitFor(() => {
            expect(api.subscribe).toHaveBeenCalledWith('test@example.com');
        });
    });

    it('success message has correct class', async () => {
        (api.subscribe as jest.Mock).mockResolvedValue({
            data: { message: 'Success' },
            error: null,
        });

        render(<SubscribeForm />);
        const input = screen.getByTestId('subscribe-input');
        const button = screen.getByTestId('subscribe-button');

        await userEvent.type(input, 'test@example.com');
        fireEvent.click(button);

        await waitFor(() => {
            const message = screen.getByTestId('subscribe-message');
            expect(message).toHaveClass('success');
        });
    });

    it('error message has correct class', async () => {
        (api.subscribe as jest.Mock).mockResolvedValue({
            data: null,
            error: 'Error',
        });

        render(<SubscribeForm />);
        const input = screen.getByTestId('subscribe-input');
        const button = screen.getByTestId('subscribe-button');

        await userEvent.type(input, 'test@example.com');
        fireEvent.click(button);

        await waitFor(() => {
            const message = screen.getByTestId('subscribe-message');
            expect(message).toHaveClass('error');
        });
    });
});
