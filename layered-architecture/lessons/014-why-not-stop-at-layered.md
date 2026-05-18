# Lesson 014: Why Not Stop At Layered?

## Objective

Explain why the current layered design is useful, where it starts to strain, and what kinds of problems could justify moving to other architectures.

## Short Answer

You absolutely could stop here for some systems.

This layered design already works. It is understandable, testable, and able to express the core workflows of the sample application. For a small team, a moderate codebase, and a single deployable service, layered architecture is often a very reasonable default.

So the right question is not:

"Why is layered bad?"

The right question is:

"What kinds of problems become awkward enough in layered architecture that another architecture might improve them?"

That is the point of the comparison.

## What Layered Is Good At

Before talking about limitations, be precise about the strengths.

### 1. It Is Easy To Explain

A new engineer can usually understand the shape quickly:

- presentation receives requests
- application orchestrates use cases
- domain holds business concepts
- infrastructure stores or integrates

That clarity is a real advantage.

### 2. It Works Well For Small And Medium Systems

When the system is not huge and the team is not large, layered architecture gives enough structure without forcing too much ceremony.

### 3. It Encourages Separation Of Concerns

Even this sample code already benefits from:

- transport logic not living in `main`
- storage not living in handlers
- use-case orchestration not living in repositories

### 4. It Is A Good Teaching Baseline

Layered architecture makes later comparisons meaningful because it gives a familiar, practical starting point.

### 5. It Can Go Surprisingly Far

A lot of real production systems are layered and stay that way for years.

So the lesson is not:

"Layered is wrong."

The lesson is:

"Layered solves many problems cheaply, but not all problems equally well."

## The Core Limitation

The main limitation is not syntax, files, or folders.

The main limitation is this:

As the business becomes richer, the application layer tends to accumulate more and more orchestration, policy coordination, and transaction responsibility.

That means:

- application services become long
- use cases become procedural
- domain objects risk becoming passive data holders plus a few small methods
- cross-cutting concerns accumulate at the service layer

In other words:

Layered architecture often works best at first, then slowly centralizes too much intelligence in the application layer unless the team is disciplined.

## Concrete Limitations In The Current Design

These are not hypothetical. They are already visible in this project.

### 1. Application Services Are Becoming Workflow Scripts

Look at services like:

- [order_service.go](/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/internal/application/order_service.go)
- [fulfillment_service.go](/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/internal/application/fulfillment_service.go)
- [return_service.go](/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/internal/application/return_service.go)

They are doing a lot:

- loading aggregates
- checking preconditions
- iterating through lines
- mutating stock
- mutating orders
- saving multiple repositories

That is normal in layered architecture, but it means the application layer becomes the place where many business workflows live procedurally.

This is fine for moderate complexity.

It becomes harder when:

- workflows branch heavily
- multiple policies interact
- invariants span several concepts
- retries, idempotency, and recovery matter more

### 2. Transaction Boundaries Are Visible But Not Well Encapsulated

You already noticed this yourself in `CreateShipment` and `CancelOrder`.

The service layer updates multiple records step by step. If a failure happens midway:

- some stock may already be changed
- the order may not yet be updated
- the shipment may not yet be saved

So the design reveals the need for a transaction boundary, but the architecture does not solve it by itself.

That leads to pressure for:

- explicit unit-of-work patterns
- infrastructure-level transactions
- compensating actions
- saga-style flow if the system becomes distributed

Layered architecture can support those things, but they are not naturally front-and-center.

### 3. Policy Logic Will Start Spreading

Right now the rules are still simple:

- `CustomBuild` requires approval
- payment must be accepted before shipment
- `Clearance` cannot be returned

But imagine adding all the canonical rules together:

- discount thresholds
- customer-tier pricing
- configurable approval conditions
- plugin-provided adjustments
- stock shortage policies
- shipping exceptions
- return windows

In a layered design, those rules often end up split across:

- application services
- domain methods
- helper services
- plugin registries

That can work, but it can become hard to answer:

"Where does this rule actually live?"

Architectures like DDD, Rules Engine, Hexagonal, or Clean often try to make those rule boundaries more deliberate.

### 4. Domain Objects Are Useful, But Still Not Very Expressive

The domain objects in this layered variant are not anemic, but they are also not especially rich.

For example:

- the `Order` knows some lifecycle transitions
- the `Quote` knows some lifecycle transitions
- but much of the real business workflow still lives outside them

This is not wrong.

It just means the model is not the main center of behavior.

If your goal is to make the domain model the primary place where invariants and decisions live, then a richer domain architecture may help more.

### 5. Read And Write Concerns Are Still Coupled To The Same Repository Shapes

Lesson `012` added a reporting service, which is good, but the reads still directly depend on write-side repositories and in-memory record shapes.

That is fine early on.

It becomes awkward when:

- reports need different storage shapes
- projections need to be denormalized
- dashboards grow independently of command workflows
- query performance matters more than write model purity

Architectures that embrace read models more explicitly can handle that tension better.

### 6. Plugin Support Is Possible, But Not Especially Natural

