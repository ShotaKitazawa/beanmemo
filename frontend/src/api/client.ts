import createClient from "openapi-fetch";
import type { paths } from "./schema.d.ts";

const baseUrl = import.meta.env.VITE_API_BASE_URL ?? "";

export const apiClient = createClient<paths>({
  baseUrl: `${baseUrl}/api`,
});

export type { paths };
