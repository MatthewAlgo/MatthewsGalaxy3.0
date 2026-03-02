/**
 * Tests for the MCP server tool definitions.
 *
 * Tests the createServer() output by invoking tools through the MCP client protocol.
 * ApiClient is mocked — we test that tools correctly forward parameters and format responses.
 */

import { describe, it, expect, vi, beforeEach } from "vitest";
import { createServer } from "../src/server.js";
import { ApiClient } from "../src/api-client.js";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { InMemoryTransport } from "@modelcontextprotocol/sdk/inMemory.js";

// ─── Test Helpers ────────────────────────────────────────────────────────────

function createMockApiClient(): ApiClient {
  return {
    login: vi.fn(),
    request: vi.fn(),
    isAuthenticated: true,
  } as unknown as ApiClient;
}

async function setupServerAndClient(mockApi: ApiClient) {
  const server = createServer(mockApi);
  const client = new Client({ name: "test-client", version: "1.0.0" });

  const [clientTransport, serverTransport] = InMemoryTransport.createLinkedPair();
  await Promise.all([
    server.connect(serverTransport),
    client.connect(clientTransport),
  ]);

  return { server, client };
}

// ─── Tests ───────────────────────────────────────────────────────────────────

describe("MCP Server Tools", () => {
  let mockApi: ApiClient;

  beforeEach(() => {
    mockApi = createMockApiClient();
  });

  // ─── Blog Management ──────────────────────────────────────────────────

  describe("list_posts", () => {
    it("returns paginated posts", async () => {
      const mockResponse = {
        data: [{ id: "1", title: "Post 1", slug: "post-1", published: true }],
        page: 1,
        limit: 20,
        total: 1,
        total_pages: 1,
      };
      vi.mocked(mockApi.request).mockResolvedValue(mockResponse);

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({ name: "list_posts", arguments: { page: 1, limit: 20 } });

      expect(mockApi.request).toHaveBeenCalledWith("GET", "/admin/posts?page=1&limit=20");
      expect(result.isError).toBeFalsy();
      const text = (result.content as Array<{ type: string; text: string }>)[0].text;
      const parsed = JSON.parse(text);
      expect(parsed.data).toHaveLength(1);
      expect(parsed.data[0].title).toBe("Post 1");
    });
  });

  describe("get_post", () => {
    it("returns a single post by slug", async () => {
      const mockPost = { id: "1", title: "Welcome", slug: "welcome", content: "Hello" };
      vi.mocked(mockApi.request).mockResolvedValue(mockPost);

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({ name: "get_post", arguments: { slug: "welcome" } });

      expect(mockApi.request).toHaveBeenCalledWith("GET", "/posts/welcome");
      expect(result.isError).toBeFalsy();
    });

    it("encodes special characters in slug", async () => {
      vi.mocked(mockApi.request).mockResolvedValue({});

      const { client } = await setupServerAndClient(mockApi);
      await client.callTool({ name: "get_post", arguments: { slug: "hello world" } });

      expect(mockApi.request).toHaveBeenCalledWith("GET", "/posts/hello%20world");
    });
  });

  describe("create_post", () => {
    it("creates a draft post with required fields", async () => {
      const mockPost = { id: "new-1", title: "New Post", slug: "new-post", published: false };
      vi.mocked(mockApi.request).mockResolvedValue(mockPost);

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({
        name: "create_post",
        arguments: { title: "New Post", content: "Some content here." },
      });

      expect(mockApi.request).toHaveBeenCalledWith("POST", "/admin/posts", {
        title: "New Post",
        content: "Some content here.",
        excerpt: undefined,
        cover_image: undefined,
        published: false,
      });
      expect(result.isError).toBeFalsy();
    });

    it("creates a published post with all optional fields", async () => {
      vi.mocked(mockApi.request).mockResolvedValue({ id: "new-2" });

      const { client } = await setupServerAndClient(mockApi);
      await client.callTool({
        name: "create_post",
        arguments: {
          title: "Full Post",
          content: "Full content here.",
          excerpt: "A summary",
          cover_image: "https://example.com/img.jpg",
          published: true,
        },
      });

      expect(mockApi.request).toHaveBeenCalledWith("POST", "/admin/posts", {
        title: "Full Post",
        content: "Full content here.",
        excerpt: "A summary",
        cover_image: "https://example.com/img.jpg",
        published: true,
      });
    });
  });

  describe("update_post", () => {
    it("sends only the provided update fields", async () => {
      vi.mocked(mockApi.request).mockResolvedValue({ id: "post-1", title: "Updated" });

      const { client } = await setupServerAndClient(mockApi);
      await client.callTool({
        name: "update_post",
        arguments: { id: "550e8400-e29b-41d4-a716-446655440000", title: "Updated" },
      });

      expect(mockApi.request).toHaveBeenCalledWith(
        "PATCH",
        "/admin/posts/550e8400-e29b-41d4-a716-446655440000",
        { title: "Updated" },
      );
    });
  });

  describe("delete_post", () => {
    it("deletes a post by ID", async () => {
      vi.mocked(mockApi.request).mockResolvedValue({ message: "Post deleted successfully" });

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({
        name: "delete_post",
        arguments: { id: "550e8400-e29b-41d4-a716-446655440000" },
      });

      expect(mockApi.request).toHaveBeenCalledWith(
        "DELETE",
        "/admin/posts/550e8400-e29b-41d4-a716-446655440000",
      );
      expect(result.isError).toBeFalsy();
    });
  });

  // ─── Dashboard & Users ────────────────────────────────────────────────

  describe("get_dashboard_stats", () => {
    it("returns aggregate statistics", async () => {
      const mockStats = {
        total_users: 10,
        total_posts: 5,
        total_comments: 20,
        total_likes: 50,
        total_subscribers: 8,
      };
      vi.mocked(mockApi.request).mockResolvedValue(mockStats);

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({ name: "get_dashboard_stats", arguments: {} });

      expect(mockApi.request).toHaveBeenCalledWith("GET", "/admin/stats");
      const parsed = JSON.parse(
        (result.content as Array<{ type: string; text: string }>)[0].text,
      );
      expect(parsed.total_posts).toBe(5);
    });
  });

  describe("list_subscribers", () => {
    it("returns active subscribers", async () => {
      const mockSubs = [
        { id: "s1", email: "a@b.c", active: true, subscribed_at: "2025-01-01T00:00:00Z" },
      ];
      vi.mocked(mockApi.request).mockResolvedValue(mockSubs);

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({ name: "list_subscribers", arguments: {} });

      expect(mockApi.request).toHaveBeenCalledWith("GET", "/admin/subscribers");
      expect(result.isError).toBeFalsy();
    });
  });

  // ─── Analytics ────────────────────────────────────────────────────────

  describe("get_post_engagement", () => {
    it("sorts posts by engagement descending", async () => {
      const mockResponse = {
        data: [
          { title: "Low", slug: "low", published: true, like_count: 1, comment_count: 0, created_at: "2025-01-01" },
          { title: "High", slug: "high", published: true, like_count: 10, comment_count: 5, created_at: "2025-01-02" },
          { title: "Mid", slug: "mid", published: true, like_count: 3, comment_count: 2, created_at: "2025-01-03" },
        ],
        page: 1,
        limit: 20,
        total: 3,
        total_pages: 1,
      };
      vi.mocked(mockApi.request).mockResolvedValue(mockResponse);

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({
        name: "get_post_engagement",
        arguments: { sort: "most_engaged" },
      });

      const parsed = JSON.parse(
        (result.content as Array<{ type: string; text: string }>)[0].text,
      );
      expect(parsed.data[0].title).toBe("High");
      expect(parsed.data[0].total_engagement).toBe(15);
      expect(parsed.data[1].title).toBe("Mid");
      expect(parsed.data[2].title).toBe("Low");
    });

    it("sorts posts by engagement ascending when least_engaged", async () => {
      const mockResponse = {
        data: [
          { title: "High", slug: "high", published: true, like_count: 10, comment_count: 5, created_at: "2025-01-01" },
          { title: "Low", slug: "low", published: true, like_count: 0, comment_count: 0, created_at: "2025-01-02" },
        ],
        page: 1,
        limit: 20,
        total: 2,
        total_pages: 1,
      };
      vi.mocked(mockApi.request).mockResolvedValue(mockResponse);

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({
        name: "get_post_engagement",
        arguments: { sort: "least_engaged" },
      });

      const parsed = JSON.parse(
        (result.content as Array<{ type: string; text: string }>)[0].text,
      );
      expect(parsed.data[0].title).toBe("Low");
      expect(parsed.data[0].total_engagement).toBe(0);
    });
  });

  describe("get_recent_email_logs", () => {
    it("returns paginated email logs", async () => {
      const mockLogs = {
        data: [
          { id: "log-1", subscriber_email: "user@test.com", sent_at: "2025-01-01T00:00:00Z", status: "sent" },
        ],
        page: 1,
        limit: 50,
        total: 1,
        total_pages: 1,
      };
      vi.mocked(mockApi.request).mockResolvedValue(mockLogs);

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({ name: "get_recent_email_logs", arguments: {} });

      expect(mockApi.request).toHaveBeenCalledWith("GET", "/admin/email-logs?page=1&limit=50");
      expect(result.isError).toBeFalsy();
    });
  });

  // ─── Error handling ───────────────────────────────────────────────────

  describe("error handling", () => {
    it("returns isError=true when API call fails", async () => {
      const { ApiClientError } = await import("../src/api-client.js");
      vi.mocked(mockApi.request).mockRejectedValue(
        new ApiClientError("Post not found", 404, "/posts/missing"),
      );

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({ name: "get_post", arguments: { slug: "missing" } });

      expect(result.isError).toBe(true);
      const text = (result.content as Array<{ type: string; text: string }>)[0].text;
      expect(text).toContain("404");
      expect(text).toContain("Post not found");
    });

    it("handles unexpected errors gracefully", async () => {
      vi.mocked(mockApi.request).mockRejectedValue(new Error("Network failure"));

      const { client } = await setupServerAndClient(mockApi);
      const result = await client.callTool({ name: "get_dashboard_stats", arguments: {} });

      expect(result.isError).toBe(true);
      const text = (result.content as Array<{ type: string; text: string }>)[0].text;
      expect(text).toContain("Network failure");
    });
  });
});
