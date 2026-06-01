# Application Architecture Comparison Repository

## Purpose

This repository implements the same business application across multiple architectural styles.

The goal is not to crown one architecture as universally best.

The goal is to make tradeoffs concrete by keeping the problem space stable while changing the design style:

- same domain
- same broad workflow
- same canonical documents
- different internal structure

The repository is meant to help you study:

- how each architecture shapes code
- what kinds of complexity each architecture handles well
- what kinds of friction appear as the system grows

## What The Sample App Covers

The sample business application includes workflows around:

- customers
- products
- inventory
- quotes
- approvals
- orders
- payments
- shipments
- returns
- reporting
- plugins

The canonical reference documents live under [docs](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/docs).

## How To Go Through The Repository

The recommended way to study an architecture is lesson by lesson.

Each lesson:

- adds one small architectural or workflow idea
- has a matching commit message convention
- is tagged in git

That means you can move through an architecture incrementally instead of reading only the final state.

Suggested flow:

1. Start with the lesson markdown in that architecture's `lessons/` folder.
2. Check out the matching tag to see the exact code state for that lesson.
3. Read the code for the new slice only.
4. Move to the next lesson.

Examples:

```powershell
git checkout layered-001
git checkout hexagonal-010
git checkout hexagonal-032
git checkout clean-001
```

To return to the latest working branch state afterward:

```powershell
git checkout -
```

The current tag/commit helper text is maintained in [docs/git-how-to.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/docs/git-how-to.md:1).

## Architectures In Scope

Planned architecture tracks are listed in [docs/architectures.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/docs/architectures.md:1).

Architectures currently implemented in this repository:

- Layered Architecture
- Hexagonal Architecture / Ports and Adapters
- Clean Architecture
- Onion Architecture

## Lesson Index

### Layered Architecture

#### Lessons

- `001` Ports-and-adapters skeleton baseline: [001-layered-skeleton.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/001-layered-skeleton.md:1)
- `002` Application service read flow: [002-application-service-read-flow.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/002-application-service-read-flow.md:1)
- `003` Domain state transition: [003-domain-state-transition.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/003-domain-state-transition.md:1)
- `004` Presentation layer: [004-presentation-layer.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/004-presentation-layer.md:1)
- `005` HTTP presentation adapter: [005-http-presentation-adapter.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/005-http-presentation-adapter.md:1)
- `006` Canonical quote inputs: [006-canonical-quote-inputs.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/006-canonical-quote-inputs.md:1)
- `007` Order conversion and reservation: [007-order-conversion-and-reservation.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/007-order-conversion-and-reservation.md:1)
- `008` Quote approval boundary: [008-quote-approval-boundary.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/008-quote-approval-boundary.md:1)
- `009` Payment and shipment gate: [009-payment-and-shipment-gate.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/009-payment-and-shipment-gate.md:1)
- `010` Cancellation and reservation release: [010-cancellation-and-reservation-release.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/010-cancellation-and-reservation-release.md:1)
- `011` Returns and restocking: [011-returns-and-restocking.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/011-returns-and-restocking.md:1)
- `012` Reporting query service: [012-reporting-query-service.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/012-reporting-query-service.md:1)
- `013` Pricing plugin extension: [013-pricing-plugin-extension.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/013-pricing-plugin-extension.md:1)
- `014` Why not stop at layered?: [014-why-not-stop-at-layered.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/layered-architecture/lessons/014-why-not-stop-at-layered.md:1)

### Hexagonal Architecture / Ports And Adapters

#### Lessons

