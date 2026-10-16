// Copyright 2026 3eAI Labs Ltd
// SPDX-License-Identifier: Apache-2.0

package diameterclient

// SubscriptionProfile mirrors the 3GPP TS 29.272 §7.3.2 Subscription-Data
// Grouped AVP (1400, vendor 10415). The backend (the HSS application)
// populates this struct on `S6aAnswer.SubProfile` for ULA and on
// `S6aEvent.SubProfile` for IDR; the gateway translates the tree into a
// Diameter Grouped AVP per TS 29.272 §7.3.
//
// This replaces the legacy `SubData string` field that carried pre-encoded
// AVP bytes — that path was wire-invalid because 1400 is Grouped, not
// OctetString. `SubData` is still parsed as a deprecated fallback, with a
// deprecation warning logged by the gateway.
//
// Field ordering follows TS 29.272 §7.3.2 Table 7.3.2/1 to keep the struct
// scannable against the spec. Pointer types are used for optional scalars
// where 0 is a valid spec value (e.g. Subscriber-Status SERVICE_GRANTED=0,
// Network-Access-Mode PACKET_AND_CIRCUIT=0, Pre-emption flags) — nil means
// "do not emit the AVP". Plain uint32 / string types are used where the
// zero-value naturally means "not set" (bitmasks, bandwidths, identifiers).
//
// Subscription-Data children that are not listed here can be added later
// without breaking the JSON envelope (additive fields).
type SubscriptionProfile struct {
	// Subscriber-Status (1424) — Enumerated. 0=SERVICE_GRANTED, 1=
	// OPERATOR_DETERMINED_BARRING. Pointer because 0 is a valid value.
	SubscriberStatus *uint32 `json:"subscriber_status,omitempty"`

	// MSISDN (701) — OctetString, TBCD-encoded MSISDN per TS 29.002.
	// Carrier-typical attach response always includes MSISDN; gateway
	// treats the bytes as opaque (TBCD encoding is the backend's job).
	MSISDN string `json:"msisdn,omitempty"` // hex-encoded TBCD bytes

	// A-MSISDN (1643) — OctetString. Additional MSISDN for SRVCC scenarios.
	AMSISDN string `json:"a_msisdn,omitempty"` // hex-encoded TBCD bytes

	// STN-SR (1433) — OctetString. Session Transfer Number for SRVCC.
	STNSR string `json:"stn_sr,omitempty"` // hex-encoded TBCD bytes

	// ICS-Indicator (1491) — Enumerated. 0=FALSE, 1=TRUE. IMS Centralized
	// Services indicator. Pointer because 0 is a valid value.
	ICSIndicator *uint32 `json:"ics_indicator,omitempty"`

	// Network-Access-Mode (1417) — Enumerated. 0=PACKET_AND_CIRCUIT,
	// 1=Reserved, 2=ONLY_PACKET. Pointer because 0 is a valid value.
	NetworkAccessMode *uint32 `json:"network_access_mode,omitempty"`

	// Operator-Determined-Barring (1425) — Unsigned32, bitmask per TS 29.272
	// §7.3.30. Bit 0 = barring of all packet oriented services, etc.
	// Zero = no barring (omit AVP).
	OperatorDeterminedBarring uint32 `json:"operator_determined_barring,omitempty"`

	// HPLMN-ODB (1418) — Unsigned32, bitmask of HPLMN-specific ODB
	// categories per TS 23.015. Zero = none.
	HPLMNODB uint32 `json:"hplmn_odb,omitempty"`

	// Regional-Subscription-Zone-Code (1446) — OctetString, 2 bytes each.
	// Up to 10 zones may be present. Empty slice = AVP not emitted.
	RegionalSubscriptionZoneCodes []string `json:"regional_subscription_zone_codes,omitempty"` // each: hex 4 chars

	// Access-Restriction-Data (1426) — Unsigned32, bitmask per TS 29.272
	// §7.3.31. Bit 0 = UTRAN-NotAllowed, bit 1 = GERAN-NotAllowed, ...
	// Zero = no restrictions (omit AVP).
	AccessRestrictionData uint32 `json:"access_restriction_data,omitempty"`

	// APN-OI-Replacement (1427) — UTF8String. Operator-controlled APN-OI
	// replacement string applied during APN construction.
	APNOIReplacement string `json:"apn_oi_replacement,omitempty"`

	// LCS-Info (1473) — Grouped. Location services subscription info.
	LCSInfo *LCSInfo `json:"lcs_info,omitempty"`

	// Teleservice-List (1486) — Grouped. CS-domain teleservice subscription.
	// Carrier-typical EPS attach omits this; included as struct stub.
	TeleserviceList *TeleserviceList `json:"teleservice_list,omitempty"`

	// 3GPP-Charging-Characteristics (13) — OctetString, 4 bytes (1 nibble
	// profile + 3 nibble behaviour bits). Carrier billing handle.
	ChargingCharacteristics string `json:"charging_characteristics,omitempty"` // hex 4 chars

	// AMBR (1435) — Grouped. UE-aggregate maximum bit rate (UE-AMBR).
	// NOT to be confused with per-APN AMBR which lives inside APN-Configuration.
	AMBR *AMBR `json:"ambr,omitempty"`

	// APN-Configuration-Profile (1429) — Grouped. The big one — list of
	// APNs the subscriber is allowed to use plus per-APN QoS / AMBR.
	APNConfigurationProfile *APNConfigurationProfile `json:"apn_configuration_profile,omitempty"`

	// RAT-Frequency-Selection-Priority-ID (1440) — Unsigned32.
	RATFrequencySelectionPriorityID uint32 `json:"rat_frequency_selection_priority_id,omitempty"`

	// CSG-Subscription-Data (1436) — Grouped, repeated. Closed Subscriber
	// Group memberships. Carrier-typical EPS attach omits unless femtocell.
	CSGSubscriptionData []CSGSubscriptionData `json:"csg_subscription_data,omitempty"`

	// Roaming-Restricted-Due-To-Unsupported-Feature (1457) — Enumerated.
	// 0=Roaming-Restricted-Due-To-Unsupported-Feature. Presence-only AVP.
	RoamingRestrictedDueToUnsupportedFeature *uint32 `json:"roaming_restricted_due_to_unsupported_feature,omitempty"`

	// Subscribed-Periodic-RAU-TAU-Timer (1619) — Unsigned32, seconds.
	// IoT use case (NB-IoT extended periodic timer).
	SubscribedPeriodicRAUTAUTimer uint32 `json:"subscribed_periodic_rau_tau_timer,omitempty"`

	// MPS-Priority (1616) — Unsigned32, bitmask. Multimedia Priority
	// Service entitlement (eMPS / public safety).
	MPSPriority uint32 `json:"mps_priority,omitempty"`

	// VPLMN-LIPA-Allowed (1617) — Enumerated. 0=LIPA_NOTALLOWED,
	// 1=LIPA_ALLOWED. Pointer because 0 is a valid value.
	VPLMNLIPAAllowed *uint32 `json:"vplmn_lipa_allowed,omitempty"`

	// Relay-Node-Indicator (1633) — Enumerated. 0=NOT_RELAY_NODE,
	// 1=RELAY_NODE. Pointer because 0 is a valid value.
	RelayNodeIndicator *uint32 `json:"relay_node_indicator,omitempty"`

	// MDT-User-Consent (1634) — Enumerated. 0=CONSENT_NOT_GIVEN,
	// 1=CONSENT_GIVEN. Pointer because 0 is a valid value.
	MDTUserConsent *uint32 `json:"mdt_user_consent,omitempty"`

	// Subscribed-VSRVCC (1636) — Enumerated. 0=VSRVCC_SUBSCRIBED.
	// Presence-only AVP.
	SubscribedVSRVCC *uint32 `json:"subscribed_vsrvcc,omitempty"`

	// ProSe-Subscription-Data (3701) — Grouped (Rel 12+). Proximity
	// Services subscription. Carrier-typical: only set for
	// public-safety / D2D-enabled subscribers.
	ProSeSubscriptionData *ProSeSubscriptionData `json:"prose_subscription_data,omitempty"`

	// Subscription-Data-Flags (1654) — Unsigned32, bitmask (Rel 14+).
	// Bit 0 = PCSCFRestoration_indicator, bit 1 = NRasSecondaryRAT, etc.
	SubscriptionDataFlags uint32 `json:"subscription_data_flags,omitempty"`

	// ---- Extended subscription data (TS 29.272 §7.3.2, TS 29.336 §8.4) ----

	// Call-Barring-Info (1488) — Grouped, multi. CS-domain SS-Code/SS-Status
	// barring per TS 22.030.
	CallBarringInfo []CallBarringInfo `json:"call_barring_info,omitempty"`

	// GPRS-Subscription-Data (1467) — Grouped. Converged-HSS 2G/3G PDP
	// context list.
	GPRSSubscriptionData *GPRSSubscriptionData `json:"gprs_subscription_data,omitempty"`

	// Trace-Data (1458, TS 29.272 §7.3.63) — Grouped. Subscriber and
	// equipment trace / MDT activation (TS 32.422).
	TraceData *TraceData `json:"trace_data,omitempty"`

	// IMSI-Group-Id (1675) — Grouped, multi (Rel-13). Group identification
	// for service-bound IMSI sets.
	IMSIGroupID []IMSIGroupID `json:"imsi_group_id,omitempty"`

	// Adjacent-Access-Restriction-Data (1673) — Grouped, multi (Rel-14).
	// Per-VPLMN restriction overrides.
	AdjacentAccessRestrictionData []AdjacentAccessRestrictionData `json:"adjacent_access_restriction_data,omitempty"`

	// V2X-Subscription-Data (1688) — Grouped (Rel-14). V2X permissions.
	V2XSubscriptionData *V2XSubscriptionData `json:"v2x_subscription_data,omitempty"`

	// Emergency-Info (1687) — Grouped (Rel-14). Emergency-bearer subscription.
	EmergencyInfo *EmergencyInfo `json:"emergency_info,omitempty"`

	// eDRX-Cycle-Length (1691) — Grouped, one per RAT type (Rel-13).
	EDRXCycleLength []EDRXCycleLength `json:"edrx_cycle_length,omitempty"`

	// Paging-Time-Window (1701) — Grouped, one per operation mode (Rel-13).
	PagingTimeWindow []PagingTimeWindow `json:"paging_time_window,omitempty"`

	// UE-Usage-Type (1680) — Unsigned32. DECOR / NB-IoT slicing tag.
	UEUsageType uint32 `json:"ue_usage_type,omitempty"`

	// DL-Buffering-Suggested-Packet-Count (1674) — Integer32 (Rel-13).
	// NB-IoT high-latency communication hint to PGW.
	DLBufferingSuggestedPacketCount int32 `json:"dl_buffering_suggested_packet_count,omitempty"`

	// Service-Gap-Time (1698) — Unsigned32 seconds (Rel-13). Power-saving.
	ServiceGapTime uint32 `json:"service_gap_time,omitempty"`
}

