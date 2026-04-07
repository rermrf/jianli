import { apiFetch } from "./api";
import { getAuthKey } from "./auth";
import type {
  CreateResumeDraftInput,
  ResumeDraftDetail,
  ResumeDraftSummary,
} from "../types/resumeDraft";

function withAuthHeaders() {
  return {
    "X-Auth-Key": getAuthKey() ?? "",
  };
}

export async function createResumeDraft(input: CreateResumeDraftInput) {
  return apiFetch<ResumeDraftDetail>("/api/resume/drafts", {
    body: JSON.stringify(input),
    headers: withAuthHeaders(),
    method: "POST",
  });
}

export async function listResumeDrafts() {
  return apiFetch<ResumeDraftSummary[]>("/api/resume/drafts", {
    headers: withAuthHeaders(),
    method: "GET",
  });
}

export async function getResumeDraft(id: number) {
  return apiFetch<ResumeDraftDetail>(`/api/resume/drafts/${id}`, {
    headers: withAuthHeaders(),
    method: "GET",
  });
}

export async function publishResumeDraft(id: number) {
  return apiFetch<{ published: boolean }>(`/api/resume/drafts/${id}/publish`, {
    headers: withAuthHeaders(),
    method: "PUT",
  });
}

export async function deleteResumeDraft(id: number) {
  return apiFetch<{ deleted: boolean }>(`/api/resume/drafts/${id}`, {
    headers: withAuthHeaders(),
    method: "DELETE",
  });
}
