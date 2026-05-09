Changelog
=========
( English | [Japanese](./CHANGELOG_ja.md) )

- Refactor internal packages into "tty8base" and "tty10base" (#23)
- Improve Open/Close state management (#24)
  - Allow safe re-open after Close
  - Prevent partial open state on Open failure
  - Make Open idempotent
- Add IsOpen() to Tty interfaces and implementations (#25)

v0.6.3
------
May 9, 2026

- tty8, tty8pe: Add retry logic on EAGAIN/EWOULDBLOCK errors. Since go-tty v0.0.8, the terminal fd is set to non-blocking mode, which can cause read operations to return EAGAIN instead of blocking (notably on macOS). This change makes tty8 and tty8pe resilient against such errors without requiring a specific version of go-tty. (#21)

v0.6.2
------
May 3, 2026

- Pin go-tty to v0.0.7 for stable blocking behavior (#18)
- Add a subpackage "fav" equivalent to "tty8pe" for Windows with Go 1.20 or earlier, "tty10pe" for others. (#19)

v0.6.1
------
Apr 29, 2026

- Fix: "tty10" and "tty10pe" did not work when standard input was redirected (#16)

v0.6.0
------
Apr 27, 2026

- Update github.com/mattn/go-tty to v2.0.0 and adapt to the API changes (#13)
- Remove `println("a")` from "tty10" and "tty10pe" on UNIX systems.

v0.5.0
------
Apr 26, 2026

- auto.Pilot.Size() is now safe when neither Open is called nor Width/Height are explicitly initialized; it falls back to a default size (80,24,nil) instead of returning (0,0,nil). (#10)
- New `ttyhook` subpackage to wrap `Tty` and intercept `GetKey` calls. (#11)

v0.4.0
------
Mar 27, 2026

- `tty8pe.Tty` and `tty10pe.Tty`: Add optional prefix hook (`OnPrefix`, `SetOnPrefix`) for experimental input handling (#4, #9)
- Rename release note files to CHANGELOG.md and CHANGELOG\_ja.md (#7)
- Improve terminal resize detection on UNIX by using SIGWINCH instead of polling (#8)

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
