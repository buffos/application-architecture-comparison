# Lesson 019: Return Query Surface Plugin

## Objective

Add an explicit read surface for return requests so callers load returns through a plugin capability instead of treating the repository as the public interface.

## Theory

The returns workflow already owns a meaningful write side:

- request
- review
- policy check
- refund and restock
- actor metadata
- idempotent review commands

But without explicit queries, callers still have an easy escape hatch:

- read the repository directly

That weakens the microkernel boundary because it makes the repository feel like the real API.

This lesson closes that gap:

- the returns plugin still owns persistence
- the plugin now exposes `GetReturnRequest`
- the plugin now exposes `ListReturnRequests`

So both write and read access go through kernel capabilities instead of leaking storage details to callers.

## Why This Matters Here

In a microkernel, plugin boundaries are not only about workflows and side effects. If reads bypass the plugin capability surface, the architecture quietly drifts toward shared storage with registration wrapped around it.

An explicit return query surface keeps the lesson honest:

- the repository remains internal plumbing
- the plugin owns the read model it chooses to expose
- callers depend on kernel capabilities, not storage details

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        RRD["kernel.ReturnReader"]
        RSA["kernel.ReturnService"]
        HST["kernel.Host"]
    end

    subgraph RTP["Returns Plugin"]
        direction TB
        RTR["returns.Repository"]
        RSS["returns.Service<br/>GetReturnRequest / ListReturnRequests"]
        RPP["returns.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MRR["Memory Return Repository"]
    end

    CLI --> HST
    HST --> RPP
    RSS --> RRD

    RTR -.used by.-> RSS
    RRD -.used by.-> CLI
    RSA -.used by.-> CLI
    RRD -.implemented by.-> RSS
    RSA -.implemented by.-> RSS
    RTR -.implemented by.-> MRR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class RRD,RSA,HST kernel;
    class RTR,RSS,RPP plugin;
    class MRR dataadapter;
    class CLI framework;
    class RRD,RSA,RTR contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or plugin registration type
- green: data adapter
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- add a kernel-owned return read capability
- expose `GetReturnRequest`
- expose `ListReturnRequests`
- support repository listing by status

Do not add order or shipment query surfaces yet.

## What To Verify

- `go test ./...` passes
- a stored return request can be loaded through the kernel capability
- return requests can be listed by status
- the demo can load and list returns without direct repository access
