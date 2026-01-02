import React from 'react';
import { render, screen } from '@testing-library/react';
import StarBackground from '../StarBackground';

describe('StarBackground', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('renders a canvas element', () => {
        render(<StarBackground />);
        const canvas = screen.getByTestId('star-background');
        expect(canvas).toBeInTheDocument();
    });

    it('applies the star-canvas class', () => {
        render(<StarBackground />);
        const canvas = screen.getByTestId('star-background');
        expect(canvas).toHaveClass('star-canvas');
    });

    it('gets 2d canvas context on mount', () => {
        render(<StarBackground />);
        expect(HTMLCanvasElement.prototype.getContext).toHaveBeenCalledWith('2d');
    });

    it('sets canvas dimensions based on window size', () => {
        Object.defineProperty(window, 'innerWidth', { value: 1920, writable: true });
        Object.defineProperty(window, 'innerHeight', { value: 1080, writable: true });

        render(<StarBackground />);
        const canvas = screen.getByTestId('star-background') as HTMLCanvasElement;

        expect(canvas.width).toBe(1920);
        expect(canvas.height).toBe(1080);
    });

    it('starts animation loop on mount', () => {
        render(<StarBackground />);
        expect(global.requestAnimationFrame).toHaveBeenCalled();
    });

    it('cleans up animation on unmount', () => {
        const { unmount } = render(<StarBackground />);
        unmount();
        expect(global.cancelAnimationFrame).toHaveBeenCalled();
    });

    it('handles window resize', () => {
        render(<StarBackground />);

        Object.defineProperty(window, 'innerWidth', { value: 800, writable: true });
        Object.defineProperty(window, 'innerHeight', { value: 600, writable: true });

        window.dispatchEvent(new Event('resize'));

        const canvas = screen.getByTestId('star-background') as HTMLCanvasElement;
        expect(canvas.width).toBe(800);
        expect(canvas.height).toBe(600);
    });
});
