import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import { VisitorsPage } from '../VisitorsPage'

describe('visitors page', () => {
  it('switches between 7-day and 30-day visitor data views', async () => {
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <VisitorsPage />
      </MemoryRouter>,
    )

    expect(screen.getByText('128')).toBeInTheDocument()
    expect(screen.getByText('3/28')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: '30天' }))

    expect(screen.getByText('462')).toBeInTheDocument()
    expect(screen.getByText('W1')).toBeInTheDocument()
  })
})
