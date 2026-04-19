import { UserManager, WebStorageStateStore } from "oidc-client-ts";

export interface OIDCSetup {
  configured: boolean;
  userManager: UserManager | null;
}

let _setup: OIDCSetup | null = null;

/**
 * Fetches OIDC configuration from GET /api/oidc-config and creates a
 * UserManager if OIDC is enabled. The result is cached so multiple callers
 * share the same UserManager instance.
 */
export async function loadOIDCSetup(): Promise<OIDCSetup> {
  if (_setup !== null) return _setup;

  try {
    const res = await fetch("/api/oidc-config");
    const data = (await res.json()) as {
      enabled: boolean;
      issuer?: string | null;
      client_id?: string | null;
      audience?: string | null;
    };

    if (!data.enabled || !data.issuer || !data.client_id) {
      _setup = { configured: false, userManager: null };
      return _setup;
    }

    const userManager = new UserManager({
      authority: data.issuer,
      client_id: data.client_id,
      redirect_uri: `${window.location.origin}/callback`,
      scope: "openid profile email",
      response_type: "code",
      userStore: new WebStorageStateStore({ store: window.localStorage }),
      automaticSilentRenew: true,
    });

    _setup = { configured: true, userManager };
    return _setup;
  } catch {
    _setup = { configured: false, userManager: null };
    return _setup;
  }
}

/** Reset the cached setup (for testing). */
export function resetOIDCSetup() {
  _setup = null;
}
