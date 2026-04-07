import createClient from "openapi-fetch";
import type { paths } from "./schema.d.ts";

const baseUrl = import.meta.env.VITE_API_BASE_URL ?? "";

type TokenProvider = () => Promise<string | null>;
let tokenProvider: TokenProvider | null = null;

/** Called by AuthProvider to inject token retrieval logic. */
export function setTokenProvider(fn: TokenProvider) {
  tokenProvider = fn;
}

export const apiClient = createClient<paths>({
  baseUrl: `${baseUrl}/api`,
  fetch: async (req) => {
    if (tokenProvider) {
      const token = await tokenProvider();
      if (token) {
        req.headers.set("Authorization", `Bearer ${token}`);
      }
    }
    return fetch(req);
  },
});

export type { paths };
