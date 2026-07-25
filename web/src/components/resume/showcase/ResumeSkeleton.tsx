export function ResumeSkeleton() {
  return (
    <div className="space-y-4" aria-busy="true" aria-live="polite">
      <div className="resume-glass h-16 animate-pulse rounded-[20px]" />
      <div className="resume-glass-strong h-40 animate-pulse rounded-[20px]" />
      <div className="resume-glass h-56 animate-pulse rounded-[20px]" />
      <span className="sr-only">加载中...</span>
    </div>
  )
}
