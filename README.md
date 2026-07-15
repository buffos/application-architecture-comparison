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

The canonical reference documents live under [docs](docs/).

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

The current tag/commit helper text is maintained in [docs/git-how-to.md](docs/git-how-to.md).

## Architectures In Scope

Planned architecture tracks are listed in [docs/architectures.md](docs/architectures.md).

Architectures currently implemented in this repository:

- Layered Architecture
- Hexagonal Architecture / Ports and Adapters
- Clean Architecture
- Onion Architecture
- Modular Monolith
- Microkernel / Plugin Architecture

## Lesson Index

### Layered Architecture

#### Lessons

- `001` Ports-and-adapters skeleton baseline: [001-layered-skeleton.md](layered-architecture/lessons/001-layered-skeleton.md)
- `002` Application service read flow: [002-application-service-read-flow.md](layered-architecture/lessons/002-application-service-read-flow.md)
- `003` Domain state transition: [003-domain-state-transition.md](layered-architecture/lessons/003-domain-state-transition.md)
- `004` Presentation layer: [004-presentation-layer.md](layered-architecture/lessons/004-presentation-layer.md)
- `005` HTTP presentation adapter: [005-http-presentation-adapter.md](layered-architecture/lessons/005-http-presentation-adapter.md)
- `006` Canonical quote inputs: [006-canonical-quote-inputs.md](layered-architecture/lessons/006-canonical-quote-inputs.md)
- `007` Order conversion and reservation: [007-order-conversion-and-reservation.md](layered-architecture/lessons/007-order-conversion-and-reservation.md)
- `008` Quote approval boundary: [008-quote-approval-boundary.md](layered-architecture/lessons/008-quote-approval-boundary.md)
- `009` Payment and shipment gate: [009-payment-and-shipment-gate.md](layered-architecture/lessons/009-payment-and-shipment-gate.md)
- `010` Cancellation and reservation release: [010-cancellation-and-reservation-release.md](layered-architecture/lessons/010-cancellation-and-reservation-release.md)
- `011` Returns and restocking: [011-returns-and-restocking.md](layered-architecture/lessons/011-returns-and-restocking.md)
- `012` Reporting query service: [012-reporting-query-service.md](layered-architecture/lessons/012-reporting-query-service.md)
- `013` Pricing plugin extension: [013-pricing-plugin-extension.md](layered-architecture/lessons/013-pricing-plugin-extension.md)
- `014` Why not stop at layered?: [014-why-not-stop-at-layered.md](layered-architecture/lessons/014-why-not-stop-at-layered.md)

### Hexagonal Architecture / Ports And Adapters

#### Lessons