// AMBR mirrors TS 29.272 §7.3.41 (also reused as 3GPP TS 29.212 QoS-
// Information sub-AVP). Used for UE-AMBR (subscription level) and as
// per-APN AMBR inside APN-Configuration. Both rates in bits per second.
// Extended-Max-Requested-BW UL/DL (RFC 8506 §10.1) carry Rel-15 high-rate
// flows where the 32-bit MaxRequestedBandwidth fields would saturate.
type AMBR struct {
	MaxRequestedBandwidthUL  uint32 `json:"max_requested_bandwidth_ul,omitempty"`
	MaxRequestedBandwidthDL  uint32 `json:"max_requested_bandwidth_dl,omitempty"`
	ExtendedMaxRequestedBWUL uint32 `json:"extended_max_requested_bw_ul,omitempty"` // AVP 554  Rel-15
	ExtendedMaxRequestedBWDL uint32 `json:"extended_max_requested_bw_dl,omitempty"` // AVP 555  Rel-15
}

// APNConfigurationProfile mirrors TS 29.272 §7.3.34 Grouped AVP (1429).
// Carrier-typical: ContextIdentifier points to the default APN whose
// configuration appears in the APNConfigurations list (the default APN's
// Context-Identifier matches the top-level field).
type APNConfigurationProfile struct {
	// Context-Identifier (1423) — Unsigned32. ID of the default APN
	// (matches one of APNConfigurations[].ContextIdentifier).
	ContextIdentifier uint32 `json:"context_identifier"`

	// All-APN-Configurations-Included-Indicator (1428) — Enumerated.
	// 0=ALL_APN_CONFIGURATIONS_INCLUDED, 1=MODIFIED/ADDED_APN_CONFIGURATIONS_
	// INCLUDED. Pointer because 0 is a valid value.
	AllAPNConfigurationsIncludedIndicator *uint32 `json:"all_apn_configurations_included_indicator,omitempty"`

	// APN-Configuration (1430) — Grouped, repeated. At least one entry
	// for a successful EPS attach (TS 29.272 §7.3.35).
	APNConfigurations []APNConfiguration `json:"apn_configurations,omitempty"`
}

