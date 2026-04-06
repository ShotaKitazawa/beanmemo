import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRecords } from "../hooks/useRecords";
import { RecordCard } from "../components/RecordCard";
import { UnfilledBanner } from "../components/UnfilledBanner";
import {
  ROAST_OPTIONS,
  BREW_OPTIONS,
  type RoastLevel,
  type BrewMethod,
} from "../lib/recordOptions";

const PAGE_SIZE = 10;

export function HomePage() {
  const navigate = useNavigate();
  const [origin, setOrigin] = useState("");
  const [roastLevel, setRoastLevel] = useState<RoastLevel | "">("");
  const [ratingMin, setRatingMin] = useState<number | "">("");
  const [brewMethod, setBrewMethod] = useState<BrewMethod | "">("");
  const [displayCount, setDisplayCount] = useState(PAGE_SIZE);

  const { records, loading, error } = useRecords({
    origin: origin || undefined,
    roast_level: roastLevel || undefined,
    rating_min: ratingMin !== "" ? ratingMin : undefined,
    brew_method: brewMethod || undefined,
  });

  const visible = records.slice(0, displayCount);
  const hasMore = records.length > displayCount;

  const handleFilterChange = () => {
    setDisplayCount(PAGE_SIZE);
  };

  return (
    <div style={{ padding: "16px", maxWidth: "600px", margin: "0 auto" }}>
      <header
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          marginBottom: "20px",
        }}
      >
        <h1
          style={{
            fontSize: "1.5rem",
            fontWeight: 800,
            color: "var(--color-brown-900)",
            margin: 0,
          }}
        >
          ☕ beanmemo
        </h1>
        <button onClick={() => navigate("/records/new")} style={fabStyle}>
          + 記録する
        </button>
      </header>

      {/* Filters */}
      <div
        style={{
          background: "var(--color-cream-200)",
          borderRadius: "10px",
          padding: "12px",
          marginBottom: "16px",
          display: "flex",
          flexWrap: "wrap",
          gap: "8px",
        }}
      >
        <input
          type="text"
          placeholder="産地で絞り込み"
          value={origin}
          onChange={(e) => {
            setOrigin(e.target.value);
            handleFilterChange();
          }}
          style={filterInputStyle}
        />
        <select
          value={roastLevel}
          onChange={(e) => {
            setRoastLevel(e.target.value as RoastLevel | "");
            handleFilterChange();
          }}
          style={filterInputStyle}
        >
          <option value="">焙煎度: すべて</option>
          {ROAST_OPTIONS.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
        <select
          value={ratingMin}
          onChange={(e) => {
            setRatingMin(e.target.value ? Number(e.target.value) : "");
            handleFilterChange();
          }}
          style={filterInputStyle}
        >
          <option value="">評価: すべて</option>
          {[1, 2, 3, 4, 5].map((r) => (
            <option key={r} value={r}>
              ★{r}以上
            </option>
          ))}
        </select>
        <select
          value={brewMethod}
          onChange={(e) => {
            setBrewMethod(e.target.value as BrewMethod | "");
            handleFilterChange();
          }}
          style={filterInputStyle}
        >
          <option value="">抽出方法: すべて</option>
          {BREW_OPTIONS.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
      </div>

      {!loading && <UnfilledBanner records={records} />}

      {loading && (
        <p style={{ color: "var(--color-brown-400)", textAlign: "center" }}>読み込み中...</p>
      )}
      {error && <p style={{ color: "#c0392b" }}>{error}</p>}

      {!loading && records.length === 0 && (
        <div style={{ textAlign: "center", padding: "48px 16px", color: "var(--color-brown-400)" }}>
          <p style={{ fontSize: "2rem", marginBottom: "12px" }}>🫘</p>
          <p style={{ fontWeight: 600, marginBottom: "8px" }}>
            {origin || roastLevel || ratingMin !== "" || brewMethod
              ? "該当する記録がありません"
              : "まだ記録がありません"}
          </p>
          {!origin && !roastLevel && ratingMin === "" && !brewMethod && (
            <p style={{ fontSize: "0.88rem" }}>最初の豆を記録してみましょう！</p>
          )}
        </div>
      )}

      <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
        {visible.map((r) => (
          <RecordCard key={r.id} record={r} />
        ))}
      </div>

      {hasMore && (
        <div style={{ textAlign: "center", marginTop: "16px" }}>
          <button onClick={() => setDisplayCount((n) => n + PAGE_SIZE)} style={moreBtnStyle}>
            もっと見る ({records.length - displayCount}件)
          </button>
        </div>
      )}
    </div>
  );
}

const fabStyle: React.CSSProperties = {
  background: "var(--color-brown-700)",
  color: "#fff",
  border: "none",
  borderRadius: "24px",
  padding: "10px 20px",
  fontSize: "0.9rem",
  fontWeight: 700,
  cursor: "pointer",
  boxShadow: "0 2px 8px rgba(111,78,55,0.4)",
};

const filterInputStyle: React.CSSProperties = {
  padding: "6px 10px",
  borderRadius: "6px",
  border: "1px solid var(--color-brown-400)",
  background: "var(--color-cream-100)",
  color: "var(--color-brown-900)",
  fontSize: "0.85rem",
  flex: "1 1 140px",
};

const moreBtnStyle: React.CSSProperties = {
  background: "transparent",
  border: "1px solid var(--color-brown-400)",
  color: "var(--color-brown-700)",
  padding: "8px 24px",
  borderRadius: "20px",
  cursor: "pointer",
  fontSize: "0.88rem",
};
