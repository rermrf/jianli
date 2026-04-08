import { useState } from "react";
import { useEffect } from "react";
import { Button } from "../components/common/Button";
import { SectionCard } from "../components/common/SectionCard";
import { AvatarUploader } from "../components/editor/AvatarUploader";
import {
  EditableListItem,
  EditableListSection,
} from "../components/editor/EditableListSection";
import { EditableTagList } from "../components/editor/EditableTagList";
import { FieldInput } from "../components/editor/FieldInput";
import { SaveDraftDialog } from "../components/editor/SaveDraftDialog";
import { SaveToast } from "../components/editor/SaveToast";
import { AppShell } from "../components/layout/AppShell";
import { TopNav } from "../components/layout/TopNav";
import { useResumeDraft } from "../hooks/useResumeDraft";
import { ApiError } from "../lib/api";
import { createResumeDraft } from "../lib/resumeDrafts";
import type {
  Award,
  EducationExperience,
  ProjectExperience,
  ResumeData,
  WorkExperience,
} from "../types/resume";

function createEmptyWorkExperience(): WorkExperience {
  return {
    company: "",
    role: "",
    startDate: "",
    endDate: "",
    description: [""],
  };
}

function createEmptyProject(): ProjectExperience {
  return {
    name: "",
    startDate: "",
    endDate: "",
    description: [""],
  };
}

function createEmptyEducation(): EducationExperience {
  return {
    school: "",
    major: "",
    degree: "",
    startDate: "",
    endDate: "",
  };
}

function createEmptyAward(): Award {
  return {
    date: "",
    title: "",
  };
}

function moveArrayItem<T>(items: T[], fromIndex: number, toIndex: number) {
  if (
    fromIndex < 0 ||
    toIndex < 0 ||
    fromIndex >= items.length ||
    toIndex >= items.length ||
    fromIndex === toIndex
  ) {
    return items;
  }

  const nextItems = [...items];
  const [item] = nextItems.splice(fromIndex, 1);
  nextItems.splice(toIndex, 0, item);
  return nextItems;
}

