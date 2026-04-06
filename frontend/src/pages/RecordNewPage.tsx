import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { StarRating } from "../components/StarRating";
import { TastingNoteInput } from "../components/TastingNoteInput";
import { createRecord } from "../hooks/useRecords";
import { useApiKey } from "../hooks/useApiKey";
import {
  ROAST_OPTIONS,
  BREW_OPTIONS,
  type RoastLevel,
  type BrewMethod,
} from "../lib/recordOptions";
import type { components } from "../api/schema.d.ts";

type CreateRecordRequest = components["schemas"]["CreateRecordRequest"];

const DRAFT_KEY = "beanmemo_record_new_draft";

interface Draft {
  step: 1 | 2;
  name: string;
  origin: string;
  roastLevel: RoastLevel | "";
  shop: string;
  price: string;
  purchasedAt: string;
  rating: number;
  tastingNote: string;
  brewMethod: BrewMethod | "";
}

function loadDraft(): Draft {
  try {
    const raw = sessionStorage.getItem(DRAFT_KEY);
    if (raw) return JSON.parse(raw) as Draft;
  } catch {}
  return {
    step: 1,
    name: "",
    origin: "",
    roastLevel: "",
    shop: "",
    price: "",
    purchasedAt: new Date().toISOString().slice(0, 10),
    rating: 0,
    tastingNote: "",
    brewMethod: "",
  };
}

function clearDraft() {
  sessionStorage.removeItem(DRAFT_KEY);
}