- `001` Ports and adapters skeleton: [001-ports-and-adapters-skeleton.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/001-ports-and-adapters-skeleton.md:1)
- `002` Second inbound adapter: [002-second-inbound-adapter.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/002-second-inbound-adapter.md:1)
- `003` Second outbound port operation: [003-second-outbound-port-operation.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/003-second-outbound-port-operation.md:1)
- `004` Second outbound port customer lookup: [004-second-outbound-port-customer-lookup.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/004-second-outbound-port-customer-lookup.md:1)
- `005` Add quote line with multiple ports: [005-add-quote-line-with-multiple-ports.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/005-add-quote-line-with-multiple-ports.md:1)
- `006` Submission and approval policy port: [006-submission-and-approval-policy-port.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/006-submission-and-approval-policy-port.md:1)
- `007` Quote to order with reservation port: [007-quote-to-order-with-reservation-port.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/007-quote-to-order-with-reservation-port.md:1)
- `008` Payment and shipment ports: [008-payment-and-shipment-ports.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/008-payment-and-shipment-ports.md:1)
- `009` Order cancellation and reservation release: [009-order-cancellation-and-reservation-release.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/009-order-cancellation-and-reservation-release.md:1)
- `010` Return request and refund port: [010-return-request-and-refund-port.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/010-return-request-and-refund-port.md:1)
- `011` Return restocking port: [011-return-restocking-port.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/011-return-restocking-port.md:1)
- `012` Return review boundary: [012-return-review-boundary.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/012-return-review-boundary.md:1)
- `013` Return eligibility policy port: [013-return-eligibility-policy-port.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/013-return-eligibility-policy-port.md:1)
- `014` Real return window policy: [014-real-return-window-policy.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/014-real-return-window-policy.md:1)
- `015` Return actor metadata: [015-return-actor-metadata.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/015-return-actor-metadata.md:1)
- `016` Return command idempotency: [016-return-command-idempotency.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/016-return-command-idempotency.md:1)
- `017` Return query surface: [017-return-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/017-return-query-surface.md:1)
- `018` Order query surface: [018-order-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/018-order-query-surface.md:1)
- `019` Shipment query surface: [019-shipment-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/019-shipment-query-surface.md:1)
- `020` Quote list query surface: [020-quote-list-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/020-quote-list-query-surface.md:1)
- `021` Product query surface: [021-product-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/021-product-query-surface.md:1)
- `022` Customer query surface: [022-customer-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/022-customer-query-surface.md:1)
- `023` Quote conversion report: [023-quote-conversion-report.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/023-quote-conversion-report.md:1)
- `024` Return rate by category report: [024-return-rate-by-category-report.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/024-return-rate-by-category-report.md:1)
- `025` Top discounted products report: [025-top-discounted-products-report.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/025-top-discounted-products-report.md:1)
- `026` Low stock items report: [026-low-stock-items-report.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/026-low-stock-items-report.md:1)
- `027` Orders awaiting approval report: [027-orders-awaiting-approval-report.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/027-orders-awaiting-approval-report.md:1)
- `028` Inventory write model: [028-inventory-write-model.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/028-inventory-write-model.md:1)
- `029` Payment review workflow: [029-payment-review-workflow.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/029-payment-review-workflow.md:1)
- `030` Partial shipment support: [030-partial-shipment-support.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/030-partial-shipment-support.md:1)
- `031` Partial returns by line: [031-partial-returns-by-line.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/031-partial-returns-by-line.md:1)
- `032` Plugin pricing extension point: [032-plugin-pricing-extension-point.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/032-plugin-pricing-extension-point.md:1)
- `033` Why not stop at hexagonal?: [033-why-not-stop-at-hexagonal.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/hexagonal-architecture/lessons/033-why-not-stop-at-hexagonal.md:1)

### Clean Architecture

#### Lessons

