import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import App from '../../App'
import { defaultResume } from '../../data/mockResume'

function mockAppFetch() {
  return vi.spyOn(globalThis, 'fetch').mockImplementation(async (input, init) => {
    const url = String(input)

    if (url === '/api/resume' && (!init || init.method === 'GET')) {
      return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    }

    if (url === '/api/auth/verify' && init?.method === 'POST') {
      return new Response(JSON.stringify({ code: 0, data: { valid: true } }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    }

    return new Response(JSON.stringify({ code: 0, data: [] }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    })
  })
}

function renderAtPath(path: string) {
  window.history.pushState({}, '', path)

  return render(<App />)
}

describe('routing smoke test', () => {
  afterEach(() => {
    sessionStorage.clear()
    vi.restoreAllMocks()
  })

  it('renders the public resume route content on the home page', async () => {
    mockAppFetch()
    renderAtPath('/')

    expect(await screen.findAllByText('温庆京')).not.toHaveLength(0)
  })

  it('redirects unauthenticated users from /edit to /login', async () => {
    mockAppFetch()
    renderAtPath('/edit')

    expect(await screen.findByText('管理后台')).toBeInTheDocument()
  })

  it('redirects back to the protected page after successful login', async () => {
    const user = userEvent.setup()
    const fetchSpy = mockAppFetch()

    renderAtPath('/edit')

    await user.type(screen.getByPlaceholderText('输入访问密钥'), 'resume-key')
    await user.click(screen.getByRole('button', { name: '验证并登录' }))

    expect(await screen.findByLabelText('姓名')).toBeInTheDocument()
    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        '/api/auth/verify',
        expect.objectContaining({ method: 'POST' }),
      ),
    )
  })

  it('navigates from the resume export action to the print page', async () => {
    const user = userEvent.setup()
    mockAppFetch()
    renderAtPath('/')

    await user.click(screen.getAllByRole('button', { name: '导出 PDF' })[0])

    expect(await screen.findByText('打印版简历')).toBeInTheDocument()
  })
})
