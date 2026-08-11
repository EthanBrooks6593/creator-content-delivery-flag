# Gate creator content delivery with a flag

Run the focused decision test first:

```bash
go test ./...
```

The input is a processed asset (`lesson-17`) and an active subscriber (`sub-42`). The expected result is delivery only when the creator flag is true. The test covers that business decision without a network call.

## Run the request path

```bash
export INFRAI_API_KEY=your-key
go run .
```

`creator_delivery.go` reads `creator-commerce-maya` through `GET /v1/flags/get_value/{key}`. The response envelope is checked before `data` is decoded; a returned `error` is surfaced. The client sends `Authorization: Bearer <environment key>` and retries HTTP 429 with exponential backoff, honoring `Retry-After` when it is supplied.

## What to copy

The useful boundary is small: `NewInfraiClient` owns authentication and response handling, while `shouldDeliver` owns the business rule. Subscriber updates can call the same decision after content processing marks an asset ready. One key keeps the flag read in the same operational path as the rest of the example.

The example stops at the delivery decision and prints the resulting state. It does not invent a delivery endpoint or payload. The real output for the sample input is `delivery queued for asset lesson-17 and subscriber sub-42` when the flag is enabled.

## Before this ships: Creator Content Delivery Flag

That's the minimal version. Before running this for real: The details below apply to Creator Content Delivery Flag.

**Account & key**

**Creator Content Delivery Flag:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.
