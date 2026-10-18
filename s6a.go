// Copyright 2026 3eAI Labs Ltd
// SPDX-License-Identifier: Apache-2.0

package diameterclient

import "time"

type DiameterHeader struct {
	TraceID     string    `json:"trace_id"`
	SourceHost  string    `json:"source_host"`
	SourceRealm string    `json:"source_realm"`
	DestRealm   string    `json:"dest_realm"`
	Timestamp   time.Time `json:"timestamp"`
	HopByHop    uint32    `json:"hop_by_hop"`
	EndToEnd    uint32    `json:"end_to_end"`
	SessionID   string    `json:"session_id"`
}

// S6aEvent covers ULR, AIR, PUR, NOR
type S6aEvent struct {
	Header DiameterHeader `json:"header"`

	// Common AVPs
	IMSI              string             `json:"imsi"`
	SupportedFeatures []SupportedFeature `json:"supported_features,omitempty"`

	// ULR AVPs
	RATType        uint32        `json:"rat_type,omitempty"`
	ULRFlags       uint32        `json:"ulr_flags,omitempty"`
	VisitedPLMN    string        `json:"visited_plmn,omitempty"`
	SGSNNumber     string        `json:"sgsn_number,omitempty"`
	TerminalInfo   *TerminalInfo `json:"terminal_info,omitempty"`
	UE_SRVCC       uint32        `json:"ue_srvcc,omitempty"`
	GMLCAddress    string        `json:"gmlc_address,omitempty"`
	ActiveAPN      []string      `json:"active_apn,omitempty"`
	HomogeneousIMS uint32        `json:"homogeneous_ims,omitempty"`

	// AIR AVPs
	AuthInfo *AuthInfo `json:"auth_info,omitempty"`

	// AuthInfoUTRAN parses Requested-UTRAN-GERAN-Authentication-Info (1408 →
	// 1409, TS 29.272 §7.3.19) for converged-HSS deployments where one MME
	// owns 4G + 3G + 2G subscribers. Same child structure as AuthInfo
	// (Number-Of-Requested-Vectors / Immediate-Response-Preferred /
	// Re-Synchronization-Info), different parent AVP. Carrier IPLMN with
	// shared subscriber base will populate one or both per the access tech
	// the UE registered on.
	AuthInfoUTRAN *AuthInfo `json:"auth_info_utran,omitempty"`

	// EquivalentPLMNList is the ULR's Equivalent-PLMN-List (1637, TS 29.272
	// §7.3.151): one entry per Visited-PLMN-Id it contains, each hex-encoded
	// 3-octet TBCD (MCC + MNC).
	EquivalentPLMNList []string `json:"equivalent_plmn_list,omitempty"`

	// MMENumberforMTSMS (1645, TS 29.272 §7.3.158) — E.164 number of the MME
	// that the HSS/SMSC should target for SMS-over-NAS MT delivery. Used by
	// SGd routing of MT short messages. Hex-encoded TBCD bytes.
	MMENumberForMTSMS string `json:"mme_number_for_mt_sms,omitempty"`

	// TraceData (1458, TS 29.272 §7.3.63) — when MME requests trace activation
	// in ULR. ULR parser populates TraceReference today (slim path); the
	// remaining Trace-Data children are not parsed from the ULR. Same struct
	// type as SubscriptionProfile.TraceData.
	TraceData *TraceData `json:"trace_data,omitempty"`

	// DRMP (301, RFC 7944) — Diameter Routing Message Priority. When MME
	// marks ULR/AIR with DRMP, gateway propagates the same priority back on
	// the answer per RFC 7944 §3.2 (proxy nodes preserve DRMP). 0 = AVP not
	// emitted (treat as default routing priority).
	DRMP uint32 `json:"drmp,omitempty"`

	// PUR AVPs
	PURFlags uint32 `json:"pur_flags,omitempty"`

	// NOR AVPs
	NORFlags      uint32 `json:"nor_flags,omitempty"`
	ContextID     uint32 `json:"context_id,omitempty"`
	ServiceSelect string `json:"service_selection,omitempty"`
	AlertReason   uint32 `json:"alert_reason,omitempty"`

	// Reverse-direction (HSS-initiated) fields. Used by IDR/CLR/RSR/DSR senders.
	// CancellationType for CLR (default 0 MME_UPDATE_PROCEDURE).
	CancellationType uint32 `json:"cancellation_type,omitempty"`
	// IDRFlags for IDR (TS 29.272 §7.3.103).
	IDRFlags uint32 `json:"idr_flags,omitempty"`

	// DSRFlags for DSR (DSR-Flags, TS 29.272 §7.3.25, Table 7.3.25/1): which
	// subscription data to withdraw. Bit 0 Regional Subscription Withdrawal,
	// 1 Complete APN Configuration Profile Withdrawal, 2 Subscribed Charging
	// Characteristics Withdrawal, 3 PDN subscription contexts Withdrawal,
	// 4 STN-SR, 5 Complete PDP context list Withdrawal, 6 PDP contexts
	// Withdrawal, 7 Roaming Restricted due to unsupported feature, 8 Trace Data
	// Withdrawal, ... DSR-Flags is mandatory in the DSR (§7.2.11); with no bit
	// set the DSR withdraws nothing.
	DSRFlags uint32 `json:"dsr_flags,omitempty"`

	// WithdrawnAPNContextIDs scopes a DSR APN withdrawal to a specific subset.
	// Used in conjunction with DSRFlags bit for APN-Configuration-Profile
	// withdrawal — empty list combined with the "complete withdrawal" bit
	// removes all APNs; non-empty list removes only the listed
	// Context-Identifiers (TS 29.272 §6.1.4 + §7.3.27).
	WithdrawnAPNContextIDs []uint32 `json:"withdrawn_apn_context_ids,omitempty"`

	// SubProfile carries the structured Subscription-Data tree for IDR
	// (TS 29.272 §7.3.2 Grouped AVP 1400); the gateway encodes the tree as
	// the Grouped AVP. Preferred over SubscriberData.
	SubProfile *SubscriptionProfile `json:"subscription_profile,omitempty"`

	// SubscriberData is the legacy hex-OctetString fallback for Subscription-
	// Data AVP. Wire-invalid (TS 29.272 §7.3.2 mandates Grouped, not
	// OctetString) — kept for one release window so producers that have
	// not migrated to SubProfile keep working with a deprecation warning
	// logged by the gateway. It will be removed in a future release. New code
	// MUST populate SubProfile instead.
	//
	// Deprecated: use SubProfile.
	SubscriberData string `json:"subscriber_data,omitempty"`

	// Further AVPs of the requests from the MME (TS 29.272 §7.2). Octet
	// strings are hex-encoded.

	// ActiveAPNs are the ULR's Active-APN AVPs (1612, §7.3.127) whole;
	// ActiveAPN keeps only their Service-Selection.
	ActiveAPNs []ActiveAPN `json:"active_apns,omitempty"`
	// SMSRegisterRequest (1648, §7.3.162): SMS_REGISTRATION_REQUIRED (0),
	// SMS_REGISTRATION_NOT_PREFERRED (1), NO_PREFERENCE (2).
	SMSRegisterRequest *uint32 `json:"sms_register_request,omitempty"`
	// SGsMMEIdentity (1664, §7.3.175).
	SGsMMEIdentity string `json:"sgs_mme_identity,omitempty"`
	// CoupledNodeDiameterID (1666, §7.3.177).
	CoupledNodeDiameterID string `json:"coupled_node_diameter_id,omitempty"`
	// AdjacentPLMNs (1672, §7.3.186): its Visited-PLMN-Ids, hex.
	AdjacentPLMNs []string `json:"adjacent_plmns,omitempty"`
	// AIRFlags (1679, §7.3.201).
	AIRFlags uint32 `json:"air_flags,omitempty"`
	// EPSLocationInformation (1496, §7.3.111), of a PUR.
	EPSLocationInformation *EPSLocationInformation `json:"eps_location_information,omitempty"`
	// MaximumUEAvailabilityTime (3329, TS 29.338 §5.3.3.22), of a NOR.
	MaximumUEAvailabilityTime *time.Time `json:"maximum_ue_availability_time,omitempty"`
	// VisitedNetworkIdentifier (600, TS 29.229), of a NOR.
	VisitedNetworkIdentifier string `json:"visited_network_identifier,omitempty"`

	// Further AVPs of the requests an application sends (TS 29.272 §7.2).
	// SupportedFeatures and DRMP above are sent too.

	// CLRFlags (1638, §7.3.152), of a CLR.
	CLRFlags uint32 `json:"clr_flags,omitempty"`
	// ResetIDs (1670, §7.3.184), of an IDR or RSR; each hex.
	ResetIDs []string `json:"reset_ids,omitempty"`
	// UserIDs (1444, §7.3.50) of an RSR beyond the one IMSI carries: each a
	// leading part of an IMSI.
	UserIDs []string `json:"user_ids,omitempty"`
	// TraceReference (1459, §7.3.64), TSCodes (1487, §7.3.100) and SSCodes
	// (1476, §7.3.87) of a DSR; hex.
	TraceReference string   `json:"trace_reference,omitempty"`
	TSCodes        []string `json:"ts_codes,omitempty"`
	SSCodes        []string `json:"ss_codes,omitempty"`
}

