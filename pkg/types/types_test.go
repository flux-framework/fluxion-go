package types

/*****************************************************************************\
 * Copyright 2024 Lawrence Livermore National Security, LLC
 * (c.f. AUTHORS, NOTICE.LLNS, LICENSE)
 *
 * This file is part of the Flux resource manager framework.
 * For details, see https://github.com/flux-framework.
 *
 * SPDX-License-Identifier: LGPL-3.0
\*****************************************************************************/

import (
	"testing"
)

func TestToString(t *testing.T) {
	type test struct {
		description string
		input       MatchType
		expected    string
	}

	tests := []test{
		{description: "unknown", input: MatchUnknown, expected: "unknown"},
		{description: "allocate", input: MatchAllocate, expected: "allocate"},
		{description: "satisfiability", input: MatchSatisfiability, expected: "satisfiability"},
		{description: "allocate_orelse_reserve", input: MatchAllocateOrElseReserve, expected: "allocate_orelse_reserve"},
		{description: "allocate_with_satisfiability", input: MatchAllocateWithSatisfiability, expected: "allocate_with_satisfiability"},
		{description: "grow_allocate", input: MatchGrowAllocation, expected: "grow_allocate"},
	}
	for _, item := range tests {
		t.Run(item.description, func(t *testing.T) {
			value := item.input.String()
			t.Logf("got %s, want %s", value, item.expected)
			if item.expected != value {
				t.Errorf("got %s, want %s", value, item.expected)
			}
		})
	}
}

func TestAsInt(t *testing.T) {
	type test struct {
		description string
		input       MatchType
		expected    int
	}

	tests := []test{
		{description: "unknown", input: MatchUnknown, expected: 0},
		{description: "allocate", input: MatchAllocate, expected: 1},
		{description: "satisfiability", input: MatchSatisfiability, expected: 5},
		{description: "allocate_orelse_reserve", input: MatchAllocateOrElseReserve, expected: 3},
		{description: "allocate_with_satisfiability", input: MatchAllocateWithSatisfiability, expected: 2},
		{description: "grow_allocate", input: MatchGrowAllocation, expected: 4},
	}
	for _, item := range tests {
		t.Run(item.description, func(t *testing.T) {
			value := int(item.input)
			t.Logf("got %d, want %d", value, item.expected)
			if item.expected != value {
				t.Errorf("got %d, want %d", value, item.expected)
			}
		})
	}
}