Lesson `013` proved that a plugin seam can be added to a layered system.

That is good.

But notice what happened:

- the application service had to own the plugin hook
- the plugin contract had to be manually introduced
- the registry had to be manually wired

Again, this works.

But extensibility is not the natural "center of gravity" of layered architecture.

Architectures like Microkernel exist precisely because plugin behavior becomes the first-class design concern rather than an add-on.

### 7. Module Boundaries Are Weak

Layered architecture gives vertical technical layers:

- presentation
- application
- domain
- infrastructure

What it does not strongly enforce by default is business module separation such as:

- quoting
- ordering
- fulfillment
- returns
- reporting

In a small app, that is acceptable.

In a large app, the result is often:

- one large application layer
- one large domain package or set of loosely related domain files
- many services that can call many other services

Modular Monolith and DDD styles exist partly to make business boundaries stronger than plain technical layering.

### 8. Dependency Direction Is Better Than Chaos, But Still Not Maximally Isolated

This project already separates concerns fairly well, but the application layer still depends directly on repository interfaces and concrete domain workflow patterns that are tightly tied to the local codebase shape.

If you want stronger protection against:

- framework leakage
- persistence influence
- transport influence
- accidental infrastructure coupling

then architectures like Hexagonal, Clean, or Onion usually push that dependency direction further and more explicitly.

### 9. It Can Encourage Service-Centric Growth

Layered systems often age like this:

1. start clean
2. add a few services
3. add more business rules
4. add helper services
5. add coordination between services
6. end up with a "service layer blob"

That is one of the classic failure modes.

The problem is not layering itself.

The problem is that layering does not strongly resist that drift unless the team actively fights it.

### 10. Architecture Intent Can Become Ambiguous

As complexity grows, a layered system can become harder to categorize:

- is the business logic supposed to live in domain objects?
- in application services?
- in policy helpers?
- in plugins?
- in repositories?

If the answer is "a little bit everywhere," the architecture is still functioning, but the conceptual clarity starts degrading.

That ambiguity is often the signal that another architecture might give you better rules of the road.

## What We Can Still Do With This Design

It is important not to exaggerate the limitations.

This design can still support:

- HTTP and CLI adapters
- transactions
- plugin hooks
- reporting
- tests
- moderate business complexity
- a single deployable modular codebase

So the limitation is not:

"We cannot implement feature X."

The limitation is more often:

"We can implement feature X, but the design may become harder to evolve, reason about, or keep clean as the system grows."

That distinction matters.

## What Kinds Of Problems Suggest Another Architecture

Here is the practical decision point.

You usually consider another architecture when one or more of these becomes important.

### Rich Domain Invariants Matter More

If you want aggregates, value objects, policies, and explicit domain services to be the center of the model, then DDD or Rich Domain Model becomes more attractive.

### Transport And Infrastructure Isolation Matter More

If you want the core use cases to be strongly isolated from HTTP, DB, CLI, and framework concerns, then Hexagonal, Clean, or Onion becomes more attractive.

### Business Modules Need Stronger Separation

If quoting, ordering, fulfillment, and returns need independent boundaries inside one codebase, then Modular Monolith or explicit component architectures become more attractive.

### Extensions Become A Primary Product Requirement

If pricing, shipping, or approvals must be plugged in or replaced regularly, then Microkernel or plugin-centered designs become more attractive.

### Rules Become Too Numerous Or Too Configurable

If business policy becomes table-driven, externally configured, or frequently changed by non-developers, then Rules Engine approaches become more attractive.

### Read And Write Models Diverge More

If reporting, dashboards, and queries need very different models than command workflows, then CQRS-like patterns or stronger query-side separation becomes more attractive.

## Why We Are Not Changing Architecture Just To Be Fancy

This is the key message for someone studying the repo.

We are not replacing layered architecture because it is obsolete.

We are exploring other architectures because they optimize for different pressures:

- stronger domain modeling
- stronger boundary isolation
- stronger modularity
- stronger extensibility
- stronger rule management
- stronger read/write separation

If those pressures are weak, layered may remain the best tradeoff.

If those pressures become strong, another architecture may reduce long-term pain.

That is the real reason to compare them.

## What Someone Should Learn From This

After reading the layered implementation, the correct conclusion should be:

1. Layered architecture is a solid baseline.
2. It handles a surprising amount of real application complexity.
3. Its biggest risk is gradual concentration of workflow and policy logic in the application layer.
4. Other architectures are useful when specific kinds of complexity become important enough to justify more structure.

That is a much better lesson than:

"Layered is simple, other architectures are advanced."

The real lesson is about tradeoffs, not prestige.

## Summary

This layered design works.

That is exactly why it is a good starting point.

But it also shows clear pressure points:

- transaction consistency across multiple repositories
- growing application-service orchestration
- scattered policy logic
- modestly expressive domain model
- weak business module boundaries
- read/write tension
- extension seams that are possible but not natural

So the reason to explore another architecture is not novelty.

It is to ask:

"Can another design handle these pressures more clearly, more safely, or with less long-term friction?"

That is the right reason to continue the comparison.
