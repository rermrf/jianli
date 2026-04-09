import { useState } from "react";
import { AppShell } from "../components/layout/AppShell";
import { TopNav } from "../components/layout/TopNav";
import { SectionCard } from "../components/common/SectionCard";
import { SaveToast } from "../components/editor/SaveToast";
import { useResumeDraft } from "../hooks/useResumeDraft";
import { ApiError } from "../lib/api";
import { updateSiteSettings } from "../lib/settings";

export function SettingsPage() {
  const { loading, siteSettings, setSiteSettings } = useResumeDraft();
  const [error, setError] = useState<string | null>(null);
  const [toastMessage, setToastMessage] = useState<string | null>(null);

  function showToast(message: string) {
    setToastMessage(message);
    window.setTimeout(() => setToastMessage(null), 1800);
  }

  async function handleTogglePdfExport(allowPdfExport: boolean) {
    const previousSettings = siteSettings;
    setSiteSettings({ allowPdfExport });
    setError(null);

    try {
      const nextSettings = await updateSiteSettings({ allowPdfExport });
      setSiteSettings(nextSettings);
      showToast(allowPdfExport ? "已开启 PDF 导出" : "已关闭 PDF 导出");
    } catch (nextError) {
      setSiteSettings(previousSettings);
      setError(
        nextError instanceof ApiError ? nextError.message : "保存站点设置失败",
      );
    }
  }

  if (loading) {
    return (
      <AppShell contentClassName="space-y-6">
        <TopNav showPdfExport={siteSettings.allowPdfExport} />
        <SectionCard className="text-sm text-slate-500">加载中...</SectionCard>
      </AppShell>
    );
  }

  return (
    <AppShell contentClassName="space-y-6">
      <TopNav showPdfExport={siteSettings.allowPdfExport} />
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">设置</h1>
        <p className="mt-2 text-sm text-slate-500">
          管理站点级行为开关，后续新增设置项也统一放在这里。
        </p>
      </div>

      <SectionCard className="space-y-4">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">站点设置</h2>
          <p className="mt-2 text-sm text-slate-500">
            关闭后，公开页和打印页都会隐藏导出按钮，后端接口也会拒绝 PDF 导出。
          </p>
        </div>
        <label className="flex items-center justify-between gap-4 rounded-2xl border border-slate-200 px-4 py-3">
          <div className="space-y-1">
            <span className="text-sm font-medium text-slate-900">
              允许访客导出 PDF
            </span>
            <p className="text-xs text-slate-500">
              关闭后访客无法下载 PDF 简历。
            </p>
          </div>
          <input
            aria-label="允许访客导出 PDF"
            checked={siteSettings.allowPdfExport}
            className="h-4 w-4"
            onChange={(event) => {
              void handleTogglePdfExport(event.target.checked);
            }}
            type="checkbox"
          />
        </label>
        {error ? (
          <p className="text-sm text-rose-600">{error}</p>
        ) : null}
      </SectionCard>

      <SaveToast
        message={toastMessage ?? undefined}
        visible={toastMessage !== null}
      />
    </AppShell>
  );
}
