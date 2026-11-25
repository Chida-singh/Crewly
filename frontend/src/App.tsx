import { Navigate, Route, Routes } from "react-router-dom";

import { LandingPage } from "./pages/LandingPage";
import { SheetPage } from "./pages/SheetPage";
import { AUTH_BYPASS } from "./lib/dev";
import { useSession } from "./lib/session";

function RequireSession({ children }: { children: React.ReactNode }) {
  const { session } = useSession();
  if (!session && !AUTH_BYPASS) return <Navigate to="/" replace />;
  return <>{children}</>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />
      <Route
        path="/sheet"
        element={
          <RequireSession>
            <SheetPage />
          </RequireSession>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
