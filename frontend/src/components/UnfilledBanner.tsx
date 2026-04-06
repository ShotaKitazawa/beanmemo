import { useNavigate } from "react-router-dom";
import type { components } from "../api/schema.d.ts";

type Record = components["schemas"]["Record"];

interface UnfilledBannerProps {
  records: Record[];
}

export function UnfilledBanner({ records }: UnfilledBannerProps) {
  const navigate = useNavigate();
  const unfilled = records.filter((r) => !r.is_note_filled);
  if (unfilled.length === 0) return null;

  const first = unfilled[0];

  return (
    <div
      onClick={() => navigate(`/records/${first.id}`)}
      style={{
        background: "#fff3cd",
        border: "1px solid #c47a1b",
        borderRadius: "10px",
        padding: "12px 16px",
        marginBottom: "16px",
        cursor: "pointer",
        display: "flex",
        alignItems: "center",
        gap: "10px",
      }}
    >
      <span style={{ fontSize: "1.3rem" }}>📝</span>
      <div>
        <div
          style={{
            fontWeight: 600,
            color: "#7a4e00",
            fontSize: "0.9rem",
          }}
        >
          テイスティングノートが未記入の豆が{unfilled.length}件あります
        </div>
        <div style={{ fontSize: "0.78rem", color: "#a06810" }}>
          「{first.name}」など — タップして追記する
        </div>
      </div>
    </div>
  );
}