// APNConfiguration mirrors TS 29.272 §7.3.35 Grouped AVP (1430).
type APNConfiguration struct {
	// Context-Identifier (1423) — Unsigned32. APN ID, mandatory.
	ContextIdentifier uint32 `json:"context_identifier"`

	// PDN-Type (1456) — Enumerated. 0=IPv4, 1=IPv6, 2=IPv4v6,
	// 3=IPv4_OR_IPv6. Pointer because 0 is a valid value.
	PDNType *uint32 `json:"pdn_type,omitempty"`

	// Service-Selection (493) — UTF8String. The APN name (e.g. "internet",
	// "ims", "sos"). Mandatory in successful APN-Configuration.
	ServiceSelection string `json:"service_selection,omitempty"`

	// EPS-Subscribed-QoS-Profile (1431) — Grouped. Default-bearer QoS.
	EPSSubscribedQoSProfile *EPSSubscribedQoSProfile `json:"eps_subscribed_qos_profile,omitempty"`

	// VPLMN-Dynamic-Address-Allowed (1432) — Enumerated. 0=NOTALLOWED,
	// 1=ALLOWED. Pointer because 0 is a valid value.
	VPLMNDynamicAddressAllowed *uint32 `json:"vplmn_dynamic_address_allowed,omitempty"`

	// Visited-Network-Identifier (600) — OctetString. Used in roaming.
	VisitedNetworkIdentifier string `json:"visited_network_identifier,omitempty"`

	// PDN-GW-Allocation-Type (1438) — Enumerated. 0=STATIC, 1=DYNAMIC.
	// Pointer because 0 is a valid value.
	PDNGWAllocationType *uint32 `json:"pdn_gw_allocation_type,omitempty"`

	// 3GPP-Charging-Characteristics (13) — OctetString, 4-hex-char string.
	// Per-APN override of the subscription-level value.
	ChargingCharacteristics string `json:"charging_characteristics,omitempty"`

	// AMBR (1435) — Grouped. APN-AMBR (per-APN). NOT UE-AMBR.
	AMBR *AMBR `json:"ambr,omitempty"`

	// APN-OI-Replacement (1427) — UTF8String. Per-APN override.
	APNOIReplacement string `json:"apn_oi_replacement,omitempty"`

	// ---- Extended APN configuration ----

	// Served-Party-IP-Address (848) — Address, multi (max 2: one v4 + one
	// v6). IETF base AVP. JSON: textual IP address.
	ServedPartyIPAddress []string `json:"served_party_ip_address,omitempty"`

	// MIP6-Agent-Info (486) — Grouped. Static PDN-GW selection (DSMIPv6).
	MIP6AgentInfo *MIP6AgentInfo `json:"mip6_agent_info,omitempty"`

	// Specific-APN-Info (1472) — Grouped, multi (Rel-10). Per-service
	// MIP6 / VisitedNetworkIdentifier overrides scoped to a service-selection.
	SpecificAPNInfo []SpecificAPNInfo `json:"specific_apn_info,omitempty"`

	// SIPTO-Permission (1613) — Unsigned32.
	SIPTOPermission uint32 `json:"sipto_permission,omitempty"`

	// LIPA-Permission (1618) — Unsigned32.
	LIPAPermission uint32 `json:"lipa_permission,omitempty"`

	// Restoration-Priority (1663) — Unsigned32. Per-APN override of
	// the top-level Restoration-Priority.
	RestorationPriority uint32 `json:"restoration_priority,omitempty"`

	// SIPTO-Local-Network-Permission (1665) — Unsigned32.
	SIPTOLocalNetworkPermission uint32 `json:"sipto_local_network_permission,omitempty"`

	// WLAN-Offloadability (1667) — Grouped (Rel-12). Per-APN E-UTRAN /
	// UTRAN WLAN offload permission.
	WLANOffloadability *WLANOffloadability `json:"wlan_offloadability,omitempty"`

	// Non-IP-PDN-Type-Indicator (1681) — Enumerated (NB-IoT).
	NonIPPDNTypeIndicator uint32 `json:"non_ip_pdn_type_indicator,omitempty"`

	// Non-IP-Data-Delivery-Mechanism (1682) — Unsigned32 (NB-IoT
	// SCEF/SGi delivery mechanism).
	NonIPDataDeliveryMechanism uint32 `json:"non_ip_data_delivery_mechanism,omitempty"`

	// SCEF-ID (3125) — DiameterIdentity. Per-APN SCEF binding.
	SCEFID string `json:"scef_id,omitempty"`

	// SCEF-Realm (1684) — DiameterIdentity.
	SCEFRealm string `json:"scef_realm,omitempty"`

	// Preferred-Data-Mode (1686) — Unsigned32 bitmask (NB-IoT data mode).
	PreferredDataMode uint32 `json:"preferred_data_mode,omitempty"`

	// PDN-Connection-Continuity (1690) — Unsigned32.
	PDNConnectionContinuity uint32 `json:"pdn_connection_continuity,omitempty"`
}

