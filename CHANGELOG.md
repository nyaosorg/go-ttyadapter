Changelog
=========
( English | [Japanese](./CHANGELOG_ja.md) )

- tty8pe.Tty: Add optional prefix hook (OnPrefix, SetOnPrefix) for experimental input handling (#4)
- Rename release note files to CHANGELOG.md and CHANGELOG\_ja.md (#7)

v0.3.0
------
Jan 29, 2026

- Add `tty8pe` and `tty10pe` packages, which treat the Escape key as a prefix to prevent misbehavior from split ESC sequences.

v0.2.0
------
Nov 5, 2025

- Added new field `auto.Pilot.OnGetKey`. (#2)
   This allows users to adjust the auto-pilot speed by setting a callback function.

v0.1.0
------
Nov 3, 2025

### Breaking changes

- Added a `height` parameter to the callback function of the `Open` method. (#1)
  Originally, this package was developed as a subpackage of `nyaosorg/go-readline-ny`, where only the terminal width (`width`) was needed for single-line input.
  Since the package is now being redesigned as a standalone and more general-purpose library, it was necessary to include the terminal height (`height`) as well.

v0.0.1
------
Nov 1, 2025

- Initial release: separated functionality from github.com/nyaosorg/go-readline-ny
