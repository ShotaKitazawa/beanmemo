import { useState, useEffect, useCallback } from "react";
import { apiClient } from "../api/client";
import type { components } from "../api/schema.d.ts";

type Record = components["schemas"]["Record"];
type RecordDetail = components["schemas"]["RecordDetail"];
type CreateRecordRequest = components["schemas"]["CreateRecordRequest"];
type UpdateRecordRequest = components["schemas"]["UpdateRecordRequest"];

interface ListParams {
  origin?: string;
  roast_level?: "light" | "medium" | "dark";
  rating_min?: number;
  brew_method?: "drip" | "espresso" | "french_press" | "other";
}

export function useRecords(params?: ListParams) {
  const [records, setRecords] = useState<Record[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRecords = useCallback(async () => {
    setLoading(true);
    setError(null);
    const { data, error: err } = await apiClient.GET("/records", {
      params: { query: params },
    });
    if (err) {
      setError("message" in err ? err.message : "取得に失敗しました");
    } else {
      setRecords(data ?? []);
    }
    setLoading(false);
  }, [params?.origin, params?.roast_level, params?.rating_min, params?.brew_method]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    void fetchRecords();
  }, [fetchRecords]);

  return { records, loading, error, refetch: fetchRecords };
}

export function useRecord(id: number) {
  const [record, setRecord] = useState<RecordDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    apiClient.GET("/records/{id}", { params: { path: { id } } }).then(({ data, error: err }) => {
      if (cancelled) return;
      if (err) {
        setError("message" in err ? err.message : "取得に失敗しました");
      } else {
        setRecord(data ?? null);
      }
      setLoading(false);
    });
    return () => {
      cancelled = true;
    };
  }, [id]);

  return { record, loading, error };
}

export async function createRecord(req: CreateRecordRequest): Promise<Record> {
  const { data, error } = await apiClient.POST("/records", { body: req });
  if (error) throw new Error("message" in error ? error.message : "作成に失敗しました");
  if (!data) throw new Error("作成に失敗しました");
  return data;
}

export async function updateRecord(id: number, req: UpdateRecordRequest): Promise<Record> {
  const { data, error } = await apiClient.PUT("/records/{id}", {
    params: { path: { id } },
    body: req,
  });
  if (error) throw new Error("message" in error ? error.message : "更新に失敗しました");
  if (!data) throw new Error("更新に失敗しました");
  return data;
}

export async function deleteRecord(id: number): Promise<void> {
  const { error } = await apiClient.DELETE("/records/{id}", {
    params: { path: { id } },
  });
  if (error) throw new Error("message" in error ? error.message : "削除に失敗しました");
}
