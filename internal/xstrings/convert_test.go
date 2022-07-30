// Copyright 2015 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package xstrings

import (
	"strings"
	"testing"
)

func TestToSnakeCaseAndToKebabCase(t *testing.T) {
	cases := _M{
		"HTTPServer":         "http_server",
		"_camelCase":         "_camel_case",
		"NoHTTPS":            "no_https",
		"Wi_thF":             "wi_th_f",
		"_AnotherTES_TCaseP": "_another_tes_t_case_p",
		"ALL":                "all",
		"_HELLO_WORLD_":      "_hello_world_",
		"HELLO_WORLD":        "hello_world",
		"HELLO____WORLD":     "hello____world",
		"TW":                 "tw",
		"_C":                 "_c",
		"http2xx":            "http_2xx",
		"HTTP2XX":            "http2_xx",
		"HTTP20xOK":          "http_20x_ok",
		"HTTP20xStatus":      "http_20x_status",
		"HTTP-20xStatus":     "http_20x_status",
		"a":                  "a",
		"Duration2m3s":       "duration_2m3s",
		"Bld4Floor3rd":       "bld4_floor_3rd",
		" _-_ ":              "_____",
		"a1b2c3d":            "a_1b2c3d",
		"A//B%%2c":           "a//b%%2c",

		"HTTP状态码404/502Error": "http_状态码404/502_error",
		"中文(字符)":              "中文(字符)",
		"混合ABCWords与123数字456": "混合_abc_words_与123_数字456",

		"  sentence case  ": "__sentence_case__",
		" Mixed-hyphen case _and SENTENCE_case and UPPER-case": "_mixed_hyphen_case__and_sentence_case_and_upper_case",
		"FROM CamelCase to snake/kebab-case":                   "from_camel_case_to_snake/kebab_case",

		"": "",
		"Abc\uFFFDE\uFFFDf\uFFFDd\uFFFD2\uFFFD00z\uFFFDZZ\uFFFDZZ": "abc_\uFFFDe\uFFFDf\uFFFDd_\uFFFD2\uFFFD00z_\uFFFDzz\uFFFDzz",
		"\uFFFD\uFFFD\uFFFD\uFFFD\uFFFD":                           "\uFFFD\uFFFD\uFFFD\uFFFD\uFFFD",
	}

	runTestCases(t, ToSnakeCase, cases)

	for k, v := range cases {
		cases[k] = strings.Replace(v, "_", "-", -1)
	}

	runTestCases(t, ToKebabCase, cases)
}

func TestToCamelCase(t *testing.T) {
	runTestCases(t, ToCamelCase, _M{
		"http_server":     "HttpServer",
		"_camel_case":     "_CamelCase",
		"no_https":        "NoHttps",
		"_complex__case_": "_Complex_Case_",
		" complex -case ": " Complex Case ",
		"all":             "All",
		"GOLANG_IS_GREAT": "GolangIsGreat",
		"GOLANG":          "Golang",
		"a":               "A",
		"好":               "好",

		"FROM CamelCase to snake/kebab-case": "FromCamelcaseToSnake/kebabCase",

		"": "",
	})
}

func TestSwapCase(t *testing.T) {
	runTestCases(t, SwapCase, _M{
		"swapCase": "SWAPcASE",
		"Θ~λa云Ξπ":  "θ~ΛA云ξΠ",
		"a":        "A",

		"": "",
	})
}

func TestFirstRuneToUpper(t *testing.T) {
	runTestCases(t, FirstRuneToUpper, _M{
		"hello, world!": "Hello, world!",
		"Hello, world!": "Hello, world!",
		"你好，世界！":        "你好，世界！",
		"a":             "A",

		"": "",
	})
}

func TestFirstRuneToLower(t *testing.T) {
	runTestCases(t, FirstRuneToLower, _M{
		"hello, world!": "hello, world!",
		"Hello, world!": "hello, world!",
		"你好，世界！":        "你好，世界！",
		"a":             "a",
		"A":             "a",

		"": "",
	})
}
