import Link from 'next/link';

export default function Footer() {
    return (
        <footer className="footer" data-testid="footer">
            <div className="container">
                <div className="footer-grid">
                    <div className="footer-brand">
                        <Link href="/" className="footer-logo">
                            <span className="footer-logo-icon">🌌</span>
                            <span className="footer-logo-text">Matthew&apos;s Galaxy</span>
                        </Link>
                        <p className="footer-tagline">
                            Exploring the frontiers of technology, one post at a time.
                        </p>
                        <p className="footer-quote">
                            <em>&quot;Ad astra per aspera&quot;</em> — Through hardships to the stars
                        </p>
                    </div>

                    <div className="footer-links">
                        <h4>Navigate</h4>
                        <Link href="/">Home</Link>
                        <Link href="/blog">Blog</Link>
                        <Link href="/about">About</Link>
                    </div>

                    <div className="footer-links">
                        <h4>Connect</h4>
                        <a href="https://www.linkedin.com/in/username" target="_blank" rel="noopener noreferrer">
                            LinkedIn
                        </a>
                        <a href="https://github.com/username" target="_blank" rel="noopener noreferrer">
                            GitHub
                        </a>
                        <a href="mailto:admin@example.com">
                            Email
                        </a>
                    </div>

                    <div className="footer-links">
                        <h4>Legal</h4>
                        <Link href="/privacy">Privacy Policy</Link>
                        <Link href="/terms">Terms of Service</Link>
                    </div>
                </div>

                <div className="footer-bottom">
                    <p>© {new Date().getFullYear()} Matthew&apos;s Galaxy. Built with ❤️ and ☕</p>
                    <div className="footer-tech">
                        <span>Next.js</span>
                        <span>•</span>
                        <span>Go</span>
                        <span>•</span>
                        <span>PostgreSQL</span>
                    </div>
                </div>
            </div>
        </footer>
    );
}
