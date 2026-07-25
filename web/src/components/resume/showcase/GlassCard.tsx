import type { PropsWithChildren } from 'react'

interface GlassCardProps {
  className?: string
  strong?: boolean
}

export function GlassCard({
  children,
  className = '',
  strong = false,
}: PropsWithChildren<GlassCardProps>) {
  return (
    <section
      className={`${strong ? 'resume-glass-strong' : 'resume-glass'} rounded-[20px] p-6 ${className}`}
    >
      {children}
    </section>
  )
}