- `000` From hexagonal to clean: [000-from-hexagonal-to-clean.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/000-from-hexagonal-to-clean.md:1)
- `001` Clean architecture skeleton: [001-clean-architecture-skeleton.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/001-clean-architecture-skeleton.md:1)
- `002` Query use case and presenter: [002-query-use-case-and-presenter.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/002-query-use-case-and-presenter.md:1)
- `003` Add quote line with gateways: [003-add-quote-line-with-gateways.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/003-add-quote-line-with-gateways.md:1)
- `004` Submit quote state transition: [004-submit-quote-state-transition.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/004-submit-quote-state-transition.md:1)
- `005` Approval policy boundary: [005-approval-policy-boundary.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/005-approval-policy-boundary.md:1)
- `006` Approve pending quote: [006-approve-pending-quote.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/006-approve-pending-quote.md:1)
- `007` Convert quote to order: [007-convert-quote-to-order.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/007-convert-quote-to-order.md:1)
- `008` Order conversion with reservation: [008-order-conversion-with-reservation.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/008-order-conversion-with-reservation.md:1)
- `009` Payment gateway and order capture: [009-payment-gateway-and-order-capture.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/009-payment-gateway-and-order-capture.md:1)
- `010` Shipment creation after payment: [010-shipment-creation-after-payment.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/010-shipment-creation-after-payment.md:1)
- `011` Order cancellation and release: [011-order-cancellation-and-release.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/011-order-cancellation-and-release.md:1)
- `012` Return request and refund boundary: [012-return-request-and-refund-boundary.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/012-return-request-and-refund-boundary.md:1)
- `013` Return restocking boundary: [013-return-restocking-boundary.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/013-return-restocking-boundary.md:1)
- `014` Return review boundary: [014-return-review-boundary.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/014-return-review-boundary.md:1)
- `015` Return eligibility policy: [015-return-eligibility-policy.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/015-return-eligibility-policy.md:1)
- `016` Real return window policy: [016-real-return-window-policy.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/016-real-return-window-policy.md:1)
- `017` Return actor metadata: [017-return-actor-metadata.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/017-return-actor-metadata.md:1)
- `018` Return command idempotency: [018-return-command-idempotency.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/018-return-command-idempotency.md:1)
- `019` Return query surface: [019-return-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/019-return-query-surface.md:1)
- `020` Order query surface: [020-order-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/020-order-query-surface.md:1)
- `021` Shipment query surface: [021-shipment-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/021-shipment-query-surface.md:1)
- `022` Quote list query surface: [022-quote-list-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/022-quote-list-query-surface.md:1)
- `023` Product query surface: [023-product-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/023-product-query-surface.md:1)
- `024` Customer query surface: [024-customer-query-surface.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/024-customer-query-surface.md:1)
- `025` Quote conversion report: [025-quote-conversion-report.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/025-quote-conversion-report.md:1)
- `026` Return rate by category report: [026-return-rate-by-category-report.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/026-return-rate-by-category-report.md:1)
- `027` Low stock items report: [027-low-stock-items-report.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/027-low-stock-items-report.md:1)
- `028` Orders awaiting approval report: [028-orders-awaiting-approval-report.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/028-orders-awaiting-approval-report.md:1)
- `029` Payment review workflow: [029-payment-review-workflow.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/029-payment-review-workflow.md:1)
- `030` Partial shipment support: [030-partial-shipment-support.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/030-partial-shipment-support.md:1)
- `031` Partial returns by line: [031-partial-returns-by-line.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/031-partial-returns-by-line.md:1)
- `032` Plugin pricing extension point: [032-plugin-pricing-extension-point.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/032-plugin-pricing-extension-point.md:1)
- `033` Why not stop at clean?: [033-why-not-stop-at-clean.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/clean-architecture/lessons/033-why-not-stop-at-clean.md:1)

### Onion Architecture

#### Lessons

- `000` From clean to onion: [000-from-clean-to-onion.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/onion-architecture/lessons/000-from-clean-to-onion.md:1)
- `001` Onion architecture skeleton: [001-onion-architecture-skeleton.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/onion-architecture/lessons/001-onion-architecture-skeleton.md:1)
- `002` Query application service: [002-query-application-service.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/onion-architecture/lessons/002-query-application-service.md:1)
- `003` Add quote line with product lookup: [003-add-quote-line-with-product-lookup.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/onion-architecture/lessons/003-add-quote-line-with-product-lookup.md:1)
- `004` Submit quote state transition: [004-submit-quote-state-transition.md](/abs/path/c:/Users/buffo/Code/architecture/01.application.architectures/onion-architecture/lessons/004-submit-quote-state-transition.md:1)

## How To Maintain This File

As new architectures and lessons are added:

- add the architecture to the implemented list
- append the new lesson to the relevant lesson section
- keep the repository-walkthrough instructions current if the git workflow changes
