import { useEffect, useRef } from 'react'
import {
  readVisitorID,
  readVisitorStartedAt,
  recordVisit,
  sendVisitDuration,
  storeVisitorID,
} from '../lib/visitors'

const HEARTBEAT_INTERVAL_MS = 15_000

/**
 * useVisitorTracking records a visit on mount (or reuses an existing
 * sessionStorage-cached visitor for cross-page sessions like / → /print)
 * and sends a duration heartbeat every 15 seconds. The heartbeat reports
 * cumulative seconds since the FIRST page in the session was opened, so
 * navigating between pages does not reset duration to 0.
 *
 * `pagehide` fires a final sendBeacon with the latest cumulative duration,
 * tightening the upper bound on lost time when an ad-blocker/CSP rule
 * happens to drop the periodic fetch heartbeat.
 *
 * @param enabled  False while the parent is still loading; setting it
 *                 true triggers recording. Setting it false stops the
 *                 heartbeat (but the row keeps its last reported value).
 */
export function useVisitorTracking(enabled: boolean) {
  const visitIDRef = useRef<number | null>(null)
  const startedAtRef = useRef<number | null>(null)

  useEffect(() => {
    if (!enabled) return

    let cancelled = false

    async function start() {
      try {
        const existing = readVisitorID()
        if (existing !== null) {
          visitIDRef.current = existing
          // Reuse the original session start time so duration accumulates
          // across page navigations within the same tab.
          startedAtRef.current = readVisitorStartedAt() ?? Date.now()
        } else {
          const id = await recordVisit()
          if (cancelled) return
          storeVisitorID(id)
          visitIDRef.current = id
          startedAtRef.current = Date.now()
        }
      } catch {
        // Recording failure is non-fatal — page should still render.
      }
    }

    function elapsedSeconds(): number | null {
      if (visitIDRef.current === null || startedAtRef.current === null) return null
      return Math.max(0, Math.round((Date.now() - startedAtRef.current) / 1000))
    }

    function heartbeat() {
      const id = visitIDRef.current
      const seconds = elapsedSeconds()
      if (id === null || seconds === null) return
      void sendVisitDuration(id, seconds)
    }

    function flushBeacon() {
      const id = visitIDRef.current
      const seconds = elapsedSeconds()
      if (id === null || seconds === null) return
      const blob = new Blob([JSON.stringify({ duration: seconds })], {
        type: 'application/json',
      })
      try {
        navigator.sendBeacon(`/api/visitors/${id}`, blob)
      } catch {
        // sendBeacon is best-effort; ignore failures.
      }
    }

    void start()
    const intervalId = window.setInterval(heartbeat, HEARTBEAT_INTERVAL_MS)
    window.addEventListener('pagehide', flushBeacon)

    return () => {
      cancelled = true
      window.clearInterval(intervalId)
      window.removeEventListener('pagehide', flushBeacon)
    }
  }, [enabled])
}
