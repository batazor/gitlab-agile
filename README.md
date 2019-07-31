# GitLab Agile

Tool for enjoy GitLab CE

### Environment

| Variable     | Description | Default |
|:-------------|:------------|:-------:|
| GITLAB_URL   | --          |   --    |
| GITLAB_TOKEN | --          |   --    |

### Todo

### 0.1.0

- [ ] Create milestoune
  - [ ] From slack? CI? UI?
    - [ ] Slack `/workflow milestoune add [Name] [from (mm/dd/yy)] [to (mm/dd/yy)]`
    - [x] Slack `/report [type] [milestoune name]` -> response table
- [ ] Statistic Export
  - [ ] Issue
    - [ ] count
    - [ ] weight (parse from title)
  - [ ] Project
    - [ ] count in milestoune
- [ ] Grafana Dashboard
  - [ ] Toolbar: Select milestoune
  - [ ] Chart by issue[count, weight]
  - [ ] SCRUM Burn Down
- [x] Process (Template)
  - [x] Read and apply `process.yaml`
    - [x] labels
    - [ ] boards
- [x] SCRUM report
  - [x] Table: `|Iteration|Plan weight|Actual weight|`
- [ ] Report for PO/TL (import to Google Docs)
  - [ ] `Planning & resources report`
- [ ] Template
  - [ ] Clone `GitLab issue/MR` to each project
