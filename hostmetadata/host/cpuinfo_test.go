/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. Licensed under the Apache License 2.0.
 * See the file "LICENSE" for details.
 */

package host

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadCPUInfo(t *testing.T) {
	info, err := readCPUInfo()
	require.NoError(t, err)
	assert.NotEmpty(t, info)
	// use package 0 if it exists, otherwise pick one that does exist
	packageId := 0
	if _, contains := info[key(keyCPUFlags)][0]; !contains {
		var k int = -1
		for k = range info[key(keyCPUFlags)] {
			break
		}
		if k == -1 {
			assert.Fail(t, "flags not found for any package")
		}
		packageId = k
	}

	assertions := map[string]func(t *testing.T){
		"FlagsAreSorted": func(t *testing.T) {
			assert.Contains(t, info[key(keyCPUFlags)], packageId)
			assert.True(t,
				sort.StringsAreSorted(strings.Split(info[key(keyCPUFlags)][packageId], ",")))
		},
		"ThreadsPerCore": func(t *testing.T) {
			assert.Contains(t, info[key(keyCPUThreadsPerCore)], packageId)
			assert.NotEmpty(t, info[key(keyCPUThreadsPerCore)][packageId])
		},
		"Caches": func(t *testing.T) {
			assert.Contains(t, info[key(keyCPUCacheL1i)], packageId)
			assert.Contains(t, info[key(keyCPUCacheL1d)], packageId)
			assert.Contains(t, info[key(keyCPUCacheL2)], packageId)
			assert.Contains(t, info[key(keyCPUCacheL3)], packageId)
			assert.NotEmpty(t, info[key(keyCPUCacheL1i)][packageId])
			assert.NotEmpty(t, info[key(keyCPUCacheL1d)][packageId])
			assert.NotEmpty(t, info[key(keyCPUCacheL2)][packageId])
			assert.NotEmpty(t, info[key(keyCPUCacheL3)][packageId])
		},
		"CachesIsANumber": func(t *testing.T) {
			assert.Contains(t, info[key(keyCPUCacheL1i)], packageId)
			_, err := strconv.Atoi(info[key(keyCPUCacheL1i)][packageId])
			require.NoError(t, err)
			assert.Contains(t, info[key(keyCPUCacheL3)], packageId)
			_, err = strconv.Atoi(info[key(keyCPUCacheL3)][packageId])
			require.NoError(t, err)
		},
		"NumCPUs": func(t *testing.T) {
			assert.Contains(t, info[key(keyCPUNumCPUs)], packageId)
			assert.NotEmpty(t, info[key(keyCPUNumCPUs)][packageId])
		},
		"CoresPerSocket": func(t *testing.T) {
			assert.Contains(t, info[key(keyCPUCoresPerSocket)], packageId)
			cps := info[key(keyCPUCoresPerSocket)][packageId]
			assert.NotEmpty(t, cps)
			i, err := strconv.Atoi(cps)
			require.NoErrorf(t, err, "%v must be parseable as a number", cps)
			assert.Greater(t, i, 0)
		},
		"OnlineCPUs": func(t *testing.T) {
			assert.Contains(t, info[key(keyCPUOnline)], packageId)
			onlines := info[key(keyCPUOnline)][packageId]
			assert.NotEmpty(t, onlines)
			ints, err := readCPURange(onlines)
			require.NoError(t, err)
			assert.NotEmpty(t, t, ints)
		},
	}
	for assertion, run := range assertions {
		t.Run(assertion, run)
	}
}

func TestOnlineCPUsFor(t *testing.T) {
	const siblings = `0-7`

	type args struct {
		coreIDs  []int
		expected string
	}
	tests := map[string]args{
		"One_CPU_Only":                 {[]int{3}, `3`},
		"A_Comma":                      {[]int{3, 5}, `3,5`},
		"A_Range":                      {[]int{0, 1, 2, 3}, `0-3`},
		"A_Range_And_Single":           {[]int{0, 1, 2, 5}, `0-2,5`},
		"Two_Ranges":                   {[]int{0, 1, 2, 5, 6, 7}, `0-2,5-7`},
		"Ranges_And_Commas":            {[]int{1, 2, 4, 6, 7}, `1-2,4,6-7`},
		"Multiple_Comma":               {[]int{1, 2, 4, 7}, `1-2,4,7`},
		"Multiple_Mixes_MultipleTimes": {[]int{0, 1, 3, 4, 6, 7}, `0-1,3-4,6-7`},
	}

	for name, test := range tests {
		c := test
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.expected, onlineCPUsFor(siblings, c.coreIDs))
		})
	}
}
