# Denote Task Format Specification

Version: 2.1.0
Date: 2026-02-17

## Version History

- 2.1.0 (2026-02-17): Added task recurrence support
  - Added `recur` field to task frontmatter
  - Supported patterns: daily, weekly, monthly, yearly, every Nd/Nw/Nm/Ny, every mon,wed,fri
  - Completing a recurring task automatically creates the next occurrence
- 2.0.1 (2025-01-14): Unified sequential ID field
  - Changed `task_id` and `project_id` to universal `index_id` field
  - Counter file now uses `next_index_id` instead of separate counters
- 2.0.0 (2025-01-13): Breaking changes - Project associations now use Denote identifiers
  - Changed `project:` field to `project_id:` using Denote timestamp
  - Added required `title:` field to frontmatter spec
  - Added `type:` field for explicit file type declaration
  - Clarified ID reference semantics
  - Added semantic versioning
- 1.0.0 (2025-07-04): Initial specification based on notes-cli implementation

## Breaking Changes in 2.0.0

Tasks now reference projects using Denote identifiers instead of string names:
- Old: `project: planning-for-lyon`
- New: `project_id: 20250627T191225`

## Breaking Changes in 2.0.1

Unified sequential ID field for both tasks and projects:
- Old: `task_id: 25` (for tasks), `project_id: 15` (for projects)
- New: `index_id: 25` (universal for both)

## Overview

This specification describes the file format and metadata structure for tasks and projects using the Denote naming convention with extended YAML frontmatter. This format is designed to be human-readable, sync-friendly, and easily parseable by multiple tools.

All files follow the core Denote principle: every file has a unique, immutable timestamp identifier.

## File Naming Convention

All task and project files follow the Denote naming convention:

```
YYYYMMDDTHHMMSS--title-slug__tag1_tag2_tag3.md
```

### Components:
- **Denote ID**: `YYYYMMDDTHHMMSS` - Timestamp of creation (e.g., `20250704T151739`)
- **Separator**: `--` - Double dash separating ID from title
- **Title**: `title-slug` - Kebab-case title (spaces become hyphens, lowercase)
- **Tag Separator**: `__` - Double underscore before tags
- **Tags**: `tag1_tag2_tag3` - Underscore-separated tags
- **Extension**: `.md` - Markdown file

### Required Tags:
- Tasks MUST include the `task` tag
- Projects MUST include the `project` tag

### Examples:
```
20250704T151739--get-a-new-front-ring-for-the-bike__task_bike_personal.md
20250624T234037--on-call-in-effect__task_itleads_active_project.md
20250627T191225--planning-for-lyon__project_travel.md
```

## YAML Frontmatter Structure

All task and project files begin with YAML frontmatter delimited by `---`:

### Task Frontmatter

```yaml
---
title: Get a new front ring for the bike  # Human-readable title
index_id: 25             # Unique sequential ID for CLI convenience
type: task               # Optional - determined by __task in filename
status: open             # Task status (see Status Values)
priority: p2             # Priority level (p1, p2, p3)
due_date: 2025-07-16     # Due date in YYYY-MM-DD format
start_date: 2025-07-01   # Start date in YYYY-MM-DD format
estimate: 5              # Time estimate (Fibonacci: 1,2,3,5,8,13)
project_id: 20250627T191225  # Denote ID of associated project
area: work               # Area of life (work, personal, home, etc.)
assignee: john-doe       # Person responsible
recur: weekly            # Recurrence pattern (see Recurrence Patterns)
tags: [bike, maintenance]  # Additional tags beyond filename tags
---
```

### Project Frontmatter

```yaml
---
title: Planning for Lyon  # Human-readable title
index_id: 15             # Unique sequential ID for CLI convenience
type: project            # Optional - determined by __project in filename
status: active           # Project status (see Status Values)
priority: p1             # Priority level (p1, p2, p3)
due_date: 2025-12-31     # Project due date
start_date: 2025-01-01   # Project start date
area: work               # Area of life
tags: [travel, conference]  # Additional tags beyond filename tags
---
```

