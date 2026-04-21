# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0](https://github.com/terry-xyz/groupie-tracker/compare/v0.0.7...v1.0.0) (2026-04-08)

### Added

- *(static/script.js)* persist filter state in sessionStorage to survive artist page navigation ([#37](https://github.com/terry-xyz/groupie-tracker/commit/4f5844f57573d51d7f93a43d32dfe88dc6ac4471))
- *(static/script.js)* replace page reload on reset with instant in-place clear and smooth slider animation ([#39](https://github.com/terry-xyz/groupie-tracker/commit/f6cca3854a1aa656a5a2aab5dd651beda27a0613))

### Fixed

- *(static)* move blue color from #applyFilters to #resetFilters and rename button to Reset Filters ([#36](https://github.com/terry-xyz/groupie-tracker/commit/f79da8ae5995db8e44ea51c35d2389ed3c919964))
- *(static/style.css)* match slider thumbs and Reset Filters to View Details color in light mode with smooth theme transition ([#38](https://github.com/terry-xyz/groupie-tracker/commit/30bac714482c007d0f808da6c1bf185010d733c6))
- *(static/script.js)* remove slider animation on reset, snap to defaults instantly ([#40](https://github.com/terry-xyz/groupie-tracker/commit/d4a31b106476f7f2575a3db1b67b677d57ba0564))
- *(handlers)* balance autocomplete suggestions ([#44](https://github.com/terry-xyz/groupie-tracker/commit/3e0c950bdee1c47fb852698dd4786f08d816d5c6))

### Documentation

- add explaining comments to all changed files per session review ([#41](https://github.com/terry-xyz/groupie-tracker/commit/387312afba71818b863e1d7dfe87a26dae7c7c39))
- *(lessons)* update lessons and README for session-2 changes ([#42](https://github.com/terry-xyz/groupie-tracker/commit/ab23913e060630d3ef9e22e6e479f5c62de25df7))
- removed lessons ([#43](https://github.com/terry-xyz/groupie-tracker/commit/450cc2c5ddaf8099aa4ece418b7e5f0fd69e7597))

## [0.0.7](https://github.com/terry-xyz/groupie-tracker/compare/v0.0.6...v0.0.7) (2026-03-21)

### Added

- *(static/script.js)* live search on typing and digits-only enforcement for year inputs ([#33](https://github.com/terry-xyz/groupie-tracker/commit/516ab9a3c54ef547498710db5ff0dfbb3bf170ee))

### Fixed

- *(handlers)* fix silent geocode parse errors, bad Content-Type on JSON errors, and year filter bypass on invalid input; cache error template ([#31](https://github.com/terry-xyz/groupie-tracker/commit/af916e300b149531d0ddd7890ad551a1ccf76408))
- *(handlers/search.go)* return empty JSON array instead of null when no artists match filters ([#32](https://github.com/terry-xyz/groupie-tracker/commit/cc5d7c2fd9922f01da2b39b257c960c992b35611))
- *(handlers)* format location names in autocomplete suggestions and match both raw and formatted location names in search ([#34](https://github.com/terry-xyz/groupie-tracker/commit/e54667d0337a29f33d2f44e26eb785f08167b29f))
- *(static/script.js)* fix bidirectional slider clamping and remove redundant Apply Filters button ([#35](https://github.com/terry-xyz/groupie-tracker/commit/67d17b02cf04864fa4105d83a1b198fdae1ea3b5))

## [0.0.6](https://github.com/terry-xyz/groupie-tracker/compare/v0.0.5...v0.0.6) (2026-03-21)

### Fixed

- *(handlers/artist)* format raw location keys and fix concert count pluralization ([#26](https://github.com/terry-xyz/groupie-tracker/commit/62565068f39710d3dae805f5c633ec69ed78226c))
- *(templates/artist.html)* adaptive map minZoom covers both axes to prevent grey bars on all screen orientations ([#27](https://github.com/terry-xyz/groupie-tracker/commit/93d8f0103ccb02609f7af852f3bea8c73c08d66f))

### Documentation

- *(handlers)* add explanatory comments across all handler files ([#29](https://github.com/terry-xyz/groupie-tracker/commit/ca2d2e8eee4c4198d2fbaa649e97987014497f48))
- *(README.md)* trim verbose sections to match idiomatic Go repo style ([#30](https://github.com/terry-xyz/groupie-tracker/commit/45c79857e57b92575b475e4142f58c31c54e4d69))

## [0.0.5](https://github.com/terry-xyz/groupie-tracker/compare/v0.0.4...v0.0.5) (2026-03-21)

### Fixed

- "sort by" select color correction ([#22](https://github.com/terry-xyz/groupie-tracker/commit/394f177c9c1f5ecfda35318eef6e15f5cad6b80c))
- *(static/style.css)* set dark text on select option elements for Windows browsers ([#23](https://github.com/terry-xyz/groupie-tracker/commit/8b7c0c16118928523a37281a53c5ededf230b19c))
- *(handlers)* log template.Execute errors instead of silently discarding them ([#25](https://github.com/terry-xyz/groupie-tracker/commit/04c419d91572fe1eb0e81845ac4c8b6dbcf99cd7))

### Documentation

- *(lessons)* update all lessons to reflect async geocoding and caching changes ([#21](https://github.com/terry-xyz/groupie-tracker/commit/3657adf4a3a5ad448d8b4f3b332d15aaaf8cb010))

## [0.0.4](https://github.com/terry-xyz/groupie-tracker/compare/v0.0.3...v0.0.4) (2026-03-20)

### Changed

- *(handlers/artist.go)* defer geocoding to async API call ([#17](https://github.com/terry-xyz/groupie-tracker/commit/c22d8a3922f5d04e759f41534b8bf8be8046a1b4))
- cache templates, add geo cache persistence and startup pre-warming ([#18](https://github.com/terry-xyz/groupie-tracker/commit/e8ab1646493b7dd93b0fce0d92bc226e84b8a68e))
- async map loading via /api/artist-geo endpoint ([#19](https://github.com/terry-xyz/groupie-tracker/commit/64434539117c0275e3fe88cd29f3989bc5eb38ab))

### Fixed

- *(templates/artist.html)* prevent grey areas and scroll lock at zoom limits ([#16](https://github.com/terry-xyz/groupie-tracker/commit/4c0a0cd2a0302dc8f062a2dac44733405f5fec85))
- *(templates/artist.html)* use singular "country" and "concert" when count is 1 ([#20](https://github.com/terry-xyz/groupie-tracker/commit/68d0a8be2755cdbec7acd17ab0e082f0e5c0cd77))

## [0.0.3](https://github.com/terry-xyz/groupie-tracker/compare/v0.0.2...v0.0.3) (2026-03-20)

### Added

- add filters, search bar, geolocalization, and visualizations extensions ([#12](https://github.com/terry-xyz/groupie-tracker/commit/d787cd91f7ec78487e4b741e0ec24a77f3a815ae))

### Changed

- transition members filter from checkboxes to sliders ([#13](https://github.com/terry-xyz/groupie-tracker/commit/1a0ee9e381da38cb21faf45e3701919a28e66af9))
- UI color changes on filters' inputs ([#14](https://github.com/terry-xyz/groupie-tracker/commit/d5479f43b4db3b2654412ff207121fa712615fe6))

### Fixed

- *(templates/artist.html)* fix map scroll zoom - smooth proportional zoom with no grey flash or after-image ([#15](https://github.com/terry-xyz/groupie-tracker/commit/44b230627ee814b846152068a6502f38e574e862))

### Documentation

- remove old documentation ([#11](https://github.com/terry-xyz/groupie-tracker/commit/2dfe77d46bbf6d590133f51bc26887600089f3e2))

## [0.0.2](https://github.com/terry-xyz/groupie-tracker/compare/v0.0.1...v0.0.2) (2026-03-12)

### Documentation

- Add documentation and license ([#6](https://github.com/terry-xyz/groupie-tracker/commit/65c0ead3cb7be3f005096c531f2b560a642c9b71))
- Add documentation and tasks ([#7](https://github.com/terry-xyz/groupie-tracker/commit/e3ef1f8f9ca29a861928706a56909f560fd6f702))
- Add comprehensive codebase learning guide ([#9](https://github.com/terry-xyz/groupie-tracker/commit/4b1068ab79ee0b2d3137f73d20724336ace283ef))
- add inline comments and Go doc comments across codebase ([#10](https://github.com/terry-xyz/groupie-tracker/commit/d7901ec86efacd23367e86caa1e86cf82fc035ff))

### Other

- Add .gitignore (ignore executables) ([#8](https://github.com/terry-xyz/groupie-tracker/commit/62706e8f4a96be86f39838c47e91001d27bdf389))

## 0.0.1 (2026-03-08)

### Added

- Add API client and data models ([#1](https://github.com/terry-xyz/groupie-tracker/commit/98aa87c76cc2c06e59770658456851bf89e93868))
- Add HTTP handlers and routing ([#2](https://github.com/terry-xyz/groupie-tracker/commit/826267b703963682430071910aea58b2d01f5e69))
- Add HTML templates for pages ([#3](https://github.com/terry-xyz/groupie-tracker/commit/a2fb3f62bb5bb737c22ca36039f0202cc5ea8c3c))
- Add glassmorphism UI styling ([#4](https://github.com/terry-xyz/groupie-tracker/commit/590c5ccad79aff2592870a9dc7bb6c5d3369ade2))
- Add search and filter functionality ([#5](https://github.com/terry-xyz/groupie-tracker/commit/865c3255f5efb8a9d1987c2c0a42163e40f792d6))
