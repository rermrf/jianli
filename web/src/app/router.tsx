import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ProtectedRoute } from "../components/layout/ProtectedRoute";
import { DraftPreviewPage } from "../pages/DraftPreviewPage";
import { DraftsPage } from "../pages/DraftsPage";
import { EditPage } from "../pages/EditPage";
import { LoginPage } from "../pages/LoginPage";
import { PrintPage } from "../pages/PrintPage";
import { ResumePage } from "../pages/ResumePage";
import { VisitorsPage } from "../pages/VisitorsPage";

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<ResumePage />} path="/" />
        <Route
          element={
            <ProtectedRoute>
              <EditPage />
            </ProtectedRoute>
          }
          path="/edit"
        />
        <Route
          element={
            <ProtectedRoute>
              <DraftsPage />
            </ProtectedRoute>
          }
          path="/drafts"
        />
        <Route
          element={
            <ProtectedRoute>
              <DraftPreviewPage />
            </ProtectedRoute>
          }
          path="/drafts/:id"
        />
        <Route
          element={
            <ProtectedRoute>
              <VisitorsPage />
            </ProtectedRoute>
          }
          path="/visitors"
        />
        <Route element={<LoginPage />} path="/login" />
        <Route element={<PrintPage />} path="/print" />
      </Routes>
    </BrowserRouter>
  );
}
