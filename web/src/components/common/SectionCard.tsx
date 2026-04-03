import type { PropsWithChildren } from 'react'

interface SectionCardProps {
  className?: string
}

export function SectionCard({
  children,
  className = '',
}: PropsWithChildren<SectionCardProps>) {
  return (
    <section
      className={`rounded-3xl border border-slate-200 bg-white p-6 shadow-[var(--shadow-card)] ${className}`}
    >
      {children}
    </section>
  )
}
