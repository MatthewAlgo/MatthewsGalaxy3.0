/**
 * MatthewsGalaxy Blog Admin MCP Server
 *
 * Defines all MCP tools for blog administration and analytics.
 * Tools are thin wrappers over the existing REST API — no new business logic.
 */

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { z } from "zod";
import { ApiClient, ApiClientError } from "./api-client.js";

// ─── Response Types ──────────────────────────────────────────────────────────

interface Post {
  id: string;
  title: string;
  slug: string;
  content: string;
  excerpt?: string;
  cover_image?: string;
  author_id: string;
  author_name?: string;
  published: boolean;
  like_count?: number;
  comment_count?: number;
  created_at: string;
  updated_at: string;
}

interface PaginatedResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

interface DashboardStats {
  total_users: number;
  total_posts: number;
  total_comments: number;
  total_likes: number;
  total_subscribers: number;
}

interface User {
  id: string;
  email: string;
  name: string;
  role: string;
  avatar_url?: string;
  bio?: string;
  created_at: string;
}

interface Subscriber {
  id: string;
  email: string;
  active: boolean;
  subscribed_at: string;
}

interface EmailLog {
  id: string;
  post_id?: string;
  subscriber_email: string;
  sent_at: string;
  status: string;
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

/**
 * Wraps a tool handler to catch ApiClientError and return MCP-formatted errors.
 * This is the single place we handle API errors for tools (DRY).
 */
function formatResult(content: unknown): { content: Array<{ type: "text"; text: string }> } {
  const text = typeof content === "string" ? content : JSON.stringify(content, null, 2);
  return { content: [{ type: "text" as const, text }] };
}

async function handleToolCall<T>(fn: () => Promise<T>): Promise<{ content: Array<{ type: "text"; text: string }>; isError?: boolean }> {
  try {
    const result = await fn();
    return formatResult(result);
  } catch (error) {
    if (error instanceof ApiClientError) {
      return {
        ...formatResult(`Error (${error.statusCode}): ${error.message} [${error.endpoint}]`),
        isError: true,
      };
    }
    const message = error instanceof Error ? error.message : String(error);
    return {
      ...formatResult(`Unexpected error: ${message}`),
      isError: true,
    };
  }
}

// ─── Server Creation ─────────────────────────────────────────────────────────

export function createServer(apiClient: ApiClient): McpServer {
  const server = new McpServer({
    name: "matthewsgalaxy-admin",
    version: "1.0.0",
  });

  // ─── Blog Management Tools ───────────────────────────────────────────────

  server.tool(
    "list_posts",
    "List all blog posts with pagination (includes drafts). Returns title, slug, published status, like/comment counts.",
    {
      page: z.number().int().min(1).default(1).describe("Page number (1-indexed)"),
      limit: z.number().int().min(1).max(100).default(20).describe("Posts per page (max 100)"),
    },
    async ({ page, limit }) =>
      handleToolCall(() =>
        apiClient.request<PaginatedResponse<Post>>("GET", `/admin/posts?page=${page}&limit=${limit}`),
      ),
  );

  server.tool(
    "get_post",
    "Get a single blog post by its slug. Returns full content, author info, and engagement metrics.",
    {
      slug: z.string().min(1).describe("The post's URL slug (e.g., 'welcome-to-matthews-galaxy')"),
    },
    async ({ slug }) =>
      handleToolCall(() => apiClient.request<Post>("GET", `/posts/${encodeURIComponent(slug)}`)),
  );

  server.tool(
    "create_post",
    "Create a new blog post. Set published=false to save as draft.",
    {
      title: z.string().min(3).describe("Post title (min 3 characters)"),
      content: z.string().min(10).describe("Post content in markdown/HTML (min 10 characters)"),
      excerpt: z.string().optional().describe("Short summary for post cards"),
      cover_image: z.string().url().optional().describe("URL to cover image"),
      published: z.boolean().default(false).describe("Publish immediately (true) or save as draft (false)"),
    },
    async (args) =>
      handleToolCall(() =>
        apiClient.request<Post>("POST", "/admin/posts", {
          title: args.title,
          content: args.content,
          excerpt: args.excerpt,
          cover_image: args.cover_image,
          published: args.published,
        }),
      ),
  );

  server.tool(
    "update_post",
    "Update an existing blog post by ID. Only include fields you want to change.",
    {
      id: z.string().uuid().describe("Post UUID"),
      title: z.string().min(3).optional().describe("New title"),
      content: z.string().min(10).optional().describe("New content"),
      excerpt: z.string().optional().describe("New excerpt"),
      cover_image: z.string().url().optional().describe("New cover image URL"),
      published: z.boolean().optional().describe("Change publish status"),
    },
    async ({ id, ...updates }) => {
      // Only send fields that were actually provided
      const body: Record<string, unknown> = {};
      if (updates.title !== undefined) body.title = updates.title;
      if (updates.content !== undefined) body.content = updates.content;
      if (updates.excerpt !== undefined) body.excerpt = updates.excerpt;
      if (updates.cover_image !== undefined) body.cover_image = updates.cover_image;
      if (updates.published !== undefined) body.published = updates.published;

      return handleToolCall(() =>
        apiClient.request<Post>("PATCH", `/admin/posts/${encodeURIComponent(id)}`, body),
      );
    },
  );

  server.tool(
    "delete_post",
    "Permanently delete a blog post by ID. This also deletes all associated comments and likes.",
    {
      id: z.string().uuid().describe("Post UUID to delete"),
    },
    async ({ id }) =>
      handleToolCall(() =>
        apiClient.request<{ message: string }>("DELETE", `/admin/posts/${encodeURIComponent(id)}`),
      ),
  );

  // ─── Dashboard & User Tools ──────────────────────────────────────────────

  server.tool(
    "get_dashboard_stats",
    "Get aggregate blog statistics: total users, posts, comments, likes, and active subscribers.",
    {},
    async () =>
      handleToolCall(() => apiClient.request<DashboardStats>("GET", "/admin/stats")),
  );

  server.tool(
    "list_users",
    "List all registered users with pagination. Shows name, email, role, and join date.",
    {
      page: z.number().int().min(1).default(1).describe("Page number"),
      limit: z.number().int().min(1).max(100).default(20).describe("Users per page"),
    },
    async ({ page, limit }) =>
      handleToolCall(() =>
        apiClient.request<PaginatedResponse<User>>("GET", `/admin/users?page=${page}&limit=${limit}`),
      ),
  );

  server.tool(
    "list_subscribers",
    "List all active email newsletter subscribers.",
    {},
    async () =>
      handleToolCall(() => apiClient.request<Subscriber[]>("GET", "/admin/subscribers")),
  );

  // ─── Analytics Tools (bundled from Issue 2) ──────────────────────────────

  server.tool(
    "get_post_engagement",
    "Get posts sorted by engagement (likes + comments). Useful for finding your most and least popular content.",
    {
      page: z.number().int().min(1).default(1).describe("Page number"),
      limit: z.number().int().min(1).max(100).default(20).describe("Posts per page"),
      sort: z
        .enum(["most_engaged", "least_engaged"])
        .default("most_engaged")
        .describe("Sort order by total engagement (likes + comments)"),
    },
    async ({ page, limit, sort }) =>
      handleToolCall(async () => {
        const result = await apiClient.request<PaginatedResponse<Post>>(
          "GET",
          `/admin/posts?page=${page}&limit=${limit}`,
        );

        // Sort by engagement (like_count + comment_count)
        const sorted = [...result.data].sort((a, b) => {
          const engagementA = (a.like_count ?? 0) + (a.comment_count ?? 0);
          const engagementB = (b.like_count ?? 0) + (b.comment_count ?? 0);
          return sort === "most_engaged"
            ? engagementB - engagementA
            : engagementA - engagementB;
        });

        return {
          ...result,
          data: sorted.map((p) => ({
            title: p.title,
            slug: p.slug,
            published: p.published,
            like_count: p.like_count ?? 0,
            comment_count: p.comment_count ?? 0,
            total_engagement: (p.like_count ?? 0) + (p.comment_count ?? 0),
            created_at: p.created_at,
          })),
        };
      }),
  );

  server.tool(
    "get_recent_email_logs",
    "View recent email delivery logs. Shows which subscribers received emails and their delivery status (sent/failed).",
    {
      page: z.number().int().min(1).default(1).describe("Page number"),
      limit: z.number().int().min(1).max(100).default(50).describe("Logs per page"),
    },
    async ({ page, limit }) =>
      handleToolCall(() =>
        apiClient.request<PaginatedResponse<EmailLog>>(
          "GET",
          `/admin/email-logs?page=${page}&limit=${limit}`,
        ),
      ),
  );

  return server;
}
