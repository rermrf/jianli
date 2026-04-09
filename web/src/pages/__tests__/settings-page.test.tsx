import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { afterEach, describe, expect, it, vi } from "vitest";
import { defaultResume } from "../../data/mockResume";
import { loginWithKey } from "../../lib/auth";
import { SettingsPage } from "../SettingsPage";

function renderSettingsPage() {
  return render(
    <MemoryRouter>
      <SettingsPage />
    </MemoryRouter>,
  );
}

describe("settings page", () => {
  afterEach(() => {
    sessionStorage.clear();
    vi.restoreAllMocks();
  });

  it("renders the pdf export toggle and saves it through the settings api", async () => {
    loginWithKey("resume-key");
    const fetchSpy = vi.spyOn(globalThis, "fetch").mockImplementation(
      async (input, init) => {
        const url = String(input);

        if (url === "/api/resume" && (!init || init.method === "GET")) {
          return new Response(
            JSON.stringify({
              code: 0,
              data: {
                resume: defaultResume,
                siteSettings: { allowPdfExport: true },
              },
            }),
            {
              status: 200,
              headers: { "Content-Type": "application/json" },
            },
          );
        }

        if (url === "/api/settings" && init?.method === "PUT") {
          return new Response(
            JSON.stringify({ code: 0, data: JSON.parse(String(init.body)) }),
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
      },
    );

    const user = userEvent.setup();
    renderSettingsPage();

    const toggle = await screen.findByLabelText("允许访客导出 PDF");
    expect(toggle).toBeChecked();

    await user.click(toggle);

    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        "/api/settings",
        expect.objectContaining({
          body: JSON.stringify({ allowPdfExport: false }),
          headers: expect.objectContaining({ "X-Auth-Key": "resume-key" }),
          method: "PUT",
        }),
      ),
    );
  });
});
