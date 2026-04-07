import { UserManager, WebStorageStateStore } from "oidc-client-ts";

const authority = import.meta.env.VITE_OIDC_AUTHORITY ?? "";
const clientId = import.meta.env.VITE_OIDC_CLIENT_ID ?? "";
const redirectUri = import.meta.env.VITE_OIDC_REDIRECT_URI ?? `${window.location.origin}/callback`;
const scope = import.meta.env.VITE_OIDC_SCOPE ?? "openid profile email";

/** true when OIDC is disabled (e.g. dev mode) */
export const oidcDisabled = authority === "";

export const userManager = oidcDisabled
  ? null
  : new UserManager({
      authority,
      client_id: clientId,
      redirect_uri: redirectUri,
      scope,
      response_type: "code",
      userStore: new WebStorageStateStore({ store: window.localStorage }),
      automaticSilentRenew: true,
    });
