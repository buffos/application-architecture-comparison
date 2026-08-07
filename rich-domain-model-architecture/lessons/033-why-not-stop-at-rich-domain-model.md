# Lesson 033: Why Not Stop At Rich Domain Model?

## Objective

Explain what this Rich Domain Model design handles well and which pressures could justify another architecture.

## Short Answer

You could stop here for many business systems.

This track now has rich aggregates, value objects, domain services, bounded-context translations, application query surfaces, reports, payment review, partial shipment, partial return, and a pricing-policy extension point.

The useful question is not whether Rich Domain Model is universally sufficient. It is which next pressure matters most for the system.

## What Rich Domain Model Is Good At

- keeping business language and invariants close to the objects that enforce them
- making aggregate boundaries and context translations explicit
- separating domain decisions from application coordination and read projections
- growing realistic workflows without turning every rule into a controller or persistence method
- giving policies and reports clear homes outside aggregate state changes

## The Core Limitation

Rich domain objects do not decide how persistence should scale, how read models should be projected, how plugins should be isolated, how rules should be authored, or when a bounded context should become a separate deployable. They also cost more modeling effort and discipline than a simpler architecture.

## Diagram

```mermaid
flowchart LR
    DOMAIN["rich domain model"] --> PRESSURE["next design pressure"]
    PRESSURE --> BOUNDARIES["Hexagonal / Clean / Onion seams"]
    PRESSURE --> MODULAR["Modular Monolith autonomy"]
    PRESSURE --> CQRS["independent read models"]
    PRESSURE --> PLUGINS["Microkernel / plugin isolation"]
    PRESSURE --> RULES["Rules Engine"]
    PRESSURE --> SERVICES["separate services"]
```

## When To Continue Comparing

### Dependency direction matters more

If the domain needs stronger protection from frameworks, databases, and delivery mechanisms, Hexagonal, Clean, or Onion Architecture adds explicit dependency seams.

### Internal module autonomy matters more

If bounded contexts need stricter ownership inside one deployable, a Modular Monolith may add useful enforcement.

### Independent reads matter more

If dashboards need their own storage and projection pipeline, CQRS is a stronger next focus.

### Extensions become the product

If third parties need discovery, isolation, versioning, and lifecycle management, Microkernel or Plugin Architecture goes deeper.

### Rules become configurable

If policies are authored outside code or composed at runtime, a Rules Engine may fit better.

### Deployment independence matters

If a context needs separate scaling or failure isolation, service-oriented decomposition may be justified.

## Why We Are Not Moving On Because Rich Domain Model Failed

Rich Domain Model has done its job here: it made business meaning, invariants, aggregate boundaries, and context translation visible. Other architectures are not automatically better; they optimize for different pressures.

## What To Verify

- lessons `000` through `033` exist
- `go test ./...` passes from this module
- the demo still runs
- every lesson has a matching commit and tag
