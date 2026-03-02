#!/usr/bin/env node

/**
 * MatthewsGalaxy Blog Admin MCP Server — Entry Point
 *
 * Reads configuration from environment variables, creates the API client
 * and MCP server, and starts the stdio transport.
 *
 * Required env vars:
 *   MG_API_URL        — Backend API base URL (default: http://localhost:8080)
 *   MG_ADMIN_EMAIL    — Admin account email
 *   MG_ADMIN_PASSWORD — Admin account password
 */

import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { ApiClient } from "./api-client.js";
import { createServer } from "./server.js";

function getRequiredEnv(name: string): string {
  const value = process.env[name];
  if (!value) {
    console.error(`FATAL: ${name} environment variable is required but not set`);
    process.exit(1);
  }
  return value;
}

async function main(): Promise<void> {
  const apiUrl = process.env.MG_API_URL ?? "http://localhost:8080";
  const adminEmail = getRequiredEnv("MG_ADMIN_EMAIL");
  const adminPassword = getRequiredEnv("MG_ADMIN_PASSWORD");

  // Create API client and verify connectivity
  const apiClient = new ApiClient(apiUrl, adminEmail, adminPassword);

  try {
    await apiClient.login();
    console.error("[matthewsgalaxy-mcp] Authenticated as admin successfully");
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    console.error(`[matthewsgalaxy-mcp] Failed to authenticate: ${message}`);
    console.error("[matthewsgalaxy-mcp] Server will attempt re-auth on first tool call");
  }

  // Create and start the MCP server
  const server = createServer(apiClient);
  const transport = new StdioServerTransport();
  await server.connect(transport);

  console.error("[matthewsgalaxy-mcp] MCP server running on stdio");
}

main().catch((error) => {
  console.error("[matthewsgalaxy-mcp] Fatal error:", error);
  process.exit(1);
});
