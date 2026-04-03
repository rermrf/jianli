import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { defaultResume } from '../../data/mockResume'
import { loginWithKey } from '../../lib/auth'
import { EditPage } from '../EditPage'

function renderEditPage() {
  return render(
    <MemoryRouter>
      <EditPage />
    </MemoryRouter>,
  )
}

describe('edit page', () => {
  afterEach(() => {
    sessionStorage.clear()
    vi.restoreAllMocks()
  })

  it('adds a new skill tag and removes an existing one', async () => {
    loginWithKey('resume-key')
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (input, init) => {
      if (String(input) === '/api/resume' && (!init || init.method === 'GET')) {
        return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        })
      }

      return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    })

    const user = userEvent.setup()
    renderEditPage()

    await screen.findByLabelText('删除技能 Go')
    await user.click(screen.getByLabelText('删除技能 Go'))
    expect(screen.queryByText('Go')).not.toBeInTheDocument()

    await user.type(screen.getByPlaceholderText('添加技能'), 'DDD')
    await user.click(screen.getByRole('button', { name: '添加技能' }))

    expect(screen.getByText('DDD')).toBeInTheDocument()
  })

  it('saves edited resume data through the backend API', async () => {
    loginWithKey('resume-key')
    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockImplementation(
      async (input, init) => {
        const url = String(input)

        if (url === '/api/resume' && (!init || init.method === 'GET')) {
          return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          })
        }

        if (url === '/api/resume' && init?.method === 'PUT') {
          return new Response(
            JSON.stringify({ code: 0, data: JSON.parse(String(init.body)) }),
            {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
            },
          )
        }

        return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        })
      },
    )

    const user = userEvent.setup()
    renderEditPage()

    const nameInput = await screen.findByLabelText('姓名')
    await user.clear(nameInput)
    await user.type(nameInput, '测试姓名')
    await user.click(screen.getAllByRole('button', { name: '保存' })[0])

    expect(await screen.findByText('已保存草稿')).toBeInTheDocument()
    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        '/api/resume',
        expect.objectContaining({
          headers: expect.objectContaining({ 'X-Auth-Key': 'resume-key' }),
          method: 'PUT',
        }),
      ),
    )
  })
})
