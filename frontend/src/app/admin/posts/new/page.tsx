'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth';
import { api } from '@/lib/api';

export default function NewPostPage() {
    const router = useRouter();
    const { isAdmin } = useAuth();
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [excerpt, setExcerpt] = useState('');
    const [coverImage, setCoverImage] = useState('');
    const [published, setPublished] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        const { data, error: apiError } = await api.createPost({
            title,
            content,
            excerpt: excerpt || undefined,
            cover_image: coverImage || undefined,
            published,
        });

        if (data) {
            router.push('/admin');
        } else {
            setError(apiError || 'Failed to create post');
        }

        setLoading(false);
    };

    if (!isAdmin) {
        return null;
    }

    return (
        <div className="admin-post-page">
            <div className="admin-post-container">
                <h1>Create New Post</h1>

                <form onSubmit={handleSubmit} className="admin-post-form">
                    {error && <div className="admin-post-error">{error}</div>}

                    <div className="admin-post-field">
                        <label htmlFor="title">Title *</label>
                        <input
                            id="title"
                            type="text"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            placeholder="Enter post title"
                            className="input-field"
                            required
                            minLength={3}
                        />
                    </div>

                    <div className="admin-post-field">
                        <label htmlFor="excerpt">Excerpt</label>
                        <textarea
                            id="excerpt"
                            value={excerpt}
                            onChange={(e) => setExcerpt(e.target.value)}
                            placeholder="Brief summary (optional)"
                            className="admin-post-textarea"
                            rows={2}
                        />
                    </div>

                    <div className="admin-post-field">
                        <label htmlFor="coverImage">Cover Image URL</label>
                        <input
                            id="coverImage"
                            type="url"
                            value={coverImage}
                            onChange={(e) => setCoverImage(e.target.value)}
                            placeholder="https://example.com/image.jpg"
                            className="input-field"
                        />
                    </div>

                    <div className="admin-post-field">
                        <label htmlFor="content">Content *</label>
                        <textarea
                            id="content"
                            value={content}
                            onChange={(e) => setContent(e.target.value)}
                            placeholder="Write your post content here... (Markdown supported)"
                            className="admin-post-content-area"
                            rows={20}
                            required
                            minLength={10}
                        />
                    </div>

                    <div className="admin-post-publish-row">
                        <label className="admin-post-checkbox">
                            <input
                                type="checkbox"
                                checked={published}
                                onChange={(e) => setPublished(e.target.checked)}
                            />
                            <span>Publish immediately</span>
                        </label>
                    </div>

                    <div className="admin-post-actions">
                        <button
                            type="button"
                            onClick={() => router.back()}
                            className="btn-secondary"
                        >
                            Cancel
                        </button>
                        <button type="submit" className="btn-primary" disabled={loading}>
                            {loading ? 'Creating...' : published ? 'Publish Post 🚀' : 'Save as Draft'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
