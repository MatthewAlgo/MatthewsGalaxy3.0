import { api } from '../api';

describe('ApiClient', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (global.fetch as jest.Mock).mockReset();
    localStorage.clear();
  });

  describe('setToken / getToken', () => {
    it('sets and gets token correctly', () => {
      api.setToken('test-token-123');
      expect(api.getToken()).toBe('test-token-123');
    });

    it('stores token in localStorage', () => {
      api.setToken('stored-token');
      expect(localStorage.setItem).toHaveBeenCalledWith('token', 'stored-token');
    });

    it('removes token when null is passed', () => {
      api.setToken(null);
      expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    });
  });

  describe('register', () => {
    it('makes POST request to correct endpoint', async () => {
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ token: 'new-token', user: { id: '1' } }),
      });

      await api.register('test@test.com', 'password123', 'Test User');

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/auth/register'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            email: 'test@test.com',
            password: 'password123',
            name: 'Test User',
          }),
        })
      );
    });

    it('returns user data on success', async () => {
      const mockUser = { id: '1', email: 'test@test.com', name: 'Test' };
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ token: 'token', user: mockUser }),
      });

      const result = await api.register('test@test.com', 'pass', 'Test');
      expect(result.data?.user).toEqual(mockUser);
    });

    it('returns error on failed registration', async () => {
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: false,
        json: () => Promise.resolve({ error: 'Email already exists' }),
      });

      const result = await api.register('test@test.com', 'pass', 'Test');
      expect(result.error).toBe('Email already exists');
    });
  });

  describe('login', () => {
    it('makes POST request to correct endpoint', async () => {
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ token: 'token', user: {} }),
      });

      await api.login('test@test.com', 'password123');

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/auth/login'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            email: 'test@test.com',
            password: 'password123',
          }),
        })
      );
    });

    it('returns user data on successful login', async () => {
      const mockUser = { id: '1', name: 'User' };
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ token: 'login-token', user: mockUser }),
      });

      const result = await api.login('test@test.com', 'pass');
      expect(result.data?.user).toEqual(mockUser);
    });
  });

  describe('getPosts', () => {
    it('makes GET request with pagination params', async () => {
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [], total: 0, page: 1, total_pages: 0 }),
      });

      await api.getPosts(2, 10);

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/posts?page=2&limit=10'),
        expect.any(Object)
      );
    });

    it('returns posts data on success', async () => {
      const mockPosts = [{ id: '1', title: 'Test Post' }];
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: mockPosts, total: 1, page: 1, total_pages: 1 }),
      });

      const result = await api.getPosts(1, 10);
      expect(result.data?.data).toEqual(mockPosts);
    });
  });

  describe('getPostBySlug', () => {
    it('makes GET request with correct slug', async () => {
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ id: '1', slug: 'test-post' }),
      });

      await api.getPostBySlug('test-post');

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/posts/test-post'),
        expect.any(Object)
      );
    });
  });

  describe('subscribe', () => {
    it('makes POST request with email', async () => {
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ message: 'Subscribed!' }),
      });

      await api.subscribe('test@example.com');

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/subscribe'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ email: 'test@example.com' }),
        })
      );
    });
  });

  describe('createComment', () => {
    it('makes POST request with content', async () => {
      api.setToken('auth-token');
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ id: '1', content: 'Great post!' }),
      });

      await api.createComment('test-slug', 'Great post!');

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/posts/test-slug/comments'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ content: 'Great post!' }),
        })
      );
    });

    it('includes authorization header', async () => {
      api.setToken('auth-token');
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({}),
      });

      await api.createComment('test-slug', 'Comment');

      expect(global.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer auth-token',
          }),
        })
      );
    });
  });

  describe('toggleLike', () => {
    it('makes POST request to like endpoint', async () => {
      api.setToken('token');
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ liked: true }),
      });

      await api.toggleLike('test-post');

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/posts/test-post/like'),
        expect.objectContaining({ method: 'POST' })
      );
    });
  });

  describe('getAdminStats', () => {
    it('makes GET request to admin stats endpoint', async () => {
      api.setToken('admin-token');
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({
          total_users: 100,
          total_posts: 50,
          total_comments: 200,
          total_likes: 500,
          total_subscribers: 25,
        }),
      });

      const result = await api.getAdminStats();

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/admin/stats'),
        expect.any(Object)
      );
      expect(result.data?.total_users).toBe(100);
    });
  });

  describe('createPost', () => {
    it('makes POST request with post data', async () => {
      api.setToken('admin-token');
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ id: '1', title: 'New Post' }),
      });

      await api.createPost({
        title: 'New Post',
        content: 'Content here',
        published: true,
      });

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/admin/posts'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('New Post'),
        })
      );
    });
  });

  describe('error handling', () => {
    it('handles network errors gracefully', async () => {
      // Suppress console.error for this test
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      
      (global.fetch as jest.Mock).mockRejectedValue(new Error('Network error'));

      const result = await api.getPosts(1, 10);
      expect(result.error).toContain('Network error');
      
      consoleSpy.mockRestore();
    });

    it('handles JSON parse errors', async () => {
      // Suppress console.error for this test
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.reject(new Error('Invalid JSON')),
      });

      const result = await api.getPosts(1, 10);
      expect(result.error).toBeDefined();
      
      consoleSpy.mockRestore();
    });
  });
});
