package internal_test

import (
	"context"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// instantiateWapcHost instantiates a test waPC host and returns it and a cleanup function.
func instantiateWapcHost(t *testing.T, r wazero.Runtime) (*wapcHost, api.Closer) {
	h := &wapcHost{t: t}
	// Export host functions (in the order defined in https://wapc.io/docs/spec/#required-host-exports)
	if host, err := r.NewModuleBuilder("wapc").
		ExportFunction("__host_call", h.hostCall,
			"__host_call", "bind_ptr", "bind_len", "ns_ptr", "ns_len", "cmd_ptr", "cmd_len", "payload_ptr", "payload_len").
		ExportFunction("__console_log", h.consoleLog,
			"__console_log", "ptr", "len").
		ExportFunction("__guest_request", h.guestRequest,
			"__guest_request", "op_ptr", "ptr").
		ExportFunction("__host_response", h.hostResponse,
			"__host_response", "ptr").
		ExportFunction("__host_response_len", h.hostResponseLen,
			"__host_response_len").
		ExportFunction("__guest_response", h.guestResponse,
			"__guest_response", "ptr", "len").
		ExportFunction("__guest_error", h.guestError,
			"__guest_error", "ptr", "len").
		ExportFunction("__host_error", h.hostError,
			"__host_error", "ptr").
		ExportFunction("__host_error_len", h.hostErrorLen,
			"__host_error_len").
		Instantiate(testCtx, r); err != nil {
		t.Errorf("Error instantiating waPC host - %v", err)
		return h, nil
	} else {
		return h, host
	}
}

type wapcHost struct {
	t                  *testing.T
	consoleLogMessages []string
}

// hostCall is the WebAssembly function export "__host_call", which initiates a host using the callHandler using
// parameters read from linear memory (wasm.Memory).
func (w *wapcHost) hostCall(ctx context.Context, m api.Module, bindPtr, bindLen, nsPtr, nsLen, cmdPtr, cmdLen, payloadPtr, payloadLen uint32) int32 {
	panic("TODO")
}

// consoleLog is the WebAssembly function export "__console_log", which logs the message stored by the guest at the
// given offset (ptr) and length (len) in linear memory (wasm.Memory).
func (w *wapcHost) consoleLog(ctx context.Context, m api.Module, ptr, len uint32) {
	msg := w.requireReadString(ctx, m.Memory(), "msg", ptr, len)
	w.consoleLogMessages = append(w.consoleLogMessages, msg)
}

// guestRequest is the WebAssembly function export "__guest_request", which writes the invokeContext.operation and
// invokeContext.guestReq to the given offsets (opPtr, ptr) in linear memory (wasm.Memory).
func (w *wapcHost) guestRequest(ctx context.Context, m api.Module, opPtr, ptr uint32) {
	panic("TODO")
}

// hostResponse is the WebAssembly function export "__host_response", which writes the invokeContext.hostResp to the
// given offset (ptr) in linear memory (wasm.Memory).
func (w *wapcHost) hostResponse(ctx context.Context, m api.Module, ptr uint32) {
	panic("TODO")
}

// hostResponse is the WebAssembly function export "__host_response_len", which returns the length of the current host
// response from invokeContext.hostResp.
func (w *wapcHost) hostResponseLen(ctx context.Context) uint32 {
	panic("TODO")
}

// guestResponse is the WebAssembly function export "__guest_response", which reads invokeContext.guestResp from the
// given offset (ptr) and length (len) in linear memory (wasm.Memory).
func (w *wapcHost) guestResponse(ctx context.Context, m api.Module, ptr, len uint32) {
	panic("TODO")
}

// guestError is the WebAssembly function export "__guest_error", which reads invokeContext.guestErr from the given
// offset (ptr) and length (len) in linear memory (wasm.Memory).
func (w *wapcHost) guestError(ctx context.Context, m api.Module, ptr, len uint32) {
	panic("TODO")
}

// hostError is the WebAssembly function export "__host_error", which writes the invokeContext.hostErr to the given
// offset (ptr) in linear memory (wasm.Memory).
func (w *wapcHost) hostError(ctx context.Context, m api.Module, ptr uint32) {
	panic("TODO")
}

// hostError is the WebAssembly function export "__host_error_len", which returns the length of the current host error
// from invokeContext.hostErr.
func (w *wapcHost) hostErrorLen(ctx context.Context) uint32 {
	panic("TODO")
}

// requireReadString is a convenience function that casts requireRead
func (w *wapcHost) requireReadString(ctx context.Context, mem api.Memory, fieldName string, offset, byteCount uint32) string {
	return string(w.requireRead(ctx, mem, fieldName, offset, byteCount))
}

// requireRead is like api.Memory except that it panics if the offset and byteCount are out of range.
func (w *wapcHost) requireRead(ctx context.Context, mem api.Memory, fieldName string, offset, byteCount uint32) []byte {
	buf, ok := mem.Read(ctx, offset, byteCount)
	if !ok {
		w.t.Fatalf("out of memory reading %s", fieldName)
	}
	return buf
}
