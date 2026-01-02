-- Matthew's Galaxy Database Schema
-- Initialize database with tables and seed admin user

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    avatar_url VARCHAR(500),
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Posts table
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(500) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    cover_image VARCHAR(500),
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    published BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Likes table
CREATE TABLE IF NOT EXISTS likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(post_id, user_id)
);

-- Subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    active BOOLEAN DEFAULT true,
    subscribed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    unsubscribed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(email)
);

-- Email notifications log
CREATE TABLE IF NOT EXISTS email_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID REFERENCES posts(id) ON DELETE SET NULL,
    subscriber_email VARCHAR(255) NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'sent'
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author_id);
CREATE INDEX IF NOT EXISTS idx_posts_published ON posts(published);
CREATE INDEX IF NOT EXISTS idx_posts_created ON posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_likes_post ON likes(post_id);
CREATE INDEX IF NOT EXISTS idx_likes_user ON likes(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_active ON subscriptions(active);

-- Insert admin user (Matei) - password is 'admin123' hashed with bcrypt
INSERT INTO users (email, password_hash, name, role, bio) VALUES (
    'admin@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.NDRJPLkUKaSZ3.lYpO',
    'Matei-Alexandru Dinu',
    'admin',
    'Ex-Software Engineer Intern @ Google | Ex-Software Engineer Intern @ Snowflake | Ex-SWE & Intern @ NXP Semiconductors | UoB CS ''25. Passionate about cloud-native technologies, backend systems, and open-source development.'
) ON CONFLICT (email) DO NOTHING;

-- Insert sample welcome post
INSERT INTO posts (title, slug, content, excerpt, author_id, published) 
SELECT 
    'Welcome to Matthew''s Galaxy',
    'welcome-to-matthews-galaxy',
    E'# Welcome to Matthew''s Galaxy! 🌌\n\nHello, and welcome to my corner of the internet! I''m **Matei-Alexandru Dinu**, a software engineer passionate about building scalable systems and exploring the frontiers of technology.\n\n## What You''ll Find Here\n\nThis blog is my space to share:\n\n- **Technical Deep Dives**: Explorations of cloud architecture, distributed systems, and backend development\n- **Career Insights**: Lessons learned from my internships at Google, Snowflake, and NXP Semiconductors\n- **Project Showcases**: Walkthroughs of personal projects and hackathon adventures\n- **Thoughts & Reflections**: Musings on tech, life, and everything in between\n\n## My Journey So Far\n\nFrom winning the NitroNLP AI Hackathon to building internal tools at Google, every step has been a learning experience. I graduated from the University of Bucharest with a 10.0/10 thesis on "Scalable Cloud Platforms for Generative AI Models."\n\n## Let''s Connect!\n\nFeel free to explore, comment, and reach out. Whether you''re a fellow developer, a recruiter, or just curious - you''re welcome here among the stars! ⭐\n\n*Ad astra per aspera* - Through hardships to the stars.',
    'Hello, and welcome to my corner of the internet! I''m Matei-Alexandru Dinu, a software engineer passionate about building scalable systems.',
    id,
    true
FROM users WHERE email = 'admin@example.com'
ON CONFLICT (slug) DO NOTHING;
