package composerVersionparser

import "regexp"
import "strings"

//import "fmt"

func substrLast(input string, lenght int) string {
	if len(input) < lenght {
		return input
	}

	return input[len(input)-lenght : len(input)]
}

func substr(input string, lenght int) string {
	if len(input) <= lenght {
		return input
	}

	return input[0:lenght]
}

func ifNotEmpty(value, defaultValue string) string {
	if value != "" {
		return value
	}

	return defaultValue
}

func expandStability(stability string) string {

	v := strings.ToLower(stability)

	if v == "a" {
		return "alpha"
	} else if v == "b" {
		return "beta"
	} else if v == "pl" || v == "p" {
		return "patch"
	} else if v == "rc" {
		return "RC"
	}

	return stability

}

// NormalizeBranch - normalizing a branch name
func NormalizeBranch(v string) string {
	branch := strings.Trim(v, " ")

	if branch == "master" || branch == "trunk" || branch == "default" {
		return Normalize(branch)
	}

	r, _ := regexp.Compile(`^v?(\d+)(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?$`)
	if r.MatchString(branch) {

		versions := r.FindStringSubmatch(branch)

		version := ""

		for i := 1; i < 5; i++ {

			if versions[i] != "" {

				x := versions[i]
				x = strings.Replace(versions[i], "*", "x", -1)
				x = strings.Replace(x, "X", "x", -1)

				version = version + x
			} else {
				version = version + ".x"
			}

		}

		return strings.Replace(version, ".x", ".9999999", -1) + "-dev"
	}

	return "dev-" + branch
}

// Normalize - normalizing a version
func Normalize(v string) string {

	version := strings.Trim(v, " ")

	// ignore aliases and just assume the alias is required instead of the source
	r, _ := regexp.Compile(`(?i)^(?:([^,\s]+)) +as +([^,\s]+)$`)
	if r.MatchString(version) {
		versions := r.FindStringSubmatch(version)
		version = versions[1]
	}

	// ignore build metadata
	r2, _ := regexp.Compile(`(?i)^(?:([^,\s+]+))\+[^\s]+$`)
	if r2.MatchString(version) {
		versions := r2.FindStringSubmatch(version)
		version = versions[1]
	}

	// match master-like branches
	r3, _ := regexp.Compile(`(?i)^(?:dev-)?(?:master|trunk|default)$`)
	if r3.MatchString(version) {
		return "9999999-dev"
	}

	if substr(strings.ToLower(version), 4) == "dev-" {

		return "dev-" + substrLast(version, len(version)-4)
	}

	// instead of index using directly the stability
	index := 0
	versions := []string{}

	// match classical versioning
	r4, _ := regexp.Compile(`(?i)^v?(\d{1,3})(\.\d+)?(\.\d+)?(\.\d+)?[._-]?(?:(stable|beta|b|RC|alpha|a|patch|pl|p)(?:[.-]?(\d+))?)?([.-]?dev)?$`)
	if r4.MatchString(version) {
		versions = r4.FindStringSubmatch(version)
		version = versions[1] + ifNotEmpty(versions[2], ".0") + ifNotEmpty(versions[3], ".0") + ifNotEmpty(versions[4], ".0")
		index = 5
	} else {
		r5, _ := regexp.Compile(`(?i)^v?(\d{4}(?:[.:-]?\d{2}){1,6}(?:[.:-]?\d{1,3})?)[._-]?(?:(stable|beta|b|RC|alpha|a|patch|pl|p)(?:[.-]?(\d+))?)?([.-]?dev)?$`)
		if r5.MatchString(version) {
			versions = r5.FindStringSubmatch(version)
			replace, _ := regexp.Compile(`\D`)
			version = replace.ReplaceAllString(versions[1], "-")
			index = 2
		} else {
			r6, _ := regexp.Compile(`(?i)^v?(\d{4,})(\.\d+)?(\.\d+)?(\.\d+)?[._-]?(?:(stable|beta|b|RC|alpha|a|patch|pl|p)(?:[.-]?(\d+))?)?([.-]?dev)?$`)
			if r6.MatchString(version) {
				versions = r6.FindStringSubmatch(version)
				version = versions[1] + ifNotEmpty(versions[2], ".0") + ifNotEmpty(versions[3], ".0") + ifNotEmpty(versions[4], ".0")
				index = 5
			}
		}
	}

	// add version modifiers if a version was matched
	if index != 0 {
		if versions[index] != "" {
			if versions[index] == "stable" {
				return version
			}

			// expand ...
			version = version + "-" + expandStability(versions[index]) + ifNotEmpty(versions[index+1], "")
		}

		if versions[index+2] != "" {
			version = version + "-dev"
		}

		return version
	}

	// match dev branches
	r7, _ := regexp.Compile(`(?i)(.*?)[.-]?dev$`)
	if r7.MatchString(version) {
		versions = r7.FindStringSubmatch(version)
		return NormalizeBranch(versions[1])

	}

	return "SOME KIND OF EXCEPTION"
}

// ParseStability returns a stability by string
func ParseStability(version string) string {

	r, _ := regexp.Compile("(?:([^#]*))")
	v := r.FindString(version)

	if "dev-" == substr(v, 4) || "-dev" == substrLast(v, 4) {
		return "dev"
	}

	r2, _ := regexp.Compile(`[._-]?(?:(stable|beta|b|rc|alpha|a|patch|pl|p)(?:[.-]?(\d+))?)?([.-]?dev)?$`)

	for _, xmatch := range r2.FindAllStringSubmatch(strings.ToLower(v), 30) {

		for _, match := range xmatch {
			if match == "beta" || match == "b" {
				return "beta"
			}

			if match == "alpha" || match == "a" {
				return "alpha"
			}

			if match == "rc" {
				return "RC"
			}

		}

	}

	return "stable"

}
