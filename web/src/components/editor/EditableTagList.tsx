import { useState } from 'react'
import { Button } from '../common/Button'
import { Tag } from '../common/Tag'

interface EditableTagListProps {
  onAdd: (skill: string) => void
  onRemove: (skill: string) => void
  skills: string[]
}

export function EditableTagList({
  onAdd,
  onRemove,
  skills,
}: EditableTagListProps) {
  const [nextSkill, setNextSkill] = useState('')

  function handleAddSkill() {
    const trimmedSkill = nextSkill.trim()

    if (!trimmedSkill) {
      return
    }

    onAdd(trimmedSkill)
    setNextSkill('')
  }

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap gap-2">
        {skills.map((skill) => (
          <span
            key={skill}
            className="inline-flex items-center gap-2 rounded-full bg-brand-50 px-3 py-1.5 text-xs font-medium text-brand-600"
          >
            <Tag className="bg-transparent px-0 py-0 text-inherit">{skill}</Tag>
            <button
              aria-label={`删除技能 ${skill}`}
              className="text-brand-400 transition hover:text-brand-600"
              onClick={() => onRemove(skill)}
              type="button"
            >
              ×
            </button>
          </span>
        ))}
      </div>
      <div className="flex flex-col gap-3 sm:flex-row">
        <input
          className="min-w-0 flex-1 rounded-2xl border border-slate-200 px-4 py-3 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-brand-500"
          onChange={(event) => setNextSkill(event.target.value)}
          placeholder="添加技能"
          value={nextSkill}
        />
        <Button className="sm:px-5" onClick={handleAddSkill} type="button">
          添加技能
        </Button>
      </div>
    </div>
  )
}
