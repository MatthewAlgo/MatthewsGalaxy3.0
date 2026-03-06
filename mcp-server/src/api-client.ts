/**
 * HTTP client for the MatthewsGalaxy backend API.
 *
 * This is the single network layer — all MCP tools call through this client.
 * Handles authentication, token refresh on 401, and typed error responses.
 */

/** Shape returned by the backend on errors */
interface ApiError {
  error: string;
}

/** Auth response from POST /auth/login */
interface AuthResponse {
  token: string;
  user: { id: string; email: string; name: string; role: string };
}

export class ApiClientError extends Error {
  constructor(
    message: string,
    public readonly statusCode: number,
    public readonly endpoint: string,
  ) {
    super(message);
    this.name = "ApiClientError";
  }
}

export class ApiClient {
  private token: string | null = null;
  private loginInProgress: Promise<void> | null = null;

  constructor(
    private readonly baseUrl: string,
    private readonly adminEmail: string,
    private readonly adminPassword: string,
  ) {}

  /**
   * Authenticate against the backend and store the JWT.
   * De-duplicates concurrent login calls.
   */
  async login(): Promise<void> {
    // If a login is already in progress, wait for it instead of firing another
    if (this.loginInProgress) {
      return this.loginInProgress;
    }

    this.loginInProgress = this.performLogin();
    try {
      await this.loginInProgress;
    } finally {
      this.loginInProgress = null;
    }
  }

  private async performLogin(): Promise<void> {
    const response = await fetch(`${this.baseUrl}/api/v1/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: this.adminEmail,
        password: this.adminPassword,
      }),
    });

    if (!response.ok) {
      const body = (await response.json().catch(() => ({}))) as ApiError;
      throw new ApiClientError(
        `Login failed: ${body.error ?? response.statusText}`,
        response.status,
        "/api/v1/auth/login",
      );
    }

    const data = (await response.json()) as AuthResponse;
    this.token = data.token;

    if (data.user.role !== "admin") {
      throw new ApiClientError(
        "Authenticated user is not an admin",
        403,
        "/api/v1/auth/login",
      );
    }
  }

  /**
   * Make an authenticated request to the backend API.
   * Automatically retries once on 401 (token expired) by re-authenticating.
   */
  async request<T>(
    method: string,
    path: string,
    body?: unknown,
  ): Promise<T> {
    // Ensure we're logged in
    if (!this.token) {
      await this.login();
    }

    const result = await this.doRequest<T>(method, path, body);

    // If 401, refresh token and retry exactly once
    if (result.status === 401) {
      this.token = null;
      await this.login();
      const retry = await this.doRequest<T>(method, path, body);
      if (retry.status !== undefined && retry.status >= 400) {
        throw new ApiClientError(
          retry.errorMessage ?? "Request failed after re-auth",
          retry.status,
          path,
        );
      }
      return retry.data as T;
    }

    if (result.status !== undefined && result.status >= 400) {
      throw new ApiClientError(
        result.errorMessage ?? "Request failed",
        result.status,
        path,
      );
    }

    return result.data as T;
  }

  private async doRequest<T>(
    method: string,
    path: string,
    body?: unknown,
  ): Promise<{ data?: T; status: number; errorMessage?: string }> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    const options: RequestInit = { method, headers };
    if (body !== undefined) {
      options.body = JSON.stringify(body);
    }

    const response = await fetch(`${this.baseUrl}/api/v1${path}`, options);

    // Handle non-JSON responses (e.g., 204 No Content)
    const contentType = response.headers.get("content-type");
    if (!contentType?.includes("application/json")) {
      return { status: response.status };
    }

    const json = await response.json();

    if (!response.ok) {
      return {
        status: response.status,
        errorMessage: (json as ApiError).error ?? response.statusText,
      };
    }

    return { data: json as T, status: response.status };
  }

  /** Check if the client has a valid token (for health checks) */
  get isAuthenticated(): boolean {
    return this.token !== null;
  }
}
