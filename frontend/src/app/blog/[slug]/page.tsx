'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import { api, Post, Comment } from '@/lib/api';
import { useAuth } from '@/lib/auth';

export default function PostPage() {
    const params = useParams();
    const slug = params.slug as string;
    const { user } = useAuth();

    const [post, setPost] = useState<Post | null>(null);
    const [comments, setComments] = useState<Comment[]>([]);
    const [loading, setLoading] = useState(true);
    const [liked, setLiked] = useState(false);
    const [likeCount, setLikeCount] = useState(0);
    const [newComment, setNewComment] = useState('');
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        const loadPost = async () => {
            const { data: postData } = await api.getPostBySlug(slug);
            if (postData) {
                setPost(postData);
                setLikeCount(postData.like_count);
            }

            const { data: commentsData } = await api.getComments(slug);
            if (commentsData) {
                setComments(commentsData);
            }

            const { data: likeData } = await api.getLikeStatus(slug);
            if (likeData) {
                setLiked(likeData.liked);
            }

            setLoading(false);
        };

        loadPost();
    }, [slug]);

    const handleLike = async () => {
        if (!user) return;
        const { data } = await api.toggleLike(slug);
        if (data) {
            setLiked(data.liked);
            setLikeCount(prev => data.liked ? prev + 1 : prev - 1);
        }
    };

    const handleComment = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user || !newComment.trim()) return;

        setSubmitting(true);
        const { data } = await api.createComment(slug, newComment);
        if (data) {
            setComments(prev => [data, ...prev]);
            setNewComment('');
        }
        setSubmitting(false);
    };

    const handleDeleteComment = async (commentId: string) => {
        if (!confirm('Delete this comment?')) return;
        const { error } = await api.deleteComment(slug, commentId);
        if (!error) {
            setComments(prev => prev.filter(c => c.id !== commentId));
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    };

    const formatContent = (content: string): string => {
        return content
            .replace(/^### (.*$)/gim, '<h3>$1</h3>')
            .replace(/^## (.*$)/gim, '<h2>$1</h2>')
            .replace(/^# (.*$)/gim, '<h1>$1</h1>')
            .replace(/\*\*(.*)\*\*/gim, '<strong>$1</strong>')
            .replace(/\*(.*)\*/gim, '<em>$1</em>')
            .replace(/\n/gim, '<br>');
    };

    if (loading) {
        return (
            <div className="post-page-loading">
                <div className="spinner"></div>
            </div>
        );
    }

    if (!post) {
        return (
            <div className="post-not-found">
                <h1>404</h1>
                <p>Post not found in this galaxy...</p>
                <Link href="/blog" className="btn-primary">Back to Blog</Link>
            </div>
        );
    }

    return (
        <article style={{ flex: 1, paddingBottom: '80px' }}>
            {/* Header */}
            <header className="post-header">
                <div className="container">
                    <Link href="/blog" className="post-back-link">← Back to Blog</Link>

                    <div className="post-meta">
                        <span className="post-meta-date">{formatDate(post.created_at)}</span>
                        <span>by {post.author_name}</span>
                    </div>

                    <h1 className="post-title">{post.title}</h1>

                    {post.excerpt && (
                        <p className="post-excerpt">{post.excerpt}</p>
                    )}

                    <div className="post-stats">
                        <button
                            onClick={handleLike}
                            className={`post-like-btn ${liked ? 'liked' : ''}`}
                            disabled={!user}
                            title={user ? (liked ? 'Unlike' : 'Like') : 'Login to like'}
                        >
                            {liked ? '❤️' : '🤍'} {likeCount}
                        </button>
                        <span className="post-stat-item">💬 {comments.length}</span>
                    </div>
                </div>
            </header>

            {/* Cover Image */}
            {post.cover_image && (
                <div className="post-cover-wrapper">
                    <img src={post.cover_image} alt={post.title} className="post-cover-image" />
                </div>
            )}

            {/* Content */}
            <div className="post-content container">
                <div className="post-prose" dangerouslySetInnerHTML={{ __html: formatContent(post.content) }} />
            </div>

            {/* Comments Section */}
            <section className="post-comments container">
                <h2 className="post-comments-title">
                    💬 Comments ({comments.length})
                </h2>

                {user ? (
                    <form onSubmit={handleComment} className="post-comment-form">
                        <textarea
                            value={newComment}
                            onChange={(e) => setNewComment(e.target.value)}
                            placeholder="Share your thoughts..."
                            className="post-comment-input"
                            rows={3}
                            required
                        />
                        <button
                            type="submit"
                            className="btn-primary"
                            disabled={submitting || !newComment.trim()}
                        >
                            {submitting ? 'Posting...' : 'Post Comment'}
                        </button>
                    </form>
                ) : (
                    <div className="post-login-prompt">
                        <p>Want to join the discussion?</p>
                        <Link href="/auth/login" className="btn-primary">Sign In to Comment</Link>
                    </div>
                )}

                <div className="post-comments-list">
                    {comments.length > 0 ? (
                        comments.map((comment) => (
                            <div key={comment.id} className="post-comment">
                                <div className="post-comment-header">
                                    <span className="post-comment-author">{comment.user_name}</span>
                                    <span className="post-comment-date">{formatDate(comment.created_at)}</span>
                                    {(user?.id === comment.user_id || user?.role === 'admin') && (
                                        <button
                                            onClick={() => handleDeleteComment(comment.id)}
                                            className="post-comment-delete"
                                        >
                                            🗑️
                                        </button>
                                    )}
                                </div>
                                <p className="post-comment-text">{comment.content}</p>
                            </div>
                        ))
                    ) : (
                        <p className="post-no-comments">No comments yet. Be the first to share your thoughts!</p>
                    )}
                </div>
            </section>
        </article>
    );
}
