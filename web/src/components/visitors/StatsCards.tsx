import { formatDuration } from '../../lib/format'
import type { VisitorStats } from '../../types/visitors'
import { SectionCard } from '../common/SectionCard'

interface StatsCardsProps {
  stats: VisitorStats
}

const cardPalette = [
  'bg-sky-50 text-sky-600',
  'bg-emerald-50 text-emerald-600',
  'bg-amber-50 text-amber-600',
  'bg-fuchsia-50 text-fuchsia-600',
]

export function StatsCards({ stats }: StatsCardsProps) {
  const cards = [
    { label: '总访问量', value: stats.totalVisits.toString() },
    { label: '今日访问', value: stats.todayVisits.toString() },
    { label: '独立访客', value: stats.uniqueVisitors.toString() },
    { label: '平均停留', value: formatDuration(stats.averageDurationSeconds) },
  ]

  return (
    <div className="grid gap-4 md:grid-cols-4">
      {cards.map((card, index) => (
        <SectionCard key={card.label} className={`p-5 ${cardPalette[index]}`}>
          <p className="text-3xl font-semibold">{card.value}</p>
          <p className="mt-2 text-sm text-slate-500">{card.label}</p>
        </SectionCard>
      ))}
    </div>
  )
}
