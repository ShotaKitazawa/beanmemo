import { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import type { User } from "oidc-client-ts";
import { oidcDisabled, userManager } from "./oidcConfig";
import { setTokenProvider } from "../api/client";

interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: User | null;
  login: () => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(!oidcDisabled);

  useEffect(() => {
    if (oidcDisabled) {
      // Dev mode: no auth needed; backend uses dummy user
      setTokenProvider(async () => null);
      return;
    }

    const mgr = userManager!;

    // Restore session from storage
    mgr.getUser().then((u) => {
      setUser(u);
      setIsLoading(false);
    });

    // Provide token to API client
    setTokenProvider(async () => {
      const u = await mgr.getUser();
      return u?.access_token ?? null;
    });

    const onUserLoaded = (u: User) => setUser(u);
    const onUserUnloaded = () => setUser(null);

    mgr.events.addUserLoaded(onUserLoaded);
    mgr.events.addUserUnloaded(onUserUnloaded);

    return () => {
      mgr.events.removeUserLoaded(onUserLoaded);
      mgr.events.removeUserUnloaded(onUserUnloaded);
    };
  }, []);

  const login = async () => {
    if (!oidcDisabled) await userManager!.signinRedirect();
  };

  const logout = async () => {
    if (!oidcDisabled) await userManager!.signoutRedirect();
  };

  const isAuthenticated = oidcDisabled || (user !== null && !user.expired);

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, user, login, logout }}>
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