// ActiveAPN (1612, TS 29.272 §7.3.127).
type ActiveAPN struct {
	ContextID                uint32            `json:"context_id"`
	ServiceSelection         string            `json:"service_selection,omitempty"`
	MIP6AgentInfo            *MIP6AgentInfo    `json:"mip6_agent_info,omitempty"`
	VisitedNetworkIdentifier string            `json:"visited_network_identifier,omitempty"`
	SpecificAPNInfo          []SpecificAPNInfo `json:"specific_apn_info,omitempty"`
}

// EPSLocationInformation (1496, TS 29.272 §7.3.111). Octet strings hex.
type EPSLocationInformation struct {
	MME  *MMELocationInformation  `json:"mme,omitempty"`
	SGSN *SGSNLocationInformation `json:"sgsn,omitempty"`
}

// MMELocationInformation (1600, §7.3.115). CurrentLocationRetrieved:
// ACTIVE-LOCATION-RETRIEVAL (0); AgeOfLocationInformation in minutes.
type MMELocationInformation struct {
	EUTRANCellGlobalIdentity string  `json:"e_utran_cell_global_identity,omitempty"`
	TrackingAreaIdentity     string  `json:"tracking_area_identity,omitempty"`
	GeographicalInformation  string  `json:"geographical_information,omitempty"`
	GeodeticInformation      string  `json:"geodetic_information,omitempty"`
	CurrentLocationRetrieved *uint32 `json:"current_location_retrieved,omitempty"`
	AgeOfLocationInformation *uint32 `json:"age_of_location_information,omitempty"`
}