export function EditPage() {
  const { draft, loading, saveDraft, setDraft } = useResumeDraft();
  const [draftDialogOpen, setDraftDialogOpen] = useState(false);
  const [draftError, setDraftError] = useState<string | null>(null);
  const [draftSaving, setDraftSaving] = useState(false);
  const [toastMessage, setToastMessage] = useState<string | null>(null);
  const [desiredCitiesInput, setDesiredCitiesInput] = useState("");

  function parseDesiredCities(value: string) {
    return value.split(/[，,]/).map((item) => item.trim()).filter(Boolean);
  }

  function showToast(message: string) {
    setToastMessage(message);
    window.setTimeout(() => setToastMessage(null), 1800);
  }

  useEffect(() => {
    setDesiredCitiesInput(draft.jobIntention.cities.join(", "));
  }, [draft.jobIntention.cities]);

  function updateProfileField<Field extends keyof typeof draft.profile>(
    field: Field,
    value: (typeof draft.profile)[Field],
  ) {
    setDraft({
      ...draft,
      profile: {
        ...draft.profile,
        [field]: value,
      },
    });
  }

  function updateSkills(nextSkills: string[]) {
    setDraft({
      ...draft,
      skills: nextSkills,
    });
  }

  function updateJobIntentionCities(value: string) {
    const cities = parseDesiredCities(value);

    setDesiredCitiesInput(value);

    setDraft({
      ...draft,
      jobIntention: {
        ...draft.jobIntention,
        cities,
      },
    });
  }

  function updateWorkExperience(index: number, nextItem: WorkExperience) {
    const nextItems = [...draft.workExperience];
    nextItems[index] = nextItem;

    setDraft({
      ...draft,
      workExperience: nextItems,
    });
  }

  function updateProject(index: number, nextItem: ProjectExperience) {
    const nextItems = [...draft.projects];
    nextItems[index] = nextItem;

    setDraft({
      ...draft,
      projects: nextItems,
    });
  }

  function updateEducation(index: number, nextItem: EducationExperience) {
    const nextItems = [...draft.education];
    nextItems[index] = nextItem;

    setDraft({
      ...draft,
      education: nextItems,
    });
  }

  function updateAward(index: number, nextItem: Award) {
    const nextItems = [...draft.awards];
    nextItems[index] = nextItem;

    setDraft({
      ...draft,
      awards: nextItems,
    });
  }

  function buildDraftForSave(): ResumeData {
    return {
      ...draft,
      jobIntention: {
        ...draft.jobIntention,
        cities: parseDesiredCities(desiredCitiesInput),
      },
    };
  }

  async function handleSaveMainResume() {
    await saveDraft(buildDraftForSave());
    showToast("已保存主简历");
  }

  async function handleSaveAsDraft(name: string, note: string) {
    if (!name) {
      setDraftError("请输入草稿名称");
      return;
    }

    try {
      setDraftSaving(true);
      await createResumeDraft({
        data: buildDraftForSave(),
        name,
        note,
      });
      setDraftDialogOpen(false);
      setDraftError(null);
      showToast("草稿已保存");
    } catch (error) {
      setDraftError(error instanceof ApiError ? error.message : "保存草稿失败");
    } finally {
      setDraftSaving(false);
    }
  }

  function openSaveDraftDialog() {
    setDraftError(null);
    setDraftDialogOpen(true);
  }

  if (loading) {
    return (
      <AppShell contentClassName="space-y-6">
        <TopNav />
        <SectionCard className="text-sm text-slate-500">加载中...</SectionCard>
      </AppShell>
    );
  }

  return (
    <AppShell contentClassName="space-y-6">
      <TopNav />
      <div className="flex items-center justify-between gap-4 md:hidden">
        <Button variant="ghost">取消</Button>
        <h1 className="text-lg font-semibold text-slate-900">编辑简历</h1>
        <div className="flex items-center gap-2">
          <Button onClick={openSaveDraftDialog} variant="secondary">
            保存为草稿
          </Button>
          <Button onClick={handleSaveMainResume}>保存主简历</Button>
        </div>
      </div>
      <div className="hidden items-center justify-between gap-4 md:flex">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900">编辑简历</h1>
          <p className="mt-2 text-sm text-slate-500">
            保存主简历会直接覆盖当前线上简历，保存为草稿则会生成一个独立版本。
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="secondary">取消</Button>
          <Button onClick={openSaveDraftDialog} variant="secondary">
            保存为草稿
          </Button>
          <Button onClick={handleSaveMainResume}>保存主简历</Button>
        </div>
      </div>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,480px)_minmax(0,1fr)]">
        <div className="space-y-6">
          <SectionCard className="space-y-4">
            <h2 className="text-lg font-semibold text-slate-900">基本信息</h2>
            <AvatarUploader
              currentUrl={draft.profile.avatarUrl}
              onUploaded={(avatarUrl) =>
                setDraft({
                  ...draft,
                  profile: {
                    ...draft.profile,
                    avatarUrl,
                  },
                })
              }
            />
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-1">
              <FieldInput
                label="姓名"
                onChange={(value) => updateProfileField("name", value)}
                value={draft.profile.name}
              />
              <FieldInput
                label="职位"
                onChange={(value) => updateProfileField("title", value)}
                value={draft.profile.title}
              />
              <FieldInput
                label="年龄"
                onChange={(value) =>
                  updateProfileField("age", Number.parseInt(value, 10) || 0)
                }
                value={String(draft.profile.age)}
              />
              <FieldInput
                label="性别"
                onChange={(value) => updateProfileField("gender", value)}
                value={draft.profile.gender}
              />
              <FieldInput
                label="学历等级"
                onChange={(value) => updateProfileField("education", value)}
                value={draft.profile.education}
              />
              <FieldInput
                label="工作年限"
                onChange={(value) => updateProfileField("experience", value)}
                value={draft.profile.experience}
              />
              <FieldInput
                label="当前所在地"
                onChange={(value) => updateProfileField("location", value)}
                value={draft.profile.location}
              />
              <FieldInput
                label="意向城市"
                onChange={updateJobIntentionCities}
                placeholder="多个城市用逗号分隔"
                value={desiredCitiesInput}
              />
              <FieldInput
                label="手机号"
                onChange={(value) => updateProfileField("phone", value)}
                value={draft.profile.phone}
              />
              <FieldInput
                label="邮箱"
                onChange={(value) => updateProfileField("email", value)}
                value={draft.profile.email}
              />
            </div>
          </SectionCard>

          <SectionCard className="space-y-4">
            <h2 className="text-lg font-semibold text-slate-900">个人技能</h2>
            <EditableTagList
              onAdd={(skill) => {
                if (!draft.skills.includes(skill)) {
                  updateSkills([...draft.skills, skill]);
                }
              }}
              onRemove={(skill) =>
                updateSkills(draft.skills.filter((item) => item !== skill))
              }
              skills={draft.skills}
            />
          </SectionCard>

          <EditableListSection
            addLabel="添加教育"
            onAdd={() =>
              setDraft({
                ...draft,
                education: [...draft.education, createEmptyEducation()],
              })
            }
            title="教育经历"
          >
            {draft.education.map((item, index) => (
              <EditableListItem
                actions={
                  <>
                    <Button
                      aria-label={`上移 教育经历 ${index + 1}`}
                      disabled={index === 0}
                      onClick={() =>
                        setDraft({
                          ...draft,
                          education: moveArrayItem(
                            draft.education,
                            index,
                            index - 1,
                          ),
                        })
                      }
                      type="button"
                      variant="ghost"
                    >
                      上移
                    </Button>
                    <Button
                      aria-label={`下移 教育经历 ${index + 1}`}
                      disabled={index === draft.education.length - 1}
                      onClick={() =>
                        setDraft({
                          ...draft,
                          education: moveArrayItem(
                            draft.education,
                            index,
                            index + 1,
                          ),
                        })
                      }
                      type="button"
                      variant="ghost"
                    >
                      下移
                    </Button>
                  </>
                }
                key={index}
                onRemove={() =>
                  setDraft({
                    ...draft,
                    education: draft.education.filter(
                      (_, current) => current !== index,
                    ),
                  })
                }
                title={`教育经历 ${index + 1}`}
              >
                <div className="grid gap-4 md:grid-cols-2">
                  <FieldInput
                    label="学校"
                    onChange={(value) =>
                      updateEducation(index, { ...item, school: value })
                    }
                    value={item.school}
                  />
                  <FieldInput
                    label="专业"
                    onChange={(value) =>
                      updateEducation(index, { ...item, major: value })
                    }
                    value={item.major}
                  />
                  <FieldInput
                    label="学历"
                    onChange={(value) =>
                      updateEducation(index, { ...item, degree: value })
                    }
                    value={item.degree}
                  />
                  <FieldInput
                    label="开始时间"
                    onChange={(value) =>
                      updateEducation(index, { ...item, startDate: value })
                    }
                    value={item.startDate}
                  />
                </div>
              </EditableListItem>
            ))}
          </EditableListSection>

          <EditableListSection
            addLabel="添加奖项"
            onAdd={() =>
              setDraft({
                ...draft,
                awards: [...draft.awards, createEmptyAward()],
              })
            }
            title="荣誉奖项"
          >
            {draft.awards.map((item, index) => (
              <EditableListItem
                key={index}
                onRemove={() =>
                  setDraft({
                    ...draft,
                    awards: draft.awards.filter(
                      (_, current) => current !== index,
                    ),
                  })
                }
                title={`奖项 ${index + 1}`}
              >
                <div className="grid gap-4 md:grid-cols-2">
                  <FieldInput
                    label="日期"
                    onChange={(value) =>
                      updateAward(index, { ...item, date: value })
                    }
                    value={item.date}
                  />
                  <FieldInput
                    label="奖项名称"
                    onChange={(value) =>
                      updateAward(index, { ...item, title: value })
                    }
                    value={item.title}
                  />
                </div>
              </EditableListItem>
            ))}
          </EditableListSection>
        </div>

        <div className="space-y-6">
          <EditableListSection
            addLabel="添加经历"
            onAdd={() =>
              setDraft({
                ...draft,
                workExperience: [
                  ...draft.workExperience,
                  createEmptyWorkExperience(),
                ],
              })
            }
            title="工作经历"
          >
            {draft.workExperience.map((item, index) => (
              <EditableListItem
                actions={
                  <>
                    <Button
                      aria-label={`上移 工作经历 ${index + 1}`}
                      disabled={index === 0}
                      onClick={() =>
                        setDraft({
                          ...draft,
                          workExperience: moveArrayItem(
                            draft.workExperience,
                            index,
                            index - 1,
                          ),
                        })
                      }
                      type="button"
                      variant="ghost"
                    >
                      上移
                    </Button>
                    <Button
                      aria-label={`下移 工作经历 ${index + 1}`}
                      disabled={index === draft.workExperience.length - 1}
                      onClick={() =>
                        setDraft({
                          ...draft,
                          workExperience: moveArrayItem(
                            draft.workExperience,
                            index,
                            index + 1,
                          ),
                        })
                      }
                      type="button"
                      variant="ghost"
                    >
                      下移
                    </Button>
                  </>
                }
                key={index}
                onRemove={() =>
                  setDraft({
                    ...draft,
                    workExperience: draft.workExperience.filter(
                      (_, current) => current !== index,
                    ),
                  })
                }
                title={`工作经历 ${index + 1}`}
              >
                <div className="grid gap-4 md:grid-cols-2">
                  <FieldInput
                    label="公司名称"
                    onChange={(value) =>
                      updateWorkExperience(index, { ...item, company: value })
                    }
                    value={item.company}
                  />
                  <FieldInput
                    label="职位"
                    onChange={(value) =>
                      updateWorkExperience(index, { ...item, role: value })
                    }
                    value={item.role}
                  />
                  <FieldInput
                    label="开始时间"
                    onChange={(value) =>
                      updateWorkExperience(index, { ...item, startDate: value })
                    }
                    value={item.startDate}
                  />
                  <FieldInput
                    label="结束时间"
                    onChange={(value) =>
                      updateWorkExperience(index, { ...item, endDate: value })
                    }
                    value={item.endDate}
                  />
                </div>
                <FieldInput
                  label="工作内容"
                  onChange={(value) =>
                    updateWorkExperience(index, {
                      ...item,
                      description: value.split("\n"),
                    })
                  }
                  textarea
                  value={item.description.join("\n")}
                />
              </EditableListItem>
            ))}
          </EditableListSection>

          <EditableListSection
            addLabel="添加项目"
            onAdd={() =>
              setDraft({
                ...draft,
                projects: [...draft.projects, createEmptyProject()],
              })
            }
            title="项目经历"
          >
            {draft.projects.map((item, index) => (
              <EditableListItem
                actions={
                  <>
                    <Button
                      aria-label={`上移 项目经历 ${index + 1}`}
                      disabled={index === 0}
                      onClick={() =>
                        setDraft({
                          ...draft,
                          projects: moveArrayItem(
                            draft.projects,
                            index,
                            index - 1,
                          ),
                        })
                      }
                      type="button"
                      variant="ghost"
                    >
                      上移
                    </Button>
                    <Button
                      aria-label={`下移 项目经历 ${index + 1}`}
                      disabled={index === draft.projects.length - 1}
                      onClick={() =>
                        setDraft({
                          ...draft,
                          projects: moveArrayItem(
                            draft.projects,
                            index,
                            index + 1,
                          ),
                        })
                      }
                      type="button"
                      variant="ghost"
                    >
                      下移
                    </Button>
                  </>
                }
                key={index}
                onRemove={() =>
                  setDraft({
                    ...draft,
                    projects: draft.projects.filter(
                      (_, current) => current !== index,
                    ),
                  })
                }
                title={`项目经历 ${index + 1}`}
              >
                <div className="grid gap-4 md:grid-cols-2">
                  <FieldInput
                    label="项目名称"
                    onChange={(value) =>
                      updateProject(index, { ...item, name: value })
                    }
                    value={item.name}
                  />
                  <FieldInput
                    label="开始时间"
                    onChange={(value) =>
                      updateProject(index, { ...item, startDate: value })
                    }
                    value={item.startDate}
                  />
                  <FieldInput
                    label="结束时间"
                    onChange={(value) =>
                      updateProject(index, { ...item, endDate: value })
                    }
                    value={item.endDate}
                  />
                </div>
                <FieldInput
                  label="项目描述"
                  onChange={(value) =>
                    updateProject(index, {
                      ...item,
                      description: value.split("\n"),
                    })
                  }
                  textarea
                  value={item.description.join("\n")}
                />
              </EditableListItem>
            ))}
          </EditableListSection>
        </div>
      </div>

      <SaveDraftDialog
        error={draftError}
        onCancel={() => setDraftDialogOpen(false)}
        onConfirm={handleSaveAsDraft}
        open={draftDialogOpen}
        saving={draftSaving}
      />
      <SaveToast
        message={toastMessage ?? undefined}
        visible={toastMessage !== null}
      />
    </AppShell>
  );
}
