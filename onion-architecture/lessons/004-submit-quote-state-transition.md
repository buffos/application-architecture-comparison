# Lesson 004: Submit Quote State Transition

## Objective

Move the first lifecycle rule into the domain core by making quote submission an explicit state transition on the `Quote` entity.

## Theory

The previous Onion lessons showed:

- create a quote
- read a quote
- add lines to a quote

But the quote was still just a mutable container with a status field.

Onion Architecture becomes more meaningful when the domain core owns actual lifecycle rules.

In this lesson:

- the application ring still loads and saves the quote
- the domain entity decides whether submission is valid
- infrastructure remains a passive persistence detail

That is an important Onion move:

- orchestration stays in the application ring
- invariants stay in the domain core

## Why This Matters Here

If the application service decides by itself whether a quote can be submitted, then the domain is still thinner than it should be.

Putting submission on the entity makes the core more central:

- only draft quotes can be submitted
- empty quotes cannot be submitted
- once submitted, the quote is no longer editable

Those are domain rules, not storage rules and not CLI rules.

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Core"]
        direction TB
        QTE["Quote Entity"]
        SUB["Submit() Rule"]
    end

    subgraph APP["Application Ring"]
        direction TB
        SQS["SubmitQuote Service"]
        QST["Quote Store"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        CLI["CLI Framework"]
        MQR["Memory Quote Repository"]
    end

    CLI --> SQS
    SQS --> QTE
    QTE --> SUB

    QST -.used by.-> SQS
    QST -.implemented by.-> MQR

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef domain fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MQR dataadapter;
    class SQS,QST app;
    class QTE,SUB domain;
    class QST contract;
```

Legend:

- blue: framework edge
- green: data adapter
- purple: application ring
- yellow: domain core
- dashed border: interface / contract
- dashed arrow: structural relationship

## Implementation Focus

Implement one lifecycle use case:

- submit quote

The code should show:

- domain submission rules on `Quote`
- an application service that loads, submits, and saves
- the existing add-line path now respecting the submitted state
- the demo creating a quote, adding a line, submitting it, and loading it again

## What To Verify

- `go test ./...` passes
- a quote with lines can be submitted
- an empty quote cannot be submitted
- a submitted quote can no longer be edited
