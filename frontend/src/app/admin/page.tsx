'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/lib/auth';
import { api, DashboardStats, Post, User, Subscription } from '@/lib/api';

export default function AdminPage() {
    const router = useRouter();
    const { user, isAdmin, loading: authLoading } = useAuth();
    const [stats, setStats] = useState<DashboardStats | null>(null);
    const [posts, setPosts] = useState<Post[]>([]);
    const [users, setUsers] = useState<User[]>([]);
    const [subscribers, setSubscribers] = useState<Subscription[]>([]);
    const [activeTab, setActiveTab] = useState<'overview' | 'posts' | 'users' | 'subscribers'>('overview');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!authLoading && !isAdmin) {
            router.push('/');
        }
    }, [authLoading, isAdmin, router]);

    useEffect(() => {
        if (isAdmin) {
            loadData();
        }
    }, [isAdmin]);

    const loadData = async () => {
        setLoading(true);

        const [statsRes, postsRes, usersRes, subsRes] = await Promise.all([
            api.getAdminStats(),
            api.getAllPosts(1, 10),
            api.getAllUsers(1, 10),
            api.getSubscribers(),
        ]);

        if (statsRes.data) setStats(statsRes.data);
        if (postsRes.data) setPosts(postsRes.data.data);
        if (usersRes.data) setUsers(usersRes.data.data);
        if (subsRes.data) setSubscribers(subsRes.data);

        setLoading(false);
    };

    const handleDeletePost = async (id: string) => {
        if (!confirm('Are you sure you want to delete this post?')) return;
        const { error } = await api.deletePost(id);
        if (!error) {
            setPosts(prev => prev.filter(p => p.id !== id));
        }
    };

    const handleDeleteUser = async (id: string) => {
        if (!confirm('Are you sure you want to delete this user?')) return;
        const { error } = await api.deleteUser(id);
        if (!error) {
            setUsers(prev => prev.filter(u => u.id !== id));
        }
    };

    if (authLoading || loading) {
        return (
            <div className="admin-loading">
                <div className="spinner"></div>
            </div>
        );
    }

    if (!isAdmin) {
        return null;
    }

    return (
        <div className="admin-page">
            <div className="admin-sidebar">
                <h2 className="admin-sidebar-title">⚡ Admin</h2>
                <nav className="admin-sidebar-nav">
                    <button
                        className={`admin-nav-item ${activeTab === 'overview' ? 'active' : ''}`}
                        onClick={() => setActiveTab('overview')}
                    >
                        📊 Overview
                    </button>
                    <button
                        className={`admin-nav-item ${activeTab === 'posts' ? 'active' : ''}`}
                        onClick={() => setActiveTab('posts')}
                    >
                        📝 Posts
                    </button>
                    <button
                        className={`admin-nav-item ${activeTab === 'users' ? 'active' : ''}`}
                        onClick={() => setActiveTab('users')}
                    >
                        👥 Users
                    </button>
                    <button
                        className={`admin-nav-item ${activeTab === 'subscribers' ? 'active' : ''}`}
                        onClick={() => setActiveTab('subscribers')}
                    >
                        ✉️ Subscribers
                    </button>
                </nav>
            </div>

            <div className="admin-content">
                {/* Overview Tab */}
                {activeTab === 'overview' && stats && (
                    <div>
                        <h1>Dashboard Overview</h1>
                        <div className="admin-stats-grid">
                            <div className="admin-stat-card">
                                <span className="admin-stat-icon">👥</span>
                                <div className="admin-stat-value">{stats.total_users}</div>
                                <div className="admin-stat-label">Total Users</div>
                            </div>
                            <div className="admin-stat-card">
                                <span className="admin-stat-icon">📝</span>
                                <div className="admin-stat-value">{stats.total_posts}</div>
                                <div className="admin-stat-label">Total Posts</div>
                            </div>
                            <div className="admin-stat-card">
                                <span className="admin-stat-icon">💬</span>
                                <div className="admin-stat-value">{stats.total_comments}</div>
                                <div className="admin-stat-label">Total Comments</div>
                            </div>
                            <div className="admin-stat-card">
                                <span className="admin-stat-icon">❤️</span>
                                <div className="admin-stat-value">{stats.total_likes}</div>
                                <div className="admin-stat-label">Total Likes</div>
                            </div>
                            <div className="admin-stat-card">
                                <span className="admin-stat-icon">✉️</span>
                                <div className="admin-stat-value">{stats.total_subscribers}</div>
                                <div className="admin-stat-label">Subscribers</div>
                            </div>
                        </div>

                        <div className="admin-quick-actions">
                            <h3>Quick Actions</h3>
                            <Link href="/admin/posts/new" className="btn-primary">
                                ✍️ Create New Post
                            </Link>
                        </div>
                    </div>
                )}

                {/* Posts Tab */}
                {activeTab === 'posts' && (
                    <div>
                        <div className="admin-tab-header">
                            <h1>Manage Posts</h1>
                            <Link href="/admin/posts/new" className="btn-primary">
                                ✍️ New Post
                            </Link>
                        </div>
                        <div className="admin-table">
                            <div className="admin-table-header">
                                <span>Title</span>
                                <span>Status</span>
                                <span>Date</span>
                                <span>Actions</span>
                            </div>
                            {posts.map((post) => (
                                <div key={post.id} className="admin-table-row">
                                    <span className="admin-post-title">{post.title}</span>
                                    <span className={`admin-status ${post.published ? 'published' : 'draft'}`}>
                                        {post.published ? 'Published' : 'Draft'}
                                    </span>
                                    <span className="admin-date">
                                        {new Date(post.created_at).toLocaleDateString()}
                                    </span>
                                    <div className="admin-actions">
                                        <Link href={`/admin/posts/${post.id}`} className="admin-edit-btn">
                                            Edit
                                        </Link>
                                        <button
                                            onClick={() => handleDeletePost(post.id)}
                                            className="admin-delete-btn"
                                        >
                                            Delete
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {/* Users Tab */}
                {activeTab === 'users' && (
                    <div className="admin-users-tab">
                        <h1>Manage Users</h1>
                        <div className="admin-table">
                            <div className="admin-table-header">
                                <span>Name</span>
                                <span>Email</span>
                                <span>Role</span>
                                <span>Joined</span>
                                <span>Actions</span>
                            </div>
                            {users.map((u) => (
                                <div key={u.id} className="admin-table-row">
                                    <span>{u.name}</span>
                                    <span className="admin-email">{u.email}</span>
                                    <span className={`admin-role ${u.role === 'admin' ? 'admin' : ''}`}>
                                        {u.role}
                                    </span>
                                    <span className="admin-date">
                                        {new Date(u.created_at).toLocaleDateString()}
                                    </span>
                                    <div className="admin-actions">
                                        {u.role !== 'admin' && u.id !== user?.id && (
                                            <button
                                                onClick={() => handleDeleteUser(u.id)}
                                                className="admin-delete-btn"
                                            >
                                                Delete
                                            </button>
                                        )}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {/* Subscribers Tab */}
                {activeTab === 'subscribers' && (
                    <div>
                        <h1>Newsletter Subscribers</h1>
                        <p className="admin-sub-count">{subscribers.length} active subscribers</p>
                        <div className="admin-subscriber-list">
                            {subscribers.map((sub) => (
                                <div key={sub.id} className="admin-subscriber">
                                    <span className="admin-sub-email">{sub.email}</span>
                                    <span className="admin-sub-date">
                                        since {new Date(sub.subscribed_at).toLocaleDateString()}
                                    </span>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
