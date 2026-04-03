import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { loginWithKey } from '../../lib/auth'
import { VisitorsPage } from '../VisitorsPage'

function visitorsResponse(totalVisits: number, labels: string[]) {
  return {
    stats: {
      totalVisits,
      todayVisits: 23,
      uniqueVisitors: 10,
      averageDurationSeconds: 120,
    },
    records: labels.map((label, index) => ({
      id: index + 1,
      ip: `112.17.45.${index + 10}`,
      city: '杭州',
      device: 'Windows',
      browser: 'Chrome',
      os: 'Windows',
      visitTime: label,
      duration: 120,
    })),
  }
}

describe('visitors page', () => {
  afterEach(() => {
    sessionStorage.clear()
    vi.restoreAllMocks()
  })

  it('switches between 7-day and 30-day visitor data views', async () => {
    loginWithKey('resume-key')
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (input) => {
      const url = String(input)

      if (url.includes('/api/visitors/stats?days=7')) {
        return new Response(
          JSON.stringify({ code: 0, data: visitorsResponse(128, ['3/28']).stats }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        )
      }

      if (url.includes('/api/visitors?days=7')) {
        return new Response(
          JSON.stringify({ code: 0, data: visitorsResponse(128, ['3/28']).records }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        )
      }

      if (url.includes('/api/visitors/stats?days=30')) {
        return new Response(
          JSON.stringify({ code: 0, data: visitorsResponse(462, ['W1']).stats }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        )
      }

      if (url.includes('/api/visitors?days=30')) {
        return new Response(
          JSON.stringify({ code: 0, data: visitorsResponse(462, ['W1']).records }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        )
      }

      return new Response(JSON.stringify({ code: 0, data: [] }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    })

    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <VisitorsPage />
      </MemoryRouter>,
    )

    expect(await screen.findByText('128')).toBeInTheDocument()
    expect((await screen.findAllByText('3/28'))[0]).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: '30天' }))

    expect(await screen.findByText('462')).toBeInTheDocument()
    expect((await screen.findAllByText('W1'))[0]).toBeInTheDocument()
  })
})
