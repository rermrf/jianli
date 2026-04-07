import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";
import App from "../../App";
import { defaultResume } from "../../data/mockResume";
import { loginWithKey } from "../../lib/auth";

const draftList = [
  {
    id: 8,
    name: "面试前调整版",
    note: "补充项目亮点",
    createdAt: "2026-04-06T10:00:00Z",
    updatedAt: "2026-04-06T11:00:00Z",
  },
  {
    id: 5,
    name: "精简项目版",
    note: "删掉无关经历",
    createdAt: "2026-04-05T10:00:00Z",
    updatedAt: "2026-04-05T11:00:00Z",
  },
];

function mockDraftFetch() {
  return vi
    .spyOn(globalThis, "fetch")
    .mockImplementation(async (input, init) => {
      const url = String(input);

      if (url === "/api/resume/drafts" && (!init || init.method === "GET")) {
        return new Response(JSON.stringify({ code: 0, data: draftList }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }

      if (url === "/api/resume/drafts/8" && (!init || init.method === "GET")) {
        return new Response(
          JSON.stringify({
            code: 0,
            data: {
              ...draftList[0],
              data: {
                ...defaultResume,
                profile: {
                  ...defaultResume.profile,
                  name: "草稿预览姓名",
                },
              },
            },
          }),
          {
            status: 200,
            headers: { "Content-Type": "application/json" },
          },
        );
      }

      if (url === "/api/resume/drafts/8/publish" && init?.method === "PUT") {
        return new Response(
          JSON.stringify({ code: 0, data: { published: true } }),
          {
            status: 200,
            headers: { "Content-Type": "application/json" },
          },
        );
      }

      if (url === "/api/resume/drafts/8" && init?.method === "DELETE") {
        return new Response(
          JSON.stringify({ code: 0, data: { deleted: true } }),
          {
            status: 200,
            headers: { "Content-Type": "application/json" },
          },
        );
      }

      if (url === "/api/resume" && (!init || init.method === "GET")) {
        return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }

      return new Response(JSON.stringify({ code: 0, data: [] }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    });
}

function renderAtPath(path: string) {
  window.history.pushState({}, "", path);

  return render(<App />);
}

describe("draft pages", () => {
  afterEach(() => {
    sessionStorage.clear();
    vi.restoreAllMocks();
  });

  it("renders the drafts list and triggers publish and delete actions", async () => {
    loginWithKey("resume-key");
    const fetchSpy = mockDraftFetch();
    const user = userEvent.setup();

    renderAtPath("/drafts");

    expect(await screen.findByText("面试前调整版")).toBeInTheDocument();
    expect(screen.getByText("补充项目亮点")).toBeInTheDocument();

    await user.click(
      screen.getByRole("button", { name: "设为主简历 面试前调整版" }),
    );
    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        "/api/resume/drafts/8/publish",
        expect.objectContaining({ method: "PUT" }),
      ),
    );

    await user.click(
      screen.getByRole("button", { name: "删除草稿 面试前调整版" }),
    );
    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        "/api/resume/drafts/8",
        expect.objectContaining({ method: "DELETE" }),
      ),
    );
  });

  it("renders one draft preview page", async () => {
    loginWithKey("resume-key");
    const fetchSpy = mockDraftFetch();
    const user = userEvent.setup();

    renderAtPath("/drafts/8");

    expect(await screen.findByText("面试前调整版")).toBeInTheDocument();
    expect((await screen.findAllByText("草稿预览姓名"))[0]).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "设为主简历" }));
    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        "/api/resume/drafts/8/publish",
        expect.objectContaining({ method: "PUT" }),
      ),
    );
  });
});
