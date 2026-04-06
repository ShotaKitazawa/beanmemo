import { useState } from "react";
import { useApiKey, type Provider } from "../hooks/useApiKey";

const PROVIDERS: Array<{ value: Provider; label: string; placeholder: string }> = [
  { value: "claude", label: "Claude", placeholder: "sk-ant-..." },
  { value: "openai", label: "OpenAI", placeholder: "sk-..." },
  { value: "google", label: "Google AI Studio", placeholder: "AIza..." },
];

const MODELS: Record<Provider, Array<{ value: string; label: string }>> = {
  claude: [
    { value: "claude-haiku-4-5", label: "Haiku 4.5" },
    { value: "claude-sonnet-4-5", label: "Sonnet 4.5" },
    { value: "claude-opus-4-6", label: "Opus 4.6" },
  ],
  openai: [
    { value: "gpt-5.4", label: "GPT-5.4" },
    { value: "gpt-5.4-mini", label: "GPT-5.4 mini" },
  ],
  google: [
    { value: "gemini-3.1-flash-lite-preview", label: "Gemini 3.1 Flash-Lite" },
    { value: "gemini-3-flash-preview", label: "Gemini 3 Flash" },
    { value: "gemini-3.1-pro-preview", label: "Gemini 3.1 Pro" },
  ],
};

export function SettingsPage() {
  const { provider, setProvider, apiKey, setApiKey, model, setModel, hasApiKey } = useApiKey();
  const [input, setInput] = useState(apiKey);
  const [saved, setSaved] = useState(false);

  const handleProviderChange = (p: Provider) => {
    setProvider(p);
    setInput("");
    setSaved(false);
  };

  const handleSave = () => {
    setApiKey(input.trim());
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  };

  const handleDelete = () => {
    setApiKey("");
    setInput("");
  };

  const currentProvider = PROVIDERS.find((p) => p.value === provider)!;

  return (
    <div style={{ padding: "16px", maxWidth: "600px", margin: "0 auto" }}>
      <h2
        style={{
          fontSize: "1.2rem",
          fontWeight: 700,
          color: "var(--color-brown-900)",
          marginBottom: "20px",
        }}
      >
        設定
      </h2>

      <div
        style={{
          background: "var(--color-cream-200)",
          borderRadius: "12px",
          padding: "16px",
        }}
      >
        <h3
          style={{
            fontSize: "0.9rem",
            fontWeight: 700,
            color: "var(--color-brown-700)",
            margin: "0 0 8px 0",
          }}
        >
          AI プロバイダ / API キー
        </h3>
        <p
          style={{
            fontSize: "0.8rem",
            color: "var(--color-brown-400)",
            marginBottom: "14px",
            lineHeight: 1.5,
          }}
        >
          AIによるテイスティングノート生成・レコメンドコメントを使用するにはAPIキーが必要です。
          キーはこのデバイスのみに保存され、サーバーには送信されません。
        </p>

        <label
          style={{
            display: "block",
            fontSize: "0.82rem",
            fontWeight: 600,
            color: "var(--color-brown-700)",
            marginBottom: "8px",
          }}
        >
          プロバイダ
        </label>
        <div style={{ display: "flex", gap: "8px", flexWrap: "wrap", marginBottom: "16px" }}>
          {PROVIDERS.map((p) => (
            <button
              key={p.value}
              onClick={() => handleProviderChange(p.value)}
              style={provider === p.value ? selectedChipStyle : chipStyle}
            >
              {p.label}
            </button>
          ))}
        </div>

        <label
          style={{
            display: "block",
            fontSize: "0.82rem",
            fontWeight: 600,
            color: "var(--color-brown-700)",
            marginBottom: "8px",
          }}
        >
          モデル
        </label>
        <div style={{ display: "flex", gap: "8px", flexWrap: "wrap", marginBottom: "16px" }}>
          {MODELS[provider].map((m) => (
            <button
              key={m.value}
              onClick={() => setModel(m.value)}
              style={model === m.value ? selectedChipStyle : chipStyle}
            >
              {m.label}
            </button>
          ))}
        </div>

        <label
          style={{
            display: "block",
            fontSize: "0.82rem",
            fontWeight: 600,
            color: "var(--color-brown-700)",
            marginBottom: "5px",
          }}
        >
          APIキー（{currentProvider.label}）
        </label>
        <input
          type="password"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder={currentProvider.placeholder}
          style={{
            display: "block",
            width: "100%",
            padding: "10px 12px",
            borderRadius: "8px",
            border: "1px solid var(--color-brown-400)",
            background: "var(--color-cream-100)",
            color: "var(--color-brown-900)",
            fontSize: "0.95rem",
            marginBottom: "12px",
            boxSizing: "border-box",
            fontFamily: "monospace",
          }}
        />

        <div style={{ display: "flex", gap: "10px" }}>
          <button
            onClick={handleSave}
            disabled={!input.trim()}
            style={{
              flex: 1,
              padding: "10px",
              borderRadius: "8px",
              border: "none",
              background: saved ? "#27ae60" : "var(--color-brown-700)",
              color: "#fff",
              fontSize: "0.9rem",
              fontWeight: 600,
              cursor: "pointer",
              transition: "background 0.2s",
            }}
          >
            {saved ? "✓ 保存しました" : "保存する"}
          </button>

          {hasApiKey && (
            <button
              onClick={handleDelete}
              style={{
                padding: "10px 16px",
                borderRadius: "8px",
                border: "1px solid #c0392b",
                background: "transparent",
                color: "#c0392b",
                fontSize: "0.9rem",
                fontWeight: 600,
                cursor: "pointer",
              }}
            >
              削除
            </button>
          )}
        </div>

        {hasApiKey && (
          <p style={{ marginTop: "10px", fontSize: "0.78rem", color: "#27ae60" }}>
            ✓ {currentProvider.label} の API キーが設定されています
          </p>
        )}
      </div>
    </div>
  );
}

const chipStyle: React.CSSProperties = {
  padding: "7px 16px",
  borderRadius: "20px",
  border: "1px solid var(--color-brown-400)",
  background: "transparent",
  color: "var(--color-brown-700)",
  fontSize: "0.85rem",
  cursor: "pointer",
};

const selectedChipStyle: React.CSSProperties = {
  ...chipStyle,
  background: "var(--color-brown-700)",
  color: "#fff",
  border: "1px solid var(--color-brown-700)",
};
