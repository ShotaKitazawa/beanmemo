import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { StarRating } from "../components/StarRating";
import { RecordCard } from "../components/RecordCard";
import { RecordForm } from "../components/RecordForm";
import { useRecord, updateRecord, deleteRecord } from "../hooks/useRecords";
import { ROAST_LABEL, BREW_LABEL, type RecordFormValues } from "../lib/recordOptions";
import type { components } from "../api/schema.d.ts";

type UpdateRecordRequest = components["schemas"]["UpdateRecordRequest"];

export function RecordDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { record, loading, error } = useRecord(Number(id));
  const [editing, setEditing] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const handleSave = async (values: RecordFormValues) => {
    if (!record) return;
    const req: UpdateRecordRequest = {
      name: values.name,
      rating: values.rating || undefined,
      origin: values.origin || null,
      roast_level: values.roastLevel || null,
      shop: values.shop || null,
      price: values.price ? parseInt(values.price, 10) : null,
      purchased_at: values.purchasedAt || null,
      tasting_note: values.tastingNote || null,
      brew_method: values.brewMethod || null,
    };
    await updateRecord(record.id, req);
    navigate(0);
  };

  const handleDelete = async () => {
    if (!record) return;
    try {
      await deleteRecord(record.id);
      navigate("/");
    } catch {
      setDeleteError("削除に失敗しました");
    }
  };

  if (loading) {
    return (
      <div style={{ padding: "32px", textAlign: "center", color: "var(--color-brown-400)" }}>
        読み込み中...
      </div>
    );
  }
  if (error || !record) {
    return (
      <div style={{ padding: "16px" }}>
        <button onClick={() => navigate(-1)} style={backBtnStyle}>
          ← 戻る
        </button>
        <p style={{ color: "#c0392b" }}>{error ?? "記録が見つかりません"}</p>
      </div>
    );
  }

  const initialFormValues: RecordFormValues = {
    name: record.name,
    rating: record.rating ?? 0,
    origin: record.origin ?? "",
    roastLevel: (record.roast_level as RecordFormValues["roastLevel"]) ?? "",
    shop: record.shop ?? "",
    price: record.price != null ? String(record.price) : "",
    purchasedAt: record.purchased_at ?? "",
    tastingNote: record.tasting_note ?? "",
    brewMethod: (record.brew_method as RecordFormValues["brewMethod"]) ?? "",
  };

  return (
    <div style={{ padding: "16px", maxWidth: "600px", margin: "0 auto" }}>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          marginBottom: "20px",
        }}
      >
        <button onClick={() => navigate(-1)} style={backBtnStyle}>
          ←
        </button>
        <div style={{ display: "flex", gap: "8px" }}>
          {!editing && (
            <button onClick={() => setEditing(true)} style={editBtnStyle}>
              編集
            </button>
          )}
          <button onClick={() => setShowDeleteConfirm(true)} style={deleteBtnStyle}>
            削除
          </button>
        </div>
      </div>

      {showDeleteConfirm && (
        <div style={overlayStyle}>
          <div style={dialogStyle}>
            <p style={{ fontWeight: 700, marginBottom: "12px" }}>
              「{record.name}」を削除しますか？
            </p>
            {deleteError && (
              <p style={{ color: "#c0392b", fontSize: "0.85rem", marginBottom: "8px" }}>
                {deleteError}
              </p>
            )}
            <div style={{ display: "flex", gap: "8px" }}>
              <button onClick={() => void handleDelete()} style={{ ...deleteBtnStyle, flex: 1 }}>
                削除する
              </button>
              <button
                onClick={() => setShowDeleteConfirm(false)}
                style={{ ...editBtnStyle, flex: 1 }}
              >
                キャンセル
              </button>
            </div>
          </div>
        </div>
      )}

      {!editing ? (
        <div>
          <h2
            style={{
              fontSize: "1.4rem",
              fontWeight: 800,
              color: "var(--color-brown-900)",
              marginBottom: "8px",
            }}
          >
            {record.name}
          </h2>
          {record.rating != null && <StarRating value={record.rating} readonly size="md" />}

          <div style={detailGridStyle}>
            {record.origin && <DetailRow label="産地" value={record.origin} />}
            {record.roast_level && (
              <DetailRow
                label="焙煎度"
                value={ROAST_LABEL[record.roast_level] ?? record.roast_level}
              />
            )}
            {record.brew_method && (
              <DetailRow
                label="抽出方法"
                value={BREW_LABEL[record.brew_method] ?? record.brew_method}
              />
            )}
            {record.shop && <DetailRow label="購入店" value={record.shop} />}
            {record.price != null && (
              <DetailRow label="価格" value={`¥${record.price.toLocaleString()}`} />
            )}
            {record.purchased_at && <DetailRow label="購入日" value={record.purchased_at} />}
          </div>

          {record.tasting_note && (
            <div style={{ marginTop: "16px" }}>
              <h3 style={subheadStyle}>テイスティングノート</h3>
              <p
                style={{
                  fontSize: "0.95rem",
                  lineHeight: 1.7,
                  color: "var(--color-brown-900)",
                  whiteSpace: "pre-wrap",
                }}
              >
                {record.tasting_note}
              </p>
            </div>
          )}

          {!record.is_note_filled && (
            <button
              onClick={() => setEditing(true)}
              style={{ ...primaryBtnStyle, marginTop: "20px", background: "#c47a1b" }}
            >
              📝 テイスティングノートを追記する
            </button>
          )}

          {record.related_records.length > 0 && (
            <div style={{ marginTop: "24px" }}>
              <h3 style={subheadStyle}>この豆の他の記録 ({record.related_records.length}件)</h3>
              <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
                {record.related_records.map((r) => (
                  <RecordCard key={r.id} record={r} />
                ))}
              </div>
            </div>
          )}
        </div>
      ) : (
        <RecordForm
          initialValues={initialFormValues}
          onSave={handleSave}
          onCancel={() => setEditing(false)}
        />
      )}
    </div>
  );
}

