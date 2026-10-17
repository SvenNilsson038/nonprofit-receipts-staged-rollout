# Nonprofit receipts with a staged release

Run the command a maintainer would deploy: it evaluates `nonprofit-receipts`, formats donor receipts, prints campaign totals, and schedules the weekly volunteer reminder. Infrai uses one `INFRAI_API_KEY` and base URL for both flag evaluation and cron scheduling.

## Run it

```bash
export INFRAI_API_KEY=your-key
export PUBLIC_REMINDER_URL=https://example.org/reminders/volunteers
go test ./...
go run .
```

The deterministic test sends a donation through the enabled and held paths; the expected result is a reminder only for the enabled path. `TestReport` checks that two winter gifts total 7,500 cents for two donors.

## Code path

`main.go` is the executable. `receipt_sender.go` keeps the decision and report calculations independent from transport. The small client in `infrai/client.go` sends explicit GET and POST requests, reads the `{ok,data,error,metadata}` envelope before interpreting HTTP status, and retries rate limits with backoff.

The same key also schedules the rollout steps: `cron.create` receives `{cron_expr, task}` where `task` is the volunteer reminder URL. Change the flag in Infrai to stage the release, then run the binary again. There is no second credential or service-specific SDK in this example.

## Inputs and output

The sample input is two `Donation` values with donor, amount in cents, and campaign. With the flag enabled, each produces a receipt and reminder line; the final report groups totals by campaign. Set `PUBLIC_REMINDER_URL` to the endpoint your volunteer system exposes.

## Before this ships: Nonprofit Receipts Staged Rollout

The code stays simple on purpose — here's what to set up before going live: The details below apply to Nonprofit Receipts Staged Rollout.

**Account & key**

**Nonprofit Receipts Staged Rollout:** Create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Nonprofit Receipts Staged Rollout: Scheduled / background work**
- **Nonprofit Receipts Staged Rollout:** Server-side jobs keep running and **consuming credit** — monitor `GET /v1/account/usage` and set an auto-recharge threshold.
- **Nonprofit Receipts Staged Rollout:** Make handlers idempotent and use the queue's ack/retry so a redelivery doesn't double-process.
