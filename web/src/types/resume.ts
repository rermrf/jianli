export interface ResumeProfile {
  avatarUrl?: string
  name: string
  title: string
  age: number
  gender: string
  education: string
  experience: string
  location: string
  hometown: string
  phone: string
  email: string
}

export interface JobIntention {
  position: string
  cities: string[]
  availability: string
}

export interface EducationExperience {
  school: string
  major: string
  degree: string
  startDate: string
  endDate: string
}

export interface WorkExperience {
  company: string
  role: string
  startDate: string
  endDate: string
  description: string[]
}

export interface ProjectExperience {
  name: string
  startDate: string
  endDate: string
  description: string[]
  url?: string
}

export interface Award {
  date: string
  title: string
}

export interface ResumeData {
  profile: ResumeProfile
  skills: string[]
  jobIntention: JobIntention
  education: EducationExperience[]
  workExperience: WorkExperience[]
  projects: ProjectExperience[]
  awards: Award[]
}
