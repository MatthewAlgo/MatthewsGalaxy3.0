'use client';

import Link from 'next/link';
import { useState } from 'react';
import { useAuth } from '@/lib/auth';

export default function Header() {
    const { user, logout, isAdmin } = useAuth();
    const [menuOpen, setMenuOpen] = useState(false);

    return (
        <header className="header" data-testid="header">
            <div className="header-container">
                <Link href="/" className="header-logo">
                    <span className="header-logo-icon">🌌</span>
                    <span className="header-logo-text">Matthew&apos;s Galaxy</span>
                </Link>

                <button
                    className="header-menu-toggle"
                    onClick={() => setMenuOpen(!menuOpen)}
                    aria-label="Toggle menu"
                    data-testid="menu-toggle"
                >
                    <span className={menuOpen ? 'header-menu-close' : 'header-menu-icon'}></span>
                </button>

                <nav className={`header-nav ${menuOpen ? 'open' : ''}`} data-testid="nav">
                    <Link href="/" className="header-nav-link" onClick={() => setMenuOpen(false)}>
                        Home
                    </Link>
                    <Link href="/blog" className="header-nav-link" onClick={() => setMenuOpen(false)}>
                        Blog
                    </Link>
                    <Link href="/about" className="header-nav-link" onClick={() => setMenuOpen(false)}>
                        About
                    </Link>

                    {user ? (
                        <>
                            {isAdmin && (
                                <Link href="/admin" className="header-nav-link" onClick={() => setMenuOpen(false)}>
                                    Admin
                                </Link>
                            )}
                            <div className="header-user-menu">
                                <span className="header-user-name">👋 {user.name}</span>
                                <button onClick={logout} className="header-logout-btn" data-testid="logout-btn">
                                    Logout
                                </button>
                            </div>
                        </>
                    ) : (
                        <Link href="/auth/login" className="header-login-btn" onClick={() => setMenuOpen(false)}>
                            Sign In
                        </Link>
                    )}
                </nav>
            </div>
        </header>
    );
}