// SGSNLocationInformation (1601, §7.3.116).
type SGSNLocationInformation struct {
	CellGlobalIdentity       string  `json:"cell_global_identity,omitempty"`
	LocationAreaIdentity     string  `json:"location_area_identity,omitempty"`
	ServiceAreaIdentity      string  `json:"service_area_identity,omitempty"`
	RoutingAreaIdentity      string  `json:"routing_area_identity,omitempty"`
	GeographicalInformation  string  `json:"geographical_information,omitempty"`
	GeodeticInformation      string  `json:"geodetic_information,omitempty"`
	CurrentLocationRetrieved *uint32 `json:"current_location_retrieved,omitempty"`
	AgeOfLocationInformation *uint32 `json:"age_of_location_information,omitempty"`
}

// EPSUserState (1495, §7.3.110). User-State (§7.3.114): DETACHED (0),
// ATTACHED_NOT_REACHABLE_FOR_PAGING (1), ATTACHED_REACHABLE_FOR_PAGING (2),
// CONNECTED_NOT_REACHABLE_FOR_PAGING (3), CONNECTED_REACHABLE_FOR_PAGING (4),
// RESERVED (5).
type EPSUserState struct {
	MMEUserState  *uint32 `json:"mme_user_state,omitempty"`
	SGSNUserState *uint32 `json:"sgsn_user_state,omitempty"`
}

