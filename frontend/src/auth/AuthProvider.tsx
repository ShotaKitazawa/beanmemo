import { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import type { UserManager } from "oidc-client-ts";
import { loadOIDCSetup } from "../oidc";
import { apiClient, setTokenProvider } from "../api/client";

interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  login: () => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [userManager, setUserManager] = useState<UserManager | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function init() {
      // 1. Fetch OIDC config from backend and create UserManager if configured.
      const setup = await loadOIDCSetup();
      if (cancelled) return;

      setUserManager(setup.userManager);

      // 2. Wire up token provider before calling /userinfo.
      if (setup.userManager) {
        const mgr = setup.userManager;
        setTokenProvider(async () => {
          const u = await mgr.getUser();
          return u?.access_token ?? null;
        });

        // Re-authenticate when the token is refreshed or revoked.
        const onLoaded = () => void checkAuth();
        const onUnloaded = () => setIsAuthenticated(false);
        mgr.events.addUserLoaded(onLoaded);
        mgr.events.addUserUnloaded(onUnloaded);
      } else {
        // OIDC disabled: no token — backend accepts all requests.
        setTokenProvider(async () => null);
      }

      // 3. Determine auth state via GET /userinfo (200 = authenticated, 401 = not).
      await checkAuth();
    }

    async function checkAuth() {
      const { data } = await apiClient.GET("/userinfo");
      if (!cancelled) {
        setIsAuthenticated(!!data);
        setIsLoading(false);
      }
    }

    void init();
    return () => {
      cancelled = true;
    };
  }, []);

  const login = async () => {
    if (userManager) await userManager.signinRedirect();
  };

  const logout = async () => {
    if (userManager) await userManager.signoutRedirect();
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuthContext(): AuthState {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuthContext must be used within AuthProvider");
  return ctx;
}
