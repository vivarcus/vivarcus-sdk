#!/usr/bin/env bash
# Map platform image tags v{ADCV}-{Assembly} to Go module tags.
#
# Platform tag (Release/binary): v26R3.3-13316
# Go module tag (go.mod / go get):  v1.26.3-3.13316
#
# github.com/vivarcus/vivarcus-sdk has no /vN import path suffix, so semver major
# must be 0 or 1. We use v1.{year}.{release}-{patch}[.{maint}...].{assembly}.

adcv_to_go_module_tag() {
	local adcv="$1"
	local assembly="$2"
	local year release patch rest prerelease m
	year="" release="" patch="" rest=""
	if [[ "$adcv" =~ ^([0-9]+)R([0-9]+)\.([0-9]+)(.*)$ ]]; then
		year="${BASH_REMATCH[1]}"
		release="${BASH_REMATCH[2]}"
		patch="${BASH_REMATCH[3]}"
		rest="${BASH_REMATCH[4]}"
	else
		return 1
	fi
	prerelease="$patch"
	if [ -n "$rest" ]; then
		rest="${rest#.}"
		local seg
		for seg in ${rest//./ }; do
			[ -n "$seg" ] || continue
			prerelease+=".${seg}"
		done
	fi
	prerelease+=".${assembly}"
	printf 'v1.%s.%s-%s' "$year" "$release" "$prerelease"
}

is_go_module_image_tag() {
	local tag="$1"
	[[ "$tag" =~ ^v1\.[0-9]+\.[0-9]+-[0-9A-Za-z.-]+$ ]]
}

parse_platform_image_tag() {
	local tag="${1#v}"
	_PARSE_ADCV=""
	_PARSE_ASSEMBLY=""
	if [[ "$tag" =~ ^(.+)-([0-9]+)$ ]] && [[ "${BASH_REMATCH[1]}" == *R* ]]; then
		_PARSE_ADCV="${BASH_REMATCH[1]}"
		_PARSE_ASSEMBLY="${BASH_REMATCH[2]}"
		return 0
	fi
	return 1
}

# v26R3.3-13316 -> v1.26.3-3.13316 ; v1.26.3-3.13316 -> unchanged
image_tag_to_go_module_tag() {
	local tag="$1"
	if is_go_module_image_tag "$tag"; then
		printf '%s\n' "$tag"
		return 0
	fi
	if parse_platform_image_tag "$tag"; then
		adcv_to_go_module_tag "$_PARSE_ADCV" "$_PARSE_ASSEMBLY"
		return 0
	fi
	return 1
}

latest_platform_image_tag() {
	local repo_root="$1"
	local best="" best_asm=0 t
	[ -n "$repo_root" ] || return 1
	git -C "$repo_root" rev-parse --git-dir >/dev/null 2>&1 || return 1
	while IFS= read -r t; do
		[ -n "$t" ] || continue
		parse_platform_image_tag "$t" || continue
		if [ "${_PARSE_ASSEMBLY}" -gt "$best_asm" ]; then
			best_asm="${_PARSE_ASSEMBLY}"
			best="$t"
		fi
	done < <(git -C "$repo_root" tag -l 'v*-*' 2>/dev/null || true)
	[ -n "$best" ] || return 1
	printf '%s\n' "$best"
}

resolve_sdk_module_ref() {
	local tag_root="${1:-}"
	if [ -n "${VIVARCUS_SDK_MODULE_REF:-}" ]; then
		printf '%s\n' "$VIVARCUS_SDK_MODULE_REF"
		return 0
	fi
	if [ -n "${VIVARCUS_PLATFORM_TAG:-}" ]; then
		printf '%s\n' "$VIVARCUS_PLATFORM_TAG"
		return 0
	fi
	if [ -z "$tag_root" ] && [ -n "${VIVARCUS_REPO_ROOT:-}" ]; then
		tag_root="$VIVARCUS_REPO_ROOT"
	fi
	local latest
	if latest="$(latest_platform_image_tag "$tag_root")"; then
		printf '%s\n' "$latest"
		return 0
	fi
	printf 'main\n'
}

# Prints the require line version for customer gosdk/go.mod (Go semver tag).
resolve_sdk_module_version() {
	local repo_root="${1:-}"
	local ref go_tag resolved
	ref="$(resolve_sdk_module_ref "$repo_root")"
	if go_tag="$(image_tag_to_go_module_tag "$ref")"; then
		printf '%s\n' "$go_tag"
		return 0
	fi
	resolved="$(GOPROXY="${GOPROXY:-direct}" go list -m "github.com/vivarcus/vivarcus-sdk@${ref}" 2>/dev/null | awk '{print $2}')" || true
	if [ -n "$resolved" ]; then
		printf '%s\n' "$resolved"
		return 0
	fi
	return 1
}

normalize_sdk_module_list_ref() {
	local ref="$1"
	local go_tag
	if go_tag="$(image_tag_to_go_module_tag "$ref")"; then
		printf '%s\n' "$go_tag"
		return 0
	fi
	printf '%s\n' "$ref"
}
