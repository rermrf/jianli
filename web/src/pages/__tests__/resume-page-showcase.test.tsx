import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { defaultResume } from '../../data/mockResume'
import { ResumePage } from '../ResumePage'

function mockFetch(resume = defaultResume, allowPdfExport = true) {
  return vi.spyOn(globalThis, 'fetch').mockImplementation(async (input, init) => {
    const url = String(input)
    if (url === '/api/resume' && (!init || init.method === 'GET')) {
      return new Response(
        JSON.stringify({
          code: 0,
          data: { resume, siteSettings: { allowPdfExport } },
        }),
        { status: 200, headers: { 'Content-Type': 'application/json' } },
      )
    }
    if (url === '/api/visitors' && init?.method === 'POST') {
      return new Response(JSON.stringify({ code: 0, data: { id: 1 } }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    }
    return new Response(JSON.stringify({ code: 0, data: {} }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    })
  })
}

describe('resume page showcase', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('uses scoped public root and capability-first headings', async () => {
    mockFetch()
    const { container } = render(
      <MemoryRouter>
        <ResumePage />
      </MemoryRouter>,
    )

    expect(await screen.findAllByText('温庆京')).not.toHaveLength(0)
    expect(container.querySelector('.resume-public')).not.toBeNull()
    expect(screen.getAllByRole('heading', { name: '项目经历' }).length).toBeGreaterThan(0)
    expect(screen.getAllByRole('heading', { name: '工作经历' }).length).toBeGreaterThan(0)
    expect(screen.queryByRole('link', { name: '编辑' })).not.toBeInTheDocument()
  })

  it('hides PDF controls when export disabled', async () => {
    mockFetch(defaultResume, false)
    render(
      <MemoryRouter>
        <ResumePage />
      </MemoryRouter>,
    )
    expect(await screen.findAllByText('温庆京')).not.toHaveLength(0)
    expect(screen.queryByRole('link', { name: '导出 PDF' })).not.toBeInTheDocument()
  })

  it('hides empty project section', async () => {
    mockFetch({ ...defaultResume, projects: [] })
    render(
      <MemoryRouter>
        <ResumePage />
      </MemoryRouter>,
    )
    expect(await screen.findAllByText('温庆京')).not.toHaveLength(0)
    expect(screen.queryByRole('heading', { name: '项目经历' })).not.toBeInTheDocument()
  })
})
