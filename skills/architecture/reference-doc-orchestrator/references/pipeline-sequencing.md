# Pipeline Sequencing

Default sequence:

1. requirements discovery
2. glossary stabilization
3. gap analysis
4. PRD drafting
5. domain modeling
6. use-case modeling
7. contract authoring
8. acceptance scenario generation
9. readiness review

## Reroute Rules

### Reroute to Discovery

When:

- the problem statement is still unstable
- actors or workflows are unknown

### Reroute to Glossary

When:

- multiple terms exist for the same concept
- the same term means different things across artifacts

### Reroute to Gap Analysis

When:

- contradictions appear
- one artifact is missing information needed by the next

### Reroute to PRD

When:

- later artifacts expose unclear scope, non-goals, or business capabilities

### Reroute to Domain Model

When:

- use cases or contracts expose weak or inconsistent business semantics

### Reroute to Use-Case Model

When:

- contract authoring reveals unstable commands or queries

### Reroute to Readiness Review

When:

- the full set appears complete and implementation is about to start
