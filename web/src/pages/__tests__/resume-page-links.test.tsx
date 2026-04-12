import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { afterEach, describe, expect, it, vi } from "vitest";
import { defaultResume } from "../../data/mockResume";
import { ResumePage } from "../ResumePage";

function mockResumePageFetch(resumeOverride = defaultResume) {
  return vi.spyOn(globalThis, "fetch").mockImplementation(async (input, init) => {
    const url = String(input);

    if (url === "/api/resume" && (!init || init.method === "GET")) {
      return new Response(
        JSON.stringify({
          code: 0,
          data: {
            resume: resumeOverride,
            siteSettings: { allowPdfExport: true },
          },
        }),
        {
          status: 200,
          headers: { "Content-Type": "application/json" },
        },
      );
    }

    if (url === "/api/visitors" && init?.method === "POST") {
      return new Response(JSON.stringify({ code: 0, data: { id: 7 } }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }

    return new Response(JSON.stringify({ code: 0, data: {} }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  });
}

describe("resume page project links", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("renders project titles as links when url is present", async () => {
    mockResumePageFetch({
      ...defaultResume,
      projects: [
        {
          ...defaultResume.projects[0],
          url: "https://github.com/example/ai-gateway",
        },
      ],
    });

    render(
      <MemoryRouter>
        <ResumePage />
      </MemoryRouter>,
    );

    const projectLinks = await screen.findAllByRole("link", {
      name: "AI Gateway - LLM 统一网关",
    });
    expect(projectLinks[0]).toHaveAttribute(
      "href",
      "https://github.com/example/ai-gateway",
    );
  });

  it("keeps project titles as plain text when url is absent", async () => {
    mockResumePageFetch({
      ...defaultResume,
      projects: [
        {
          ...defaultResume.projects[0],
          url: undefined,
        },
      ],
    });

    render(
      <MemoryRouter>
        <ResumePage />
      </MemoryRouter>,
    );

    expect(
      await screen.findAllByText("AI Gateway - LLM 统一网关"),
    ).not.toHaveLength(0);
    expect(
      screen.queryByRole("link", { name: "AI Gateway - LLM 统一网关" }),
    ).not.toBeInTheDocument();
  });
});