- `001` Ports and adapters skeleton: [001-ports-and-adapters-skeleton.md](hexagonal-architecture/lessons/001-ports-and-adapters-skeleton.md)
- `002` Second inbound adapter: [002-second-inbound-adapter.md](hexagonal-architecture/lessons/002-second-inbound-adapter.md)
- `003` Second outbound port operation: [003-second-outbound-port-operation.md](hexagonal-architecture/lessons/003-second-outbound-port-operation.md)
- `004` Second outbound port customer lookup: [004-second-outbound-port-customer-lookup.md](hexagonal-architecture/lessons/004-second-outbound-port-customer-lookup.md)
- `005` Add quote line with multiple ports: [005-add-quote-line-with-multiple-ports.md](hexagonal-architecture/lessons/005-add-quote-line-with-multiple-ports.md)
- `006` Submission and approval policy port: [006-submission-and-approval-policy-port.md](hexagonal-architecture/lessons/006-submission-and-approval-policy-port.md)
- `007` Quote to order with reservation port: [007-quote-to-order-with-reservation-port.md](hexagonal-architecture/lessons/007-quote-to-order-with-reservation-port.md)
- `008` Payment and shipment ports: [008-payment-and-shipment-ports.md](hexagonal-architecture/lessons/008-payment-and-shipment-ports.md)
- `009` Order cancellation and reservation release: [009-order-cancellation-and-reservation-release.md](hexagonal-architecture/lessons/009-order-cancellation-and-reservation-release.md)
- `010` Return request and refund port: [010-return-request-and-refund-port.md](hexagonal-architecture/lessons/010-return-request-and-refund-port.md)
- `011` Return restocking port: [011-return-restocking-port.md](hexagonal-architecture/lessons/011-return-restocking-port.md)
- `012` Return review boundary: [012-return-review-boundary.md](hexagonal-architecture/lessons/012-return-review-boundary.md)
- `013` Return eligibility policy port: [013-return-eligibility-policy-port.md](hexagonal-architecture/lessons/013-return-eligibility-policy-port.md)
- `014` Real return window policy: [014-real-return-window-policy.md](hexagonal-architecture/lessons/014-real-return-window-policy.md)
- `015` Return actor metadata: [015-return-actor-metadata.md](hexagonal-architecture/lessons/015-return-actor-metadata.md)
- `016` Return command idempotency: [016-return-command-idempotency.md](hexagonal-architecture/lessons/016-return-command-idempotency.md)
- `017` Return query surface: [017-return-query-surface.md](hexagonal-architecture/lessons/017-return-query-surface.md)
- `018` Order query surface: [018-order-query-surface.md](hexagonal-architecture/lessons/018-order-query-surface.md)
- `019` Shipment query surface: [019-shipment-query-surface.md](hexagonal-architecture/lessons/019-shipment-query-surface.md)
- `020` Quote list query surface: [020-quote-list-query-surface.md](hexagonal-architecture/lessons/020-quote-list-query-surface.md)
- `021` Product query surface: [021-product-query-surface.md](hexagonal-architecture/lessons/021-product-query-surface.md)
- `022` Customer query surface: [022-customer-query-surface.md](hexagonal-architecture/lessons/022-customer-query-surface.md)
- `023` Quote conversion report: [023-quote-conversion-report.md](hexagonal-architecture/lessons/023-quote-conversion-report.md)
- `024` Return rate by category report: [024-return-rate-by-category-report.md](hexagonal-architecture/lessons/024-return-rate-by-category-report.md)
- `025` Top discounted products report: [025-top-discounted-products-report.md](hexagonal-architecture/lessons/025-top-discounted-products-report.md)
- `026` Low stock items report: [026-low-stock-items-report.md](hexagonal-architecture/lessons/026-low-stock-items-report.md)
- `027` Orders awaiting approval report: [027-orders-awaiting-approval-report.md](hexagonal-architecture/lessons/027-orders-awaiting-approval-report.md)
- `028` Inventory write model: [028-inventory-write-model.md](hexagonal-architecture/lessons/028-inventory-write-model.md)
- `029` Payment review workflow: [029-payment-review-workflow.md](hexagonal-architecture/lessons/029-payment-review-workflow.md)
- `030` Partial shipment support: [030-partial-shipment-support.md](hexagonal-architecture/lessons/030-partial-shipment-support.md)
- `031` Partial returns by line: [031-partial-returns-by-line.md](hexagonal-architecture/lessons/031-partial-returns-by-line.md)
- `032` Plugin pricing extension point: [032-plugin-pricing-extension-point.md](hexagonal-architecture/lessons/032-plugin-pricing-extension-point.md)
- `033` Why not stop at hexagonal?: [033-why-not-stop-at-hexagonal.md](hexagonal-architecture/lessons/033-why-not-stop-at-hexagonal.md)

### Clean Architecture

#### Lessons

