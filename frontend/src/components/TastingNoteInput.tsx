import { useState } from "react";
import { Link } from "react-router-dom";
import {
  generateTastingQuestions,
  generateTastingNote,
  type TastingQuestion,
  type TastingAnswers,
} from "../api/ai";
import type { Provider } from "../hooks/useApiKey";

interface TastingNoteInputProps {
  value: string;
  onChange: (value: string) => void;
  provider: Provider;
  apiKey: string;
  model: string;
}

type Step = "input" | "questions" | "done";

const FREE = "__free__";

export function TastingNoteInput({
  value,
  onChange,
  provider,
  apiKey,
  model,
}: TastingNoteInputProps) {
  const [step, setStep] = useState<Step>("input");
  // capturedText: "AIに深掘りしてもらう" を押した瞬間の入力テキスト（質問生成・ノート生成に使用）
  const [capturedText, setCapturedText] = useState("");
  const [questions, setQuestions] = useState<TastingQuestion[]>([]);
  const [answers, setAnswers] = useState<TastingAnswers>({});
  const [freeInputs, setFreeInputs] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [prevValue, setPrevValue] = useState<string>("");

  const resetAI = (restoreValue?: string) => {
    if (restoreValue !== undefined) onChange(restoreValue);
    setQuestions([]);
    setAnswers({});
    setFreeInputs({});
    setError(null);
    setStep("input");
  };

  const toggleAnswer = (question: string, val: string) => {
    setAnswers((prev) => ({ ...prev, [question]: prev[question] === val ? "" : val }));
  };

  const handleDeepen = async (baseText?: string) => {
    const text = baseText ?? value;
    if (!text.trim()) return;
    setCapturedText(text);
    setLoading(true);
    setError(null);
    try {
      const qs = await generateTastingQuestions(provider, apiKey, model, text);
      setQuestions(qs);
      setAnswers({});
      setFreeInputs({});
      setStep("questions");
    } catch (e) {
      setError(e instanceof Error ? e.message : "エラーが発生しました");
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async () => {
    if (loading) return;
    setLoading(true);
    setError(null);
    try {
      const effectiveAnswers: TastingAnswers = {};
      for (const q of questions) {
        const ans = answers[q.question];
        if (!ans) continue;
        if (ans === FREE) {
          const text = freeInputs[q.question]?.trim();
          if (text) effectiveAnswers[q.question] = text;
        } else {
          effectiveAnswers[q.question] = ans;
        }
      }
      setPrevValue(value);
      const note = await generateTastingNote(
        provider,
        apiKey,
        model,
        capturedText,
        effectiveAnswers,
      );
      onChange(note);
      setStep("done");
    } catch (e) {
      setError(e instanceof Error ? e.message : "エラーが発生しました");
    } finally {
      setLoading(false);
    }
  };

  if (!apiKey) {
    return (
      <div>
        <textarea
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder="テイスティングノートを入力..."
          rows={4}
          style={textareaStyle}
        />
        <p style={{ fontSize: "0.78rem", color: "var(--color-brown-400)", marginTop: "6px" }}>
          <Link to="/settings" style={{ color: "var(--color-brown-700)" }}>
            APIキーを設定
          </Link>
          するとAIがノートを自動生成します
        </p>
      </div>
    );
  }

  if (step === "input") {
    return (
      <div>
        <textarea
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder="飲んだ感想を自由に入力してください（例：なんか酸っぱくてフルーティーだった）"
          rows={3}
          style={textareaStyle}
        />
        {error && <p style={{ color: "#c0392b", fontSize: "0.82rem" }}>{error}</p>}
        <button
          onClick={() => void handleDeepen()}
          disabled={loading || !value.trim()}
          style={aiBtnStyle}
        >
          {loading ? "生成中..." : "✨ AIに深掘りしてもらう"}
        </button>
      </div>
    );
  }

  if (step === "questions") {
    return (
      <div>
        <p
          style={{
            fontSize: "0.88rem",
            color: "var(--color-brown-700)",
            marginBottom: "12px",
            fontWeight: 600,
          }}
        >
          もう少し教えてください（未選択は無回答として扱います）:
        </p>
        {questions.map((q, i) => {
          const sel = answers[q.question];
          const chip = (val: string, muted = false): React.CSSProperties => ({
            padding: "6px 14px",
            borderRadius: "20px",
            border: `1px solid ${sel === val ? "var(--color-brown-700)" : "var(--color-brown-400)"}`,
            background: sel === val ? "var(--color-brown-700)" : "transparent",
            color:
              sel === val ? "#fff" : muted ? "var(--color-brown-400)" : "var(--color-brown-700)",
            fontSize: "0.82rem",
            cursor: "pointer",
          });
          return (
            <div key={i} style={{ marginBottom: "14px" }}>
              <p
                style={{
                  fontSize: "0.88rem",
                  fontWeight: 600,
                  marginBottom: "6px",
                  color: "var(--color-brown-900)",
                }}
              >
                {q.question}
              </p>
              <div style={{ display: "flex", flexWrap: "wrap", gap: "8px" }}>
                {q.choices.map((choice) => (
                  <button
                    key={choice}
                    onClick={() => toggleAnswer(q.question, choice)}
                    style={chip(choice)}
                  >
                    {choice}
                  </button>
                ))}
                <button onClick={() => toggleAnswer(q.question, FREE)} style={chip(FREE, true)}>
                  自由記述
                </button>
              </div>
              {sel === FREE && (
                <input
                  type="text"
                  value={freeInputs[q.question] ?? ""}
                  onChange={(e) =>
                    setFreeInputs((prev) => ({ ...prev, [q.question]: e.target.value }))
                  }
                  placeholder="自由に記述してください"
                  style={{
                    display: "block",
                    width: "100%",
                    marginTop: "8px",
                    padding: "8px 10px",
                    borderRadius: "8px",
                    border: "1px solid var(--color-brown-400)",
                    background: "var(--color-cream-100)",
                    color: "var(--color-brown-900)",
                    fontSize: "0.88rem",
                    boxSizing: "border-box",
                  }}
                />
              )}
            </div>
          );
        })}
        {error && <p style={{ color: "#c0392b", fontSize: "0.82rem" }}>{error}</p>}
        <button
          onClick={() => void handleGenerate()}
          disabled={loading}
          style={
            loading
              ? { ...saveBtnStyle, background: "var(--color-brown-400)", cursor: "not-allowed" }
              : saveBtnStyle
          }
        >
          {loading ? "生成中..." : "📝 ノートを生成する"}
        </button>
      </div>
    );
  }

  // step === "done"
  return (
    <div>
      <textarea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        rows={5}
        style={textareaStyle}
      />
      <p style={{ fontSize: "0.78rem", color: "var(--color-brown-400)", marginTop: "4px" }}>
        ✨ AIが生成しました。自由に編集できます。
      </p>
      <button onClick={() => void handleDeepen(value)} style={{ ...aiBtnStyle, marginTop: "8px" }}>
        ✨ もう一度AI深掘りする
      </button>
      <button
        onClick={() => resetAI(prevValue)}
        style={{
          ...aiBtnStyle,
          marginTop: "6px",
          color: "var(--color-brown-400)",
          borderColor: "var(--color-brown-400)",
        }}
      >
        キャンセル
      </button>
    </div>
  );
}

const textareaStyle: React.CSSProperties = {
  width: "100%",
  padding: "10px 12px",
  borderRadius: "8px",
  border: "1px solid var(--color-brown-400)",
  background: "var(--color-cream-100)",
  color: "var(--color-brown-900)",
  fontSize: "0.95rem",
  resize: "vertical",
  fontFamily: "inherit",
  boxSizing: "border-box",
};

const aiBtnStyle: React.CSSProperties = {
  display: "block",
  width: "100%",
  marginTop: "10px",
  padding: "10px",
  borderRadius: "8px",
  border: "1px solid #c47a1b",
  background: "transparent",
  color: "#c47a1b",
  fontSize: "0.9rem",
  fontWeight: 600,
  cursor: "pointer",
};

const saveBtnStyle: React.CSSProperties = {
  display: "block",
  width: "100%",
  marginTop: "10px",
  padding: "10px",
  borderRadius: "8px",
  border: "none",
  background: "var(--color-brown-700)",
  color: "#fff",
  fontSize: "0.9rem",
  fontWeight: 600,
  cursor: "pointer",
};