## Field Specifications

### Required Fields

#### title
- Type: String
- Required: Yes
- Description: Human-readable title for display
- Note: May differ from filename slug for readability

#### type
- Type: String (enum)
- Required: No
- Values: `task`, `project`
- Description: Optional explicit declaration of file type
- Note: File type is determined by `__task` or `__project` tags in filename

### Common Fields

#### index_id
- Type: Integer
- Required: Yes
- Description: Universal sequential ID for CLI convenience (e.g., `task 25` or `project 15`)
- Note: The Denote ID (timestamp) is the canonical identifier

#### status
- Type: String (enum)
- Required: No (default: "open" for tasks, "active" for projects)
- Task values: `open`, `done`, `paused`, `delegated`, `dropped`
- Project values: `active`, `completed`, `paused`, `cancelled`

#### priority
- Type: String (enum)
- Required: No
- Values: `p1` (highest), `p2` (medium), `p3` (low)
- Display: Often shown as [P1], [P2], [P3]

#### due_date / start_date
- Type: String (date)
- Required: No
- Format: `YYYY-MM-DD`
- Example: `2025-07-16`

#### area
- Type: String
- Required: No
- Description: Life area or context
- Common values: `work`, `personal`, `home`, `health`, `finance`

#### tags
- Type: Array of strings
- Required: No
- Description: Additional tags beyond those in filename
- Example: `[urgent, waiting-for-review]`

### Task-Specific Fields

#### estimate
- Type: Integer
- Required: No
- Values: Fibonacci sequence (1, 2, 3, 5, 8, 13)
- Description: Time/effort estimate

#### project_id
- Type: String (Denote ID)
- Required: No
- Format: `YYYYMMDDTHHMMSS`
- Description: Denote ID of associated project file
- Example: `20250627T191225`
- Note: Use YAML comments for human context if needed

#### assignee
- Type: String
- Required: No
- Description: Person responsible for the task

#### recur
- Type: String
- Required: No
- Description: Recurrence pattern for repeating tasks
- Requires: `due_date` must be set when `recur` is set
- Behavior: When a recurring task is marked done, a new task is automatically created with the next due date calculated from the original due date (fixed schedule)
- Values:
  - Simple: `daily`, `weekly`, `monthly`, `yearly`
  - Interval: `every <N>d`, `every <N>w`, `every <N>m`, `every <N>y` (e.g., `every 2w` for biweekly)
  - Day-of-week: `every monday`, `every mon,wed,fri`
- Note: If the computed next date would be in the past (late completion), it advances until the next future occurrence

## Content Structure

After the YAML frontmatter, the file contains Markdown content:

```markdown
---
title: Fix slow database queries
task_id: 35
type: task
status: open
project_id: 20250627T191225  # Website Optimization
---

Main task description goes here.

## Notes
Additional notes and details.

## Log Entries
[2025-07-04] Investigated initial approach
[2025-07-05] Waiting for feedback from team
```

### Log Entry Format
- Format: `[YYYY-MM-DD] Entry text`
- Location: Added chronologically after frontmatter
- Purpose: Timestamped progress updates

## File Organization

### Directory Structure
```
notes/
├── tasks/           # Task files (can be same as notes dir)
├── projects/        # Project files (can be same as notes dir)
└── .denote-task-counter.json  # Sequential ID counter
```

### Counter File

#### .denote-task-counter.json
Tracks next available sequential ID:
```json
{
  "next_index_id": 73,
  "spec_version": "2.0.1"
}
```

Note: Previously named `.notes-cli-id-counter.json` in v1.0.0

## ID Reference Guidelines

### Primary Identifier
The Denote ID (timestamp) is the canonical, immutable identifier for all files.