// LocalTimeZone (1649, §7.3.156). TimeZone is "±n", n quarter hours from
// UTC (§7.3.163: "+8" is two hours ahead). DaylightSavingTime (§7.3.164):
// NO_ADJUSTMENT (0), PLUS_ONE_HOUR_ADJUSTMENT (1), PLUS_TWO_HOURS_ADJUSTMENT (2).
type LocalTimeZone struct {
	TimeZone           string `json:"time_zone"`
	DaylightSavingTime uint32 `json:"daylight_saving_time"`
}

type SupportedFeature struct {
	VendorID      uint32 `json:"vendor_id"`
	FeatureListID uint32 `json:"feature_list_id"`
	FeatureList   uint32 `json:"feature_list"`
}

type TerminalInfo struct {
	IMEI string `json:"imei,omitempty"`
	SV   string `json:"software_version,omitempty"`
	MEID string `json:"meid,omitempty"` // 3GPP2-MEID (1471, TS 29.272 §7.3.6), hex
}

type AuthInfo struct {
	RequestedVectors uint32 `json:"requested_vectors"`
	ImmediateRes     bool   `json:"immediate_response"`
	ReSyncInfo       string `json:"resync_info,omitempty"` // OctetString
}

// S13Event covers ECR
type S13Event struct {
	Header       DiameterHeader `json:"header"`
	IMSI         string         `json:"imsi"`
	TerminalInfo *TerminalInfo  `json:"terminal_info"`
}

