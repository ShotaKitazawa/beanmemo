import { useState, useEffect } from "react";
import { apiClient } from "../api/client";
import type { components } from "../api/schema.d.ts";

type StatsSummary = components["schemas"]["StatsSummary"];
type FlavorWord = components["schemas"]["FlavorWord"];

export function useStatsSummary() {
  const [summary, setSummary] = useState<StatsSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    apiClient.GET("/stats/summary").then(({ data, error: err }) => {
      if (err) {
        setError("message" in err ? err.message : "取得に失敗しました");
      } else {
        setSummary(data ?? null);
      }
      setLoading(false);
    });
  }, []);

  return { summary, loading, error };
}

export function useFlavorWords() {
  const [words, setWords] = useState<FlavorWord[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiClient.GET("/stats/flavor-words").then(({ data }) => {
      setWords(data ?? []);
      setLoading(false);
    });
  }, []);

  return { words, loading };
}
