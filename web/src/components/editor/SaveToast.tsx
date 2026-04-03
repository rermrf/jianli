interface SaveToastProps {
  visible: boolean
}

export function SaveToast({ visible }: SaveToastProps) {
  if (!visible) {
    return null
  }

  return (
    <div
      aria-live="polite"
      className="fixed bottom-6 left-1/2 -translate-x-1/2 rounded-full bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-lg"
    >
      已保存草稿
    </div>
  )
}