// EPSSubscribedQoSProfile mirrors TS 29.272 §7.3.37 Grouped AVP (1431).
// Default-bearer QoS for the APN's default EPS bearer.
type EPSSubscribedQoSProfile struct {
	// QoS-Class-Identifier (1028) — Enumerated. 1..9 standardized QCIs.
	// Carrier-typical default-bearer = 9 (non-GBR best-effort internet).
	// Pointer because... actually QCI 0 is reserved/invalid, so plain
	// uint32 with 0 = omit is safe and matches the gx_fwd_ccri convention.
	QoSClassIdentifier uint32 `json:"qos_class_identifier,omitempty"`

	// Allocation-Retention-Priority (1034) — Grouped.
	AllocationRetentionPriority *AllocationRetentionPriority `json:"allocation_retention_priority,omitempty"`
}

// AllocationRetentionPriority mirrors TS 29.272 §7.3.40 Grouped AVP (1034).
// Reused by Gx (TS 29.212 §5.3.45). Layout matches gx.AllocationRetentionPriority
// but lives here to keep S6a self-contained.
type AllocationRetentionPriority struct {
	// Priority-Level (1046) — Unsigned32, 1..15 (1=highest, 15=lowest).
	// Carrier-typical default-bearer ARP = 9.
	PriorityLevel uint32 `json:"priority_level"`

	// Pre-emption-Capability (1047) — Enumerated. 0=ENABLED, 1=DISABLED.
	// Pointer because 0 is a valid value.
	PreEmptionCapability *uint32 `json:"pre_emption_capability,omitempty"`

	// Pre-emption-Vulnerability (1048) — Enumerated. 0=ENABLED, 1=DISABLED.
	// Pointer because 0 is a valid value.
	PreEmptionVulnerability *uint32 `json:"pre_emption_vulnerability,omitempty"`
}

