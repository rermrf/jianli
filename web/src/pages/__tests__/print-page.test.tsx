import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { defaultResume } from '../../data/mockResume'
import { PrintPage } from '../PrintPage'

describe('print page', () => {
  afterEach(() => {
    vi.restoreAllMocks()
    delete (window as { __printReady?: boolean }).__printReady
  })

  it('sets window.__printReady to true once the resume has loaded so the chromedp exporter can proceed', async () => {
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (input) => {
      const url = String(input)

      if (url === '/api/resume') {
        return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        })
      }

      return new Response(JSON.stringify({ code: 0, data: {} }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    })

    expect((window as { __printReady?: boolean }).__printReady).toBeUndefined()

    render(
      <MemoryRouter>
        <PrintPage />
      </MemoryRouter>,
    )

    expect(await screen.findByText(defaultResume.profile.name)).toBeInTheDocument()

    await waitFor(() => {
      expect((window as { __printReady?: boolean }).__printReady).toBe(true)
    })
  })

  it('renders profile facts in the header plus education and awards from the published resume', async () => {
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (input) => {
      const url = String(input)

      if (url === '/api/resume') {
        return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        })
      }

      return new Response(JSON.stringify({ code: 0, data: {} }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    })

    render(
      <MemoryRouter>
        <PrintPage />
      </MemoryRouter>,
    )

    expect(await screen.findByText(defaultResume.profile.name)).toBeInTheDocument()
    expect(screen.queryByText('打印版简历')).not.toBeInTheDocument()
    expect(screen.queryByText('个人基本信息')).not.toBeInTheDocument()
    expect(screen.getByText('25岁 / 男 / 本科 / 0.9年 / 籍贯：江西赣州')).toBeInTheDocument()
    expect(screen.getByText('教育经历')).toBeInTheDocument()
    expect(screen.getByText(defaultResume.education[0].school)).toBeInTheDocument()
    expect(screen.getByText('荣誉奖项')).toBeInTheDocument()
    expect(screen.getAllByText(defaultResume.awards[0].title)).not.toHaveLength(0)
  })

  it('downloads a PDF file from the backend when the download button is clicked', async () => {
    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockImplementation(async (input) => {
      const url = String(input)

      if (url === '/api/resume') {
        return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        })
      }

      if (url === '/api/resume/pdf') {
        return new Response(new Uint8Array([1, 2, 3]), {
          status: 200,
          headers: { 'Content-Type': 'application/pdf' },
        })
      }

      return new Response(JSON.stringify({ code: 0, data: {} }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    })

    const createObjectURL = vi.fn(() => 'blob:resume-pdf')
    const revokeObjectURL = vi.fn()
    vi.stubGlobal('URL', {
      ...URL,
      createObjectURL,
      revokeObjectURL,
    })

    const appendChild = vi.spyOn(document.body, 'appendChild')
    const removeChild = vi.spyOn(document.body, 'removeChild')
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <PrintPage />
      </MemoryRouter>,
    )

    expect(await screen.findByText(defaultResume.profile.name)).toBeInTheDocument()
    expect(screen.queryByText('打印版简历')).not.toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: '下载 PDF' }))

    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        '/api/resume/pdf',
        expect.objectContaining({ method: 'GET' }),
      ),
    )

    expect(createObjectURL).toHaveBeenCalled()
    expect(appendChild).toHaveBeenCalled()
    expect(removeChild).toHaveBeenCalled()
    expect(revokeObjectURL).toHaveBeenCalledWith('blob:resume-pdf')
  })

  it('hides the download button when pdf export is disabled', async () => {
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (input) => {
      const url = String(input)

      if (url === '/api/resume') {
        return new Response(
          JSON.stringify({
            code: 0,
            data: {
              resume: defaultResume,
              siteSettings: { allowPdfExport: false },
            },
          }),
          {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          },
        )
      }

      return new Response(JSON.stringify({ code: 0, data: {} }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    })

    render(
      <MemoryRouter>
        <PrintPage />
      </MemoryRouter>,
    )

    expect(await screen.findByText(defaultResume.profile.name)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '下载 PDF' })).not.toBeInTheDocument()
  })
})
