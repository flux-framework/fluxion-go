package fluxcli

import (
	"fmt"
	"os"
	"testing"

	"github.com/flux-framework/fluxion-go/pkg/types"
)

func TestFluxcliContext(t *testing.T) {
	cli := NewReapiClient()
	if !cli.HasContext() {
		t.Errorf("Context is null")
	}

}

// NewClient creates a new client for testing
func NewClient(jgf string) (*ReapiClient, error) {
	cli := NewReapiClient()
	err := cli.InitContext(jgf, "{\"prune-filters\": \"ALL:core,ALL:gpu,ALL:memory\"}")
	if err != nil {
		fmt.Printf("Error initializing jobspec context for ReapiClient: %v\n", err)
		fmt.Printf("Errors so far: %s\n", cli.GetErrMsg())
	}
	return cli, err
}

func TestFind(t *testing.T) {

	// This job takes up the entire graph
	jgf, err := os.ReadFile("./data/find/tiny.json")
	if err != nil {
		t.Log("Error reading JGF file")
	}
	cli, err := NewClient(string(jgf))
	if err != nil {
		t.Errorf("creating new client %v\n", err)
	}

	// This takes up the entire graph
	jobspec, err := os.ReadFile("./data/find/full.yaml")
	if err != nil {
		t.Log("Error reading cancel initial jobspec file")
	}

	t.Log("allocate initial job to take up entire graph")
	reserved, allocated, at, _, initialJobid, err := cli.MatchAllocate(false, string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("matchAllocate reservation false %v\n", err)
	}
	printOutput(reserved, allocated, at, initialJobid, err)
	if allocated == "" {
		t.Error("find full allocation should succeed")
	}

	t.Log("find status=up")
	out, err := cli.Find("status=up")
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("find status=up %v\n", err)
	}
	t.Log(out)

}

