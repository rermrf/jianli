import type { VisitorTrendPoint } from '../../types/visitors'
import { SectionCard } from '../common/SectionCard'

interface TrendChartProps {
  points: VisitorTrendPoint[]
}

/**
 * Bar chart of daily visit counts. Pure SVG to avoid pulling in a chart
 * library (~80KB+). The points array is expected to already be padded by
 * the caller so it represents a continuous N-day window even when some
 * days have count=0.
 */
export function TrendChart({ points }: TrendChartProps) {
  if (points.length === 0) {
    return (
      <SectionCard className="space-y-4">
        <h2 className="text-lg font-semibold text-slate-900">近期访问趋势</h2>
        <div className="rounded-2xl bg-slate-50 p-6 text-sm text-slate-400">暂无数据</div>
      </SectionCard>
    )
  }

  const maxValue = Math.max(...points.map((point) => point.count), 1)
  // Bar geometry within a 100x100 viewBox: 4px gap on each side, even
  // distribution across the remaining width. We render bars from the
  // bottom (y=100) upward.
  const innerWidth = 96
  const xOffset = 2
  const slotWidth = innerWidth / points.length
  const barWidth = Math.max(1, slotWidth * 0.7)
  const barXOffset = (slotWidth - barWidth) / 2

  return (
    <SectionCard className="space-y-4">
      <h2 className="text-lg font-semibold text-slate-900">近期访问趋势</h2>
      <div className="space-y-3 rounded-2xl bg-slate-50 p-4">
        <svg
          aria-label="访问趋势柱状图"
          className="h-52 w-full"
          preserveAspectRatio="none"
          viewBox="0 0 100 100"
          role="img"
        >
          {points.map((point, index) => {
            if (point.count === 0) {
              // Render a near-zero placeholder so axis stays aligned but
              // the bar is visually invisible.
              return null
            }
            const height = (point.count / maxValue) * 90 // leave 10% headroom
            const x = xOffset + index * slotWidth + barXOffset
            const y = 100 - height
            return (
              <rect
                key={point.date}
                fill="#4a90d9"
                height={height}
                rx={0.5}
                width={barWidth}
                x={x}
                y={y}
              >
                <title>{`${point.date}：${point.count} 次`}</title>
              </rect>
            )
          })}
        </svg>
        <div
          className="grid gap-1 text-[10px] text-slate-400"
          style={{ gridTemplateColumns: `repeat(${points.length}, minmax(0, 1fr))` }}
        >
          {points.map((point) => (
            <span key={point.date} className="truncate text-center">
              {formatBarLabel(point.date)}
            </span>
          ))}
        </div>
      </div>
    </SectionCard>
  )
}

/**
 * `2026-04-27` → `04-27`. Robust to malformed input: falls back to the raw
 * string so a date the backend returns in some other format still shows
 * something rather than NaN.
 */
function formatBarLabel(isoDate: string): string {
  if (isoDate.length >= 10 && isoDate[4] === '-' && isoDate[7] === '-') {
    return isoDate.slice(5)
  }
  return isoDate
}
