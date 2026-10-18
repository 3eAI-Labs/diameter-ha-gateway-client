// Copyright 2026 3eAI Labs Ltd
// SPDX-License-Identifier: Apache-2.0

package diameterclient

// Sub-AVP families of the extended Subscription-Data (TS 29.272 §7.3.x and
// TS 29.336 §8.4).

// SpecificAPNInfo mirrors TS 29.272 §7.3.82 Grouped AVP (1472, Rel-10).
// Repeated within APNConfiguration to scope a particular service-selection
// (APN string) to a specific MIP6 / VisitedNetworkIdentifier override —
// used in DSMIPv6 + per-APN home agent allocation.
type SpecificAPNInfo struct {
	ServiceSelection         string         `json:"service_selection"`                    // AVP 493
	MIP6AgentInfo            *MIP6AgentInfo `json:"mip6_agent_info,omitempty"`            // AVP 486
	VisitedNetworkIdentifier string         `json:"visited_network_identifier,omitempty"` // AVP 600
}

// MIP6AgentInfo mirrors TS 29.272 §7.3.45 (and RFC 5447) Grouped AVP (486)
// for DSMIPv6 home agent assignment. Dictionary symbol asymmetry: AVP 334
// = `MIPHomeAgentAddress` (no "6"), AVP 348 = `MIPHomeAgentHost` (no "6"),
// AVP 125 = `MIP6HomeLinkPrefix` (with "6"), AVP 486 = `MIP6AgentInfo`.
// Field names keep "6" per TS 29.272 §7.3.45 spec readability; encoder
// resolves the AVP code by symbol regardless.
type MIP6AgentInfo struct {
	MIP6HomeAgentAddress []string `json:"mip6_home_agent_address,omitempty"` // AVP 334
	MIP6HomeAgentHost    *struct {
		DestinationHost  string `json:"destination_host"`  // AVP 293
		DestinationRealm string `json:"destination_realm"` // AVP 283
	} `json:"mip6_home_agent_host,omitempty"` // AVP 348
	// MIP6HomeLinkPrefix (AVP 125, RFC 5447 §4.2.4) is hex: the prefix
	// length octet followed by the 16-octet prefix.
	MIP6HomeLinkPrefix string `json:"mip6_home_link_prefix,omitempty"`
}

// WLANOffloadability mirrors TS 29.272 §7.3.123 Grouped AVP (1667, Rel-12).
// Per-APN WLAN offload permissions for E-UTRAN and UTRAN.
type WLANOffloadability struct {
	WLANOffloadabilityEUTRAN uint32 `json:"wlan_offloadability_eutran,omitempty"` // AVP 1668
	WLANOffloadabilityUTRAN  uint32 `json:"wlan_offloadability_utran,omitempty"`  // AVP 1669
}

// ExternalClient mirrors TS 29.272 §7.3.30 Grouped AVP (1479) — LCS detail
// expansion of LCSPrivacyException allowing per-client GMLC restrictions
// and notification-to-UE-user behavior.
type ExternalClient struct {
	ClientIdentity       string `json:"client_identity"`                   // AVP 1480
	GMLCRestriction      uint32 `json:"gmlc_restriction,omitempty"`        // AVP 1481
	NotificationToUEUser uint32 `json:"notification_to_ue_user,omitempty"` // AVP 1478
}

// ServiceType mirrors TS 29.272 §7.3.32 Grouped AVP (1483). LCS detail —
// service-type-coded permissions paired with per-type GMLC restriction.
type ServiceType struct {
	ServiceTypeCode      uint32 `json:"service_type_code"`                 // AVP 1484
	GMLCRestriction      uint32 `json:"gmlc_restriction,omitempty"`        // AVP 1481
	NotificationToUEUser uint32 `json:"notification_to_ue_user,omitempty"` // AVP 1478
}

// CallBarringInfo mirrors TS 29.272 §7.3.31 Grouped AVP (1488). CS-domain
// supplementary-service barring per SS-Code per TS 22.030. Repeated.
type CallBarringInfo struct {
	SSCode   []byte `json:"ss_code"`   // AVP 1476  OctetString
	SSStatus []byte `json:"ss_status"` // AVP 1477  OctetString
}

