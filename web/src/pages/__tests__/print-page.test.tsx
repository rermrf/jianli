import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { defaultResume } from '../../data/mockResume'
import { PrintPage } from '../PrintPage'

describe('print page', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('renders profile facts, education, and awards from the published resume', async () => {
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
    expect(screen.getByText(String(defaultResume.profile.age))).toBeInTheDocument()
    expect(screen.getByText(defaultResume.profile.gender)).toBeInTheDocument()
    expect(screen.getByText(defaultResume.profile.education)).toBeInTheDocument()
    expect(screen.getByText(defaultResume.profile.experience)).toBeInTheDocument()
    expect(screen.getByText(defaultResume.profile.hometown)).toBeInTheDocument()
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

    expect(await screen.findByText('打印版简历')).toBeInTheDocument()
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
})
