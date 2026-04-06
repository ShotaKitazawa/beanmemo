import { useState } from "react";
import { StarRating } from "./StarRating";
import { TastingNoteInput } from "./TastingNoteInput";
import { useApiKey } from "../hooks/useApiKey";
import {
  ROAST_OPTIONS,
  BREW_OPTIONS,
  type RoastLevel,
  type BrewMethod,
  type RecordFormValues,
} from "../lib/recordOptions";

interface RecordFormProps {
  initialValues: RecordFormValues;
  onSave: (values: RecordFormValues) => Promise<void>;
  onCancel: () => void;
}

export function RecordForm({ initialValues, onSave, onCancel }: RecordFormProps) {
  const { provider, apiKey, model } = useApiKey();
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [name, setName] = useState(initialValues.name);
  const [rating, setRating] = useState(initialValues.rating);
  const [origin, setOrigin] = useState(initialValues.origin);
  const [roastLevel, setRoastLevel] = useState<RoastLevel | "">(initialValues.roastLevel);
  const [shop, setShop] = useState(initialValues.shop);
  const [price, setPrice] = useState(initialValues.price);
  const [purchasedAt, setPurchasedAt] = useState(initialValues.purchasedAt);
  const [tastingNote, setTastingNote] = useState(initialValues.tastingNote);
  const [brewMethod, setBrewMethod] = useState<BrewMethod | "">(initialValues.brewMethod);

  const handleSave = async () => {
    if (!name.trim()) return;
    setSaving(true);
    setError(null);
    try {
      await onSave({
        name: name.trim(),
        rating,
        origin,
        roastLevel,
        shop,
        price,
        purchasedAt,
        tastingNote,
        brewMethod,
      });
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div>
      <label style={labelStyle}>豆の名前</label>
      <input
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
        style={inputStyle}
      />

      <label style={labelStyle}>産地</label>
      <input
        type="text"
        value={origin}
        onChange={(e) => setOrigin(e.target.value)}
        style={inputStyle}
      />

      <label style={labelStyle}>焙煎度</label>
      <div style={{ display: "flex", gap: "8px", flexWrap: "wrap", marginBottom: "14px" }}>
        {ROAST_OPTIONS.map((o) => (
          <button
            key={o.value}
            onClick={() => setRoastLevel((p) => (p === o.value ? "" : o.value))}
            style={roastLevel === o.value ? selectedChipStyle : chipStyle}
          >
            {o.label}
          </button>
        ))}
      </div>

      <label style={labelStyle}>購入店</label>
      <input
        type="text"
        value={shop}
        onChange={(e) => setShop(e.target.value)}
        style={inputStyle}
      />

      <label style={labelStyle}>価格（円）</label>
      <input
        type="number"
        value={price}
        onChange={(e) => setPrice(e.target.value)}
        min={0}
        style={inputStyle}
      />

      <label style={labelStyle}>購入日</label>
      <input
        type="date"
        value={purchasedAt}
        onChange={(e) => setPurchasedAt(e.target.value)}
        style={inputStyle}
      />

      <label style={labelStyle}>総合評価</label>
      <div style={{ marginBottom: "14px" }}>
        <StarRating value={rating} onChange={setRating} size="lg" />
      </div>

      <label style={labelStyle}>抽出方法</label>
      <div style={{ display: "flex", gap: "8px", flexWrap: "wrap", marginBottom: "14px" }}>
        {BREW_OPTIONS.map((o) => (
          <button
            key={o.value}
            onClick={() => setBrewMethod((p) => (p === o.value ? "" : o.value))}
            style={brewMethod === o.value ? selectedChipStyle : chipStyle}
          >
            {o.label}
          </button>
        ))}
      </div>

      <label style={labelStyle}>テイスティングノート</label>
      <TastingNoteInput
        value={tastingNote}
        onChange={setTastingNote}
        provider={provider}
        apiKey={apiKey}
        model={model}
      />

      {error && <p style={{ color: "#c0392b", fontSize: "0.85rem", marginTop: "8px" }}>{error}</p>}

      <div style={{ display: "flex", gap: "10px", marginTop: "16px" }}>
        <button
          onClick={() => void handleSave()}
          disabled={saving || !name.trim()}
          style={{ ...primaryBtnStyle, flex: 1 }}
        >
          {saving ? "保存中..." : "保存する"}
        </button>
        <button onClick={onCancel} style={{ ...secondaryBtnStyle, flex: 1 }}>
          キャンセル
        </button>
      </div>
    </div>
  );
}

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

const primaryBtnStyle: React.CSSProperties = {
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

const secondaryBtnStyle: React.CSSProperties = {
  display: "block",
  width: "100%",
  padding: "12px",
  borderRadius: "10px",
  border: "1px solid var(--color-brown-400)",
  background: "transparent",
  color: "var(--color-brown-700)",
  fontSize: "1rem",
  fontWeight: 600,
  cursor: "pointer",
};

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
