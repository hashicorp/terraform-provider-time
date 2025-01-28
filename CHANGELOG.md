## 0.12.1 (September 11, 2024)

NOTES:

* all: This release introduces no functional changes. It does however include dependency updates which address upstream CVEs. ([#351](https://github.com/hashicorp/terraform-provider-time/issues/351))

## 0.12.0 (July 17, 2024)

ENHANCEMENTS:

* resource/time_static: If the `rfc3339` value is set in config and known at plan-time, all other attributes will also be known during plan. ([#255](https://github.com/hashicorp/terraform-provider-time/issues/255))

## 0.11.2 (May 28, 2024)

NOTES:

* This release introduces no functional changes. It does however include dependency updates which address upstream CVEs. ([#327](https://github.com/hashicorp/terraform-provider-time/issues/327))

## 0.11.1 (March 11, 2024)

NOTES:

* No functional changes from v0.11.0. Minor documentation fixes. ([#299](https://github.com/hashicorp/terraform-provider-time/issues/299))

## 0.11.0 (March 11, 2024)

FEATURES:

* functions/rfc3339_parse: Added a new `rfc3339_parse` function that parses an RFC3339 timestamp string and returns an object representation. ([#280](https://github.com/hashicorp/terraform-provider-time/issues/280))

## 0.10.0 (December 06, 2023)

BUG FIXES:

* resource/time_offset: Fix bug preventing multiple offset arguments from being set ([#189](https://github.com/hashicorp/terraform-provider-time/issues/189))

## 0.9.2 (November 28, 2023)

NOTES:

* This release introduces no functional changes. It does however include dependency updates which address upstream CVEs. ([#263](https://github.com/hashicorp/terraform-provider-time/issues/263))

## 0.9.1 (November 2, 2022)

BUG FIXES:

* resource/time_rotating: Correctly retrieve the Terraform state during update ([#132](https://github.com/hashicorp/terraform-provider-time/pull/132))

## 0.9.0 (October 11, 2022)

NOTES:

* provider: Rewritten to use the new [`terraform-plugin-framework`](https://www.terraform.io/plugin/framework) ([#112](https://github.com/hashicorp/terraform-provider-time/pull/112))

## 0.8.0 (August 10, 2022)

BUG FIXES:

* documentation: Changed wording from "Conflicts with other `offset_`/`rotation_` arguments." to "At least one of the `offset_`/`rotation_` arguments must be configured." to correctly reflect the use of `AtLeastOneOf` ([#105](https://github.com/hashicorp/terraform-provider-time/pull/105)) 

NOTES:

* provider: Upgrade Go version to 1.18 ([#114](https://github.com/hashicorp/terraform-provider-time/pull/114))
* provider: Enable `golangci-lint` ([#105](https://github.com/hashicorp/terraform-provider-time/pull/105))

## 0.7.2 (July 01, 2021)

BUG FIXES:

* resource/time_sleep: Prevent `context deadline exceeded` error when timeout duration is configured above 20 minutes ([#45](https://github.com/hashicorp/terraform-provider-time/issues/45))

## 0.7.1 (May 04, 2021)

BUG FIXES:

* provider: Ensure `darwin/arm64` platform is included in releases

## 0.7.0 (February 19, 2021)

Binary releases of this provider now include the darwin-arm64 platform. This version contains no further changes.

## 0.6.0 (October 04, 2020)

BREAKING CHANGES:

* Dropped support for Terraform 0.11 and lower ([#16](https://github.com/hashicorp/terraform-provider-time/issues/16))

ENHANCEMENTS

* Made `time_sleep` context aware, allowing easier early cancellation ([#16](https://github.com/hashicorp/terraform-provider-time/issues/16))

## 0.5.0 (May 13, 2020)

FEATURES

* **New Resource:** `time_sleep` ([#12](https://github.com/hashicorp/terraform-provider-time/issues/12))

# 0.4.0 (April 21, 2020)

BREAKING CHANGES:

* resource/time_offset: `keepers` argument renamed to `triggers`
* resource/time_offset: Remove non-RFC3339 RFC and `unixdate` attribute
* resource/time_rotating: `keepers` argument renamed to `triggers`
* resource/time_rotating: Remove non-RFC3339 RFC and `unixdate` attributes
* resource/time_static: `keepers` argument renamed to `triggers`
* resource/time_static: Remove non-RFC3339 RFC and `unixdate` attributes

# v0.3.0

ENHANCEMENTS:

* resource/time_offset: Add `keepers` argument
* resource/time_rotating: Add `keepers` argument
* resource/time_static: Add `keepers` argument

BUG FIXES:

* resource/time_offset: Ensure `base_rfc3339` is always set in Terraform state during creation, even if unconfigured

# v0.2.0

BREAKING CHANGES:

* resource/time_static: The `expiration_` arguments have been moved to the new `time_rotating` resource as `rotation_` arguments.

FEATURES:

* **New Resource:** `time_offset`
* **New Resource:** `time_rotating`

# v0.1.0

FEATURES:

* **New Resource:** `time_static`
