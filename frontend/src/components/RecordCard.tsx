import { useNavigate } from "react-router-dom";
import { StarRating } from "./StarRating";
import type { components } from "../api/schema.d.ts";

type CoffeeRecord = components["schemas"]["Record"];

interface RecordCardProps {
  record: CoffeeRecord;
}

const ROAST_LABEL: { [key: string]: string } = {
  light: "浅煎り",
  medium: "中煎り",
  dark: "深煎り",
};

export function RecordCard({ record }: RecordCardProps) {
  const navigate = useNavigate();

  return (
    <div
      onClick={() => navigate(`/records/${record.id}`)}
      style={{
        background: "var(--color-cream-200)",
        borderRadius: "12px",
        padding: "16px",
        cursor: "pointer",
        boxShadow: "0 1px 4px rgba(59,31,10,0.08)",
        position: "relative",
        transition: "box-shadow 0.15s",
      }}
    >
      {!record.is_note_filled && (
        <span
          style={{
            position: "absolute",
            top: "10px",
            right: "10px",
            background: "#c47a1b",
            color: "#fff",
            fontSize: "0.65rem",
            padding: "2px 7px",
            borderRadius: "20px",
            fontWeight: 600,
          }}
        >
          未追記
        </span>
      )}
      <div
        style={{
          fontWeight: 700,
          fontSize: "1rem",
          color: "var(--color-brown-900)",
          marginBottom: "4px",
        }}
      >
        {record.name}
      </div>
      {record.origin && (
        <div
          style={{
            fontSize: "0.82rem",
            color: "var(--color-brown-400)",
            marginBottom: "6px",
          }}
        >
          {record.origin}
          {record.roast_level && ` · ${ROAST_LABEL[record.roast_level] ?? record.roast_level}`}
        </div>
      )}
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
        }}
      >
        <StarRating value={record.rating} readonly size="sm" />
        <span style={{ fontSize: "0.75rem", color: "var(--color-brown-400)" }}>
          {record.purchased_at
            ? record.purchased_at
            : new Date(record.created_at).toLocaleDateString("ja-JP")}
        </span>
      </div>
    </div>
  );
}