### Reference Methods
Tasks/projects can be referenced by:
1. **Denote ID** (canonical): `20250704T151739`
2. **Sequential ID** (CLI convenience): `task 35` or `project 15`
3. **Partial Denote ID**: `0704T1517` (must be unambiguous)
4. **Title** (human search): "Fix kitchen sink"

### Association Example
```yaml
# Task file: 20250705T093000--implement-search__task.md
project_id: 20250627T191225  # Lyon Planning project

# Project file: 20250627T191225--lyon-planning__project.md
title: Lyon Planning
```

## Tool Implementation Guidelines

### Project Association
1. Store project's Denote ID in task's `project_id` field
2. Never store project names/titles for association
3. Use YAML comments for human readability when editing
4. Build indexes/caches using Denote IDs as keys

### UI Considerations
1. Always display human-readable titles, never Denote IDs
2. Provide project selection interfaces (don't require ID entry)
3. Show project context when viewing tasks
4. Cache project metadata for performance

### Migration from v1.0.0
Tasks using string-based project associations need migration:
```yaml
# Old (v1.0.0)
project: planning-for-lyon

# New (v2.0.0)
project_id: 20250627T191225  # Planning for Lyon
```

### Migration from v2.0.0 to v2.0.1
Unified sequential ID field:
```yaml
# Old (v2.0.0)
task_id: 25      # For tasks
project_id: 15   # For projects

# New (v2.0.1)
index_id: 25     # Universal for both tasks and projects
```

## Sync Considerations

1. **Counter File**: Store as `.denote-task-counter.json` in task directory
2. **Conflict Resolution**: If counter missing, scan for highest index_id
3. **Denote IDs**: Include microseconds to minimize collision risk
4. **Project References**: Denote IDs are stable across renames/moves

## Best Practices

1. Let tools manage ID assignment (both Denote and sequential)
2. Use YAML comments to add human context to IDs
3. Keep titles concise but descriptive
4. Use consistent tag vocabulary
5. Store related tasks and projects in same directory tree
6. Add log entries for significant progress
7. Cache project lookups for performance in UIs

## Example Files

### Task File
```
Filename: 20250704T151739--fix-homepage-layout__task_website.md
```

```yaml
---
title: Fix homepage layout
index_id: 50
type: task
status: open
priority: p2
due_date: 2025-07-10
area: work
project_id: 20250615T120000  # Website Redesign
estimate: 5
---

The homepage layout breaks on mobile devices. Need to fix responsive CSS.

## Checklist
- [ ] Identify breakpoint issues
- [ ] Fix navigation menu
- [ ] Test on multiple devices
- [ ] Update browser compatibility

[2025-07-04] Reported by customer
[2025-07-05] Reproduced on iPhone 12
```

### Project File
```
Filename: 20250615T120000--website-redesign__project_work.md
```

```yaml
---
title: Website Redesign
index_id: 8
type: project
status: active
priority: p1
start_date: 2025-06-15
due_date: 2025-09-01
area: work
tags: [q3-goals, customer-facing]
---

Complete redesign of company website with new branding and improved UX.

## Objectives
- Modernize visual design
- Improve mobile experience  
- Increase conversion rate
- Enhance accessibility

## Milestones
- [x] Design mockups approved
- [ ] Frontend implementation
- [ ] Backend integration
- [ ] User testing
- [ ] Launch

[2025-06-15] Project kickoff meeting
[2025-06-20] Selected design agency
[2025-07-01] Mockups delivered
```

## Semantic Versioning

This specification follows [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** version for incompatible changes (e.g., 2.0.0)
- **MINOR** version for backwards-compatible features (e.g., 2.1.0)
- **PATCH** version for backwards-compatible fixes (e.g., 2.0.1)

Tools implementing this spec should:
1. Declare which spec version they support
2. Check `spec_version` in counter file if present
3. Provide migration tools for major version changes
4. Accept files from same major version