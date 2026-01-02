'use client';

import { useState } from 'react';
import { api } from '@/lib/api';

export default function SubscribeForm() {
    const [email, setEmail] = useState('');
    const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');
    const [message, setMessage] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setStatus('loading');

        const { data, error } = await api.subscribe(email);

        if (data) {
            setStatus('success');
            setMessage(data.message);
            setEmail('');
        } else {
            setStatus('error');
            setMessage(error || 'Something went wrong');
        }

        setTimeout(() => {
            setStatus('idle');
            setMessage('');
        }, 5000);
    };

    return (
        <div className="subscribe-container" data-testid="subscribe-form">
            <div className="subscribe-content">
                <div className="subscribe-icon">✉️</div>
                <h3 className="subscribe-title">Stay in Orbit</h3>
                <p className="subscribe-description">
                    Get notified when new posts launch. No spam, only stellar content.
                </p>

                <form onSubmit={handleSubmit} className="subscribe-form">
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="your@email.com"
                        className="subscribe-input"
                        required
                        disabled={status === 'loading'}
                        data-testid="subscribe-input"
                    />
                    <button
                        type="submit"
                        className="subscribe-button"
                        disabled={status === 'loading'}
                        data-testid="subscribe-button"
                    >
                        {status === 'loading' ? '...' : 'Subscribe 🚀'}
                    </button>
                </form>

                {message && (
                    <p className={`subscribe-message ${status}`} data-testid="subscribe-message">
                        {status === 'success' ? '✨ ' : '⚠️ '}{message}
                    </p>
                )}
            </div>
        </div>
    );
}