// S6aAnswer is the cross-command S6a NATS reply envelope. The backend (the
// HSS application answering diameter.s6a.{air,ulr,pur,nor}) populates the
// subset of fields relevant to the answered command; the gateway selects them
// by command code when it builds the Diameter answer.
//
// Crypto material is carried as lowercase hex strings (see EUTRANVector) so
// the JSON envelope stays self-describing.
type S6aAnswer struct {
	ResultCode uint32 `json:"result_code"`

	// ExperimentalResultCode is a 3GPP (vendor 10415) result code, such as
	// DIAMETER_ERROR_USER_UNKNOWN (5001), carried in Experimental-Result with
	// no Result-Code (TS 29.272 §7.4.3). In a reply to the gateway, a non-zero
	// value makes the gateway answer the MME with Experimental-Result and
	// ResultCode is ignored. In a reply from the gateway to an
	// application-initiated command, exactly one of ResultCode and
	// ExperimentalResultCode is non-zero.
	ExperimentalResultCode uint32 `json:"experimental_result_code,omitempty"`

	// ULA fields
	ULAFlags uint32 `json:"ula_flags,omitempty"`

	// SubProfile carries the structured Subscription-Data tree for ULA
	// (TS 29.272 §7.3.2 Grouped AVP 1400); the gateway encodes the tree as
	// the Grouped AVP. Preferred over SubData.
	SubProfile *SubscriptionProfile `json:"subscription_profile,omitempty"`

	// SubData is the legacy hex-OctetString fallback for Subscription-Data
	// AVP. Wire-invalid (TS 29.272 §7.3.2 mandates Grouped, not OctetString) —
	// kept for one release window so producers that have not migrated to
	// SubProfile keep working with a deprecation warning logged by the
	// gateway. It will be removed in a future release. New code MUST populate
	// SubProfile instead.
	//
	// Deprecated: use SubProfile.
	SubData string `json:"subscriber_data,omitempty"`

	// AIA fields — at least one E-UTRAN-Vector required for a successful
	// AIA per TS 29.272 §7.3. Each vector populates a child Authentication-
	// Info group AVP (1413) containing nested E-UTRAN-Vector group (1414).
	AuthVectors []EUTRANVector `json:"auth_vectors,omitempty"`

	// PUA fields
	PUAFlags uint32 `json:"pua_flags,omitempty"`

	// ResetIDs are the Reset-ID AVPs (1670, TS 29.272 §7.3.184) of a ULA: the
	// HSS resources this subscriber's data depends on. Each entry is an
	// OctetString carried verbatim. Of the answers the gateway builds, only the
	// ULA carries Reset-ID (§7.2.4); it is ignored on the others.
	ResetIDs []string `json:"reset_ids,omitempty"` // each: hex-encoded bytes

	// SupportedFeatures echoes the HSS-side feature support per TS 29.229
	// §6.2 + TS 29.272 §7.3.10 (Rel-12+ feature negotiation: Multiple-PDN-
	// Connectivity, ARP-Override, Mission-Critical-Service, etc.). When set,
	// gateway emits one Supported-Features Grouped AVP (628) per entry on
	// ULA/AIA/PUA/NOA answers. Producer is responsible for the intersection
	// logic (HSS bits AND incoming MME bits); gateway just transports.
	SupportedFeatures []SupportedFeature `json:"supported_features,omitempty"`

	// FailedAVPs carries the offending AVP(s) when ResultCode signals an
	// AVP-related error (5004 INVALID_AVP_VALUE, 5005 MISSING_AVP, 5007
	// CONTRADICTING_AVPS, 5008 AVP_NOT_ALLOWED, 5009 AVP_OCCURS_TOO_MANY_TIMES,
	// 5014 INVALID_AVP_LENGTH). Per RFC 6733 §7.5 these answers SHOULD
	// include a Failed-AVP grouped AVP (279) wrapping the offending AVP so
	// peers can debug. The gateway encodes the entries the backend supplies.
	FailedAVPs []FailedAVPEntry `json:"failed_avps,omitempty"`

	// Redirect-Host indication (RFC 6733 §6.5). Populated when ResultCode
	// (or ExperimentalResultCode) is 3006 (DIAMETER_REDIRECT_INDICATION)
	// and the backend wants the gateway to encode Redirect-Host AVPs into
	// the outbound answer.
	//
	// RedirectHosts is the ordered list of DiameterURI targets — at least
	// one entry is required when ResultCode == 3006. RedirectHostUsage
	// follows §6.13 (0=DONT_CACHE, 1=ALL_SESSION, ..., 6=ALL_USER); when
	// omitted, peer treats as DONT_CACHE. RedirectMaxCacheTime is in
	// seconds (§6.12); zero ⇒ AVP not emitted (peer applies its default).
	RedirectHosts        []string `json:"redirect_hosts,omitempty"`
	RedirectHostUsage    uint32   `json:"redirect_host_usage,omitempty"`
	RedirectMaxCacheTime uint32   `json:"redirect_max_cache_time,omitempty"`

	// Further AVPs of an application's ULA or AIA (TS 29.272 §7.2).

	// ErrorDiagnostic (1614, §7.3.128), with an Experimental-Result.
	ErrorDiagnostic *uint32 `json:"error_diagnostic,omitempty"`
	// UTRANVectors and GERANVectors (1415, 1416) of an AIA, beside
	// AuthVectors (E-UTRAN).
	UTRANVectors []UTRANVector `json:"utran_vectors,omitempty"`
	GERANVectors []GERANVector `json:"geran_vectors,omitempty"`
	// UEUsageType (1680, §7.3.202), of an AIA.
	UEUsageType *uint32 `json:"ue_usage_type,omitempty"`

	// The MME's answer to a CLR, IDR, DSR or RSR, in the gateway's reply to
	// the application; SupportedFeatures and FailedAVPs above carry the
	// answer's own. Octet strings hex.

	// IDAFlags (1441, §7.3.47) and DSAFlags (1422, §7.3.26).
	IDAFlags uint32 `json:"ida_flags,omitempty"`
	DSAFlags uint32 `json:"dsa_flags,omitempty"`
	// IMSVoiceOverPSSessionsSupported (1492, §7.3.106): NOT_SUPPORTED (0),
	// SUPPORTED (1).
	IMSVoiceOverPSSessionsSupported *uint32 `json:"ims_voice_over_ps_sessions_supported,omitempty"`
	// LastUEActivityTime (1494, §7.3.108).
	LastUEActivityTime *time.Time `json:"last_ue_activity_time,omitempty"`
	// RATType (1032, TS 29.212 §5.3.31) of the IDA.
	RATType                uint32                  `json:"rat_type,omitempty"`
	EPSUserState           *EPSUserState           `json:"eps_user_state,omitempty"`
	EPSLocationInformation *EPSLocationInformation `json:"eps_location_information,omitempty"`
	LocalTimeZone          *LocalTimeZone          `json:"local_time_zone,omitempty"`
}