func TestCancel(t *testing.T) {

	// This job takes up the entire graph
	// data/resource/jgfs/elastic/tiny-partial-cancel.json
	jgf, err := os.ReadFile("./data/cancel/tiny-partial-cancel.json")
	if err != nil {
		t.Log("Error reading JGF file")
	}
	cli, err := NewClient(string(jgf))
	if err != nil {
		t.Errorf("creating new client %v\n", err)
	}

	// This takes up the entire graph
	jobspec, err := os.ReadFile("./data/cancel/test022.yaml")
	if err != nil {
		t.Log("Error reading cancel initial jobspec file")
	}

	// We will ask for this job after requesting cancel (needs one node)
	cancelJobspec, err := os.ReadFile("./data/cancel/test023.yaml")
	if err != nil {
		t.Log("Error reading one node jobspec file")
	}

	// this requests to cancel 1/2 nodes
	cancelRequest, err := os.ReadFile("./data/cancel/node-1-partial-cancel.json")
	if err != nil {
		t.Log("Error reading partial cancel request file")
	}

	// Allocate the initial job - this takes up the entire graph
	t.Log("allocate initial job to take up entire graph")
	reserved, allocated, at, _, initialJobid, err := cli.MatchAllocate(false, string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("matchAllocate reservation false %v\n", err)
	}
	printOutput(reserved, allocated, at, initialJobid, err)
	if allocated == "" {
		t.Log("initial allocation failed (is empty)")
		t.Error("initial allocation failed (is empty)")
	}

	// Partial cancel 1/2 nodes
	t.Log("partial cancel 1/2 nodes")
	totalCancel, err := cli.PartialCancel(initialJobid, string(cancelRequest), false)
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("partial cancel one node %v\n", err)
	}
	t.Logf("Partial cancel output: total_removal: %v, err: %v\n", totalCancel, err)

	// The same request for 2 nodes should still fail
	t.Log("Match allocate for 2 nodes (all of graph resources) should now fail")
	reserved, allocated, at, _, jobid, err := cli.MatchAllocate(false, string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("matchAllocate 2 nodes when graph is 1/2 allocated: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)
	if allocated != "" {
		t.Errorf("matchAllocate should have failed, not enough resources.")
	}

	// We should be able to request one node
	t.Log("Match allocate for 1 node should succeed")
	reserved, allocated, at, _, oneNodeJobid, err := cli.MatchAllocate(false, string(cancelJobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("matchAllocate for one node after partial cancel: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)
	if allocated == "" {
		t.Error("matchAllocate for one node after partial cancel should have succeeded")
	}

	// Cancel should work of the initial job
	t.Log("cancel for partially cancelled job")
	err = cli.Cancel(initialJobid, false)
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("cancel for partially cancelled job: %v\n", err)
	}

	// And the one node job
	t.Log("cancel for one node job")
	err = cli.Cancel(oneNodeJobid, false)
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("cancel for single node job job: %v\n", err)
	}

	t.Log("Match allocate for 2 nodes (all of graph resources) should now succeed")
	reserved, allocated, at, _, jobid, err = cli.MatchAllocate(false, string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("matchAllocate 2 nodes when graph is empty: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)
	if allocated == "" {
		t.Errorf("matchAllocate should not have failed, we have two nodes again.")
	}

}

func TestMatchAllocate(t *testing.T) {

	// This job is just for one node
	jgf, err := os.ReadFile("./data/tiny.json")
	if err != nil {
		t.Error("Error reading JGF file")
	}
	cli, err := NewClient(string(jgf))
	if err != nil {
		t.Errorf("creating new client %v\n", err)
	}

	// This job is just for one node
	jobspec, err := os.ReadFile("./data/match/jobspec.yaml")
	if err != nil {
		t.Error("Error reading Jobspec file")
	}
	t.Logf("Jobspec:\n %s\n", jobspec)

	// Testing reservation false
	reserved, allocated, at, overhead, jobid, err := cli.MatchAllocate(false, string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("matchAllocate reservation false %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	// Testing reservation true
	reserved, allocated, at, overhead, jobid1, err := cli.MatchAllocate(true, string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("matchAllocate reservation true %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	// Same call, but directly to Match specifying the match operator type
	reserved, allocated, at, overhead, jobid, err = cli.Match(types.MatchAllocate, string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("match with match_op MatchAllocate %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	// Same, but with match_op match allocate or else reserve
	reserved, allocated, at, overhead, jobid, err = cli.Match(types.MatchAllocateOrElseReserve, string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("match with match_op MatchAllocateElseReserve %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	// Same, but with satisfy
	reserved, allocated, at, overhead, jobid, err = cli.Match(types.MatchAllocateWithSatisfiability, string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printOutput(reserved, allocated, at, jobid, err)

	// Just satisfy (boolean response)
	sat, overhead, err := cli.MatchSatisfy(string(jobspec))
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("Error in ReapiClient MatchAllocate: %v\n", err)
	}
	printSatOutput(sat, err)

	// Cancel of one of the jobid (not ok if does not exist)
	err = cli.Cancel(jobid1, false)
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("info for not-cancelled job %v\n", err)
	}

	// test info for cancelled job
	reserved, at, overhead, mode, err := cli.Info(jobid1)
	if err != nil {
		t.Logf("Fluxion errors %s\n", cli.GetErrMsg())
		t.Errorf("info for cancelled job %v\n", err)
	}
	t.Logf("Info output jobid 1: %t, %d, %f, %s, %v\n", reserved, at, overhead, mode, err)
}

func printOutput(reserved bool, allocated string, at int64, jobid int64, err error) {
	fmt.Println("\n\t----Match Allocate output---")
	isAllocated := allocated != ""
	fmt.Printf("jobid: %d\nreserved: %t\nallocated: %v\nat: %d\nerror: %v\n", jobid, reserved, isAllocated, at, err)
}

func printSatOutput(sat bool, err error) {
	fmt.Println("\n\t----Match Satisfy output---")
	fmt.Printf("satisfied: %t\nerror: %v\n", sat, err)
}
