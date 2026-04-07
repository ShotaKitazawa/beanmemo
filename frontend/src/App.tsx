import { BrowserRouter, Routes, Route, NavLink, Navigate } from "react-router-dom";
import { AuthProvider } from "./auth/AuthProvider";
import { useAuth } from "./hooks/useAuth";
import { HomePage } from "./pages/HomePage";
import { RecordNewPage } from "./pages/RecordNewPage";
import { RecordDetailPage } from "./pages/RecordDetailPage";
import { RecommendPage } from "./pages/RecommendPage";
import { ProfilePage } from "./pages/ProfilePage";
import { SettingsPage } from "./pages/SettingsPage";
import { LoginPage } from "./pages/LoginPage";
import { CallbackPage } from "./pages/CallbackPage";
import type { ReactNode } from "react";

export function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
}

function AppRoutes() {
  return (
    <div style={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}>
      <main style={{ flex: 1, paddingBottom: "80px" }}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/callback" element={<CallbackPage />} />
          <Route
            path="/*"
            element={
              <RequireAuth>
                <AuthenticatedRoutes />
              </RequireAuth>
            }
          />
        </Routes>
      </main>
    </div>
  );
}

function RequireAuth({ children }: { children: ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div
        style={{
          minHeight: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <p>読み込み中...</p>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}

function AuthenticatedRoutes() {
  return (
    <>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/records/new" element={<RecordNewPage />} />
        <Route path="/records/:id" element={<RecordDetailPage />} />
        <Route path="/recommend" element={<RecommendPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/settings" element={<SettingsPage />} />
      </Routes>

      {/* Bottom navigation */}
      <nav
        style={{
          position: "fixed",
          bottom: 0,
          left: 0,
          right: 0,
          background: "var(--color-cream-100)",
          borderTop: "1px solid var(--color-cream-200)",
          display: "flex",
          justifyContent: "space-around",
          padding: "8px 0 calc(8px + env(safe-area-inset-bottom))",
          boxShadow: "0 -2px 12px rgba(59,31,10,0.06)",
          zIndex: 50,
        }}
      >
        <NavItem to="/" icon="🏠" label="ホーム" />
        <NavItem to="/recommend" icon="🎯" label="おすすめ" />
        <NavItem to="/profile" icon="☕" label="好み" />
        <NavItem to="/settings" icon="⚙️" label="設定" />
      </nav>
    </>
  );
}

function NavItem({ to, icon, label }: { to: string; icon: string; label: string }) {
  return (
    <NavLink
      to={to}
      end={to === "/"}
      style={({ isActive }) => ({
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: "2px",
        textDecoration: "none",
        color: isActive ? "var(--color-brown-700)" : "var(--color-brown-400)",
        fontSize: "0.65rem",
        fontWeight: isActive ? 700 : 400,
        minWidth: "48px",
      })}
    >
      <span style={{ fontSize: "1.3rem" }}>{icon}</span>
      {label}
    </NavLink>
  );
}
