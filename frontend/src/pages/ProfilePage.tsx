import { useState } from "react";
import { useStatsSummary, useFlavorWords } from "../hooks/useStats";
import { RoastBar } from "../components/RoastBar";
import { StarRating } from "../components/StarRating";
import { useApiKey } from "../hooks/useApiKey";
import { generateProfileComment } from "../api/ai";

export function ProfilePage() {
  const { summary, loading } = useStatsSummary();
  const { words } = useFlavorWords();
  const { provider, apiKey, model } = useApiKey();
  const [aiComment, setAiComment] = useState<string | null>(null);
  const [commentLoading, setCommentLoading] = useState(false);

  const topOrigins = summary?.by_origin.slice(0, 3).map((o) => o.label) ?? [];
  const topRoast = summary?.by_roast_level.sort((a, b) => b.count - a.count)[0]?.label ?? "";
  const topFlavorWords = words.slice(0, 10).map((w) => w.word);

  const handleGenerateComment = async () => {
    setCommentLoading(true);
    try {
      const comment = await generateProfileComment(
        provider,
        apiKey,
        model,
        topOrigins,
        topRoast,
        topFlavorWords,
      );
      setAiComment(comment);
    } catch {
      setAiComment("コメントの生成に失敗しました");
    } finally {
      setCommentLoading(false);
    }
  };

  const hasData = (summary?.total_records ?? 0) > 0;

  return (
    <div style={{ padding: "16px", maxWidth: "600px", margin: "0 auto" }}>
      <h2
        style={{
          fontSize: "1.2rem",
          fontWeight: 700,
          color: "var(--color-brown-900)",
          marginBottom: "6px",
        }}
      >
        好みカード
      </h2>
      <p style={{ fontSize: "0.82rem", color: "var(--color-brown-400)", marginBottom: "20px" }}>
        店頭で店員に見せてお気に入りを伝えましょう
      </p>

      {loading && (
        <p style={{ color: "var(--color-brown-400)", textAlign: "center" }}>読み込み中...</p>
      )}

      {!loading && !hasData && (
        <div
          style={{
            background: "var(--color-cream-200)",
            borderRadius: "12px",
            padding: "32px",
            textAlign: "center",
            color: "var(--color-brown-400)",
          }}
        >
          <p style={{ fontSize: "1.5rem", marginBottom: "12px" }}>📊</p>
          <p>記録を増やすと傾向が表示されます</p>
          <p style={{ fontSize: "0.82rem" }}>現在 {summary?.total_records ?? 0} 件</p>
        </div>
      )}

      {!loading && hasData && (
        <>
          {/* AI Comment */}
          {apiKey && (
            <div style={cardStyle}>
              {aiComment ? (
                <p
                  style={{
                    fontSize: "0.95rem",
                    lineHeight: 1.6,
                    color: "var(--color-brown-900)",
                    fontStyle: "italic",
                  }}
                >
                  ✨ {aiComment}
                </p>
              ) : (
                <button
                  onClick={() => void handleGenerateComment()}
                  disabled={commentLoading}
                  style={genBtnStyle}
                >
                  {commentLoading ? "生成中..." : "✨ AIに好みを分析してもらう"}
                </button>
              )}
            </div>
          )}

          {/* Top Origins */}
          <div style={cardStyle}>
            <h3 style={cardTitleStyle}>好きな産地 TOP3</h3>
            {summary!.by_origin.length === 0 ? (
              <p style={emptyStyle}>産地データなし</p>
            ) : (
              <div>
                {summary!.by_origin.slice(0, 3).map((o, i) => (
                  <div
                    key={o.label}
                    style={{
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "space-between",
                      marginBottom: "8px",
                    }}
                  >
                    <span
                      style={{
                        fontSize: "0.9rem",
                        color: "var(--color-brown-900)",
                        fontWeight: i === 0 ? 700 : 400,
                      }}
                    >
                      {i + 1}. {o.label}
                    </span>
                    <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                      <StarRating value={Math.round(o.avg_rating)} readonly size="sm" />
                      <span style={{ fontSize: "0.78rem", color: "var(--color-brown-400)" }}>
                        ({o.count}件)
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Roast Level */}
          <div style={cardStyle}>
            <h3 style={cardTitleStyle}>焙煎度の傾向</h3>
            <RoastBar byRoastLevel={summary!.by_roast_level} />
          </div>

          {/* Brew Methods */}
          <div style={cardStyle}>
            <h3 style={cardTitleStyle}>よく使う抽出方法</h3>
            {summary!.by_brew_method.length === 0 ? (
              <p style={emptyStyle}>データなし</p>
            ) : (
              <div style={{ display: "flex", flexWrap: "wrap", gap: "8px" }}>
                {summary!.by_brew_method.map((b) => {
                  const BREW_LABEL: Record<string, string> = {
                    drip: "ドリップ",
                    espresso: "エスプレッソ",
                    french_press: "フレンチプレス",
                    other: "その他",
                  };
                  const total = summary!.by_brew_method.reduce((s, x) => s + x.count, 0);
                  const pct = Math.round((b.count / total) * 100);
                  return (
                    <span
                      key={b.label}
                      style={{
                        background: "var(--color-brown-700)",
                        color: "#fff",
                        borderRadius: "20px",
                        padding: "5px 12px",
                        fontSize: "0.82rem",
                      }}
                    >
                      {BREW_LABEL[b.label] ?? b.label} {pct}%
                    </span>
                  );
                })}
              </div>
            )}
          </div>

          {/* Flavor Words */}
          <div style={cardStyle}>
            <h3 style={cardTitleStyle}>好きなフレーバーワード</h3>
            {words.length === 0 ? (
              <p style={emptyStyle}>テイスティングノートを追加すると表示されます</p>
            ) : (
              <div style={{ display: "flex", flexWrap: "wrap", gap: "8px" }}>
                {words.slice(0, 12).map((w) => (
                  <span
                    key={w.word}
                    style={{
                      background: "var(--color-cream-200)",
                      color: "var(--color-brown-900)",
                      border: "1px solid var(--color-brown-400)",
                      borderRadius: "20px",
                      padding: "4px 12px",
                      fontSize: "0.82rem",
                    }}
                  >
                    {w.word}
                  </span>
                ))}
              </div>
            )}
          </div>

          <p
            style={{
              textAlign: "center",
              fontSize: "0.78rem",
              color: "var(--color-brown-400)",
              marginTop: "8px",
            }}
          >
            {summary!.total_records}件の記録から生成
          </p>
        </>
      )}
    </div>
  );
}

const cardStyle: React.CSSProperties = {
  background: "var(--color-cream-200)",
  borderRadius: "12px",
  padding: "16px",
  marginBottom: "14px",
};

const cardTitleStyle: React.CSSProperties = {
  fontSize: "0.85rem",
  fontWeight: 700,
  color: "var(--color-brown-700)",
  margin: "0 0 12px 0",
};

const emptyStyle: React.CSSProperties = {
  fontSize: "0.82rem",
  color: "var(--color-brown-400)",
};

const genBtnStyle: React.CSSProperties = {
  display: "block",
  width: "100%",
  padding: "10px",
  borderRadius: "8px",
  border: "none",
  background: "var(--color-brown-700)",
  color: "#fff",
  fontSize: "0.9rem",
  fontWeight: 600,
  cursor: "pointer",
};
