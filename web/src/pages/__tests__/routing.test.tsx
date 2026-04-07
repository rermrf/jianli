import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";
import App from "../../App";
import { defaultResume } from "../../data/mockResume";
import { loginWithKey } from "../../lib/auth";

function mockAppFetch(resumeOverride = defaultResume) {
  return vi
    .spyOn(globalThis, "fetch")
    .mockImplementation(async (input, init) => {
      const url = String(input);

      if (url === "/api/resume" && (!init || init.method === "GET")) {
        return new Response(JSON.stringify({ code: 0, data: resumeOverride }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }

      if (url === "/api/auth/verify" && init?.method === "POST") {
        return new Response(
          JSON.stringify({ code: 0, data: { valid: true } }),
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

      if (url === "/api/resume/drafts" && (!init || init.method === "GET")) {
        return new Response(JSON.stringify({ code: 0, data: [] }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }

      if (url === "/api/resume/drafts/8" && (!init || init.method === "GET")) {
        return new Response(
          JSON.stringify({
            code: 0,
            data: {
              id: 8,
              name: "面试前调整版",
              note: "补充项目亮点",
              data: resumeOverride,
              createdAt: "2026-04-06T10:00:00Z",
              updatedAt: "2026-04-06T10:00:00Z",
            },
          }),
          {
            status: 200,
            headers: { "Content-Type": "application/json" },
          },
        );
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

describe("routing smoke test", () => {
  afterEach(() => {
    sessionStorage.clear();
    vi.restoreAllMocks();
  });

  it("shows only 简历 in nav when the user is not authenticated", async () => {
    mockAppFetch();
    renderAtPath("/");

    await screen.findAllByText("温庆京");
    expect(screen.getByText("简历")).toBeInTheDocument();
    expect(screen.queryByText("编辑")).not.toBeInTheDocument();
    expect(screen.queryByText("草稿")).not.toBeInTheDocument();
    expect(screen.queryByText("访客")).not.toBeInTheDocument();
  });

  it("redirects unauthenticated users from /edit to /login", async () => {
    mockAppFetch();
    renderAtPath("/edit");

    expect(await screen.findByText("管理后台")).toBeInTheDocument();
  });

  it("redirects unauthenticated users from /drafts to /login", async () => {
    mockAppFetch();
    renderAtPath("/drafts");

    expect(await screen.findByText("管理后台")).toBeInTheDocument();
  });

  it("redirects unauthenticated users from /drafts/8 to /login", async () => {
    mockAppFetch();
    renderAtPath("/drafts/8");

    expect(await screen.findByText("管理后台")).toBeInTheDocument();
  });

  it("enters login page after triple-clicking the avatar within 1.5 seconds", async () => {
    const user = userEvent.setup();
    mockAppFetch();
    renderAtPath("/");

    const avatar = (
      await screen.findAllByRole("button", { name: "头像入口" })
    )[0];
    await user.click(avatar);
    await user.click(avatar);
    await user.click(avatar);

    expect(await screen.findByText("管理后台")).toBeInTheDocument();
  });

  it("redirects back to the protected page after successful login", async () => {
    const user = userEvent.setup();
    const fetchSpy = mockAppFetch();

    renderAtPath("/edit");

    await user.type(screen.getByPlaceholderText("输入访问密钥"), "resume-key");
    await user.click(screen.getByRole("button", { name: "验证并登录" }));

    expect(await screen.findByLabelText("姓名")).toBeInTheDocument();
    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        "/api/auth/verify",
        expect.objectContaining({ method: "POST" }),
      ),
    );
  });

  it("shows 编辑, 草稿 and 访客 in nav when the user is authenticated", async () => {
    loginWithKey("resume-key");
    mockAppFetch();
    renderAtPath("/");

    await screen.findAllByText("温庆京");
    expect(screen.getByText("简历")).toBeInTheDocument();
    expect(screen.getByText("编辑")).toBeInTheDocument();
    expect(screen.getByText("草稿")).toBeInTheDocument();
    expect(screen.getByText("访客")).toBeInTheDocument();
  });

  it("renders uploaded avatar on the public resume page", async () => {
    mockAppFetch({
      ...defaultResume,
      profile: {
        ...defaultResume.profile,
        avatarUrl: "/uploads/avatars/avatar-1.png",
      },
    });
    renderAtPath("/");

    const avatars = await screen.findAllByAltText("温庆京头像");
    expect(avatars[0]).toHaveAttribute("src", "/uploads/avatars/avatar-1.png");
  });

  it("navigates from the resume export action to the print page", async () => {
    const user = userEvent.setup();
    mockAppFetch({
      ...defaultResume,
      profile: {
        ...defaultResume.profile,
        avatarUrl: "/uploads/avatars/avatar-1.png",
      },
    });
    renderAtPath("/");

    await user.click(screen.getAllByRole("button", { name: "导出 PDF" })[0]);

    expect(await screen.findByText("打印版简历")).toBeInTheDocument();
    const avatars = await screen.findAllByAltText("温庆京头像");
    expect(avatars[0]).toHaveAttribute("src", "/uploads/avatars/avatar-1.png");
  });
});
