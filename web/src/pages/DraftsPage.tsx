import { Link } from "react-router-dom";
import { Button } from "../components/common/Button";
import { EmptyHint } from "../components/common/EmptyHint";
import { SectionCard } from "../components/common/SectionCard";
import { AppShell } from "../components/layout/AppShell";
import { TopNav } from "../components/layout/TopNav";
import { useResumeDraft } from "../hooks/useResumeDraft";
import { useResumeDraftList } from "../hooks/useResumeDraftList";

export function DraftsPage() {
  const { loading: settingsLoading, siteSettings } = useResumeDraft();
  const showPdfExport = !settingsLoading && siteSettings.allowPdfExport;
  const { drafts, error, loading, publishDraft, removeDraft } =
    useResumeDraftList();

  if (loading) {
    return (
      <AppShell contentClassName="space-y-6">
        <TopNav showPdfExport={showPdfExport} />
        <SectionCard className="text-sm text-slate-500">加载中...</SectionCard>
      </AppShell>
    );
  }

  return (
    <AppShell contentClassName="space-y-6">
      <TopNav showPdfExport={showPdfExport} />
      <SectionCard className="space-y-2">
        <h1 className="text-2xl font-semibold text-slate-900">草稿版本</h1>
        <p className="text-sm text-slate-500">
          在这里预览、发布或清理管理员保存的简历草稿。
        </p>
      </SectionCard>

      {error ? (
        <SectionCard className="text-sm text-rose-600">{error}</SectionCard>
      ) : null}

      {drafts.length === 0 ? (
        <EmptyHint
          description="先去编辑页保存一个草稿版本，再回来管理它们。"
          title="还没有保存过草稿"
        />
      ) : (
        <div className="space-y-4">
          {drafts.map((draft) => (
            <SectionCard className="space-y-4" key={draft.id}>
              <div className="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
                <div className="space-y-1">
                  <h2 className="text-lg font-semibold text-slate-900">
                    {draft.name}
                  </h2>
                  <p className="text-sm text-slate-500">{draft.note}</p>
                  <p className="text-xs text-slate-400">
                    更新时间：{draft.updatedAt}
                  </p>
                </div>
                <div className="flex flex-wrap items-center gap-2">
                  <Link to={`/drafts/${draft.id}`}>
                    <Button variant="secondary">预览</Button>
                  </Link>
                  <Button
                    aria-label={`设为主简历 ${draft.name}`}
                    onClick={() => void publishDraft(draft.id)}
                  >
                    设为主简历
                  </Button>
                  <Button
                    aria-label={`删除草稿 ${draft.name}`}
                    onClick={() => void removeDraft(draft.id)}
                    variant="ghost"
                  >
                    删除草稿
                  </Button>
                </div>
              </div>
            </SectionCard>
          ))}
        </div>
      )}
    </AppShell>
  );
}
