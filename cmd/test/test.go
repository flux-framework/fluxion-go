/*****************************************************************************\
 * Copyright 2023 Lawrence Livermore National Security, LLC
 * (c.f. AUTHORS, NOTICE.LLNS, LICENSE)
 *
 * This file is part of the Flux resource manager framework.
 * For details, see https://github.com/flux-framework.
 *
 * SPDX-License-Identifier: LGPL-3.0
\*****************************************************************************/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/flux-framework/fluxion-go/pkg/fluxcli"
	"github.com/flux-framework/fluxion-go/pkg/types"
)

// readFile, panic and fail on error
func readFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading JGF file")
	}
	return string(content)
}

// checkErrors looks to see if the error is not an empty
// string. If not, we print it
func checkErrors(cli *fluxcli.ReapiClient) {
	msg := cli.GetErrMsg()
	if msg != "" {
		fmt.Printf("Errors so far: %s\n", msg)
	}
}

func main() {
	jgfPtr := flag.String("jgf", "", "path to jgf")
	jgfGrowPtr := flag.String("grow", "", "path to grow jgf")
	jobspecPtr := flag.String("jobspec", "", "path to jobspecs")
	flag.Parse()

	// JobSpecs for testing:
	// 1. Vanilla "hello world" from node down to core
	jobspec := readFile(filepath.Join(*jobspecPtr, "test001.yaml"))

	// 4 nodes, will work after grow but not before it
	growJobspec := readFile(filepath.Join(*jobspecPtr, "grow.yaml"))

	jgf := readFile(*jgfPtr)
	growJgf := readFile(*jgfGrowPtr)

	// Client creation
	cli := fluxcli.NewReapiClient()
	err := cli.InitContext(jgf, "{}")
	if err != nil {
		log.Fatalf("Error initializing jobspec context for ReapiClient: %v\n", err)
	}
	checkErrors(cli)

	fmt.Printf("     Jobspec:\n %s\n", jobspec)
	fmt.Printf("Grow Jobspec:\n %s\n", growJobspec)

	// Match Allocate without and with reservation
	reserve := false
	reserved, allocated, at, overhead, jobid, err := cli.MatchAllocate(reserve, jobspec)
	if err != nil {
		log.Fatalf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	reserve = true
	reserved, allocated, at, overhead, jobid, err = cli.MatchAllocate(reserve, jobspec)
	checkErrors(cli)
	if err != nil {
		log.Fatalf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	// Match derivatives
	reserved, allocated, at, overhead, jobid, err = cli.Match(types.MatchAllocate, jobspec)
	checkErrors(cli)
	if err != nil {
		log.Fatalf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	// ...or else Reserve
	reserved, allocated, at, overhead, jobid, err = cli.Match(types.MatchAllocateOrElseReserve, jobspec)
	checkErrors(cli)
	if err != nil {
		log.Fatalf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	// ...with Satisfy
	reserved, allocated, at, overhead, jobid, err = cli.Match(types.MatchAllocateWithSatisfiability, jobspec)
	checkErrors(cli)
	if err != nil {
		log.Fatalf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	// Match Satisfy
	sat, overhead, err := cli.MatchSatisfy(jobspec)
	checkErrors(cli)
	if err != nil {
		log.Fatalf("Error in ReapiClient MatchSatisfy: %v\n", err)
	}
	printSatOutput(sat, err)

	// Cancel
	err = cli.Cancel(1, false)
	if err != nil {
		log.Fatalf("Error in ReapiClient Cancel: %v\n", err)
	}
	fmt.Printf("Cancel output: %v\n", err)

	// Info
	reserved, at, overhead, mode, err := cli.Info(1)
	if err != nil {
		log.Fatalf("Error in ReapiClient Info: %v\n", err)
	}
	fmt.Printf("Info output jobid 1: %t, %d, %f, %s, %v\n", reserved, at, overhead, mode, err)

	reserved, at, overhead, mode, err = cli.Info(2)
	if err != nil {
		log.Fatalf("Error in ReapiClient Info: %v\n", err)
	}
	fmt.Printf("Info output jobid 2: %t, %d, %f, %v\n", reserved, at, overhead, err)

	// Test grow - we first will ask "tiny" for resources it doesn't have (4 nodes)
	// then we will grow, and the request should succeed.
	fmt.Println("Asking to MatchSatisfy 4 nodes (not possible)")
	sat, overhead, err = cli.MatchSatisfy(growJobspec)
	checkErrors(cli)
	if err != nil {
		log.Fatalf("Error in ReapiClient MatchSatisfy: %v\n", err)
	}
	printSatOutput(sat, err)

	// This should not be satisfied!
	if sat {
		log.Fatalf("Error in ReapiClient MatchSatisfy - asking for 4 nodes with only 2 should fail: %v\n", err)
	}

	// Add nodes: note the data file comes originally from:
	// https://github.com/milroy/flux-sched/tree/grow-api/t/data/resource/jgfs/elastic
	fmt.Println("🍔 Asking to Grow from 2 to 4 Nodes")
	err = cli.Grow(growJgf)
	if err != nil {
		log.Fatalf("Error in ReapiClient Grow: %s %s\n", err, cli.GetErrMsg())
	}
	fmt.Printf("Grow request return value: %v\n", err)

	// Ask again - it should succeed this time.
	fmt.Println("Asking to MatchSatisfy 4 nodes (now IS possible)")
	sat, overhead, err = cli.MatchSatisfy(growJobspec)
	checkErrors(cli)

	if err != nil {
		log.Fatalf("Error in ReapiClient MatchSatisfy: %v\n", err)
	}
	printSatOutput(sat, err)

	// This should not be satisfied!
	if !sat {
		log.Fatalf("Error in ReapiClient MatchSatisfy - asking for 4 nodes should now succeed: %v\n", err)
	}

	// Shrink (remove subgraph) for node2
	fmt.Println("🥕 Asking to Shrink from 4 to 3 Nodes")
	err = cli.Shrink("/tiny0/rack0/node2")
	if err != nil {
		log.Fatalf("Error in ReapiClient Shrink: %s %s\n", err, cli.GetErrMsg())
	}
	fmt.Printf("Shrink request return value: %v\n", err)

	fmt.Println("Asking to MatchSatisfy 4 nodes (again, not possible)")
	sat, overhead, err = cli.MatchSatisfy(growJobspec)
	checkErrors(cli)
	if err != nil {
		log.Fatalf("Error in ReapiClient MatchSatisfy: %v\n", err)
	}
	printSatOutput(sat, err)
	if sat {
		log.Fatalf("Error in ReapiClient MatchSatisfy - asking for 4 nodes with only 3 should fail: %v\n", err)
	}

	fmt.Println("Asking to ShrinkMulti by 2 nodes")
	err = cli.ShrinkMulti([]string{"/tiny0/rack0/node3", "/tiny0/rack0/node1"})
	if err != nil {
		log.Fatalf("Error in ReapiClient ShrinkMulti: %s %s\n", err, cli.GetErrMsg())
	}
	fmt.Printf("Shrink request return value: %v\n", err)
}

func printOutput(reserved bool, allocated string, at int64, jobid uint64, err error) {
	fmt.Println("\n\t----Match Allocate output---")
	fmt.Printf("jobid: %d\nreserved: %t\nallocated: %s\nat: %d\nerror: %v\n", jobid, reserved, allocated, at, err)
}

func printSatOutput(sat bool, err error) {
	fmt.Println("\n\t----Match Satisfy output---")
	fmt.Printf("satisfied: %t\nerror: %v\n", sat, err)
}
