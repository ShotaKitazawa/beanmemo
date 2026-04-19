import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { loadOIDCSetup } from "../oidc";

export function CallbackPage() {
  const navigate = useNavigate();

  useEffect(() => {
    loadOIDCSetup()
      .then(({ userManager }) => {
        if (!userManager) {
          navigate("/", { replace: true });
          return;
        }
        return userManager.signinRedirectCallback().then(() => navigate("/", { replace: true }));
      })
      .catch((err) => {
        console.error("OIDC callback error:", err);
        navigate("/login", { replace: true });
      });
  }, [navigate]);

  return (
    <div
      style={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <p>認証中...</p>
    </div>
  );
}