// TraceData mirrors TS 29.272 §7.3.63 Grouped AVP (1458): subscriber and
// equipment trace and Minimization-of-Drive-Tests (MDT) activation parameters
// (TS 32.422).
type TraceData struct {
	TraceReference        []byte            `json:"trace_reference"`                // AVP 1459  OctetString (3 octets per TS 32.422)
	TraceDepth            uint32            `json:"trace_depth"`                    // AVP 1462  Enumerated
	TraceNETypeList       []byte            `json:"trace_ne_type_list"`             // AVP 1463  OctetString bitstring (TS 32.422 §5.1)
	TraceInterfaceList    []byte            `json:"trace_interface_list,omitempty"` // AVP 1464  OctetString bitstring
	TraceEventList        []byte            `json:"trace_event_list"`               // AVP 1465  OctetString bitstring
	OMCID                 string            `json:"omc_id,omitempty"`               // AVP 1466
	TraceCollectionEntity string            `json:"trace_collection_entity"`        // AVP 1452  Address (IP of TCE)
	MDTConfiguration      *MDTConfiguration `json:"mdt_configuration,omitempty"`    // AVP 1622  Grouped (Rel-10)
}

// AreaScope mirrors TS 29.272 §7.3.138 Grouped AVP (1624). MDT geographic
// scope filter — each entity-id slice is multi (one wire AVP per entry).
type AreaScope struct {
	CellGlobalIdentity       [][]byte `json:"cell_global_identity,omitempty"`        // AVP 1604
	EUTRANCellGlobalIdentity [][]byte `json:"eutran_cell_global_identity,omitempty"` // AVP 1602
	RoutingAreaIdentity      [][]byte `json:"routing_area_identity,omitempty"`       // AVP 1605
	LocationAreaIdentity     [][]byte `json:"location_area_identity,omitempty"`      // AVP 1606
	TrackingAreaIdentity     [][]byte `json:"tracking_area_identity,omitempty"`      // AVP 1603
}

// MDTConfiguration mirrors TS 29.272 §7.3.136 Grouped AVP (1622). MDT
// measurement parameters; AreaScope is a nested Grouped (§7.3.138).
type MDTConfiguration struct {
	JobType                 uint32     `json:"job_type"`                             // AVP 1623  Enumerated
	AreaScope               *AreaScope `json:"area_scope,omitempty"`                 // AVP 1624  Grouped
	ListOfMeasurements      uint32     `json:"list_of_measurements"`                 // AVP 1625  Unsigned32 bitmask (TS 32.422 §5.1 — opaque)
	ReportingTrigger        uint32     `json:"reporting_trigger,omitempty"`          // AVP 1626
	ReportInterval          uint32     `json:"report_interval,omitempty"`            // AVP 1627
	ReportAmount            uint32     `json:"report_amount,omitempty"`              // AVP 1628
	EventThresholdRSRP      uint32     `json:"event_threshold_rsrp,omitempty"`       // AVP 1629
	EventThresholdRSRQ      uint32     `json:"event_threshold_rsrq,omitempty"`       // AVP 1630
	LoggingInterval         uint32     `json:"logging_interval,omitempty"`           // AVP 1631
	LoggingDuration         uint32     `json:"logging_duration,omitempty"`           // AVP 1632
	MeasurementPeriodLTE    uint32     `json:"measurement_period_lte,omitempty"`     // AVP 1655
	MeasurementPeriodUMTS   uint32     `json:"measurement_period_umts,omitempty"`    // AVP 1656
	CollectionPeriodRRMLTE  uint32     `json:"collection_period_rrm_lte,omitempty"`  // AVP 1657
	CollectionPeriodRRMUMTS uint32     `json:"collection_period_rrm_umts,omitempty"` // AVP 1658
	PositioningMethod       []byte     `json:"positioning_method,omitempty"`         // AVP 1659  OctetString
	MeasurementQuantity     []byte     `json:"measurement_quantity,omitempty"`       // AVP 1660  OctetString
	EventThresholdEvent1F   int32      `json:"event_threshold_event_1f,omitempty"`   // AVP 1661  Integer32
	EventThresholdEvent1I   int32      `json:"event_threshold_event_1i,omitempty"`   // AVP 1662  Integer32
}

// GPRSSubscriptionData mirrors TS 29.272 §7.3.71 Grouped AVP (1467). Used by
// converged HSS deployments that own 2G/3G + 4G subscription data in one
// realm — PDPContext list is the legacy GPRS analogue of APNConfiguration.
type GPRSSubscriptionData struct {
	CompleteDataListIncludedIndicator uint32       `json:"complete_data_list_included_indicator,omitempty"` // AVP 1468
	PDPContext                        []PDPContext `json:"pdp_context"`                                     // AVP 1469
}

