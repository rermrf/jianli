import type { VisitorTrendPoint } from '../../types/visitors'
import { SectionCard } from '../common/SectionCard'

interface TrendChartProps {
  points: VisitorTrendPoint[]
}

export function TrendChart({ points }: TrendChartProps) {
  const maxValue = Math.max(...points.map((point) => point.value), 1)
  const chartPoints = points
    .map((point, index) => {
      const x = (index / Math.max(points.length - 1, 1)) * 100
      const y = 100 - (point.value / maxValue) * 70

      return `${x},${y}`
    })
    .join(' ')

  return (
    <SectionCard className="space-y-4">
      <h2 className="text-lg font-semibold text-slate-900">近况访问趋势</h2>
      <div className="space-y-4 rounded-2xl bg-slate-50 p-4">
        <svg
          className="h-52 w-full"
          preserveAspectRatio="none"
          viewBox="0 0 100 100"
        >
          <polyline
            fill="none"
            points={chartPoints}
            stroke="#4a90d9"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
          />
        </svg>
        <div className="grid grid-cols-4 gap-3 text-xs text-slate-400 md:grid-cols-7">
          {points.map((point) => (
            <span key={point.label}>{point.label}</span>
          ))}
        </div>
      </div>
    </SectionCard>
  )
}
