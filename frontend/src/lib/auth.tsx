'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { api, User } from './api';

interface AuthContextType {
    user: User | null;
    loading: boolean;
    login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
    register: (email: string, password: string, name: string) => Promise<{ success: boolean; error?: string }>;
    logout: () => void;
    isAdmin: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export { AuthContext };

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Check for existing token
        const token = api.getToken();
        if (token) {
            api.getCurrentUser().then(({ data, error }) => {
                if (data) {
                    setUser(data);
                } else {
                    api.setToken(null);
                }
                setLoading(false);
            });
        } else {
            setLoading(false);
        }
    }, []);

    const login = async (email: string, password: string) => {
        const { data, error } = await api.login(email, password);
        if (data) {
            api.setToken(data.token);
            setUser(data.user);
            return { success: true };
        }
        return { success: false, error };
    };

    const register = async (email: string, password: string, name: string) => {
        const { data, error } = await api.register(email, password, name);
        if (data) {
            api.setToken(data.token);
            setUser(data.user);
            return { success: true };
        }
        return { success: false, error };
    };

    const logout = () => {
        api.setToken(null);
        setUser(null);
    };

    const isAdmin = user?.role === 'admin';

    return (
        <AuthContext.Provider value={{ user, loading, login, register, logout, isAdmin }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}
