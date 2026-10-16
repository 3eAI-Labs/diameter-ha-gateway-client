// Copyright 2026 3eAI Labs Ltd
// SPDX-License-Identifier: Apache-2.0

// Package diameterclient is the NATS wire contract of diameter-ha-gateway: the subjects
// and the JSON message schema of the S6a (MME ↔ HSS, 3GPP TS 29.272) and S13 (MME ↔ EIR)
// interfaces the gateway bridges to NATS.
//
// This module is licensed under the Apache License 2.0, separately from the gateway
// server (AGPL-3.0-only), so an HSS, an EIR or any other application can link it without
// taking on the server's licence. The dependency runs one way: the server imports this
// module for the wire contract, this module never imports the server.
//
// # Wire model
//
// Every subject is NATS request/reply with a JSON body.
//
//   - Network-initiated commands (ULR, AIR, PUR, NOR, ECR): the gateway sends an
//     [S6aEvent] (or [S13Event]) as a NATS request and builds the Diameter answer from
//     the [S6aAnswer] (or [S13Answer]) in the reply. See [SubjectULR] and the others.
//   - Application-initiated commands (CLR, IDR, DSR, RSR): the application sends an
//     [S6aEvent] whose Header.SourceHost names the target MME; the reply is an
//     [S6aAnswer] with the Result-Code of the MME's answer. See [SubjectCLR] and the others.
//
// Result codes are Diameter numbers: [S6aAnswer.ResultCode] for base protocol codes
// (RFC 6733 §7.1), [S6aAnswer.ExperimentalResultCode] for 3GPP codes under vendor 10415
// (TS 29.272 §7.4). The two number spaces overlap — 5001 is DIAMETER_AVP_UNSUPPORTED as a
// Result-Code and DIAMETER_ERROR_USER_UNKNOWN as an Experimental-Result — so put a 3GPP
// code in ExperimentalResultCode, never in ResultCode.
//
// # Scope
//
// Only the S6a and S13 contract is published here. The gateway's Gx, Gy, Ro and SGd NATS
// interfaces are not part of this module yet.
//
// # Security
//
// The [SubjectAIR] reply carries authentication vectors ([EUTRANVector]: RAND, XRES,
// AUTN, KASME). Restrict it with NATS permissions so that only the gateway can publish
// requests and only the authentication centre can subscribe, and run NATS with TLS.
// Never log these payloads.
package diameterclient
