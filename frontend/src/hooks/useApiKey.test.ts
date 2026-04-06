import { describe, it, expect, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useApiKey, DEFAULT_MODELS } from "./useApiKey";

beforeEach(() => {
  localStorage.clear();
});

describe("useApiKey", () => {
  it('初期状態でproviderが"claude"になる', () => {
    const { result } = renderHook(() => useApiKey());
    expect(result.current.provider).toBe("claude");
  });

  it("初期状態でapiKeyが空文字になる", () => {
    const { result } = renderHook(() => useApiKey());
    expect(result.current.apiKey).toBe("");
    expect(result.current.hasApiKey).toBe(false);
  });

  it("localStorageにapiKeyがある場合は初期値として読み込む", () => {
    localStorage.setItem("beanmemo_ai_key_claude", "sk-test-key");
    const { result } = renderHook(() => useApiKey());
    expect(result.current.apiKey).toBe("sk-test-key");
    expect(result.current.hasApiKey).toBe(true);
  });

  it("localStorageにproviderがある場合は初期値として読み込む", () => {
    localStorage.setItem("beanmemo_ai_provider", "openai");
    const { result } = renderHook(() => useApiKey());
    expect(result.current.provider).toBe("openai");
  });

  it("setApiKey でapiKeyが更新されlocalStorageに保存される", () => {
    const { result } = renderHook(() => useApiKey());
    act(() => {
      result.current.setApiKey("new-key");
    });
    expect(result.current.apiKey).toBe("new-key");
    expect(localStorage.getItem("beanmemo_ai_key_claude")).toBe("new-key");
  });

  it("setApiKey に空文字を渡すとlocalStorageから削除される", () => {
    localStorage.setItem("beanmemo_ai_key_claude", "existing-key");
    const { result } = renderHook(() => useApiKey());
    act(() => {
      result.current.setApiKey("");
    });
    expect(result.current.apiKey).toBe("");
    expect(localStorage.getItem("beanmemo_ai_key_claude")).toBeNull();
  });

  it("setProvider でproviderが切り替わる", () => {
    const { result } = renderHook(() => useApiKey());
    act(() => {
      result.current.setProvider("openai");
    });
    expect(result.current.provider).toBe("openai");
    expect(localStorage.getItem("beanmemo_ai_provider")).toBe("openai");
  });

  it("setProvider で切り替え先のapiKeyが読み込まれる", () => {
    localStorage.setItem("beanmemo_ai_key_openai", "openai-key");
    const { result } = renderHook(() => useApiKey());
    act(() => {
      result.current.setProvider("openai");
    });
    expect(result.current.apiKey).toBe("openai-key");
  });

  it("setModel でmodelが更新されlocalStorageに保存される", () => {
    const { result } = renderHook(() => useApiKey());
    act(() => {
      result.current.setModel("claude-opus-4-6");
    });
    expect(result.current.model).toBe("claude-opus-4-6");
    expect(localStorage.getItem("beanmemo_ai_model_claude")).toBe("claude-opus-4-6");
  });

  it("デフォルトモデルが各プロバイダに設定されている", () => {
    expect(DEFAULT_MODELS.claude).toBeTruthy();
    expect(DEFAULT_MODELS.openai).toBeTruthy();
    expect(DEFAULT_MODELS.google).toBeTruthy();
  });
});
