/*****************************************************************************\
 * Copyright 2023 Lawrence Livermore National Security, LLC
 * (c.f. AUTHORS, NOTICE.LLNS, LICENSE)
 *
 * This file is part of the Flux resource manager framework.
 * For details, see https://github.com/flux-framework.
 *
 * SPDX-License-Identifier: LGPL-3.0
\*****************************************************************************/

package fluxmodule

/*
#include "resource/reapi/bindings/c/reapi_module.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/flux-framework/fluxion-go/pkg/types"
)

type (
	ReapiCtx    C.struct_reapi_module_ctx_t
	ReapiModule struct {
		ctx ReapiCtx
	}
)

// NewReapiModule creates a new resource API module
// reapi_module_ctx_t *reapi_module_new ();
func NewReapiModule() *ReapiModule {
	ctx := (*ReapiCtx)(C.reapi_module_new())
	return &ReapiModule{ctx: ctx}
}

// HasContext exposes the private ctx, telling the caller if it is set
func (m *ReapiModule) HasContext() bool {
	return m.ctx != nil
}

// Given an integer return code, convert to go error
// Also provide a meaningful string to the developer user
func retvalToError(code int, message string) error {
	if code == 0 {
		return nil
	}
	return fmt.Errorf(message+" %d", code)
}

// Destroy destroys the resource API context
// void reapi_module_destroy (reapi_module_ctx_t *ctx);
func (m *ReapiModule) Destroy() {
	C.reapi_module_destroy((*C.struct_reapi_module_ctx)(m.ctx))
}

// Match matches and allocates resources
// int reapi_module_match_allocate (reapi_module_ctx_t *ctx, char *match_op,
// at: is the scheduled time "at"
func (m *ReapiModule) Match(
	match_op types.MatchType,
	jobspec string,
	jobid int,
) (reserved bool, allocated string, at int64, overhead float64, err error) {
	// var atlong C.long = (C.long)(at)
	var r = C.CString("teststring")

	// Jobspec as a CString
	spec := C.CString(jobspec)

	fluxerr := (int)(C.reapi_module_match((*C.struct_reapi_module_ctx)(m.ctx),
		C.match_op_t(match_op),
		spec,
		(C.ulong)(jobid),
		(*C.bool)(&reserved),
		&r,
		(*C.long)(&at),
		(*C.double)(&overhead)))

	allocated = C.GoString(r)

	defer C.free(unsafe.Pointer(r))
	defer C.free(unsafe.Pointer(spec))

	err = retvalToError(fluxerr, "issue matching allocation for resource api")
	return reserved, allocated, at, overhead, err
}

// MatchAllocate matches and allocates resources
func (m *ReapiModule) MatchAllocate(
	orelse_reserve bool,
	jobspec string,
	jobid int,
) (reserved bool, allocated string, at int64, overhead float64, err error) {
	var match_op string

	if orelse_reserve {
		match_op = types.MatchAllocateOrElseReserve
	} else {
		match_op = types.MatchAllocate
	}

	return m.Match(match_op, jobspec, jobid)
}

// MatchSatisfy runs satisfiability check on jobspec
func (m *ReapiModule) MatchSatisfy(
	jobspec string,
	jobid int,
) (reserved bool, allocated string, at int64, overhead float64, err error) {
	match_op := types.MatchSatisfiability
	return m.Match(match_op, jobspec, jobid)
}

// Info gets the information on the allocation or reservation corresponding
// to jobid.
//
//	\param ctx       reapi_module_ctx_t context object
//	\param jobid     const jobid of the uint64_t type.
//	\param reserved  Boolean into which to return true if this job has been
//	                 reserved instead of allocated.
//	\param at        If allocated, 0 is returned; if reserved, actual time
//	                 at which the job is reserved.
//	\param ov        Double into which to return performance overhead
//	                 in terms of elapse time needed to complete
//	                 the match operation.
//	\return          0 on success; -1 on error.
//
// int reapi_module_info (reapi_module_ctx_t *ctx, const uint64_t jobid,
//
//	bool *reserved, int64_t *at, double *ov);
func (m *ReapiModule) Info(ctx *ReapiCtx, jobid int64) (reserved bool, at int64, overhead float64, err error) {
	fluxerr := (int)(C.reapi_module_info((*C.struct_reapi_module_ctx)(m.ctx),
		(C.ulong)(jobid),
		(*C.bool)(&reserved),
		(*C.long)(&at),
		(*C.double)(&overhead)))

	err = retvalToError(fluxerr, "issue getting module info")
	return reserved, at, overhead, err
}

// Cancel cancels the allocation or reservation corresponding to jobid.
//
//	\param ctx       reapi_module_ctx_t context object
//	\param jobid     jobid of the uint64_t type.
//	\param noent_ok  don't return an error on nonexistent jobid
//	\return          0 on success; -1 on error.
//
// int reapi_module_cancel (reapi_module_ctx_t *ctx,
//
//	const uint64_t jobid, bool noent_ok);*/
func (m *ReapiModule) Cancel(jobid int64, noent_ok bool) (err error) {
	fluxerr := (int)(C.reapi_module_cancel((*C.struct_reapi_module_ctx)(m.ctx),
		(C.ulong)(jobid),
		(C.bool)(noent_ok)))
	return retvalToError(fluxerr, "issue with cancel")
}
