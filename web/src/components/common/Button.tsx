import type { ButtonHTMLAttributes, PropsWithChildren } from 'react'

type ButtonVariant = 'primary' | 'secondary' | 'ghost'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
}

const variantClassName: Record<ButtonVariant, string> = {
  primary:
    'bg-brand-500 text-white shadow-sm hover:bg-brand-600 focus-visible:outline-brand-500',
  secondary:
    'border border-slate-200 bg-white text-slate-700 hover:border-slate-300 hover:bg-slate-50 focus-visible:outline-slate-400',
  ghost: 'bg-transparent text-slate-500 hover:bg-slate-100 focus-visible:outline-slate-300',
}

export function Button({
  children,
  className = '',
  type = 'button',
  variant = 'primary',
  ...props
}: PropsWithChildren<ButtonProps>) {
  return (
    <button
      className={`inline-flex items-center justify-center rounded-xl px-4 py-2 text-sm font-medium transition focus-visible:outline-2 focus-visible:outline-offset-2 ${variantClassName[variant]} ${className}`}
      type={type}
      {...props}
    >
      {children}
    </button>
  )
}
