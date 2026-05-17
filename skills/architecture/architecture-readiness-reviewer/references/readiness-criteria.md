# Readiness Criteria

The artifact set is ready for architecture implementation when most of the following are true:

- the product scope is bounded
- the domain vocabulary is stable
- the main workflows are explicit
- important rules and invariants are modeled
- command/query intent is stable
- the external contract is stable enough to compare implementations
- acceptance scenarios exist
- there are no unresolved contradictions that would change behavior

Not-ready signals:

- the same workflow is described differently across documents
- key statuses are unstable
- major failure paths are missing
- architecture-specific assumptions leak into requirements
- the external contract exposes behavior not justified by the product or use-case model
- tests cannot be derived from the artifacts
