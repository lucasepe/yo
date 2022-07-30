package template

import (
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
)

type builtinFunction struct {
	fn      interface{}
	summary string
	usage   string
}

var functionsMap = map[string]builtinFunction{
	// Date functions
	"date": {
		fn:      date,
		summary: "Formats a date.",
		usage:   `now | date "2006-01-02"`,
	},
	"dateInZone": {
		fn:      dateInZone,
		summary: "Formats a date with a timezone.",
		usage:   `dateInZone "2006-01-02" (now) "UTC"`,
	},
	"dateModify": {
		fn:      dateModify,
		summary: "Modify a date.",
		usage:   `now | dateModify "-1.5h"`,
	},
	"duration": {
		fn:      duration,
		summary: "Formats seconds as a time.Duration.",
		usage:   "duration 95",
	},
	"now": {
		fn:      time.Now,
		summary: "The current date/time.",
		usage:   `now | date "2006-01-02"`,
	},
	"toDate": {
		fn:      toDate,
		summary: "Converts a string to a date.",
		usage:   `toDate "2006-01-02" "2017-12-31"`,
	},
	"unixEpoch": {
		fn:      unixEpoch,
		summary: "Returns the seconds since the unix epoch.",
		usage:   "now | unixEpoch",
	},

	// String functions
	"trim": {
		fn:      strings.TrimSpace,
		summary: "Removes space from either side of a string.",
		usage:   `trim "   hello    "`,
	},
	"trimSuffix": {
		fn:      func(a, b string) string { return strings.TrimSuffix(b, a) },
		summary: "Trim just the suffix from a string.",
		usage:   `trimSuffix "-" "hello-"`,
	},
	"trimPrefix": {
		fn:      func(a, b string) string { return strings.TrimPrefix(b, a) },
		summary: "Trim just the prefix from a string.",
		usage:   `trimPrefix "-" "-hello"`,
	},
	"upper": {
		fn:      strings.ToUpper,
		summary: "Convert the entire string to uppercase.",
		usage:   `upper "hello"`,
	},
	"lower": {
		fn:      strings.ToLower,
		summary: "Convert the entire string to lowercase.",
		usage:   `lower "HELLO"`,
	},
	"title": {
		fn:      strings.Title,
		summary: "Convert to title case.",
		usage:   `title "hello world"`,
	},

	// Switch order so that "foo" | repeat 5
	"repeat": {
		fn:      func(count int, str string) string { return strings.Repeat(str, count) },
		summary: "Repeat a string multiple times.",
		usage:   `repeat 3 "hello"`,
	},
	"substr": {
		fn:      substring,
		summary: "Get a substring from a string.",
		usage:   `substr 0 5 "hello world"`,
	},

	"contains": {
		fn:      func(substr string, str string) bool { return strings.Contains(str, substr) },
		summary: "Test to see if one string is contained inside of another.",
		usage:   `contains "cat" "catch"`,
	},
	"hasPrefix": {
		fn:      func(substr string, str string) bool { return strings.HasPrefix(str, substr) },
		summary: "Test whether a string has a given prefix.",
		usage:   `hasPrefix "cat" "catch"`,
	},
	"hasSuffix": {
		fn:      func(substr string, str string) bool { return strings.HasSuffix(str, substr) },
		summary: "Test whether a string has a given suffix.",
		usage:   `hasSuffix "tch" "catch"`,
	},

	"split": {
		fn:      func(sep, orig string) []string { return strings.Split(orig, sep) },
		summary: "Split a string into a list of strings.",
		usage:   `split "$" "foo$bar$baz$bar"`,
	},
	"join": {
		fn:      join,
		summary: "Join a list of strings into a single string, with the given separator.",
		usage:   `join "-" .Names`,
	},

	"replace": {
		fn:      replace,
		summary: "Perform simple string replacement.",
		usage:   `"I Am Henry VIII" | replace " " "-"`,
	},

	"regexMatch": {
		fn:      regexMatch,
		summary: "Test if the input string matches a regular expression.",
		usage:   `regexMatch "dog$" "bulldog"`,
	},
	"regexFindAll": {
		fn:      regexFindAll,
		summary: "Returns a slice of all matches of the regular expression.",
		usage:   `regexFindAll "[2,4,6,8]" "123456789" -1`,
	},
	"regexFind": {
		fn:      regexFind,
		summary: "Return the first match of the regular expression.",
		usage:   `regexFind "[a-zA-Z][1-9]" "abcd1234"`,
	},

	// OS:
	"env": {
		fn:      os.Getenv,
		summary: "Reads an environment variable.",
		usage:   `env "HOME"`,
	},
	"expandenv": {
		fn:      os.ExpandEnv,
		summary: "Substitute environment variables in a string.",
		usage:   `expandenv "Your path is set to $PATH"`,
	},
	"pathSep": {
		fn:      func() string { return string(os.PathSeparator) },
		summary: "Returns OS-specific path separator.",
		usage:   `pathSep`,
	},
	"pathListSep": {
		fn:      func() string { return string(os.PathListSeparator) },
		summary: "Returns OS-specific path list separator.",
		usage:   `pathSep`,
	},
	"tempDir": {
		fn:      os.TempDir,
		summary: "Returns the default directory to use for temporary files.",
		usage:   `tempDir`,
	},

	// UUIDs:
	"uuid": {
		fn:      uuidv4,
		summary: "Generate UUID v4 universally unique IDs.",
		usage:   `uuid`,
	},

	// Digests
	"sha1sum": {
		fn:      sha1sum,
		summary: "Computes the SHA1 digest of a specified string",
		usage:   `sha1sum "Hello world!"`,
	},
	"sha256sum": {
		fn:      sha256sum,
		summary: "Computes the SHA256 digest of a specified string",
		usage:   `sha256sum "Hello world!"`,
	},

	// Encoding:
	"b64enc": {
		fn:      base64encode,
		summary: "Encode or decode with Base64",
		usage:   `b64enc "Hello world!"`,
	},
	"b64dec": {
		fn:      base64decode,
		summary: "Decode with Base64",
		usage:   `b64dec "Hello world!"`,
	},

	// Defaults
	"default": {
		fn:      dfault,
		summary: "Set a simple default value.",
		usage:   `default "foo" .Bar`,
	},
	"coalesce": {
		fn:      coalesce,
		summary: "Takes a list of values and returns the first non-empty one.",
		usage:   `coalesce .name .parent.name "Scarlett"`,
	},

	// Counters and Randoms
	"incr": {
		fn:      func(i interface{}) int64 { return toInt64(i) + 1 },
		summary: "Increment an integer value by one.",
		usage:   `incr 7`,
	},
	"decr": {
		fn:      func(i interface{}) int64 { return toInt64(i) - 1 },
		summary: "Decrement an integer value by one.",
		usage:   `decr 9`,
	},
	"rand": {
		fn:      randInt,
		summary: "Returns a random integer value from min (inclusive) to max (exclusive).",
		usage:   `rand 8 16`,
	},

	// Lists
	"has": {
		fn:      has,
		summary: "Test to see if a list has a particular element.",
		usage:   `has 4 $myList`,
	},
	"uniq": {
		fn:      uniq,
		summary: "Generate a list with all of the duplicates removed.",
		usage:   `split "$" "foo$bar$baz$bar" | uniq`,
	},
}

// Names returns the builtin functions names.
func Names() []string {
	var res []string
	for k := range functionsMap {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

// Summary return the specific builtin function description.
func Summary(name string) string {
	res, ok := functionsMap[name]
	if ok {
		return res.summary
	}
	return ""
}

// Usage return the specific builtin function usage.
func Usage(name string) string {
	res, ok := functionsMap[name]
	if ok {
		return res.usage
	}
	return ""
}

// TxtFuncMap returns a 'text/template'.FuncMap
func TxtFuncMap() template.FuncMap {
	fm := make(map[string]interface{}, len(functionsMap))
	for k, v := range functionsMap {
		fm[k] = v.fn
	}

	return template.FuncMap(fm)
}
