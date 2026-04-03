import type { VisitorRecord } from '../../types/visitors'
import { SectionCard } from '../common/SectionCard'

interface VisitorListProps {
  records: VisitorRecord[]
}

export function VisitorList({ records }: VisitorListProps) {
  return (
    <div className="space-y-3 md:hidden">
      {records.map((record) => (
        <SectionCard key={record.id} className="flex items-start gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-slate-100 text-slate-400">
            💻
          </div>
          <div className="space-y-1">
            <p className="text-sm font-medium text-slate-800">
              {record.city} · {record.browser} · {record.device}
            </p>
            <p className="text-xs text-slate-400">
              {record.visitTime} · IP {record.ip}
            </p>
          </div>
        </SectionCard>
      ))}
    </div>
  )
}
