import { Link, useParams } from "react-router-dom";
import { Button } from "../components/common/Button";
import { SectionCard } from "../components/common/SectionCard";
import { AppShell } from "../components/layout/AppShell";
import { TopNav } from "../components/layout/TopNav";
import { ResumeDesktopLayout } from "../components/resume/ResumeDesktopLayout";
import { ResumeMobileLayout } from "../components/resume/ResumeMobileLayout";
import { useResumeDraftDetail } from "../hooks/useResumeDraftDetail";

export function DraftPreviewPage() {
  const params = useParams();
  const draftID = params.id ? Number(params.id) : null;
  const { draft, error, loading, publishDraft } = useResumeDraftDetail(
    Number.isFinite(draftID) ? draftID : null,
  );

  if (loading) {
    return (
      <AppShell contentClassName="space-y-6">
        <TopNav />
        <SectionCard className="text-sm text-slate-500">加载中...</SectionCard>
      </AppShell>
    );
  }

  if (!draft || error) {
    return (
      <AppShell contentClassName="space-y-6">
        <TopNav />
        <SectionCard className="text-sm text-rose-600">
          {error ?? "草稿不存在"}
        </SectionCard>
      </AppShell>
    );
  }

  return (
    <AppShell contentClassName="space-y-6">
      <TopNav />
      <SectionCard className="space-y-3">
        <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
          <div className="space-y-1">
            <h1 className="text-2xl font-semibold text-slate-900">
              {draft.name}
            </h1>
            <p className="text-sm text-slate-500">{draft.note}</p>
            <p className="text-xs text-slate-400">
              更新时间：{draft.updatedAt}
            </p>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <Link to="/drafts">
              <Button variant="secondary">返回草稿列表</Button>
            </Link>
            <Button onClick={() => void publishDraft()}>设为主简历</Button>
          </div>
        </div>
      </SectionCard>
      <ResumeMobileLayout resume={draft.data} />
      <ResumeDesktopLayout resume={draft.data} />
    </AppShell>
  );
}
