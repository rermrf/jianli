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
