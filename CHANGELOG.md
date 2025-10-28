# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0](https://github.com/mas2020-golang/raptor/releases/tag/v0.4.0) - 2025-xx-xx

### Added
- add `nav` command to open secret URLs and copy passwords to the clipboard (#93)
- the commands are grouped into categories (#124)

### Changed
- change the default pwd lenght for the `create password` cmd to 14 (#111)
- refactor secret lookup to reuse shared helper across commands
- The `get secret` command has been replaced by `get` command

### Fixed
- prevent secret edit from exiting when the target secret is missing (#93)

## [0.2.0](https://github.com/mas2020-golang/raptor/releases/tag/v0.2.0) - 2024-02-12

### Added

- 66816fd (#39) add the box name to the end of the secret ls command
- f783c94 (#50) add the CRYPTEX_FOLDER env variable to set a box folder other than the default one
- 8eda61c change version to 0.2.0-dev
- fdfafb5 feature/67 add the interactive box mode (#72)
- 3098de0 feature/77 install goreleaser and test it locally (#79)
- 2e678e6 feature/78 open a box giving the path (#78)
- 7416dc5 fix: print command interrupts the app when the given secret is not existing in the box (#76)
- 17c6229 fix: verbose mode gets reset after executing the print command (#74)
- 321c6bb refactor(box): replace proto buffer with the YAML box representation (#66)
- b69ea1e set version to 0.2.0
- dec8b01 use the goutils module for the cryptex output

## [0.1.0](https://github.com/mas2020-golang/raptor/releases/tag/v0.1.0-rc.1) - 2022-05-14

### Added
- Add the 'box create' command
- Add the 'box list' command
- Add the 'secret create' command
- Add the encryption layer to the box
- Add the 'secret list' command
- Get a sensitive data from a secret
- Add the '--items' flag to the 'secret ls' command
- Add the login field to the box
- Add the command 'secret print' to show information related to a secret

### Changed

### Removed

## Types of changes
- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.
