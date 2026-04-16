import { renderHook, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";

vi.mock("../api/client", () => ({
  apiClient: { GET: vi.fn() },
  setTokenProvider: vi.fn(),
}));

import { useUserinfo } from "./useUserinfo";
import { apiClient } from "../api/client";

describe("useUserinfo", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("returns userinfo on successful fetch", async () => {
    vi.mocked(apiClient.GET).mockResolvedValue({
      data: { sub: "user|123", name: "Alice", email: "alice@example.com", picture: null },
      error: undefined,
      response: new Response(),
    });

    const { result } = renderHook(() => useUserinfo());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.userinfo).toEqual({
      sub: "user|123",
      name: "Alice",
      email: "alice@example.com",
      picture: null,
    });
    expect(result.current.error).toBeNull();
  });

  it("sets error when fetch fails", async () => {
    vi.mocked(apiClient.GET).mockResolvedValue({
      data: undefined,
      error: { message: "unauthorized" },
      response: new Response(null, { status: 401 }),
    });

    const { result } = renderHook(() => useUserinfo());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.userinfo).toBeNull();
    expect(result.current.error).toBeTruthy();
  });
});
