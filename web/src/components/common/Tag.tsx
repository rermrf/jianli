import type { PropsWithChildren } from 'react'

interface TagProps {
  className?: string
}

export function Tag({ children, className = '' }: PropsWithChildren<TagProps>) {
  return (
    <span
      className={`inline-flex items-center rounded-full bg-brand-50 px-3 py-1 text-xs font-medium text-brand-600 ${className}`}
    >
      {children}
    </span>
  )
}
