import { renderHook, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import type { ReactNode } from "react";
import { createElement } from "react";

// Mock loadOIDCSetup to avoid real fetch calls.
vi.mock("../oidc", () => ({
  loadOIDCSetup: vi.fn(),
  resetOIDCSetup: vi.fn(),
}));

// Mock apiClient so /userinfo doesn't hit the network.
vi.mock("../api/client", () => ({
  setTokenProvider: vi.fn(),
  apiClient: {
    GET: vi.fn(),
  },
}));

import { AuthProvider } from "../auth/AuthProvider";
import { useAuth } from "./useAuth";
import { loadOIDCSetup } from "../oidc";
import { apiClient } from "../api/client";

const wrapper = ({ children }: { children: ReactNode }) =>
  createElement(AuthProvider, null, children);

const mockUserManager = {
  getUser: vi.fn(),
  signinRedirect: vi.fn(),
  signoutRedirect: vi.fn(),
  events: {
    addUserLoaded: vi.fn(),
    addUserUnloaded: vi.fn(),
  },
};

describe("useAuth", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("is not authenticated when /userinfo returns no data", async () => {
    vi.mocked(loadOIDCSetup).mockResolvedValue({ configured: false, userManager: null });
    vi.mocked(apiClient.GET).mockResolvedValue({
      data: undefined,
      error: { message: "unauthorized" },
      response: new Response(null, { status: 401 }),
    });

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    expect(result.current.isAuthenticated).toBe(false);
  });

  it("is authenticated when /userinfo returns data", async () => {
    vi.mocked(loadOIDCSetup).mockResolvedValue({ configured: false, userManager: null });
    vi.mocked(apiClient.GET).mockResolvedValue({
      data: { sub: "local", name: "dev" },
      error: undefined,
      response: new Response(),
    });

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    expect(result.current.isAuthenticated).toBe(true);
  });

  it("login calls signinRedirect when OIDC is configured", async () => {
    vi.mocked(loadOIDCSetup).mockResolvedValue({
      configured: true,
      userManager: mockUserManager as never,
    });
    vi.mocked(mockUserManager.getUser).mockResolvedValue(null);
    vi.mocked(apiClient.GET).mockResolvedValue({
      data: undefined,
      error: { message: "unauthorized" },
      response: new Response(null, { status: 401 }),
    });
    vi.mocked(mockUserManager.signinRedirect).mockResolvedValue(undefined);

    const { result } = renderHook(() => useAuth(), { wrapper });
    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await result.current.login();

    expect(mockUserManager.signinRedirect).toHaveBeenCalledOnce();
  });

  it("logout calls signoutRedirect when OIDC is configured", async () => {
    vi.mocked(loadOIDCSetup).mockResolvedValue({
      configured: true,
      userManager: mockUserManager as never,
    });
    vi.mocked(mockUserManager.getUser).mockResolvedValue(null);
    vi.mocked(apiClient.GET).mockResolvedValue({
      data: undefined,
      error: { message: "unauthorized" },
      response: new Response(null, { status: 401 }),
    });
    vi.mocked(mockUserManager.signoutRedirect).mockResolvedValue(undefined);

    const { result } = renderHook(() => useAuth(), { wrapper });
    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await result.current.logout();

    expect(mockUserManager.signoutRedirect).toHaveBeenCalledOnce();
  });
});
