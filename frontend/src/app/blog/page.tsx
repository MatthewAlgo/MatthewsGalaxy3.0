'use client';

import { useEffect, useState } from 'react';
import { api, Post } from '@/lib/api';
import PostCard from '@/components/PostCard';

export default function BlogPage() {
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);

    const loadPosts = async (pageNum: number) => {
        setLoading(true);
        const { data } = await api.getPosts(pageNum, 9);
        if (data) {
            setPosts(data.data);
            setTotalPages(data.total_pages);
        }
        setLoading(false);
    };

    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        loadPosts(page);
    }, [page]);

    return (
        <div style={{ flex: 1 }}>
            {/* Hero */}
            <section className="blog-hero">
                <div className="container">
                    <span className="blog-emoji">📚</span>
                    <h1 className="blog-title">Blog Posts</h1>
                    <p className="blog-subtitle">
                        Thoughts, tutorials, and insights from my journey through tech.
                    </p>
                </div>
            </section>

            {/* Posts */}
            <section className="blog-section container">
                {loading ? (
                    <div className="blog-loading">
                        <div className="spinner"></div>
                    </div>
                ) : posts.length > 0 ? (
                    <>
                        <div className="blog-posts-grid">
                            {posts.map((post) => (
                                <PostCard key={post.id} post={post} />
                            ))}
                        </div>

                        {totalPages > 1 && (
                            <div className="blog-pagination">
                                <button
                                    onClick={() => setPage(p => Math.max(1, p - 1))}
                                    disabled={page === 1}
                                    className="blog-page-btn"
                                >
                                    ← Previous
                                </button>
                                <span className="blog-page-info">
                                    Page {page} of {totalPages}
                                </span>
                                <button
                                    onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                                    disabled={page === totalPages}
                                    className="blog-page-btn"
                                >
                                    Next →
                                </button>
                            </div>
                        )}
                    </>
                ) : (
                    <div className="blog-no-posts">
                        <div className="blog-no-posts-icon">🌌</div>
                        <h3>No Posts Yet</h3>
                        <p>The cosmos is vast, but this space is empty... for now.</p>
                        <p>Check back soon for stellar content!</p>
                    </div>
                )}
            </section>
        </div>
    );
}
