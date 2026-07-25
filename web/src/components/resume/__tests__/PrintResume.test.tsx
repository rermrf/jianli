import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import type { ResumeData } from '../../../types/resume'
import { defaultResume } from '../../../data/mockResume'
import { PrintResume } from '../PrintResume'

describe('PrintResume', () => {
  it('renders sections in print-optimized order: skills, work, projects, education, awards', () => {
    const { container } = render(<PrintResume resume={defaultResume} />)

    const sectionTitles = Array.from(
      container.querySelectorAll('.print-section-title'),
    ).map((node) => node.textContent?.trim())

    expect(sectionTitles).toEqual([
      '个人技能',
      '工作经历',
      '项目经历',
      '教育经历',
      '荣誉奖项',
    ])
  })

  it('hides a section entirely (heading included) when its array is empty', () => {
    const noAwards: ResumeData = { ...defaultResume, awards: [] }
    render(<PrintResume resume={noAwards} />)

    expect(screen.queryByText('荣誉奖项')).not.toBeInTheDocument()
  })

  it('hides every optional section when all data arrays are empty, keeping only the header', () => {
    const sparse: ResumeData = {
      ...defaultResume,
      skills: [],
      workExperience: [],
      projects: [],
      education: [],
      awards: [],
    }
    const { container } = render(<PrintResume resume={sparse} />)

    expect(screen.queryByText('个人技能')).not.toBeInTheDocument()
    expect(screen.queryByText('工作经历')).not.toBeInTheDocument()
    expect(screen.queryByText('项目经历')).not.toBeInTheDocument()
    expect(screen.queryByText('教育经历')).not.toBeInTheDocument()
    expect(screen.queryByText('荣誉奖项')).not.toBeInTheDocument()

    // Header stays
    expect(screen.getByText(defaultResume.profile.name)).toBeInTheDocument()

    // No items either
    expect(container.querySelectorAll('.print-item').length).toBe(0)
  })

  it('marks every work / project / education / award entry with the print-item class so smart pagination can hold them together', () => {
    const { container } = render(<PrintResume resume={defaultResume} />)

    const expected =
      defaultResume.workExperience.length +
      defaultResume.projects.length +
      defaultResume.education.length +
      defaultResume.awards.length

    expect(container.querySelectorAll('.print-item').length).toBe(expected)
  })

  it('exposes an id="print-root" sentinel for the chromedp PDF exporter to wait on', () => {
    const { container } = render(<PrintResume resume={defaultResume} />)
    expect(container.querySelector('#print-root')).not.toBeNull()
  })

  it('renders project links as anchors using the brand color and an external indicator', () => {
    render(<PrintResume resume={defaultResume} />)

    const linkedProject = defaultResume.projects.find((p) => p.url)!
    const link = screen.getByRole('link', { name: new RegExp(linkedProject.name) })

    expect(link).toHaveAttribute('href', linkedProject.url!)
    expect(link).toHaveAttribute('target', '_blank')
    expect(link.className).toContain('text-brand-600')
  })

  it('joins skills with a middle-dot separator instead of rendering pill tags', () => {
    render(<PrintResume resume={defaultResume} />)
    expect(screen.getByText(defaultResume.skills.join(' · '))).toBeInTheDocument()
  })
})
