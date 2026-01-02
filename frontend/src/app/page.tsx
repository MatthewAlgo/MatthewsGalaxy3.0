'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { api, Post } from '@/lib/api';
import PostCard from '@/components/PostCard';
import SubscribeForm from '@/components/SubscribeForm';

export default function Home() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getPosts(1, 4).then(({ data }) => {
      if (data) {
        setPosts(data.data);
      }
      setLoading(false);
    });
  }, []);

  return (
    <div style={{ flex: 1 }}>
      {/* Hero Section */}
      <section className="home-hero">
        <div className="home-hero-content">
          <div className="home-hero-emoji">🚀</div>
          <h1 className="home-hero-title fade-in-up">
            Welcome to <span className="gradient-text">Matthew&apos;s Galaxy</span>
          </h1>
          <p className="home-hero-subtitle fade-in-up fade-in-delay-1">
            Exploring the frontiers of technology, cloud architecture, and software engineering.
            Join me on this cosmic journey through code.
          </p>
          <div className="home-hero-cta fade-in-up fade-in-delay-2">
            <Link href="/blog" className="btn-primary">
              Explore Posts 🌟
            </Link>
            <Link href="/about" className="btn-secondary">
              About Me
            </Link>
          </div>
        </div>

        <div className="home-hero-stats">
          <div className="home-stat">
            <span className="home-stat-icon">💼</span>
            <span className="home-stat-label">Experience at</span>
            <span className="home-stat-value">Google, Snowflake, NXP</span>
          </div>
          <div className="home-stat">
            <span className="home-stat-icon">🎓</span>
            <span className="home-stat-label">Thesis Grade</span>
            <span className="home-stat-value">10.0 / 10</span>
          </div>
          <div className="home-stat">
            <span className="home-stat-icon">🏆</span>
            <span className="home-stat-label">NitroNLP Hackathon</span>
            <span className="home-stat-value">1st Place</span>
          </div>
        </div>
      </section>

      {/* Featured Posts */}
      <section className="home-section container">
        <div className="home-section-header">
          <h2>Latest Posts</h2>
          <Link href="/blog" className="home-view-all">View All →</Link>
        </div>

        {loading ? (
          <div className="home-loading">
            <div className="spinner"></div>
          </div>
        ) : posts.length > 0 ? (
          <div className="home-posts-grid">
            {posts.map((post, index) => (
              <PostCard key={post.id} post={post} featured={index === 0} />
            ))}
          </div>
        ) : (
          <div className="home-no-posts">
            <p>No posts yet. Check back soon! ✨</p>
          </div>
        )}
      </section>

      {/* About Preview */}
      <section className="home-about-preview">
        <div className="container">
          <div className="home-about-grid">
            <div className="home-about-content">
              <span className="home-about-label">👋 Hi, I&apos;m Matei</span>
              <h2>Software Engineer & Tech Explorer</h2>
              <p>
                I&apos;m a results-oriented software engineer with experience at top-tier companies
                including <strong>Google</strong>, <strong>Snowflake</strong>, and <strong>NXP Semiconductors</strong>.
                I graduated from the University of Bucharest with a focus on cloud-native technologies,
                backend systems, and distributed computing.
              </p>
              <p>
                Currently based in Warsaw, Poland, I&apos;m passionate about building scalable systems
                and sharing knowledge through this blog.
              </p>
              <Link href="/about" className="btn-primary">
                Learn More About Me →
              </Link>
            </div>
            <div className="home-about-skills">
              <h3>Tech Stack</h3>
              <div className="home-skill-tags">
                <span className="home-skill-tag">Go</span>
                <span className="home-skill-tag">Python</span>
                <span className="home-skill-tag">TypeScript</span>
                <span className="home-skill-tag">Kubernetes</span>
                <span className="home-skill-tag">Docker</span>
                <span className="home-skill-tag">PostgreSQL</span>
                <span className="home-skill-tag">React</span>
                <span className="home-skill-tag">Next.js</span>
                <span className="home-skill-tag">Linux</span>
                <span className="home-skill-tag">C++</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Newsletter */}
      <section className="home-section container">
        <SubscribeForm />
      </section>
    </div>
  );
}