- `000` From hexagonal to clean: [000-from-hexagonal-to-clean.md](clean-architecture/lessons/000-from-hexagonal-to-clean.md)
- `001` Clean architecture skeleton: [001-clean-architecture-skeleton.md](clean-architecture/lessons/001-clean-architecture-skeleton.md)
- `002` Query use case and presenter: [002-query-use-case-and-presenter.md](clean-architecture/lessons/002-query-use-case-and-presenter.md)
- `003` Add quote line with gateways: [003-add-quote-line-with-gateways.md](clean-architecture/lessons/003-add-quote-line-with-gateways.md)
- `004` Submit quote state transition: [004-submit-quote-state-transition.md](clean-architecture/lessons/004-submit-quote-state-transition.md)
- `005` Approval policy boundary: [005-approval-policy-boundary.md](clean-architecture/lessons/005-approval-policy-boundary.md)
- `006` Approve pending quote: [006-approve-pending-quote.md](clean-architecture/lessons/006-approve-pending-quote.md)
- `007` Convert quote to order: [007-convert-quote-to-order.md](clean-architecture/lessons/007-convert-quote-to-order.md)
- `008` Order conversion with reservation: [008-order-conversion-with-reservation.md](clean-architecture/lessons/008-order-conversion-with-reservation.md)
- `009` Payment gateway and order capture: [009-payment-gateway-and-order-capture.md](clean-architecture/lessons/009-payment-gateway-and-order-capture.md)
- `010` Shipment creation after payment: [010-shipment-creation-after-payment.md](clean-architecture/lessons/010-shipment-creation-after-payment.md)
- `011` Order cancellation and release: [011-order-cancellation-and-release.md](clean-architecture/lessons/011-order-cancellation-and-release.md)
- `012` Return request and refund boundary: [012-return-request-and-refund-boundary.md](clean-architecture/lessons/012-return-request-and-refund-boundary.md)
- `013` Return restocking boundary: [013-return-restocking-boundary.md](clean-architecture/lessons/013-return-restocking-boundary.md)
- `014` Return review boundary: [014-return-review-boundary.md](clean-architecture/lessons/014-return-review-boundary.md)
- `015` Return eligibility policy: [015-return-eligibility-policy.md](clean-architecture/lessons/015-return-eligibility-policy.md)
- `016` Real return window policy: [016-real-return-window-policy.md](clean-architecture/lessons/016-real-return-window-policy.md)
- `017` Return actor metadata: [017-return-actor-metadata.md](clean-architecture/lessons/017-return-actor-metadata.md)
- `018` Return command idempotency: [018-return-command-idempotency.md](clean-architecture/lessons/018-return-command-idempotency.md)
- `019` Return query surface: [019-return-query-surface.md](clean-architecture/lessons/019-return-query-surface.md)
- `020` Order query surface: [020-order-query-surface.md](clean-architecture/lessons/020-order-query-surface.md)
- `021` Shipment query surface: [021-shipment-query-surface.md](clean-architecture/lessons/021-shipment-query-surface.md)
- `022` Quote list query surface: [022-quote-list-query-surface.md](clean-architecture/lessons/022-quote-list-query-surface.md)
- `023` Product query surface: [023-product-query-surface.md](clean-architecture/lessons/023-product-query-surface.md)
- `024` Customer query surface: [024-customer-query-surface.md](clean-architecture/lessons/024-customer-query-surface.md)
- `025` Quote conversion report: [025-quote-conversion-report.md](clean-architecture/lessons/025-quote-conversion-report.md)
- `026` Return rate by category report: [026-return-rate-by-category-report.md](clean-architecture/lessons/026-return-rate-by-category-report.md)
- `027` Low stock items report: [027-low-stock-items-report.md](clean-architecture/lessons/027-low-stock-items-report.md)
- `028` Orders awaiting approval report: [028-orders-awaiting-approval-report.md](clean-architecture/lessons/028-orders-awaiting-approval-report.md)
- `029` Payment review workflow: [029-payment-review-workflow.md](clean-architecture/lessons/029-payment-review-workflow.md)
- `030` Partial shipment support: [030-partial-shipment-support.md](clean-architecture/lessons/030-partial-shipment-support.md)
- `031` Partial returns by line: [031-partial-returns-by-line.md](clean-architecture/lessons/031-partial-returns-by-line.md)
- `032` Plugin pricing extension point: [032-plugin-pricing-extension-point.md](clean-architecture/lessons/032-plugin-pricing-extension-point.md)
- `033` Why not stop at clean?: [033-why-not-stop-at-clean.md](clean-architecture/lessons/033-why-not-stop-at-clean.md)

### Onion Architecture

#### Lessons

