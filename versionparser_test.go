package composerVersionparser

import "testing"

func Test_substr(t *testing.T) {

	if substr("test", 3) != "tes" {
		t.Error("tes should be " + substr("test", 3))
	}

	if substr("test", 4) != "test" {
		t.Error("test should be " + substr("test", 4))
	}

	if substr("test", 5) != "test" {
		t.Error("test should be " + substr("test", 5))
	}

	if substr("test", 1) != "t" {
		t.Error("t should be " + substr("test", 1))
	}

	if substr("test", 10) != "test" {
		t.Error("test should be " + substr("test", 10))
	}
}

func Test_substrLast(t *testing.T) {

	if substrLast("test", 3) != "est" {
		t.Error("est should be " + substrLast("test", 3))
	}

	if substrLast("test", 1) != "t" {
		t.Error("t should be " + substrLast("test", 1))
	}

	if substrLast("test", 4) != "test" {
		t.Error("test should be " + substrLast("test", 4))
	}

	if substrLast("test", 5) != "test" {
		t.Error("test should be " + substrLast("test", 5))
	}

	if substrLast("test", 15) != "test" {
		t.Error("test should be " + substrLast("test", 15))
	}
}

func Test_Normalized(t *testing.T) {
	branchNames := [][]string{

		[]string{"1.0.0", "1.0.0.0"},                             // none
		[]string{"1.2.3.4", "1.2.3.4"},                           // none/2
		[]string{"1.0.0RC1dev", "1.0.0.0-RC1-dev"},               // parses state
		[]string{"1.0.0-rC15-dev", "1.0.0.0-RC15-dev"},           // CI parsing
		[]string{"1.0.0.RC.15-dev", "1.0.0.0-RC15-dev"},          // delimiters
		[]string{"1.0.0-rc1", "1.0.0.0-RC1"},                     // RC uppercase
		[]string{"1.0.0.pl3-dev", "1.0.0.0-patch3-dev"},          // patch replace
		[]string{"1.0-dev", "1.0.0.0-dev"},                       // forces w.x.y.z
		[]string{"0", "0.0.0.0"},                                 // forces w.x.y.z/2
		[]string{"10.4.13-beta", "10.4.13.0-beta"},               // parses long
		[]string{"10.4.13beta2", "10.4.13.0-beta2"},              // parses long/2
		[]string{"10.4.13beta.2", "10.4.13.0-beta2"},             // parses long/semver
		[]string{"10.4.13-b", "10.4.13.0-beta"},                  // expand shorthand
		[]string{"10.4.13-b5", "10.4.13.0-beta5"},                // expand shorthand2
		[]string{"v1.0.0", "1.0.0.0"},                            // strips leading v
		[]string{"v20100102", "20100102"},                        // strips v/datetime
		[]string{"2010.01", "2010-01"},                           // parses dates y-m
		[]string{"2010.01.02", "2010-01-02"},                     // parses dates w/ .
		[]string{"2010-01-02", "2010-01-02"},                     // parses dates w/ -
		[]string{"2010-01-02.5", "2010-01-02-5"},                 // parses numbers
		[]string{"2010.1.555", "2010.1.555.0"},                   // parses dates y.m.Y
		[]string{"20100102-203040", "20100102-203040"},           // parses datetime
		[]string{"20100102203040-10", "20100102203040-10"},       // parses dt+number
		[]string{"20100102-203040-p1", "20100102-203040-patch1"}, // parses dt+patch
		[]string{"dev-master", "9999999-dev"},                    // parses master
		[]string{"dev-trunk", "9999999-dev"},                     // parses trunk
		[]string{"1.x-dev", "1.9999999.9999999.9999999-dev"},     // parses branches
		[]string{"dev-feature-foo", "dev-feature-foo"},           // parses arbitrary
		[]string{"DEV-FOOBAR", "dev-FOOBAR"},                     // parses arbitrary2
		[]string{"dev-feature/foo", "dev-feature/foo"},           // parses arbitrary3
		[]string{"dev-master as 1.0.0", "9999999-dev"},           // ignores aliases
		[]string{"dev-master+foo.bar", "9999999-dev"},            // semver metadata
		[]string{"1.0.0-beta.5+foo", "1.0.0.0-beta5"},            // semver metadata/2
		[]string{"1.0.0+foo", "1.0.0.0"},                         // semver metadata/3
		[]string{"1.0.0+foo as 2.0", "1.0.0.0"},                  // metadata w/ alias

	}

	for _, d := range branchNames {
		if Normalize(d[0]) != d[1] {
			t.Error(d[0] + " is " + Normalize(d[0]) + " should be " + d[1])
		}
	}
}

func Test_ParseSability(t *testing.T) {

	stabilities := [][]string{
		[]string{"stable", "1.0"},
		[]string{"stable", "3.2.1"},
		[]string{"stable", "v3.2.1"},
		[]string{"dev", "v2.0.x-dev"},
		[]string{"stable", "1"},
		[]string{"dev", "v2.0.x-dev"},
		[]string{"dev", "v2.0.x-dev"},
		[]string{"dev", "v2.0.x-dev#abc123"},
		[]string{"dev", "v2.0.x-dev#trunk/@123"},
		[]string{"RC", "3.0-RC2"},
		[]string{"dev", "dev-master"},
		[]string{"dev", "3.1.2-dev"},
		[]string{"stable", "3.1.2-pl2"},
		[]string{"stable", "3.1.2-patch"},
		[]string{"alpha", "3.1.2-alpha5"},
		[]string{"beta", "3.1.2-beta"},
		[]string{"beta", "2.0B1"},
		[]string{"alpha", "1.2.0a1"},
		[]string{"alpha", "1.2_a1"},
		[]string{"RC", "2.0.0rc1"},
	}

	for _, d := range stabilities {
		if ParseStability(d[1]) != d[0] {
			t.Error(d[1] + " is " + ParseStability(d[1]) + " should be " + d[0])
		}
	}

}

func Test_successfulNormalizedBranches(t *testing.T) {
	branchNames := [][]string{
		[]string{"v1.x", "1.9999999.9999999.9999999-dev"},
		[]string{"v1.*", "1.9999999.9999999.9999999-dev"},
		[]string{"v1.0", "1.0.9999999.9999999-dev"},
		[]string{"2.0", "2.0.9999999.9999999-dev"},
		[]string{"v1.0.x", "1.0.9999999.9999999-dev"},
		[]string{"v1.0.3.*", "1.0.3.9999999-dev"},
		[]string{"v2.4.0", "2.4.0.9999999-dev"},
		[]string{"2.4.4", "2.4.4.9999999-dev"},
		[]string{"master", "9999999-dev"},
		[]string{"trunk", "9999999-dev"},
		[]string{"feature-a", "dev-feature-a"},
		[]string{"FOOBAR", "dev-FOOBAR"},
	}

	for _, d := range branchNames {
		if NormalizeBranch(d[0]) != d[1] {
			t.Error(d[0] + " is " + NormalizeBranch(d[0]) + " should be " + d[1])
		}
	}
}
