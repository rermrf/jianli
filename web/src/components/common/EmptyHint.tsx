interface EmptyHintProps {
  title: string
  description: string
}

export function EmptyHint({ title, description }: EmptyHintProps) {
  return (
    <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 px-4 py-5 text-center">
      <p className="text-sm font-medium text-slate-700">{title}</p>
      <p className="mt-2 text-sm leading-6 text-slate-500">{description}</p>
    </div>
  )
}