- `000` From clean to onion: [000-from-clean-to-onion.md](onion-architecture/lessons/000-from-clean-to-onion.md)
- `001` Onion architecture skeleton: [001-onion-architecture-skeleton.md](onion-architecture/lessons/001-onion-architecture-skeleton.md)
- `002` Query application service: [002-query-application-service.md](onion-architecture/lessons/002-query-application-service.md)
- `003` Add quote line with product lookup: [003-add-quote-line-with-product-lookup.md](onion-architecture/lessons/003-add-quote-line-with-product-lookup.md)
- `004` Submit quote state transition: [004-submit-quote-state-transition.md](onion-architecture/lessons/004-submit-quote-state-transition.md)
- `005` Approval policy boundary: [005-approval-policy-boundary.md](onion-architecture/lessons/005-approval-policy-boundary.md)
- `006` Approve pending quote: [006-approve-pending-quote.md](onion-architecture/lessons/006-approve-pending-quote.md)
- `007` Convert quote to order: [007-convert-quote-to-order.md](onion-architecture/lessons/007-convert-quote-to-order.md)
- `008` Order conversion with reservation: [008-order-conversion-with-reservation.md](onion-architecture/lessons/008-order-conversion-with-reservation.md)
- `009` Payment gateway and order capture: [009-payment-gateway-and-order-capture.md](onion-architecture/lessons/009-payment-gateway-and-order-capture.md)
- `010` Shipment creation after payment: [010-shipment-creation-after-payment.md](onion-architecture/lessons/010-shipment-creation-after-payment.md)
- `011` Order cancellation and release: [011-order-cancellation-and-release.md](onion-architecture/lessons/011-order-cancellation-and-release.md)
- `012` Return request and refund boundary: [012-return-request-and-refund-boundary.md](onion-architecture/lessons/012-return-request-and-refund-boundary.md)
- `013` Return restocking boundary: [013-return-restocking-boundary.md](onion-architecture/lessons/013-return-restocking-boundary.md)
- `014` Return review boundary: [014-return-review-boundary.md](onion-architecture/lessons/014-return-review-boundary.md)
- `015` Return eligibility policy: [015-return-eligibility-policy.md](onion-architecture/lessons/015-return-eligibility-policy.md)
- `016` Real return window policy: [016-real-return-window-policy.md](onion-architecture/lessons/016-real-return-window-policy.md)
- `017` Return actor metadata: [017-return-actor-metadata.md](onion-architecture/lessons/017-return-actor-metadata.md)
- `018` Return command idempotency: [018-return-command-idempotency.md](onion-architecture/lessons/018-return-command-idempotency.md)
- `019` Return query surface: [019-return-query-surface.md](onion-architecture/lessons/019-return-query-surface.md)
- `020` Order query surface: [020-order-query-surface.md](onion-architecture/lessons/020-order-query-surface.md)
- `021` Shipment query surface: [021-shipment-query-surface.md](onion-architecture/lessons/021-shipment-query-surface.md)
- `022` Quote list query surface: [022-quote-list-query-surface.md](onion-architecture/lessons/022-quote-list-query-surface.md)
- `023` Product query surface: [023-product-query-surface.md](onion-architecture/lessons/023-product-query-surface.md)
- `024` Customer query surface: [024-customer-query-surface.md](onion-architecture/lessons/024-customer-query-surface.md)
- `025` Quote conversion report: [025-quote-conversion-report.md](onion-architecture/lessons/025-quote-conversion-report.md)
- `026` Return rate by category report: [026-return-rate-by-category-report.md](onion-architecture/lessons/026-return-rate-by-category-report.md)
- `027` Low stock items report: [027-low-stock-items-report.md](onion-architecture/lessons/027-low-stock-items-report.md)
- `028` Orders awaiting approval report: [028-orders-awaiting-approval-report.md](onion-architecture/lessons/028-orders-awaiting-approval-report.md)
- `029` Payment review workflow: [029-payment-review-workflow.md](onion-architecture/lessons/029-payment-review-workflow.md)
- `030` Partial shipment support: [030-partial-shipment-support.md](onion-architecture/lessons/030-partial-shipment-support.md)
- `031` Partial returns by line: [031-partial-returns-by-line.md](onion-architecture/lessons/031-partial-returns-by-line.md)
- `032` Plugin pricing extension point: [032-plugin-pricing-extension-point.md](onion-architecture/lessons/032-plugin-pricing-extension-point.md)
- `033` Why not stop at onion?: [033-why-not-stop-at-onion.md](onion-architecture/lessons/033-why-not-stop-at-onion.md)

### Modular Monolith

#### Lessons

