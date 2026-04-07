import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { afterEach, describe, expect, it, vi } from "vitest";
import { defaultResume } from "../../data/mockResume";
import { loginWithKey } from "../../lib/auth";
import { EditPage } from "../EditPage";

function renderEditPage() {
  return render(
    <MemoryRouter>
      <EditPage />
    </MemoryRouter>,
  );
}

function mockResumeFetch(resumeOverride = defaultResume) {
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

      if (url === "/api/resume" && init?.method === "PUT") {
        return new Response(
          JSON.stringify({ code: 0, data: JSON.parse(String(init.body)) }),
          {
            status: 200,
            headers: { "Content-Type": "application/json" },
          },
        );
      }

      if (url === "/api/resume/drafts" && init?.method === "POST") {
        const body = JSON.parse(String(init.body));

        return new Response(
          JSON.stringify({
            code: 0,
            data: {
              id: 8,
              name: body.name,
              note: body.note,
              data: body.data,
              createdAt: "2026-04-07T10:00:00Z",
              updatedAt: "2026-04-07T10:00:00Z",
            },
          }),
          {
            status: 200,
            headers: { "Content-Type": "application/json" },
          },
        );
      }

      return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    });
}

describe("edit page", () => {
  afterEach(() => {
    sessionStorage.clear();
    vi.restoreAllMocks();
  });

  it("adds a new skill tag and removes an existing one", async () => {
    loginWithKey("resume-key");
    mockResumeFetch();

    const user = userEvent.setup();
    renderEditPage();

    await screen.findByLabelText("删除技能 Go");
    await user.click(screen.getByLabelText("删除技能 Go"));
    expect(screen.queryByText("Go")).not.toBeInTheDocument();

    await user.type(screen.getByPlaceholderText("添加技能"), "DDD");
    await user.click(screen.getByRole("button", { name: "添加技能" }));

    expect(screen.getByText("DDD")).toBeInTheDocument();
  });

  it("keeps focus while typing in the award name field", async () => {
    loginWithKey("resume-key");
    mockResumeFetch();

    const user = userEvent.setup();
    renderEditPage();

    const awardNameInput = (await screen.findAllByLabelText("奖项名称"))[0];
    await user.click(awardNameInput);
    await user.type(awardNameInput, "A");

    const currentAwardNameInput = screen.getAllByLabelText("奖项名称")[0];
    expect(currentAwardNameInput).toHaveFocus();
  });

  it("disables move controls at the list boundaries", async () => {
    loginWithKey("resume-key");
    mockResumeFetch();

    renderEditPage();

    expect(
      await screen.findByRole("button", { name: "上移 教育经历 1" }),
    ).toBeDisabled();
    expect(
      screen.getByRole("button", { name: "下移 项目经历 3" }),
    ).toBeDisabled();
  });

  it("moves work experience items down in the visible order", async () => {
    loginWithKey("resume-key");
    mockResumeFetch({
      ...defaultResume,
      workExperience: [
        defaultResume.workExperience[0],
        {
          company: "第二家公司",
          role: "后端工程师",
          startDate: "2024.01",
          endDate: "2025.01",
          description: ["负责第二段经历"],
        },
      ],
    });

    const user = userEvent.setup();
    renderEditPage();

    expect((await screen.findAllByLabelText("公司名称"))[0]).toHaveValue(
      "杭州云徕科技有限公司",
    );

    await user.click(screen.getByRole("button", { name: "下移 工作经历 1" }));

    const companyInputs = screen.getAllByLabelText("公司名称");
    expect(companyInputs[0]).toHaveValue("第二家公司");
    expect(companyInputs[1]).toHaveValue("杭州云徕科技有限公司");
  });

  it("uploads a cropped avatar and includes avatarUrl in resume save", async () => {
    loginWithKey("resume-key");
    const fetchSpy = vi
      .spyOn(globalThis, "fetch")
      .mockImplementation(async (input, init) => {
        const url = String(input);

        if (url === "/api/resume" && (!init || init.method === "GET")) {
          return new Response(
            JSON.stringify({ code: 0, data: defaultResume }),
            {
              status: 200,
              headers: { "Content-Type": "application/json" },
            },
          );
        }

        if (url === "/api/upload/avatar" && init?.method === "POST") {
          return new Response(
            JSON.stringify({
              code: 0,
              data: { url: "/uploads/avatars/avatar-1.png" },
            }),
            {
              status: 200,
              headers: { "Content-Type": "application/json" },
            },
          );
        }

        if (url === "/api/resume" && init?.method === "PUT") {
          return new Response(
            JSON.stringify({ code: 0, data: JSON.parse(String(init.body)) }),
            {
              status: 200,
              headers: { "Content-Type": "application/json" },
            },
          );
        }

        return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      });
    vi.stubGlobal("URL", {
      ...URL,
      createObjectURL: vi.fn(() => "blob:avatar-preview"),
      revokeObjectURL: vi.fn(),
    });

    const user = userEvent.setup();
    renderEditPage();

    const file = new File([new Uint8Array([1, 2, 3])], "avatar.png", {
      type: "image/png",
    });

    await user.upload(await screen.findByLabelText("选择头像"), file);
    expect(await screen.findByAltText("头像预览")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "确认上传头像" }));
    await waitFor(() => expect(fetchSpy.mock.calls.length).toBeGreaterThan(1));
    await user.click(screen.getAllByRole("button", { name: "保存主简历" })[0]);

    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        "/api/upload/avatar",
        expect.objectContaining({ method: "POST" }),
      ),
    );

    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        "/api/resume",
        expect.objectContaining({
          body: expect.stringContaining("/uploads/avatars/avatar-1.png"),
          method: "PUT",
        }),
      ),
    );
  });

  it("saves reordered project order through the main resume api", async () => {
    loginWithKey("resume-key");
    const fetchSpy = mockResumeFetch();

    const user = userEvent.setup();
    renderEditPage();

    await user.click(
      await screen.findByRole("button", { name: "下移 项目经历 1" }),
    );
    await user.click(screen.getAllByRole("button", { name: "保存主简历" })[0]);

    await waitFor(() => {
      const saveCall = fetchSpy.mock.calls.find(
        ([url, init]) =>
          String(url) === "/api/resume" && init?.method === "PUT",
      );
      expect(saveCall).toBeDefined();
      const body = JSON.parse(String(saveCall?.[1]?.body));
      expect(body.projects[0].name).toBe("Kafka 运维消息服务");
      expect(body.projects[1].name).toBe("AI Gateway - LLM 统一网关");
    });
  });

  it("saves reordered project order through the draft api", async () => {
    loginWithKey("resume-key");
    const fetchSpy = mockResumeFetch();

    const user = userEvent.setup();
    renderEditPage();

    await user.click(
      await screen.findByRole("button", { name: "下移 项目经历 1" }),
    );
    await user.click(screen.getAllByRole("button", { name: "保存为草稿" })[0]);
    await user.type(screen.getByLabelText("草稿名称"), "顺序调整草稿");
    await user.type(screen.getByLabelText("草稿备注"), "验证项目顺序");
    await user.click(screen.getByRole("button", { name: "确认保存草稿" }));

    await waitFor(() => {
      const saveCall = fetchSpy.mock.calls.find(
        ([url, init]) =>
          String(url) === "/api/resume/drafts" && init?.method === "POST",
      );
      expect(saveCall).toBeDefined();
      const body = JSON.parse(String(saveCall?.[1]?.body));
      expect(body.data.projects[0].name).toBe("Kafka 运维消息服务");
      expect(body.data.projects[1].name).toBe("AI Gateway - LLM 统一网关");
    });
  });

  it("saves edited resume data through the backend API", async () => {
    loginWithKey("resume-key");
    const fetchSpy = mockResumeFetch();

    const user = userEvent.setup();
    renderEditPage();

    const nameInput = await screen.findByLabelText("姓名");
    await user.clear(nameInput);
    await user.type(nameInput, "测试姓名");
    await user.click(screen.getAllByRole("button", { name: "保存主简历" })[0]);

    expect(await screen.findByText("已保存主简历")).toBeInTheDocument();
    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        "/api/resume",
        expect.objectContaining({
          headers: expect.objectContaining({ "X-Auth-Key": "resume-key" }),
          method: "PUT",
        }),
      ),
    );
  });
});
