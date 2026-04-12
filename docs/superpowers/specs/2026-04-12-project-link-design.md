# Project Link Field Design

## Goal

Add an optional project URL field to project experience items so personal projects on the public resume page can expose a clickable title that opens the related site or code repository.

## Scope

- Add an optional `url` field to each project experience item.
- Show project titles as links on the public resume page only when `url` is present.
- Keep draft preview and print/PDF output unchanged for now.
- Preserve compatibility with existing resume and draft data that do not contain the new field.

## Non-Goals

- No multi-link support per project in this iteration.
- No link rendering changes on draft preview, print preview, or exported PDF.
- No database schema migration that splits project data into separate relational columns.

## Current Architecture

The project stores the full resume as JSON in the `resume.data` column and full draft payloads as JSON in `resume_drafts.data`. Project experience is therefore part of the resume JSON document rather than its own database table.

This means adding `url` is a JSON shape extension, not a SQL table change.

## Recommended Data Model

Extend `ProjectExperience` with:

```ts
interface ProjectExperience {
  name: string
  startDate: string
  endDate: string
  description: string[]
  url?: string
}
```

Why this model:

- It matches the current requirement directly.
- It keeps the editing experience simple.
- It is backward-compatible because the new field is optional.
- It leaves room to evolve into a richer `links` model later if needed.

## UI Behavior

### Edit Page

In the project experience editor, add a new input labeled `项目地址`.

Behavior:

- The field is optional.
- Empty input means no link behavior on the public resume.
- No extra placeholder UI is shown on the public page when absent.

### Public Resume Page

For each project item:

- If `url` exists and is non-empty, render the project title as a link.
- If `url` is absent or empty, render the title exactly as plain text.

Link behavior recommendation:

- open in a new tab
- include `rel="noreferrer"` or equivalent safe external-link handling

### Draft Preview And Print/PDF

Do not change visible behavior in this iteration.

This keeps scope small and avoids widening the requirement before the public-page behavior is validated.

## Validation Rules

Recommended lightweight validation for this iteration:

- treat the field as optional
- trim surrounding whitespace before save
- if provided, store the trimmed string

Avoid strict URL parsing in the first version unless current form validation patterns already support it cleanly. A malformed URL can be improved later without blocking the feature delivery.

## Backend Impact

Although the backend mostly stores raw resume JSON, any backend code that decodes project items into a typed struct should be extended with the optional `url` field so the data is not silently dropped during read/transform flows.

Known area to update:

- PDF export resume decoding structures

Even though PDF output will not render the link yet, the decode struct should still accept the field for compatibility.

## Deployment And Data Migration

This project is already deployed, but this change does **not** require a database table migration.

Why:

- `resume` and `resume_drafts` store JSON blobs
- adding `url` changes the JSON payload shape, not the SQL schema
- existing rows remain valid because `url` is optional

Deployment steps should be:

1. deploy the updated application build
2. restart or rebuild the app container/process
3. have the admin edit projects and fill `项目地址` where needed
4. save the main resume or drafts normally

No manual `ALTER TABLE` or backfill is required.

## Risks

- If a typed backend decode path is missed, the field may be dropped in downstream transforms.
- If link rendering is added without safe target/rel handling, external navigation behavior may be undesirable.
- If validation is too strict too early, legitimate project URLs may be rejected unnecessarily.

## Testing Strategy

Frontend:

- type definition coverage for the new optional field
- edit page test proving `项目地址` is included in saved resume payload
- public resume page test proving linked title renders when `url` exists
- public resume page test proving plain title remains when `url` is absent

Backend:

- confirm resume save/load round-trips the new field
- confirm any typed project decode path, especially PDF export decode, accepts the field without failure

## Acceptance Criteria

- Admin can fill an optional `项目地址` field for each project.
- Saved resume JSON preserves the field.
- Public resume project title is clickable only when a URL is present.
- Existing resumes and drafts without `url` continue to load and render correctly.
- Production deployment requires application rollout only, not SQL schema migration.
