interface RoastBarProps {
  byRoastLevel: Array<{ label: string; count: number; avg_rating: number }>;
}

const ROAST_ORDER = ["light", "medium", "dark"];
const ROAST_LABEL: Record<string, string> = {
  light: "浅煎り",
  medium: "中煎り",
  dark: "深煎り",
};

export function RoastBar({ byRoastLevel }: RoastBarProps) {
  const totalCount = byRoastLevel.reduce((s, r) => s + r.count, 0);
  if (totalCount === 0) {
    return <p style={{ color: "var(--color-brown-400)", fontSize: "0.88rem" }}>データなし</p>;
  }

  const sorted = ROAST_ORDER.map((key) => {
    const found = byRoastLevel.find((r) => r.label === key);
    return { label: key, count: found?.count ?? 0, avg: found?.avg_rating ?? 0 };
  });

  return (
    <div>
      <div
        style={{
          display: "flex",
          height: "16px",
          borderRadius: "8px",
          overflow: "hidden",
          marginBottom: "8px",
        }}
      >
        {sorted.map(({ label, count }) => {
          const pct = totalCount > 0 ? (count / totalCount) * 100 : 0;
          if (pct === 0) return null;
          const colors: Record<string, string> = {
            light: "#f5c98a",
            medium: "#c47a1b",
            dark: "#3b1f0a",
          };
          return (
            <div
              key={label}
              style={{
                width: `${pct}%`,
                background: colors[label] ?? "#888",
              }}
              title={`${ROAST_LABEL[label] ?? label}: ${count}件`}
            />
          );
        })}
      </div>
      <div style={{ display: "flex", gap: "16px", flexWrap: "wrap" }}>
        {sorted
          .filter((r) => r.count > 0)
          .map(({ label, count, avg }) => (
            <span key={label} style={{ fontSize: "0.8rem", color: "var(--color-brown-700)" }}>
              {ROAST_LABEL[label] ?? label}: {count}件 (★{avg.toFixed(1)})
            </span>
          ))}
      </div>
    </div>
  );
}
