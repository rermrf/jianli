interface FieldInputProps {
  label: string
  onChange: (value: string) => void
  placeholder?: string
  textarea?: boolean
  value: string
}

export function FieldInput({
  label,
  onChange,
  placeholder,
  textarea = false,
  value,
}: FieldInputProps) {
  const className =
    'w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-brand-500'

  return (
    <label className="flex flex-col gap-2 text-sm font-medium text-slate-700">
      <span>{label}</span>
      {textarea ? (
        <textarea
          className={`${className} min-h-28 resize-y`}
          onChange={(event) => onChange(event.target.value)}
          placeholder={placeholder}
          value={value}
        />
      ) : (
        <input
          className={className}
          onChange={(event) => onChange(event.target.value)}
          placeholder={placeholder}
          value={value}
        />
      )}
    </label>
  )
}
