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

func main() {
	jgfPtr := flag.String("jgf", "", "path to jgf")
	jobspecPtr := flag.String("jobspecs", "", "directory path of jobspecs")
	cancelRequestsPtr := flag.String("cancel", "", "directory to cancel requests")
	reserve := false
	flag.Parse()

	// This job is just for one node
	jgf, err := os.ReadFile(*jgfPtr)
	if err != nil {
		log.Panicf("Error reading JGF file")
	}

	// This job will request all 2 nodes so we can cancel one
	jobspec, err := os.ReadFile(filepath.Join(*jobspecPtr, "test001.yaml"))
	if err != nil {
		log.Panicf("Error reading jobspec file")
	}
	cancelJobspec, err := os.ReadFile(filepath.Join(*jobspecPtr, "test002.yaml"))
	if err != nil {
		log.Panicf("Error reading cancel jobspec file")
	}
	cancelRequest, err := os.ReadFile(filepath.Join(*cancelRequestsPtr, "partial-cancel.json"))
	if err != nil {
		log.Panicf("Error reading partial-cancel.json")
	}
	fullCancelRequest, err := os.ReadFile(filepath.Join(*cancelRequestsPtr, "cancel.json"))
	if err != nil {
		log.Panicf("Error reading partial-cancel.json")
	}
	fmt.Printf("Jobspec:\n %s\n", jobspec)

	cli := fluxcli.NewReapiClient()
	err = cli.InitContext(string(jgf), "{}")
	if err != nil {
		log.Panicf("Error initializing jobspec context for ReapiClient: %v\n", err)
	}
	fmt.Printf("Errors so far: %s\n", cli.GetErrMsg())

	reserved, allocated, at, overhead, jobid, err := cli.MatchAllocate(reserve, string(jobspec))
	if err != nil {
		log.Panicf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	reserved, allocated, at, overhead, jobid, err = cli.MatchAllocate(true, string(jobspec))
	fmt.Println("Errors so far: \n", cli.GetErrMsg())

	if err != nil {
		log.Panicf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	reserved, allocated, at, overhead, jobid, err = cli.Match(types.MatchAllocate, string(jobspec))
	fmt.Println("Errors so far: \n", cli.GetErrMsg())

	if err != nil {
		log.Panicf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	reserved, allocated, at, overhead, jobid, err = cli.Match(types.MatchAllocateOrElseReserve, string(jobspec))
	fmt.Println("Errors so far: \n", cli.GetErrMsg())

	if err != nil {
		log.Panicf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	reserved, allocated, at, overhead, jobid, err = cli.Match(types.MatchAllocateWithSatisfiability, string(jobspec))
	fmt.Println("Errors so far: \n", cli.GetErrMsg())

	if err != nil {
		log.Panicf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	sat, overhead, err := cli.MatchSatisfy(string(jobspec))
	fmt.Println("Errors so far: \n", cli.GetErrMsg())

	if err != nil {
		log.Panicf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printSatOutput(sat, err)

	err = cli.Cancel(1, false)
	if err != nil {
		log.Panicf("Error in ReapiClient Cancel: %v\n", err)
	}
	fmt.Printf("Cancel output: %v\n", err)

	reserved, at, overhead, mode, err := cli.Info(1)
	if err != nil {
		log.Panicf("Error in ReapiClient Info: %v\n", err)
	}
	fmt.Printf("Info output jobid 1: %t, %d, %f, %s, %v\n", reserved, at, overhead, mode, err)

	reserved, at, overhead, mode, err = cli.Info(2)
	if err != nil {
		log.Panicf("Error in ReapiClient Info: %v\n", err)
	}
	fmt.Printf("Info output jobid 2: %t, %d, %f, %v\n", reserved, at, overhead, err)

	// Finally, request two nodes
	fmt.Printf("Match allocate for 2 nodes (all of graph resources)")
	reserved, allocated, at, overhead, cancelJobid, err := cli.MatchAllocate(reserve, string(cancelJobspec))
	if err != nil {
		log.Panicf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, cancelJobid, err)

	// And partial cancel one
	totalCancel, err := cli.PartialCancel(int64(cancelJobid), string(cancelRequest), false)
	if err != nil {
		log.Panicf("Error in ReapiClient PartialCancel: %v\n", err)
	}
	fmt.Printf("Partial cancel output: total_removal: %v, err: %v\n", totalCancel, err)

	// This should now fail
	fmt.Printf("Match allocate for 2 nodes (all of graph resources) should now fail")
	reserved, allocated, at, overhead, jobid, err = cli.MatchAllocate(false, string(cancelJobspec))
	if err != nil {
		log.Panicf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)
	if allocated != "" {
		log.Panicf("MatchAllocate should have failed, not enough resources.")
	}
	err = cli.Cancel(int64(cancelJobid), false)
	if err != nil {
		log.Panicf("Error in ReapiClient cancel for partially cancelled job: %v\n", err)
	}
	fmt.Println(fullCancelRequest)
}

func printOutput(reserved bool, allocated string, at int64, jobid uint64, err error) {
	fmt.Println("\n\t----Match Allocate output---")
	fmt.Printf("jobid: %d\nreserved: %t\nallocated: %s\nat: %d\nerror: %v\n", jobid, reserved, allocated, at, err)
}

func printSatOutput(sat bool, err error) {
	fmt.Println("\n\t----Match Satisfy output---")
	fmt.Printf("satisfied: %t\nerror: %v\n", sat, err)
}