- `000` From onion to modular monolith: [000-from-onion-to-modular-monolith.md](modular-monolith/lessons/000-from-onion-to-modular-monolith.md)
- `001` Modular monolith skeleton: [001-modular-monolith-skeleton.md](modular-monolith/lessons/001-modular-monolith-skeleton.md)
- `002` Quote query through module API: [002-quote-query-through-module-api.md](modular-monolith/lessons/002-quote-query-through-module-api.md)
- `003` Add quote line with product module: [003-add-quote-line-with-product-module.md](modular-monolith/lessons/003-add-quote-line-with-product-module.md)
- `004` Submit quote state transition: [004-submit-quote-state-transition.md](modular-monolith/lessons/004-submit-quote-state-transition.md)
- `005` Approval policy module: [005-approval-policy-module.md](modular-monolith/lessons/005-approval-policy-module.md)
- `006` Approve pending quote: [006-approve-pending-quote.md](modular-monolith/lessons/006-approve-pending-quote.md)
- `007` Convert quote to order: [007-convert-quote-to-order.md](modular-monolith/lessons/007-convert-quote-to-order.md)
- `008` Order conversion with reservation: [008-order-conversion-with-reservation.md](modular-monolith/lessons/008-order-conversion-with-reservation.md)
- `009` Payment gateway and order capture: [009-payment-gateway-and-order-capture.md](modular-monolith/lessons/009-payment-gateway-and-order-capture.md)
- `010` Shipment creation after payment: [010-shipment-creation-after-payment.md](modular-monolith/lessons/010-shipment-creation-after-payment.md)
- `011` Order cancellation and release: [011-order-cancellation-and-release.md](modular-monolith/lessons/011-order-cancellation-and-release.md)
- `012` Return request and refund boundary: [012-return-request-and-refund-boundary.md](modular-monolith/lessons/012-return-request-and-refund-boundary.md)
- `013` Return restocking boundary: [013-return-restocking-boundary.md](modular-monolith/lessons/013-return-restocking-boundary.md)
- `014` Return review boundary: [014-return-review-boundary.md](modular-monolith/lessons/014-return-review-boundary.md)
- `015` Return eligibility policy: [015-return-eligibility-policy.md](modular-monolith/lessons/015-return-eligibility-policy.md)
- `016` Real return window policy: [016-real-return-window-policy.md](modular-monolith/lessons/016-real-return-window-policy.md)
- `017` Return actor metadata: [017-return-actor-metadata.md](modular-monolith/lessons/017-return-actor-metadata.md)
- `018` Return command idempotency: [018-return-command-idempotency.md](modular-monolith/lessons/018-return-command-idempotency.md)
- `019` Return query surface: [019-return-query-surface.md](modular-monolith/lessons/019-return-query-surface.md)
- `020` Order query surface: [020-order-query-surface.md](modular-monolith/lessons/020-order-query-surface.md)
- `021` Shipment query surface: [021-shipment-query-surface.md](modular-monolith/lessons/021-shipment-query-surface.md)
- `022` Quote list query surface: [022-quote-list-query-surface.md](modular-monolith/lessons/022-quote-list-query-surface.md)
- `023` Product query surface: [023-product-query-surface.md](modular-monolith/lessons/023-product-query-surface.md)
- `024` Customer query surface: [024-customer-query-surface.md](modular-monolith/lessons/024-customer-query-surface.md)
- `025` Quote conversion report: [025-quote-conversion-report.md](modular-monolith/lessons/025-quote-conversion-report.md)
- `026` Return rate by category report: [026-return-rate-by-category-report.md](modular-monolith/lessons/026-return-rate-by-category-report.md)
- `027` Low stock items report: [027-low-stock-items-report.md](modular-monolith/lessons/027-low-stock-items-report.md)
- `028` Orders awaiting approval report: [028-orders-awaiting-approval-report.md](modular-monolith/lessons/028-orders-awaiting-approval-report.md)
- `029` Payment review workflow: [029-payment-review-workflow.md](modular-monolith/lessons/029-payment-review-workflow.md)
- `030` Partial shipment support: [030-partial-shipment-support.md](modular-monolith/lessons/030-partial-shipment-support.md)
- `031` Partial returns by line: [031-partial-returns-by-line.md](modular-monolith/lessons/031-partial-returns-by-line.md)
- `032` Plugin pricing extension point: [032-plugin-pricing-extension-point.md](modular-monolith/lessons/032-plugin-pricing-extension-point.md)
- `033` Why not stop at modular monolith?: [033-why-not-stop-at-modular-monolith.md](modular-monolith/lessons/033-why-not-stop-at-modular-monolith.md)

