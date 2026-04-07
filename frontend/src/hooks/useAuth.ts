import { useAuthContext } from "../auth/AuthProvider";

/**
 * Returns auth state and actions.
 * - isAuthenticated: true when the user is logged in (or when OIDC is disabled)
 * - isLoading: true while the initial session is being restored
 * - login: redirects to the OIDC provider login page
 * - logout: redirects to the OIDC provider logout page
 */
export function useAuth() {
  return useAuthContext();
}
