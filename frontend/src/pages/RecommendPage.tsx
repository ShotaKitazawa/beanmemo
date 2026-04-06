import { useState } from "react";
import { useRecommend } from "../hooks/useRecommend";
import { StarRating } from "../components/StarRating";
import { useApiKey } from "../hooks/useApiKey";
import { generateRecommendComment } from "../api/ai";

export function RecommendPage() {
  const [origin, setOrigin] = useState("");
  const [name, setName] = useState("");
  const { result, loading, error, search } = useRecommend();
  const { provider, apiKey, model } = useApiKey();
  const [aiComment, setAiComment] = useState<string | null>(null);
  const [commentLoading, setCommentLoading] = useState(false);

  const handleSearch = async () => {
    setAiComment(null);
    await search({ origin: origin || undefined, name: name || undefined });
  };

  const handleAiComment = async (score: number) => {
    setCommentLoading(true);
    try {
      const comment = await generateRecommendComment(provider, apiKey, model, score, origin, name);
      setAiComment(comment);
    } catch {
      setAiComment("コメントの生成に失敗しました");
    } finally {
      setCommentLoading(false);
    }
  };

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
        レコメンド
      </h2>

      <div style={searchCardStyle}>
        <label style={labelStyle}>産地</label>
        <input
          type="text"
          value={origin}
          onChange={(e) => setOrigin(e.target.value)}
          placeholder="例: エチオピア"
          style={inputStyle}
        />

        <label style={labelStyle}>豆の名前</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="例: イルガチェフェ"
          style={inputStyle}
        />

        <button
          onClick={() => void handleSearch()}
          disabled={loading || (!origin && !name)}
          style={searchBtnStyle}
        >
          {loading ? "検索中..." : "好みスコアを確認する"}
        </button>
      </div>

      {error && <p style={{ color: "#c0392b" }}>{error}</p>}

      {result && (
        <div style={resultCardStyle}>
          {result.locked ? (
            <div style={{ textAlign: "center", padding: "16px" }}>
              <p style={{ fontSize: "2rem", marginBottom: "12px" }}>🔒</p>
              <p
                style={{
                  fontWeight: 700,
                  color: "var(--color-brown-900)",
                  marginBottom: "8px",
                }}
              >
                まだレコメンドを使えません
              </p>
              <p style={{ fontSize: "0.88rem", color: "var(--color-brown-400)" }}>
                あと{result.records_needed}件記録するとレコメンドが使えます
              </p>
              <p style={{ fontSize: "0.78rem", color: "var(--color-brown-400)", marginTop: "8px" }}>
                現在 {result.total_records} / 5 件
              </p>
            </div>
          ) : (
            <div>
              <p
                style={{
                  fontSize: "0.82rem",
                  color: "var(--color-brown-400)",
                  marginBottom: "8px",
                }}
              >
                好みスコア ({result.total_records}件の記録から算出)
              </p>
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "16px",
                  marginBottom: "16px",
                }}
              >
                <span
                  style={{
                    fontSize: "2.5rem",
                    fontWeight: 800,
                    color: "var(--color-brown-700)",
                  }}
                >
                  {result.score != null ? result.score.toFixed(1) : "-"}
                </span>
                <div>
                  <StarRating
                    value={result.score != null ? Math.round(result.score) : 0}
                    readonly
                    size="lg"
                  />
                  <p
                    style={{
                      fontSize: "0.78rem",
                      color: "var(--color-brown-400)",
                      margin: "4px 0 0",
                    }}
                  >
                    / 5.0
                  </p>
                </div>
              </div>

              <div
                style={{
                  display: "flex",
                  flexDirection: "column",
                  gap: "6px",
                  marginBottom: "16px",
                }}
              >
                {result.origin_avg != null && (
                  <div style={breakdownRowStyle}>
                    <span>産地平均</span>
                    <span style={{ fontWeight: 600 }}>★{result.origin_avg.toFixed(1)}</span>
                  </div>
                )}
                {result.name_match_avg != null && (
                  <div style={breakdownRowStyle}>
                    <span>この豆の平均</span>
                    <span style={{ fontWeight: 600 }}>★{result.name_match_avg.toFixed(1)}</span>
                  </div>
                )}
              </div>

              {result.score == null && (
                <p style={{ fontSize: "0.85rem", color: "var(--color-brown-400)" }}>
                  産地や名前に一致する記録がありませんでした
                </p>
              )}

              {apiKey && result.score != null && (
                <>
                  {aiComment ? (
                    <p
                      style={{
                        fontSize: "0.9rem",
                        lineHeight: 1.6,
                        color: "var(--color-brown-900)",
                        fontStyle: "italic",
                        padding: "12px",
                        background: "var(--color-cream-100)",
                        borderRadius: "8px",
                      }}
                    >
                      ✨ {aiComment}
                    </p>
                  ) : (
                    <button
                      onClick={() => void handleAiComment(result.score!)}
                      disabled={commentLoading}
                      style={commentBtnStyle}
                    >
                      {commentLoading ? "生成中..." : "✨ AIコメントを生成"}
                    </button>
                  )}
                </>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

const searchCardStyle: React.CSSProperties = {
  background: "var(--color-cream-200)",
  borderRadius: "12px",
  padding: "16px",
  marginBottom: "16px",
};

const resultCardStyle: React.CSSProperties = {
  background: "var(--color-cream-200)",
  borderRadius: "12px",
  padding: "16px",
};

const labelStyle: React.CSSProperties = {
  display: "block",
  fontSize: "0.82rem",
  fontWeight: 600,
  color: "var(--color-brown-700)",
  marginBottom: "5px",
};

const inputStyle: React.CSSProperties = {
  display: "block",
  width: "100%",
  padding: "10px 12px",
  borderRadius: "8px",
  border: "1px solid var(--color-brown-400)",
  background: "var(--color-cream-100)",
  color: "var(--color-brown-900)",
  fontSize: "0.95rem",
  marginBottom: "14px",
  boxSizing: "border-box",
};

const searchBtnStyle: React.CSSProperties = {
  display: "block",
  width: "100%",
  padding: "12px",
  borderRadius: "10px",
  border: "none",
  background: "var(--color-brown-700)",
  color: "#fff",
  fontSize: "1rem",
  fontWeight: 700,
  cursor: "pointer",
};

const commentBtnStyle: React.CSSProperties = {
  display: "block",
  width: "100%",
  padding: "10px",
  borderRadius: "8px",
  border: "none",
  background: "var(--color-brown-700)",
  color: "#fff",
  fontSize: "0.88rem",
  fontWeight: 600,
  cursor: "pointer",
};

const breakdownRowStyle: React.CSSProperties = {
  display: "flex",
  justifyContent: "space-between",
  fontSize: "0.85rem",
  color: "var(--color-brown-700)",
  padding: "6px 10px",
  background: "var(--color-cream-100)",
  borderRadius: "6px",
};