// LCSInfo mirrors TS 29.272 §7.3.43 Grouped AVP (1473). Slim — full LCS
// schema is large; carrier-typical EPS attach omits LCS-Info entirely
// (location services are pulled on-demand via SLh, not pushed at attach).
// This struct is a placeholder so producers can populate it once needed
// without a schema migration.
type LCSInfo struct {
	GMLCNumber              []string              `json:"gmlc_number,omitempty"` // hex TBCD bytes, repeated
	LCSPrivacyExceptionList []LCSPrivacyException `json:"lcs_privacy_exception_list,omitempty"`
	MOLR                    []MOLR                `json:"mo_lr,omitempty"`
}

// LCSPrivacyException — TS 29.272 §7.3.45.
type LCSPrivacyException struct {
	SSCode               []byte  `json:"ss_code,omitempty"`   // AVP 1476  OctetString (TS 29.002)
	SSStatus             []byte  `json:"ss_status,omitempty"` // AVP 1477  OctetString (TS 29.002)
	NotificationToUEUser *uint32 `json:"notification_to_ue_user,omitempty"`

	ExternalClient []ExternalClient `json:"external_client,omitempty"` // AVP 1479
	PLMNClient     []uint32         `json:"plmn_client,omitempty"`     // AVP 1482
	ServiceType    []ServiceType    `json:"service_type,omitempty"`    // AVP 1483
}

// MOLR — TS 29.272 §7.3.49 placeholder.
type MOLR struct {
	SSCode   []byte `json:"ss_code,omitempty"`   // AVP 1476  OctetString
	SSStatus []byte `json:"ss_status,omitempty"` // AVP 1477  OctetString
}

// TeleserviceList — TS 29.272 §7.3.79. CS-domain teleservices the
// subscriber is provisioned for (legacy GSM/UMTS). Slim placeholder.
type TeleserviceList struct {
	TSCodes [][]byte `json:"ts_codes,omitempty"` // TS-Code (1487)  OctetString (TS 29.002)
}

// CSGSubscriptionData mirrors TS 29.272 §7.3.78 Grouped AVP (1436).
// Closed Subscriber Group memberships. Carrier-typical: only set for
// femtocell-enabled subscribers.
type CSGSubscriptionData struct {
	CSGId             uint32   `json:"csg_id"`
	ExpirationDate    int64    `json:"expiration_date,omitempty"`    // unix seconds, sent as a Time AVP; 0 = no expiry
	ServiceSelections []string `json:"service_selections,omitempty"` // associated APN names
}

// ProSeSubscriptionData mirrors TS 29.272 §7.3.137 Grouped AVP (3701).
// Proximity Services (Rel 12+). Slim placeholder.
type ProSeSubscriptionData struct {
	ProSePermission uint32 `json:"prose_permission,omitempty"` // bitmask
}
