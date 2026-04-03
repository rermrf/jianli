import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { defaultResume } from '../../data/mockResume'
import { ResumePage } from '../ResumePage'

describe('resume page tracking', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('tracks a visit on mount and sends duration on pagehide', async () => {
    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockImplementation(async (input, init) => {
      const url = String(input)

      if (url === '/api/resume') {
        return new Response(JSON.stringify({ code: 0, data: defaultResume }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        })
      }

      if (url === '/api/visitors' && init?.method === 'POST') {
        return new Response(JSON.stringify({ code: 0, data: { id: 7 } }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        })
      }

      return new Response(JSON.stringify({ code: 0, data: {} }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    })
    const sendBeacon = vi.fn<(url: string, data?: BodyInit | null) => boolean>(() => true)
    vi.stubGlobal('navigator', {
      ...navigator,
      sendBeacon,
    })

    render(
      <MemoryRouter>
        <ResumePage />
      </MemoryRouter>,
    )

    expect(await screen.findAllByText('温庆京')).not.toHaveLength(0)
    await waitFor(() =>
      expect(fetchSpy).toHaveBeenCalledWith(
        '/api/visitors',
        expect.objectContaining({ method: 'POST' }),
      ),
    )

    window.dispatchEvent(new Event('pagehide'))

    expect(sendBeacon).toHaveBeenCalled()
    const firstCall = sendBeacon.mock.calls[0]
    expect(firstCall?.[0]).toBe('/api/visitors/7')
  })
})
