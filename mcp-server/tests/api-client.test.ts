/**
 * Tests for the ApiClient HTTP layer.
 *
 * All tests use mocked fetch — no real network calls.
 */

import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { ApiClient, ApiClientError } from "../src/api-client.js";

// ─── Mock fetch ──────────────────────────────────────────────────────────────

function mockFetch(responses: Array<{ status: number; body?: unknown }>): void {
  let callIndex = 0;
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => {
      const resp = responses[callIndex++];
      if (!resp) throw new Error(`Unexpected fetch call #${callIndex}`);
      return {
        ok: resp.status >= 200 && resp.status < 300,
        status: resp.status,
        statusText: `Status ${resp.status}`,
        headers: {
          get: (name: string) =>
            name.toLowerCase() === "content-type" ? "application/json" : null,
        },
        json: async () => resp.body ?? {},
      } as unknown as Response;
    }),
  );
}

describe("ApiClient", () => {
  const BASE = "http://localhost:8080";

  beforeEach(() => {
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  // ─── Login ───────────────────────────────────────────────────────────────

  it("login() stores token on success", async () => {
    mockFetch([
      { status: 200, body: { token: "jwt-123", user: { id: "1", email: "a@b.c", name: "Admin", role: "admin" } } },
    ]);

    const client = new ApiClient(BASE, "admin@test.com", "pass");
    await client.login();
    expect(client.isAuthenticated).toBe(true);
  });

  it("login() throws if credentials are invalid", async () => {
    mockFetch([{ status: 401, body: { error: "Invalid credentials" } }]);

    const client = new ApiClient(BASE, "bad@test.com", "bad");
    await expect(client.login()).rejects.toThrow(ApiClientError);
  });

  it("login() throws if user is not admin", async () => {
    mockFetch([
      { status: 200, body: { token: "jwt-123", user: { id: "1", email: "u@b.c", name: "User", role: "user" } } },
    ]);

    const client = new ApiClient(BASE, "user@test.com", "pass");
    await expect(client.login()).rejects.toThrow("not an admin");
  });

  it("login() de-duplicates concurrent calls", async () => {
    mockFetch([
      { status: 200, body: { token: "jwt-1", user: { id: "1", email: "a@b.c", name: "A", role: "admin" } } },
    ]);

    const client = new ApiClient(BASE, "a@b.c", "pass");
    // Fire two logins concurrently — should only make one fetch call
    await Promise.all([client.login(), client.login()]);
    expect(fetch).toHaveBeenCalledTimes(1);
  });

  // ─── Authenticated requests ──────────────────────────────────────────────

  it("request() includes Authorization header", async () => {
    mockFetch([
      { status: 200, body: { token: "jwt-abc", user: { id: "1", email: "a@b.c", name: "Admin", role: "admin" } } },
      { status: 200, body: { total_posts: 5 } },
    ]);

    const client = new ApiClient(BASE, "a@b.c", "pass");
    const stats = await client.request<{ total_posts: number }>("GET", "/admin/stats");
    expect(stats.total_posts).toBe(5);

    // Verify the second call had the auth header
    const calls = vi.mocked(fetch).mock.calls;
    const secondCall = calls[1];
    const headers = (secondCall[1] as RequestInit).headers as Record<string, string>;
    expect(headers["Authorization"]).toBe("Bearer jwt-abc");
  });

  it("request() auto-refreshes token on 401", async () => {
    mockFetch([
      // Initial login
      { status: 200, body: { token: "token-1", user: { id: "1", email: "a@b.c", name: "Admin", role: "admin" } } },
      // First request fails with 401 (token expired)
      { status: 401, body: { error: "Token expired" } },
      // Re-login
      { status: 200, body: { token: "token-2", user: { id: "1", email: "a@b.c", name: "Admin", role: "admin" } } },
      // Retry succeeds
      { status: 200, body: { data: [] } },
    ]);

    const client = new ApiClient(BASE, "a@b.c", "pass");
    const result = await client.request<{ data: unknown[] }>("GET", "/admin/posts");
    expect(result.data).toEqual([]);
    expect(fetch).toHaveBeenCalledTimes(4); // login + fail + re-login + retry
  });

  it("request() throws on non-401 error", async () => {
    mockFetch([
      { status: 200, body: { token: "jwt-1", user: { id: "1", email: "a@b.c", name: "Admin", role: "admin" } } },
      { status: 500, body: { error: "Internal server error" } },
    ]);

    const client = new ApiClient(BASE, "a@b.c", "pass");
    await expect(client.request("GET", "/admin/bad")).rejects.toThrow(ApiClientError);
  });

  it("request() sends JSON body for POST", async () => {
    mockFetch([
      { status: 200, body: { token: "jwt-1", user: { id: "1", email: "a@b.c", name: "Admin", role: "admin" } } },
      { status: 201, body: { id: "post-1", title: "Hello" } },
    ]);

    const client = new ApiClient(BASE, "a@b.c", "pass");
    const post = await client.request<{ id: string; title: string }>("POST", "/admin/posts", {
      title: "Hello",
      content: "World",
    });
    expect(post.title).toBe("Hello");

    const calls = vi.mocked(fetch).mock.calls;
    const postCall = calls[1];
    expect((postCall[1] as RequestInit).body).toBe(JSON.stringify({ title: "Hello", content: "World" }));
  });
});