// PDPContext mirrors TS 29.272 §7.3.71 child Grouped AVP (1469). Per-APN
// 2G/3G PDP context configuration; ExtPDPType/ExtPDPAddress are Rel-9 IPv4v6
// PDP support.
type PDPContext struct {
	ContextIdentifier          uint32   `json:"context_identifier"`                      // AVP 1423
	PDPType                    []byte   `json:"pdp_type"`                                // AVP 1470  OctetString — TS 29.060 PDP type
	PDPAddress                 []string `json:"pdp_address,omitempty"`                   // AVP 1227  Address (multi, TS 32.299)
	QoSSubscribed              []byte   `json:"qos_subscribed,omitempty"`                // AVP 1404  OctetString — TS 24.008 §10.5.6.5
	VPLMNDynamicAddressAllowed uint32   `json:"vplmn_dynamic_address_allowed,omitempty"` // AVP 1432
	ServiceSelection           string   `json:"service_selection"`                       // AVP 493   APN string
	ChargingCharacteristics    string   `json:"charging_characteristics,omitempty"`      // AVP 13
	ExtPDPType                 []byte   `json:"ext_pdp_type,omitempty"`                  // AVP 1620  Rel-9 IPv4v6
	ExtPDPAddress              []string `json:"ext_pdp_address,omitempty"`               // AVP 1621  Address (multi)
	AMBR                       *AMBR    `json:"ambr,omitempty"`                          // AVP 1435
	APNOIReplacement           string   `json:"apn_oi_replacement,omitempty"`            // AVP 1427
	SIPTOPermission            uint32   `json:"sipto_permission,omitempty"`              // AVP 1613
	LIPAPermission             uint32   `json:"lipa_permission,omitempty"`               // AVP 1618
	RestorationPriority        uint32   `json:"restoration_priority,omitempty"`          // AVP 1663
}

// IMSIGroupID mirrors TS 29.272 §7.3.122 Grouped AVP (1675, Rel-13). Group
// identification for service-bound IMSI sets; repeated.
type IMSIGroupID struct {
	GroupServiceID uint32 `json:"group_service_id"`         // AVP 1676
	GroupPLMNID    []byte `json:"group_plmn_id,omitempty"`  // AVP 1677  OctetString (TBCD PLMN)
	LocalGroupID   []byte `json:"local_group_id,omitempty"` // AVP 1678  OctetString
}

// AdjacentAccessRestrictionData mirrors TS 29.272 §7.3.137 Grouped AVP (1673,
// Rel-14). Per-VPLMN access-restriction overrides for adjacent PLMN scenarios.
type AdjacentAccessRestrictionData struct {
	VisitedPLMNID         []byte `json:"visited_plmn_id"`         // AVP 1407  OctetString (TBCD)
	AccessRestrictionData uint32 `json:"access_restriction_data"` // AVP 1426
}

// V2XSubscriptionData mirrors TS 29.272 §7.3.212 Grouped AVP (1688, Rel-14).
// Vehicle-to-everything (V2X) subscription permissions.
type V2XSubscriptionData struct {
	V2XPermission uint32 `json:"v2x_permission,omitempty"` // AVP 1689  Bitmask
}

// EmergencyInfo mirrors TS 29.272 §7.3.131 Grouped AVP (1687, Rel-14).
// Emergency-bearer-services subscription parameters.
type EmergencyInfo struct {
	MIP6AgentInfo *MIP6AgentInfo `json:"mip6_agent_info,omitempty"` // AVP 486
}

// EDRXCycleLength mirrors TS 29.272 §7.3.215 Grouped AVP (1691): the
// extended DRX cycle subscribed for one RAT type.
type EDRXCycleLength struct {
	RATType uint32 `json:"rat_type"` // AVP 1032 (TS 29.212), e.g. 1004 EUTRAN, 1005 EUTRAN-NB-IoT
	Value   []byte `json:"value"`    // AVP 1692 eDRX-Cycle-Length-Value, 1 octet (§7.3.216)
}

// PagingTimeWindow mirrors TS 29.272 §7.3.226 Grouped AVP (1701): the paging
// time window subscribed for one operation mode.
type PagingTimeWindow struct {
	OperationMode uint32 `json:"operation_mode"` // AVP 1702: 1 Iu, 2 WB-S1, 3 NB-S1 (§7.3.227)
	Length        []byte `json:"length"`         // AVP 1703 Paging-Time-Window-Length, 1 octet (§7.3.228)
}
