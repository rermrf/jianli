import type { PropsWithChildren, ReactNode } from 'react'
import { Button } from '../common/Button'
import { SectionCard } from '../common/SectionCard'

interface EditableListSectionProps {
  addLabel: string
  children: ReactNode
  onAdd: () => void
  title: string
}

interface EditableListItemProps {
  children: ReactNode
  onRemove: () => void
  title: string
}

export function EditableListSection({
  addLabel,
  children,
  onAdd,
  title,
}: EditableListSectionProps) {
  return (
    <SectionCard className="space-y-4">
      <div className="flex items-center justify-between gap-3">
        <h2 className="text-lg font-semibold text-slate-900">{title}</h2>
        <Button onClick={onAdd} type="button" variant="secondary">
          {addLabel}
        </Button>
      </div>
      <div className="space-y-4">{children}</div>
    </SectionCard>
  )
}

export function EditableListItem({
  children,
  onRemove,
  title,
}: PropsWithChildren<EditableListItemProps>) {
  return (
    <div className="space-y-4 rounded-2xl border border-slate-200 bg-slate-50 p-4">
      <div className="flex items-center justify-between gap-3">
        <h3 className="text-sm font-semibold text-slate-800">{title}</h3>
        <Button onClick={onRemove} type="button" variant="ghost">
          删除
        </Button>
      </div>
      {children}
    </div>
  )
}
