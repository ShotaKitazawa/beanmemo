import { useState } from "react";

export type Provider = "claude" | "openai" | "google";

const PROVIDER_KEY = "beanmemo_ai_provider";
const keyFor = (p: Provider) => `beanmemo_ai_key_${p}`;
const modelKeyFor = (p: Provider) => `beanmemo_ai_model_${p}`;

export const DEFAULT_MODELS: Record<Provider, string> = {
  claude: "claude-haiku-4-5",
  openai: "gpt-5.4",
  google: "gemini-3.1-flash-lite-preview",
};

function loadProvider(): Provider {
  return (localStorage.getItem(PROVIDER_KEY) as Provider) ?? "claude";
}

function loadModel(p: Provider): string {
  return localStorage.getItem(modelKeyFor(p)) ?? DEFAULT_MODELS[p];
}

export function useApiKey() {
  const [provider, setProviderState] = useState<Provider>(loadProvider);
  const [apiKey, setApiKeyState] = useState<string>(
    () => localStorage.getItem(keyFor(loadProvider())) ?? "",
  );
  const [model, setModelState] = useState<string>(() => loadModel(loadProvider()));

  const setProvider = (p: Provider) => {
    localStorage.setItem(PROVIDER_KEY, p);
    setProviderState(p);
    setApiKeyState(localStorage.getItem(keyFor(p)) ?? "");
    setModelState(loadModel(p));
  };

  const setApiKey = (key: string) => {
    if (key) {
      localStorage.setItem(keyFor(provider), key);
    } else {
      localStorage.removeItem(keyFor(provider));
    }
    setApiKeyState(key);
  };

  const setModel = (m: string) => {
    localStorage.setItem(modelKeyFor(provider), m);
    setModelState(m);
  };

  return { provider, setProvider, apiKey, setApiKey, model, setModel, hasApiKey: !!apiKey };
}
