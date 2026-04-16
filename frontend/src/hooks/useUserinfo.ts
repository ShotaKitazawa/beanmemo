import { useEffect, useState } from "react";
import { apiClient } from "../api/client";
import type { components } from "../api/schema.d.ts";

type UserinfoResponse = components["schemas"]["UserinfoResponse"];

interface UseUserinfoResult {
  userinfo: UserinfoResponse | null;
  loading: boolean;
  error: string | null;
}

export function useUserinfo(): UseUserinfoResult {
  const [userinfo, setUserinfo] = useState<UserinfoResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    apiClient
      .GET("/userinfo")
      .then(({ data, error: apiError }) => {
        if (cancelled) return;
        if (apiError || !data) {
          setError(apiError ? String(apiError) : "failed to fetch userinfo");
          setUserinfo(null);
        } else {
          setUserinfo(data);
          setError(null);
        }
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        setError(String(err));
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, []);

  return { userinfo, loading, error };
}
