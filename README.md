# diameter-ha-gateway-client

NATS wire contract for [diameter-ha-gateway](https://github.com/3eAI-Labs/diameter-ha-gateway), a
Diameter gateway that bridges an MME's S6a and S13 interfaces (3GPP TS 29.272) to applications over
NATS.

Licensed under the Apache License 2.0, separately from the gateway (AGPL-3.0-only), so an HSS, an
EIR or any other application can link it under any licence.

```bash
go get github.com/3eAI-Labs/diameter-ha-gateway-client
```

## What it covers

The subjects and JSON envelopes of S6a and S13. Every subject is core NATS request/reply. Keep
these subjects out of JetStream streams: a stream that captures a subject answers every request with
its own publish acknowledgement, so the requester never sees the reply.

| Diameter command | Direction | Subject | Request → reply |
|------------------|-----------|---------|-----------------|
| Update-Location (ULR/ULA) | MME → application | `diameter.s6a.ulr` | `S6aEvent` → `S6aAnswer` |
| Authentication-Information (AIR/AIA) | MME → application | `diameter.s6a.air` | `S6aEvent` → `S6aAnswer` |
| Purge-UE (PUR/PUA) | MME → application | `diameter.s6a.pur` | `S6aEvent` → `S6aAnswer` |
| Notify (NOR/NOA) | MME → application | `diameter.s6a.nor` | `S6aEvent` → `S6aAnswer` |
| ME-Identity-Check (ECR/ECA, S13) | MME → application | `diameter.s13.ecr` | `S13Event` → `S13Answer` |
| Cancel-Location (CLR/CLA) | application → MME | `diameter.s6a.clr` | `S6aEvent` → `S6aAnswer` |
| Insert-Subscriber-Data (IDR/IDA) | application → MME | `diameter.s6a.idr` | `S6aEvent` → `S6aAnswer` |
| Delete-Subscriber-Data (DSR/DSA) | application → MME | `diameter.s6a.dsr` | `S6aEvent` → `S6aAnswer` |
| Reset (RSR/RSA) | application → MME | `diameter.s6a.rsr` | `S6aEvent` → `S6aAnswer` |

For application-initiated commands, `S6aEvent.Header.SourceHost` names the target MME, only the
gateway subscribes to the subject, and the reply carries `ResultCode` or, when the MME answered with
an Experimental-Result, `ExperimentalResultCode`.

Purge-UE and Notify are sent by the MME to the HSS only (TS 29.272 §7.2.13, §7.2.17), so there is no
application-initiated PUR or NOR.

The gateway's Gx, Gy, Ro and SGd NATS interfaces are not part of this module yet.

## Usage

An HSS answering Update-Location:

```go
nc, err := nats.Connect(natsURL)
if err != nil { /* ... */ }

_, err = nc.QueueSubscribe(diameterclient.SubjectULR, "hss", func(m *nats.Msg) {
	var ev diameterclient.S6aEvent
	if err := json.Unmarshal(m.Data, &ev); err != nil {
		return // the gateway times out and answers the MME with an error
	}
	ans := diameterclient.S6aAnswer{ResultCode: 2001} // DIAMETER_SUCCESS
	if !known(ev.IMSI) {
		// 3GPP result codes go in ExperimentalResultCode (TS 29.272 §7.4), never in ResultCode.
		ans = diameterclient.S6aAnswer{ExperimentalResultCode: 5001} // DIAMETER_ERROR_USER_UNKNOWN
	} else {
		ans.SubProfile = profileFor(ev.IMSI)
	}
	data, _ := json.Marshal(ans)
	_ = m.Respond(data)
})
```

Subscribe with a queue group so that one instance of your application answers each request.

## Security

The `diameter.s6a.air` reply carries authentication vectors (RAND, XRES, AUTN, KASME). Restrict the
subject with NATS permissions — only the gateway publishes requests, only the authentication centre
subscribes — run NATS with TLS, and never log these payloads.

## Versioning

Below 1.0, a minor version may change the schema; patch versions do not. `TestWireSchema` pins
every field name and JSON type of the envelopes.

## Licence

[Apache License 2.0](LICENSE). Contributions require the
[3eAI Labs Contributor Licence Agreement](CLA.md).
