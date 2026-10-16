// Copyright 2026 3eAI Labs Ltd
// SPDX-License-Identifier: Apache-2.0

package diameterclient

// NATS subjects of the S6a and S13 interfaces. Every subject is NATS request/reply with a
// JSON body.
//
// Network-initiated (gateway → application). The gateway publishes a request when an MME
// sends the Diameter command and waits for the reply to build the Diameter answer.
const (
	// SubjectULR carries Update-Location-Request: S6aEvent → S6aAnswer.
	SubjectULR = "diameter.s6a.ulr"
	// SubjectAIR carries Authentication-Information-Request: S6aEvent → S6aAnswer.
	// The answer carries authentication vectors (key material); restrict this subject
	// with NATS permissions to the gateway and the authentication centre.
	SubjectAIR = "diameter.s6a.air"
	// SubjectPUR carries Purge-UE-Request: S6aEvent → S6aAnswer.
	SubjectPUR = "diameter.s6a.pur"
	// SubjectNOR carries Notify-Request: S6aEvent → S6aAnswer.
	SubjectNOR = "diameter.s6a.nor"
	// SubjectECR carries ME-Identity-Check-Request (S13): S13Event → S13Answer.
	SubjectECR = "diameter.s13.ecr"
)

// Application-initiated (application → gateway → MME). The application sends an S6aEvent
// whose Header.SourceHost names the target MME; the reply is an S6aAnswer carrying the
// outcome of the MME's answer: ResultCode, or ExperimentalResultCode (with ResultCode 0)
// when the MME answered with an Experimental-Result. Only the gateway subscribes to these
// subjects.
const (
	// SubjectCLR sends Cancel-Location-Request.
	SubjectCLR = "diameter.s6a.clr"
	// SubjectIDR sends Insert-Subscriber-Data-Request.
	SubjectIDR = "diameter.s6a.idr"
	// SubjectDSR sends Delete-Subscriber-Data-Request.
	SubjectDSR = "diameter.s6a.dsr"
	// SubjectRSR sends Reset-Request.
	SubjectRSR = "diameter.s6a.rsr"
)