export function RecordNewPage() {
  const navigate = useNavigate();
  const { provider, apiKey, model } = useApiKey();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const initial = loadDraft();
  const [step, setStep] = useState<1 | 2>(initial.step);
  const [name, setName] = useState(initial.name);
  const [origin, setOrigin] = useState(initial.origin);
  const [roastLevel, setRoastLevel] = useState<RoastLevel | "">(initial.roastLevel);
  const [shop, setShop] = useState(initial.shop);
  const [price, setPrice] = useState(initial.price);
  const [purchasedAt, setPurchasedAt] = useState(initial.purchasedAt);
  const [rating, setRating] = useState(initial.rating);
  const [tastingNote, setTastingNote] = useState(initial.tastingNote);
  const [brewMethod, setBrewMethod] = useState<BrewMethod | "">(initial.brewMethod);

  useEffect(() => {
    const draft: Draft = {
      step,
      name,
      origin,
      roastLevel,
      shop,
      price,
      purchasedAt,
      rating,
      tastingNote,
      brewMethod,
    };
    sessionStorage.setItem(DRAFT_KEY, JSON.stringify(draft));
  }, [step, name, origin, roastLevel, shop, price, purchasedAt, rating, tastingNote, brewMethod]);

  const buildRequest = (withStep2: boolean): CreateRecordRequest => {
    const req: CreateRecordRequest = { name: name.trim() };
    if (origin) req.origin = origin;
    if (roastLevel) req.roast_level = roastLevel;
    if (shop) req.shop = shop;
    if (price) req.price = parseInt(price, 10);
    if (purchasedAt) req.purchased_at = purchasedAt;
    if (withStep2) {
      if (rating) req.rating = rating;
      if (tastingNote) req.tasting_note = tastingNote;
      if (brewMethod) req.brew_method = brewMethod;
    }
    return req;
  };

  const handleSaveLater = async () => {
    if (!name.trim()) return;
    setSubmitting(true);
    setError(null);
    try {
      await createRecord(buildRequest(false));
      clearDraft();
      navigate("/");
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存に失敗しました");
    } finally {
      setSubmitting(false);
    }
  };

  const handleSaveWithStep2 = async () => {
    if (!name.trim()) return;
    setSubmitting(true);
    setError(null);
    try {
      const created = await createRecord(buildRequest(true));
      clearDraft();
      navigate(`/records/${created.id}`);
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存に失敗しました");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div style={{ padding: "16px", maxWidth: "600px", margin: "0 auto" }}>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          gap: "12px",
          marginBottom: "20px",
        }}
      >
        <button onClick={() => navigate(-1)} style={backBtnStyle}>
          ←
        </button>
        <h2
          style={{
            fontSize: "1.2rem",
            fontWeight: 700,
            color: "var(--color-brown-900)",
            margin: 0,
          }}
        >
          新しく記録する
        </h2>
      </div>

      {/* Step 1 */}
      <section style={sectionStyle}>
        <h3 style={sectionTitleStyle}>Step 1 — 基本情報</h3>

        <label style={labelStyle}>豆の名前 *</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="例: エチオピア イルガチェフェ"
          style={inputStyle}
        />

        <label style={labelStyle}>産地</label>
        <input
          type="text"
          value={origin}
          onChange={(e) => setOrigin(e.target.value)}
          placeholder="例: エチオピア"
          style={inputStyle}
        />

        <label style={labelStyle}>焙煎度</label>
        <div style={{ display: "flex", gap: "8px", flexWrap: "wrap", marginBottom: "16px" }}>
          {ROAST_OPTIONS.map((o) => (
            <button
              key={o.value}
              onClick={() => setRoastLevel((prev) => (prev === o.value ? "" : o.value))}
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
          placeholder="例: 〇〇コーヒー"
          style={inputStyle}
        />

        <label style={labelStyle}>価格（円）</label>
        <input
          type="number"
          value={price}
          onChange={(e) => setPrice(e.target.value)}
          placeholder="例: 1200"
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

        {error && step === 1 && <p style={{ color: "#c0392b", fontSize: "0.85rem" }}>{error}</p>}

        {step === 1 && (
          <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
            <button onClick={() => setStep(2)} disabled={!name.trim()} style={primaryBtnStyle}>
              Step2もまとめて入力する
            </button>
            <button
              onClick={() => void handleSaveLater()}
              disabled={!name.trim() || submitting}
              style={secondaryBtnStyle}
            >
              {submitting ? "保存中..." : "あとで追記する"}
            </button>
          </div>
        )}
      </section>

      {/* Step 2 */}
      {step === 2 && (
        <section style={sectionStyle}>
          <h3 style={sectionTitleStyle}>Step 2 — テイスティング（任意）</h3>

          <label style={labelStyle}>総合評価</label>
          <div style={{ marginBottom: "16px" }}>
            <StarRating value={rating} onChange={setRating} size="lg" />
            {rating === 0 && (
              <span
                style={{ fontSize: "0.78rem", color: "var(--color-brown-400)", marginLeft: "8px" }}
              >
                タップして評価
              </span>
            )}
          </div>

          <label style={labelStyle}>抽出方法</label>
          <div
            style={{
              display: "flex",
              gap: "8px",
              flexWrap: "wrap",
              marginBottom: "16px",
            }}
          >
            {BREW_OPTIONS.map((o) => (
              <button
                key={o.value}
                onClick={() => setBrewMethod((prev) => (prev === o.value ? "" : o.value))}
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

          {error && <p style={{ color: "#c0392b", fontSize: "0.85rem" }}>{error}</p>}

          <button
            onClick={() => void handleSaveWithStep2()}
            disabled={!name.trim() || submitting}
            style={{ ...primaryBtnStyle, marginTop: "16px" }}
          >
            {submitting ? "保存中..." : "保存する"}
          </button>
        </section>
      )}
    </div>
  );
}

const sectionStyle: React.CSSProperties = {
  background: "var(--color-cream-200)",
  borderRadius: "12px",
  padding: "16px",
  marginBottom: "16px",
};

const sectionTitleStyle: React.CSSProperties = {
  fontSize: "0.9rem",
  fontWeight: 700,
  color: "var(--color-brown-700)",
  marginBottom: "14px",
  margin: "0 0 14px 0",
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

const backBtnStyle: React.CSSProperties = {
  background: "transparent",
  border: "none",
  fontSize: "1.2rem",
  color: "var(--color-brown-700)",
  cursor: "pointer",
  padding: "4px 8px",
};
