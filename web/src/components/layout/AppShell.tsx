import type { PropsWithChildren } from 'react'

interface AppShellProps {
  contentClassName?: string
}

export function AppShell({
  children,
  contentClassName = '',
}: PropsWithChildren<AppShellProps>) {
  return (
    <main className="min-h-screen bg-white text-slate-900 md:bg-slate-50">
      <div className={`mx-auto max-w-6xl px-4 py-6 md:px-6 md:py-10 ${contentClassName}`}>
        {children}
      </div>
    </main>
  )
}
