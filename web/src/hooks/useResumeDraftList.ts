import { useEffect, useState } from "react";
import {
  deleteResumeDraft,
  listResumeDrafts,
  publishResumeDraft,
} from "../lib/resumeDrafts";
import type { ResumeDraftSummary } from "../types/resumeDraft";

export function useResumeDraftList() {
  const [drafts, setDrafts] = useState<ResumeDraftSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  async function refresh() {
    try {
      const nextDrafts = await listResumeDrafts();
      setDrafts(nextDrafts);
      setError(null);
    } catch (loadError) {
      setError(loadError instanceof Error ? loadError.message : "加载草稿失败");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void refresh();
  }, []);

  async function publishDraft(id: number) {
    await publishResumeDraft(id);
  }

  async function removeDraft(id: number) {
    await deleteResumeDraft(id);
    setDrafts((current) => current.filter((draft) => draft.id !== id));
  }

  return {
    drafts,
    error,
    loading,
    publishDraft,
    refresh,
    removeDraft,
  };
}
