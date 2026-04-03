import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from '../../App'

function renderAtPath(path: string) {
  window.history.pushState({}, '', path)

  return render(<App />)
}

describe('routing smoke test', () => {
  it('renders the public resume route content on the home page', () => {
    renderAtPath('/')

    expect(screen.getAllByText('温庆京')[0]).toBeInTheDocument()
  })

  it('redirects unauthenticated users from /edit to /login', async () => {
    renderAtPath('/edit')

    expect(await screen.findByText('管理后台')).toBeInTheDocument()
  })

  it('redirects back to the protected page after successful login', async () => {
    const user = userEvent.setup()

    renderAtPath('/edit')

    await user.type(screen.getByPlaceholderText('输入访问密钥'), 'resume-key')
    await user.click(screen.getByRole('button', { name: '验证并登录' }))

    expect(await screen.findByLabelText('姓名')).toBeInTheDocument()
  })

  it('navigates from the resume export action to the print page', async () => {
    const user = userEvent.setup()

    renderAtPath('/')

    await user.click(screen.getAllByRole('button', { name: '导出 PDF' })[0])

    expect(await screen.findByText('打印版简历')).toBeInTheDocument()
  })
})
