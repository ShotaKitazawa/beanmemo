import { useAuth } from "../hooks/useAuth";

export function LoginPage() {
  const { login } = useAuth();

  return (
    <div
      style={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        gap: "24px",
        background: "var(--color-cream-50)",
        padding: "24px",
      }}
    >
      <h1 style={{ fontSize: "2rem", color: "var(--color-brown-800)" }}>☕ beanmemo</h1>
      <p style={{ color: "var(--color-brown-600)", textAlign: "center" }}>
        コーヒー豆の記録を始めましょう
      </p>
      <button
        onClick={() => void login()}
        style={{
          padding: "12px 32px",
          fontSize: "1rem",
          background: "var(--color-brown-700)",
          color: "#fff",
          border: "none",
          borderRadius: "8px",
          cursor: "pointer",
        }}
      >
        ログイン
      </button>
    </div>
  );
}
