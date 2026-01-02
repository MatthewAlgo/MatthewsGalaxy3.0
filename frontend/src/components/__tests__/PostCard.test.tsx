import React from 'react';
import { render, screen } from '@testing-library/react';
import PostCard from '../PostCard';
import { Post } from '@/lib/api';

const mockPost: Post = {
    id: '1',
    title: 'Test Post Title',
    slug: 'test-post-title',
    content: 'This is the test post content.',
    excerpt: 'This is a short excerpt for the test post.',
    cover_image: 'https://example.com/image.jpg',
    author_id: 'author-1',
    author_name: 'John Doe',
    published: true,
    like_count: 42,
    comment_count: 15,
    created_at: '2026-01-01T12:00:00Z',
    updated_at: '2026-01-01T12:00:00Z',
};

describe('PostCard', () => {
    it('renders the post card', () => {
        render(<PostCard post={mockPost} />);
        expect(screen.getByTestId('post-card')).toBeInTheDocument();
    });

    it('displays the post title', () => {
        render(<PostCard post={mockPost} />);
        expect(screen.getByText('Test Post Title')).toBeInTheDocument();
    });

    it('displays the post excerpt', () => {
        render(<PostCard post={mockPost} />);
        expect(screen.getByText('This is a short excerpt for the test post.')).toBeInTheDocument();
    });

    it('displays the author name', () => {
        render(<PostCard post={mockPost} />);
        expect(screen.getByText('by John Doe')).toBeInTheDocument();
    });

    it('displays the formatted date', () => {
        render(<PostCard post={mockPost} />);
        expect(screen.getByText('January 1, 2026')).toBeInTheDocument();
    });

    it('displays the like count', () => {
        render(<PostCard post={mockPost} />);
        expect(screen.getByText('❤️ 42')).toBeInTheDocument();
    });

    it('displays the comment count', () => {
        render(<PostCard post={mockPost} />);
        expect(screen.getByText('💬 15')).toBeInTheDocument();
    });

    it('has a Read More link', () => {
        render(<PostCard post={mockPost} />);
        expect(screen.getByText('Read More →')).toBeInTheDocument();
    });

    it('links to the correct post URL', () => {
        render(<PostCard post={mockPost} />);
        const readMoreLink = screen.getByText('Read More →').closest('a');
        expect(readMoreLink).toHaveAttribute('href', '/blog/test-post-title');
    });

    it('title links to the correct post URL', () => {
        render(<PostCard post={mockPost} />);
        const titleLink = screen.getByText('Test Post Title').closest('a');
        expect(titleLink).toHaveAttribute('href', '/blog/test-post-title');
    });

    it('displays the cover image when provided', () => {
        render(<PostCard post={mockPost} />);
        const image = screen.getByRole('img');
        expect(image).toHaveAttribute('src', 'https://example.com/image.jpg');
        expect(image).toHaveAttribute('alt', 'Test Post Title');
    });

    it('does not display image when cover_image is not provided', () => {
        const postWithoutImage = { ...mockPost, cover_image: undefined };
        render(<PostCard post={postWithoutImage} />);
        expect(screen.queryByRole('img')).not.toBeInTheDocument();
    });

    it('applies featured class when featured prop is true', () => {
        render(<PostCard post={mockPost} featured={true} />);
        expect(screen.getByTestId('post-card')).toHaveClass('featured');
    });

    it('does not apply featured class by default', () => {
        render(<PostCard post={mockPost} />);
        expect(screen.getByTestId('post-card')).not.toHaveClass('featured');
    });

    it('handles post without excerpt', () => {
        const postWithoutExcerpt = { ...mockPost, excerpt: undefined };
        render(<PostCard post={postWithoutExcerpt} />);
        expect(screen.getByTestId('post-card')).toBeInTheDocument();
        expect(screen.queryByText('This is a short excerpt')).not.toBeInTheDocument();
    });

    it('handles zero likes and comments', () => {
        const postWithZeroStats = { ...mockPost, like_count: 0, comment_count: 0 };
        render(<PostCard post={postWithZeroStats} />);
        expect(screen.getByText('❤️ 0')).toBeInTheDocument();
        expect(screen.getByText('💬 0')).toBeInTheDocument();
    });

    it('formats different dates correctly', () => {
        const postWithDifferentDate = { ...mockPost, created_at: '2025-06-15T10:30:00Z' };
        render(<PostCard post={postWithDifferentDate} />);
        expect(screen.getByText('June 15, 2025')).toBeInTheDocument();
    });
});
