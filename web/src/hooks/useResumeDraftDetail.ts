import { useEffect, useState } from "react";
import {
  deleteResumeDraft,
  getResumeDraft,
  publishResumeDraft,
} from "../lib/resumeDrafts";
import type { ResumeDraftDetail } from "../types/resumeDraft";

export function useResumeDraftDetail(id: number | null) {
  const [draft, setDraft] = useState<ResumeDraftDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (id === null) {
      setDraft(null);
      setError("无效草稿编号");
      setLoading(false);
      return;
    }

    const draftID = id;
    let active = true;

    async function loadDraft() {
      try {
        const nextDraft = await getResumeDraft(draftID);
        if (!active) {
          return;
        }
        setDraft(nextDraft);
        setError(null);
      } catch (loadError) {
        if (!active) {
          return;
        }
        setError(
          loadError instanceof Error ? loadError.message : "加载草稿失败",
        );
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    }

    void loadDraft();

    return () => {
      active = false;
    };
  }, [id]);

  async function publishDraft() {
    if (id === null) {
      return;
    }

    await publishResumeDraft(id);
  }

  async function removeDraft() {
    if (id === null) {
      return;
    }

    await deleteResumeDraft(id);
  }

  return {
    draft,
    error,
    loading,
    publishDraft,
    removeDraft,
  };
}
