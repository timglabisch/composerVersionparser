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

func wouldMatch(regex, value string) bool {
	r, _ := regexp.Compile(regex)
	return r.MatchString(value)
}

func match(regex, value string) (matched bool, values []string) {
	r, _ := regexp.Compile(regex)
	if !r.MatchString(value) {
		return false, []string{}
	}

	return true, r.FindStringSubmatch(value)
}

// NormalizeBranch - normalizing a branch name
func NormalizeBranch(v string) string {
	branch := strings.Trim(v, " ")

	if branch == "master" || branch == "trunk" || branch == "default" {
		_, normalized := Normalize(branch)
		return normalized
	}

	if matched, versions := match(`^v?(\d+)(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?$`, branch); matched { //r.MatchString(branch) {
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
func Normalize(v string) (ok bool, NormalizedVersion string) {

	version := strings.Trim(v, " ")

	// ignore aliases and just assume the alias is required instead of the source
	if matched, versions := match(`(?i)^(?:([^,\s]+)) +as +([^,\s]+)$`, version); matched {
		version = versions[1]
	}

	// ignore build metadata
	if matched, versions := match(`(?i)^(?:([^,\s+]+))\+[^\s]+$`, version); matched {
		version = versions[1]
	}

	// match master-like branches
	if wouldMatch(`(?i)^(?:dev-)?(?:master|trunk|default)$`, version) {
		return true, "9999999-dev"
	}

	if substr(strings.ToLower(version), 4) == "dev-" {
		return true, "dev-" + substrLast(version, len(version)-4)
	}

	// instead of index using directly the stability
	index := 0
	matched := false
	versions := []string{}

	if matched, versions = match(`(?i)^v?(\d{1,3})(\.\d+)?(\.\d+)?(\.\d+)?[._-]?(?:(stable|beta|b|RC|alpha|a|patch|pl|p)(?:[.-]?(\d+))?)?([.-]?dev)?$`, version); matched {
		version = versions[1] + ifNotEmpty(versions[2], ".0") + ifNotEmpty(versions[3], ".0") + ifNotEmpty(versions[4], ".0")
		index = 5
	} else if matched, versions = match(`(?i)^v?(\d{4}(?:[.:-]?\d{2}){1,6}(?:[.:-]?\d{1,3})?)[._-]?(?:(stable|beta|b|RC|alpha|a|patch|pl|p)(?:[.-]?(\d+))?)?([.-]?dev)?$`, version); matched {
		replace, _ := regexp.Compile(`\D`)
		version = replace.ReplaceAllString(versions[1], "-")
		index = 2
	} else if matched, versions = match(`(?i)^v?(\d{4,})(\.\d+)?(\.\d+)?(\.\d+)?[._-]?(?:(stable|beta|b|RC|alpha|a|patch|pl|p)(?:[.-]?(\d+))?)?([.-]?dev)?$`, version); matched {
		version = versions[1] + ifNotEmpty(versions[2], ".0") + ifNotEmpty(versions[3], ".0") + ifNotEmpty(versions[4], ".0")
		index = 5
	}

	// add version modifiers if a version was matched
	if matched && index != 0 {
		if versions[index] == "stable" {
			return true, version
		} else if versions[index] != "" {
			version = version + "-" + expandStability(versions[index]) + ifNotEmpty(versions[index+1], "")
		}

		if versions[index+2] != "" {
			version = version + "-dev"
		}

		return true, version
	}

	// match dev branches
	if matched, versions := match(`(?i)(.*?)[.-]?dev$`, version); matched {
		return true, NormalizeBranch(versions[1])
	}

	return false, ""
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