// UTRANVector (1415, TS 29.272 §7.3.19); octet strings hex.
type UTRANVector struct {
	Item uint32 `json:"item"`
	RAND string `json:"rand"`
	XRES string `json:"xres"`
	AUTN string `json:"autn"`
	CK   string `json:"ck"` // Confidentiality-Key (625, TS 29.229)
	IK   string `json:"ik"` // Integrity-Key (626, TS 29.229)
}

// GERANVector (1416, TS 29.272 §7.3.20); octet strings hex.
type GERANVector struct {
	Item uint32 `json:"item"`
	RAND string `json:"rand"`
	SRES string `json:"sres"`
	Kc   string `json:"kc"`
}

// FailedAVPEntry describes a single AVP to wrap inside a Failed-AVP grouped
// AVP (RFC 6733 §7.5). Value is the raw payload bytes the producer wants
// echoed; for OctetString-typed AVPs this is the bytes verbatim, for
// Unsigned32 it's a 4-byte big-endian encoding, etc. Producers that don't
// have access to the original wire bytes can supply an empty Value — the
// encoder still emits the AVP code so peers see the Failed-AVP shape.
//
// VendorID is 0 for base IETF AVPs (User-Name, Session-Id, ...) and 10415
// for 3GPP vendor-specific (RAT-Type, ULR-Flags, ...). The encoder sets
// the V-bit when VendorID != 0.
type FailedAVPEntry struct {
	Code     uint32 `json:"code"`
	VendorID uint32 `json:"vendor_id,omitempty"`
	// Value is the raw AVP payload as a hex-encoded string. Empty value
	// means "echo the AVP code only" (peer learns the offending field
	// even without seeing the rejected payload).
	Value string `json:"value,omitempty"`
}

// EUTRANVector is one EPS authentication vector (TS 33.401 / TS 29.272
// §7.3.18). Crypto values use lowercase hex strings without 0x prefix —
// the gateway hex-decodes them into OctetString AVP payloads. Hex (not
// raw bytes) keeps the JSON envelope self-describing for ops debugging.
type EUTRANVector struct {
	Item  uint32 `json:"item"`  // 1..N — Item-Number AVP (1419)
	RAND  string `json:"rand"`  // hex 32 chars (16 bytes)
	XRES  string `json:"xres"`  // hex 16 chars (8 bytes)
	AUTN  string `json:"autn"`  // hex 32 chars (16 bytes)
	KASME string `json:"kasme"` // hex 64 chars (32 bytes)
}

type S13Answer struct {
	ResultCode uint32 `json:"result_code"`
	// ExperimentalResultCode is a 3GPP result such as
	// DIAMETER_ERROR_EQUIPMENT_UNKNOWN (5422, TS 29.272 §7.4.3.5); when
	// non-zero the gateway answers with Experimental-Result and ResultCode is
	// ignored.
	ExperimentalResultCode uint32 `json:"experimental_result_code,omitempty"`
	// Status is the Equipment-Status (TS 29.272 §7.3.51): 0 PERMITTEDLISTED,
	// 1 PROHIBITEDLISTED, 2 TRACKINGLISTED. Sent only with ResultCode 2001.
	Status uint32 `json:"equipment_status"`
}
