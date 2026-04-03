import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it } from 'vitest'
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
    localStorage.clear()
  })

  it('adds a new skill tag and removes an existing one', async () => {
    const user = userEvent.setup()

    renderEditPage()

    await user.click(screen.getByLabelText('删除技能 Go'))
    expect(screen.queryByText('Go')).not.toBeInTheDocument()

    await user.type(screen.getByPlaceholderText('添加技能'), 'DDD')
    await user.click(screen.getByRole('button', { name: '添加技能' }))

    expect(screen.getByText('DDD')).toBeInTheDocument()
  })

  it('saves edited resume data to localStorage', async () => {
    const user = userEvent.setup()

    renderEditPage()

    const nameInput = screen.getByLabelText('姓名')
    await user.clear(nameInput)
    await user.type(nameInput, '测试姓名')
    await user.click(screen.getAllByRole('button', { name: '保存' })[0])

    expect(screen.getByText('已保存草稿')).toBeInTheDocument()
    expect(localStorage.getItem('resume:draft')).toContain('测试姓名')
  })
})
