// API client for Matthew's Galaxy backend

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

interface ApiResponse<T> {
  data?: T;
  error?: string;
}

class ApiClient {
  private token: string | null = null;

  constructor() {
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('token');
    }
  }

  setToken(token: string | null) {
    this.token = token;
    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem('token', token);
      } else {
        localStorage.removeItem('token');
      }
    }
  }

  getToken(): string | null {
    return this.token;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    };

    if (this.token) {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${this.token}`;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1${endpoint}`, {
        ...options,
        headers,
      });

      const data = await response.json();

      if (!response.ok) {
        return { error: data.error || 'An error occurred' };
      }

      return { data };
    } catch (error) {
      console.error('API Error:', error);
      return { error: 'Network error. Please try again.' };
    }
  }

  // Auth
  async register(email: string, password: string, name: string) {
    return this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, name }),
    });
  }

  async login(email: string, password: string) {
    return this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  async getCurrentUser() {
    return this.request<User>('/me');
  }

  async updateProfile(data: { name?: string; bio?: string; avatar_url?: string }) {
    return this.request<User>('/me', {
      method: 'PATCH',
      body: JSON.stringify(data),
    });
  }

  // Posts
  async getPosts(page = 1, limit = 10) {
    return this.request<PaginatedResponse<Post>>(`/posts?page=${page}&limit=${limit}`);
  }

  async getPostBySlug(slug: string) {
    return this.request<Post>(`/posts/${slug}`);
  }

  // Comments
  async getComments(slug: string) {
    return this.request<Comment[]>(`/posts/${slug}/comments`);
  }

  async createComment(slug: string, content: string) {
    return this.request<Comment>(`/posts/${slug}/comments`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    });
  }

  async deleteComment(slug: string, commentId: string) {
    return this.request<void>(`/posts/${slug}/comments/${commentId}`, {
      method: 'DELETE',
    });
  }

  // Likes
  async getLikeStatus(slug: string) {
    return this.request<{ count: number; liked: boolean }>(`/posts/${slug}/likes`);
  }

  async toggleLike(slug: string) {
    return this.request<{ liked: boolean }>(`/posts/${slug}/like`, {
      method: 'POST',
    });
  }

  // Subscription
  async subscribe(email: string) {
    return this.request<{ message: string }>('/subscribe', {
      method: 'POST',
      body: JSON.stringify({ email }),
    });
  }

  // Admin
  async getAdminStats() {
    return this.request<DashboardStats>('/admin/stats');
  }

  async getAllPosts(page = 1, limit = 20) {
    return this.request<PaginatedResponse<Post>>(`/admin/posts?page=${page}&limit=${limit}`);
  }

  async createPost(post: CreatePostRequest) {
    return this.request<Post>('/admin/posts', {
      method: 'POST',
      body: JSON.stringify(post),
    });
  }

  async updatePost(id: string, post: Partial<CreatePostRequest>) {
    return this.request<Post>(`/admin/posts/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(post),
    });
  }

  async deletePost(id: string) {
    return this.request<void>(`/admin/posts/${id}`, {
      method: 'DELETE',
    });
  }

  async getAllUsers(page = 1, limit = 20) {
    return this.request<PaginatedResponse<User>>(`/admin/users?page=${page}&limit=${limit}`);
  }

  async deleteUser(id: string) {
    return this.request<void>(`/admin/users/${id}`, {
      method: 'DELETE',
    });
  }

  async getSubscribers() {
    return this.request<Subscription[]>('/admin/subscribers');
  }
}

// Types
export interface User {
  id: string;
  email: string;
  name: string;
  role: string;
  avatar_url?: string;
  bio?: string;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface Post {
  id: string;
  title: string;
  slug: string;
  content: string;
  excerpt?: string;
  cover_image?: string;
  author_id: string;
  author_name: string;
  author_avatar_url?: string;
  published: boolean;
  like_count: number;
  comment_count: number;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: string;
  post_id: string;
  user_id: string;
  user_name: string;
  user_avatar_url?: string;
  content: string;
  created_at: string;
}

export interface Subscription {
  id: string;
  email: string;
  active: boolean;
  subscribed_at: string;
}

export interface DashboardStats {
  total_users: number;
  total_posts: number;
  total_comments: number;
  total_likes: number;
  total_subscribers: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  excerpt?: string;
  cover_image?: string;
  published: boolean;
}

export const api = new ApiClient();
