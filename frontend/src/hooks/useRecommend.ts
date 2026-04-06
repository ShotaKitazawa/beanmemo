import { useState } from "react";
import { apiClient } from "../api/client";
import type { components } from "../api/schema.d.ts";

type RecommendResult = components["schemas"]["RecommendResult"];

export function useRecommend() {
  const [result, setResult] = useState<RecommendResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const search = async (params: { origin?: string; name?: string }) => {
    setLoading(true);
    setError(null);
    const { data, error: err } = await apiClient.GET("/recommend", {
      params: { query: params },
    });
    if (err) {
      setError("message" in err ? err.message : "取得に失敗しました");
    } else {
      setResult(data ?? null);
    }
    setLoading(false);
  };

  return { result, loading, error, search };
}
