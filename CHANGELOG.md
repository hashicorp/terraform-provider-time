# v0.4.0

BREAKING CHANGES:

* resource/time_offset: `keepers` argument renamed to `triggers`
* resource/time_rotating: `keepers` argument renamed to `triggers`
* resource/time_static: `keepers` argument renamed to `triggers`

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
