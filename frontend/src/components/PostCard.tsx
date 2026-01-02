import Link from 'next/link';
import { Post } from '@/lib/api';

interface PostCardProps {
    post: Post;
    featured?: boolean;
}

export default function PostCard({ post, featured = false }: PostCardProps) {
    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    };

    return (
        <article className={`post-card ${featured ? 'featured' : ''}`} data-testid="post-card">
            {post.cover_image && (
                <div className="post-card-image-wrapper">
                    <img src={post.cover_image} alt={post.title} className="post-card-image" />
                </div>
            )}

            <div className="post-card-content">
                <div className="post-card-meta">
                    <span className="post-card-date">{formatDate(post.created_at)}</span>
                    <span>by {post.author_name}</span>
                </div>

                <Link href={`/blog/${post.slug}`}>
                    <h3 className="post-card-title">{post.title}</h3>
                </Link>

                {post.excerpt && (
                    <p className="post-card-excerpt">{post.excerpt}</p>
                )}

                <div className="post-card-footer">
                    <div className="post-card-stats">
                        <span className="post-card-stat">❤️ {post.like_count}</span>
                        <span className="post-card-stat">💬 {post.comment_count}</span>
                    </div>

                    <Link href={`/blog/${post.slug}`} className="post-card-read-more">
                        Read More →
                    </Link>
                </div>
            </div>
        </article>
    );
}
