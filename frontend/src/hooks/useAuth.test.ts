import { renderHook } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import type { ReactNode } from "react";
import { createElement } from "react";

// Mock oidcConfig before importing AuthProvider
vi.mock("../auth/oidcConfig", () => ({
  oidcDisabled: false,
  userManager: {
    getUser: vi.fn(),
    signinRedirect: vi.fn(),
    signoutRedirect: vi.fn(),
    events: {
      addUserLoaded: vi.fn(),
      addUserUnloaded: vi.fn(),
      removeUserLoaded: vi.fn(),
      removeUserUnloaded: vi.fn(),
    },
  },
}));

// Mock setTokenProvider so it doesn't need a real client
vi.mock("../api/client", () => ({
  setTokenProvider: vi.fn(),
  apiClient: {},
}));

import { AuthProvider } from "../auth/AuthProvider";
import { useAuth } from "./useAuth";
import { userManager } from "../auth/oidcConfig";

const wrapper = ({ children }: { children: ReactNode }) =>
  createElement(AuthProvider, null, children);

describe("useAuth", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("is not authenticated when no user is stored", async () => {
    vi.mocked(userManager!.getUser).mockResolvedValue(null);

    const { result } = renderHook(() => useAuth(), { wrapper });

    // isLoading starts true, wait for resolution
    await vi.waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.user).toBeNull();
  });

  it("is authenticated when a valid user is stored", async () => {
    const fakeUser = { access_token: "tok", expired: false } as never;
    vi.mocked(userManager!.getUser).mockResolvedValue(fakeUser);

    const { result } = renderHook(() => useAuth(), { wrapper });

    await vi.waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.isAuthenticated).toBe(true);
  });

  it("is not authenticated when stored user is expired", async () => {
    const expiredUser = { access_token: "tok", expired: true } as never;
    vi.mocked(userManager!.getUser).mockResolvedValue(expiredUser);

    const { result } = renderHook(() => useAuth(), { wrapper });

    await vi.waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.isAuthenticated).toBe(false);
  });

  it("login calls signinRedirect", async () => {
    vi.mocked(userManager!.getUser).mockResolvedValue(null);
    vi.mocked(userManager!.signinRedirect).mockResolvedValue();

    const { result } = renderHook(() => useAuth(), { wrapper });
    await vi.waitFor(() => expect(result.current.isLoading).toBe(false));

    await result.current.login();

    expect(userManager!.signinRedirect).toHaveBeenCalledOnce();
  });

  it("logout calls signoutRedirect", async () => {
    vi.mocked(userManager!.getUser).mockResolvedValue(null);
    vi.mocked(userManager!.signoutRedirect).mockResolvedValue();

    const { result } = renderHook(() => useAuth(), { wrapper });
    await vi.waitFor(() => expect(result.current.isLoading).toBe(false));

    await result.current.logout();

    expect(userManager!.signoutRedirect).toHaveBeenCalledOnce();
  });
});
