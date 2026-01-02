import React from 'react';
import { render, screen } from '@testing-library/react';
import Footer from '../Footer';

describe('Footer', () => {
    it('renders the footer element', () => {
        render(<Footer />);
        expect(screen.getByTestId('footer')).toBeInTheDocument();
    });

    it('displays the logo', () => {
        render(<Footer />);
        expect(screen.getByText("Matthew's Galaxy")).toBeInTheDocument();
    });

    it('displays the tagline', () => {
        render(<Footer />);
        expect(screen.getByText(/Exploring the frontiers of technology/)).toBeInTheDocument();
    });

    it('displays the Latin quote', () => {
        render(<Footer />);
        expect(screen.getByText(/Ad astra per aspera/)).toBeInTheDocument();
    });

    it('has navigation links section', () => {
        render(<Footer />);
        expect(screen.getByText('Navigate')).toBeInTheDocument();
    });

    it('has connect links section', () => {
        render(<Footer />);
        expect(screen.getByText('Connect')).toBeInTheDocument();
    });

    it('has legal links section', () => {
        render(<Footer />);
        expect(screen.getByText('Legal')).toBeInTheDocument();
    });

    it('displays current year in copyright', () => {
        render(<Footer />);
        const currentYear = new Date().getFullYear().toString();
        expect(screen.getByText(new RegExp(currentYear))).toBeInTheDocument();
    });

    it('has correct navigation links', () => {
        render(<Footer />);
        const homeLinks = screen.getAllByText('Home');
        expect(homeLinks.length).toBeGreaterThan(0);
        expect(screen.getByText('Blog')).toBeInTheDocument();
        expect(screen.getByText('About')).toBeInTheDocument();
    });

    it('has social media links with correct targets', () => {
        render(<Footer />);
        const linkedinLink = screen.getByText('LinkedIn');
        expect(linkedinLink).toHaveAttribute('target', '_blank');
        expect(linkedinLink).toHaveAttribute('rel', 'noopener noreferrer');
    });

    it('has GitHub link', () => {
        render(<Footer />);
        const githubLink = screen.getByText('GitHub');
        expect(githubLink).toHaveAttribute('href', expect.stringContaining('github.com'));
    });

    it('has email link', () => {
        render(<Footer />);
        const emailLink = screen.getByText('Email');
        expect(emailLink).toHaveAttribute('href', expect.stringContaining('mailto:'));
    });

    it('displays tech stack', () => {
        render(<Footer />);
        expect(screen.getByText('Next.js')).toBeInTheDocument();
        expect(screen.getByText('Go')).toBeInTheDocument();
        expect(screen.getByText('PostgreSQL')).toBeInTheDocument();
    });

    it('has privacy policy link', () => {
        render(<Footer />);
        const privacyLink = screen.getByText('Privacy Policy');
        expect(privacyLink).toHaveAttribute('href', '/privacy');
    });

    it('has terms of service link', () => {
        render(<Footer />);
        const termsLink = screen.getByText('Terms of Service');
        expect(termsLink).toHaveAttribute('href', '/terms');
    });
});
