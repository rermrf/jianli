import type { ResumeData } from "./resume";

export interface ResumeDraftSummary {
  id: number;
  name: string;
  note: string;
  createdAt: string;
  updatedAt: string;
}

export interface ResumeDraftDetail extends ResumeDraftSummary {
  data: ResumeData;
}

export interface CreateResumeDraftInput {
  name: string;
  note: string;
  data: ResumeData;
}