### Microkernel / Plugin Architecture

#### Lessons

- `000` From modular monolith to microkernel: [000-from-modular-monolith-to-microkernel.md](microkernel-architecture/lessons/000-from-modular-monolith-to-microkernel.md)
- `001` Microkernel skeleton: [001-microkernel-skeleton.md](microkernel-architecture/lessons/001-microkernel-skeleton.md)
- `002` Quote query through kernel capability: [002-quote-query-through-kernel-capability.md](microkernel-architecture/lessons/002-quote-query-through-kernel-capability.md)
- `003` Add quote line with product plugin: [003-add-quote-line-with-product-plugin.md](microkernel-architecture/lessons/003-add-quote-line-with-product-plugin.md)
- `004` Submit quote state transition: [004-submit-quote-state-transition.md](microkernel-architecture/lessons/004-submit-quote-state-transition.md)
- `005` Approval policy plugin: [005-approval-policy-plugin.md](microkernel-architecture/lessons/005-approval-policy-plugin.md)
- `006` Approve pending quote: [006-approve-pending-quote.md](microkernel-architecture/lessons/006-approve-pending-quote.md)
- `007` Convert quote to order: [007-convert-quote-to-order.md](microkernel-architecture/lessons/007-convert-quote-to-order.md)
- `008` Order conversion with reservation: [008-order-conversion-with-reservation.md](microkernel-architecture/lessons/008-order-conversion-with-reservation.md)
- `009` Payment gateway and order capture: [009-payment-gateway-and-order-capture.md](microkernel-architecture/lessons/009-payment-gateway-and-order-capture.md)
- `010` Shipment creation after payment: [010-shipment-creation-after-payment.md](microkernel-architecture/lessons/010-shipment-creation-after-payment.md)
- `011` Order cancellation and release: [011-order-cancellation-and-release.md](microkernel-architecture/lessons/011-order-cancellation-and-release.md)
- `012` Return request and refund plugin: [012-return-request-and-refund-plugin.md](microkernel-architecture/lessons/012-return-request-and-refund-plugin.md)
- `013` Return restocking plugin: [013-return-restocking-plugin.md](microkernel-architecture/lessons/013-return-restocking-plugin.md)
- `014` Return review plugin: [014-return-review-plugin.md](microkernel-architecture/lessons/014-return-review-plugin.md)
- `015` Return eligibility plugin: [015-return-eligibility-plugin.md](microkernel-architecture/lessons/015-return-eligibility-plugin.md)
- `016` Real return window plugin: [016-real-return-window-plugin.md](microkernel-architecture/lessons/016-real-return-window-plugin.md)
- `017` Return actor metadata plugin: [017-return-actor-metadata-plugin.md](microkernel-architecture/lessons/017-return-actor-metadata-plugin.md)
- `018` Return command idempotency plugin: [018-return-command-idempotency-plugin.md](microkernel-architecture/lessons/018-return-command-idempotency-plugin.md)
- `019` Return query surface plugin: [019-return-query-surface-plugin.md](microkernel-architecture/lessons/019-return-query-surface-plugin.md)
- `020` Order query surface plugin: [020-order-query-surface-plugin.md](microkernel-architecture/lessons/020-order-query-surface-plugin.md)
- `021` Shipment query surface plugin: [021-shipment-query-surface-plugin.md](microkernel-architecture/lessons/021-shipment-query-surface-plugin.md)
- `022` Quote list query surface plugin: [022-quote-list-query-surface-plugin.md](microkernel-architecture/lessons/022-quote-list-query-surface-plugin.md)

## How To Maintain This File

As new architectures and lessons are added:

- add the architecture to the implemented list
- append the new lesson to the relevant lesson section
- keep the repository-walkthrough instructions current if the git workflow changes

