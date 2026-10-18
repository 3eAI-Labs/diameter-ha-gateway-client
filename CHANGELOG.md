# Changelog

All notable changes to `diameter-ha-gateway-client` are documented in this file.
Format loosely follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versioning follows [SemVer](https://semver.org/) (pre-1.0: a minor version may be breaking).

## [Unreleased]

## [0.2.0] — 2026-10-18

### Added

All additions are new optional fields; no field changes name or type. Octet strings are hex.

- `S6aEvent`, from the MME's requests: `active_apns` (Active-APN whole, with Context-Identifier),
  `sms_register_request`, `sgs_mme_identity`, `coupled_node_diameter_id`, `adjacent_plmns` (ULR),
  `air_flags` (AIR), `eps_location_information` (PUR), `maximum_ue_availability_time` and
  `visited_network_identifier` (NOR); `terminal_info.meid`.
- `S6aEvent`, for the requests an application sends: `clr_flags` (CLR), `reset_ids` (IDR, RSR),
  `user_ids` (RSR, beyond `imsi`), `trace_reference`, `ts_codes`, `ss_codes` (DSR).
  `supported_features` and `drmp` are now sent in these requests too.
- `S6aAnswer`, in an application's ULA or AIA: `error_diagnostic`, `utran_vectors`,
  `geran_vectors`, `ue_usage_type`.
- `S6aAnswer`, in the gateway's reply to a CLR, IDR, DSR or RSR — the MME's answer, which until now
  was reduced to its result: `ida_flags`, `dsa_flags`, `ims_voice_over_ps_sessions_supported`,
  `last_ue_activity_time`, `rat_type`, `eps_user_state`, `eps_location_information`,
  `local_time_zone`, and its `supported_features` and `failed_avps`.
- Types `ActiveAPN`, `EPSLocationInformation`, `MMELocationInformation`,
  `SGSNLocationInformation`, `EPSUserState`, `LocalTimeZone`, `UTRANVector`, `GERANVector`.

### Changed

- `mip6_home_link_prefix` is hex: the prefix length octet and the 16-octet prefix (RFC 5447
  §4.2.4). The gateway wrote the string's own bytes, which cannot hold that binary value in JSON.
- From a NOR, `service_selection` is the NOR's own Service-Selection only; the gateway filled it
  with a ULR's first Active-APN's (see `active_apns`).

### Fixed

- The comment on `TerminalInfo.SV` called it the 3GPP2 MEID; it is the Software-Version.

## [0.1.0] — 2026-10-16

### Added

- First release of the wire contract as its own module: the S6a and S13 envelopes (`S6aEvent`,
  `S6aAnswer`, `S13Event`, `S13Answer`, `SubscriptionProfile` and its sub-types) and the subject
  constants. The types were previously in the gateway module, at
  `github.com/3eAI-Labs/diameter-ha-gateway/pkg/api/diameter/v1`; the package is now
  `github.com/3eAI-Labs/diameter-ha-gateway-client` (package `diameterclient`).
- Compared with the types in gateway 0.9.4: the extended Subscription-Data families (call barring,
  GPRS subscription data, IMSI group, adjacent access restriction, V2X, emergency info, eDRX,
  specific APN info, MIP6 agent info, WLAN offloadability, LCS detail, Rel-15 extended bandwidth).
- `TestWireSchema`: a golden file pinning every JSON field name and type.

### Changed

- **Breaking:** `S6aEvent.trace_data` is the full `TraceData` object (TS 29.272 §7.3.63); in 0.9.4
  it was `{"trace_reference": "<string>"}`. `trace_reference` is now base64-encoded octets.
- `S13Answer.experimental_result_code`: an EIR reports DIAMETER_ERROR_EQUIPMENT_UNKNOWN (5422) and
  other 3GPP results here; the gateway sends them as Experimental-Result, and sends
  `equipment_status` only with `result_code` 2001. S13 itself did not work in gateway 0.9.4 (the
  gateway could not parse an ECR).
- Replies to application-initiated commands carry `experimental_result_code` when the MME answered
  with an Experimental-Result; gateway 0.9.4 reported such an answer as `result_code` 2001.
- Subscription-Data checked against TS 29.272 Table 7.3.1/1 (gateway side: codes, M-bit and types
  on the wire). Changes to the types:
  - **Breaking:** `header_compression_configuration` is gone: TS 29.272 has no such AVP, and the
    codes the gateway used for it (1698-1702) are Service-Gap-Time, Broadcast-Location-Assistance-
    Data-Types, Paging-Time-Window and Operation-Mode.
  - **Breaking:** `edrx_subscription_data` is replaced by `edrx_cycle_length` (eDRX-Cycle-Length,
    §7.3.215, one per RAT type) and `paging_time_window` (Paging-Time-Window, §7.3.226, one per
    operation mode), the shape TS 29.272 defines. The old AVP (1711-1713) was not in the
    specification; those codes are UE-PC5-QoS and its children.
  - **Breaking:** `lcs_info.*.ss_code`, `ss_status` and `teleservice_list.ts_codes` are octets
    (base64), as TS 29.272 types SS-Code, SS-Status and TS-Code and as `call_barring_info` already
    had them; they were integers sent as Unsigned32.
  - **Breaking:** `S6aAnswer.equivalent_plmn_list` is gone. Equivalent-PLMN-List is part of the ULR
    only (§7.2.3); `S6aEvent.equivalent_plmn_list` now lists the Visited-PLMN-Ids inside the
    Grouped AVP (gateway 0.9.4 delivered the whole AVP body as one hex string).
  - `reset_ids` is sent on a ULA only, the one answer whose ABNF carries Reset-ID.
  - **Breaking:** `aese_communication_pattern` and `monitoring_event_configuration` are gone until
    they can be checked against TS 29.336. What the gateway sent broke the ABNF TS 29.272 gives
    them: AESE-Communication-Pattern needs a SCEF reference ID, which the type had no field for;
    Monitoring-Event-Configuration is 3122 (it was sent as 3142), and 3122 is what Monitoring-Type
    was sent as; over S6a it needs SCEF-Realm.
- Gateway 0.9.4 also subscribed to `diameter.s6a.pur` and `diameter.s6a.nor` to send PUR/NOR
  towards an MME; it no longer does. Both subjects carry an MME's request to the application only.

[Unreleased]: https://github.com/3eAI-Labs/diameter-ha-gateway-client/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/3eAI-Labs/diameter-ha-gateway-client/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/3eAI-Labs/diameter-ha-gateway-client/releases/tag/v0.1.0