function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <div style={{ display: "contents" }}>
      <dt style={{ fontSize: "0.8rem", color: "var(--color-brown-400)", fontWeight: 600 }}>
        {label}
      </dt>
      <dd style={{ fontSize: "0.9rem", color: "var(--color-brown-900)", margin: 0 }}>{value}</dd>
    </div>
  );
}

const detailGridStyle: React.CSSProperties = {
  display: "grid",
  gridTemplateColumns: "auto 1fr",
  gap: "8px 16px",
  marginTop: "16px",
  padding: "14px",
  background: "var(--color-cream-200)",
  borderRadius: "10px",
};

const subheadStyle: React.CSSProperties = {
  fontSize: "0.85rem",
  fontWeight: 700,
  color: "var(--color-brown-700)",
  margin: "0 0 8px 0",
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

const backBtnStyle: React.CSSProperties = {
  background: "transparent",
  border: "none",
  fontSize: "1.2rem",
  color: "var(--color-brown-700)",
  cursor: "pointer",
  padding: "4px 8px",
};

const editBtnStyle: React.CSSProperties = {
  padding: "8px 16px",
  borderRadius: "20px",
  border: "1px solid var(--color-brown-700)",
  background: "transparent",
  color: "var(--color-brown-700)",
  fontSize: "0.85rem",
  fontWeight: 600,
  cursor: "pointer",
};

const deleteBtnStyle: React.CSSProperties = {
  padding: "8px 16px",
  borderRadius: "20px",
  border: "1px solid #c0392b",
  background: "transparent",
  color: "#c0392b",
  fontSize: "0.85rem",
  fontWeight: 600,
  cursor: "pointer",
};

const overlayStyle: React.CSSProperties = {
  position: "fixed",
  inset: 0,
  background: "rgba(0,0,0,0.4)",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  zIndex: 100,
};

const dialogStyle: React.CSSProperties = {
  background: "var(--color-cream-100)",
  borderRadius: "14px",
  padding: "24px",
  width: "90%",
  maxWidth: "360px",
};
