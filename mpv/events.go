package mpv

// #include "utils.h"
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

type Event struct {
	ID      uint64
	Error   Error
	Data    interface{}
	EventID EventID
}

type EProperty struct {
	Name     string
	Format   Format
	Property interface{}
}

type ELogMessage struct {
	Prefix, Text string
	Level        LogLevel
}

type EClientMessage []string

type EStartFile int64

type EEndFile struct {
	PlaylistEntryID       int64
	PlaylistInsertID      int64
	PlaylistInsertEntries int
	Error                 Error
	Reason                EndFileReason
}

type EHook struct {
	Name string
	ID   uint64
}

type ECommandReply *Node


// RequestEvent
//
// status = true means enabled, otherwise disabled
func (m Mpv) RequestEvent(event EventID, status bool) error {
	var cstatus C.int = 0
	if status {
		cstatus = 1
	}

	code := C.mpv_request_event(m.ctx, C.mpv_event_id(event), cstatus)
	return Error(code).Err()
}

func (m Mpv) RequestLogMessages(level LogLevel) error {
	clevel := C.CString(level.String())
	defer C.free(unsafe.Pointer(clevel))

	code := C.mpv_request_log_messages(m.ctx, clevel)
	return Error(code).Err()
}

func (m Mpv) Wakeup() {
	C.mpv_wakeup(m.ctx)
}

func (m Mpv) SetWakeupCallback() {
	panic("[SetWakeupCallback] function not implemented.")
}

func (m Mpv) WaitAsyncRequests() {
	C.mpv_wait_async_requests(m.ctx)
}

func (m Mpv) HookAdd(name string, priority int, id uint64) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	code := C.mpv_hook_add(m.ctx, C.ulong(id), cname, C.int(priority))
	return Error(code).Err()
}

func (m Mpv) HookContinue(id uint64) error {
	code := C.mpv_hook_continue(m.ctx, C.ulong(id))
	return Error(code).Err()
}

func (m Mpv) EventWait(timeout float64) *Event {
	cevent := C.mpv_wait_event(m.ctx, C.double(timeout))
	return decodeCEvent(cevent)
}

func decodeCEvent(e *C.mpv_event) *Event {
	result := Event{ID: uint64(e.reply_userdata), Error: Error(e.error), Data: nil, EventID: EventID(e.event_id)}
	if result.Error == 0 && result.EventID == EventNone {
		return &result
	}

	switch result.EventID {
	case EventGetPropertyReply, EventPropertyChange:
		obj := *(*C.mpv_event_property)(e.data)
		format := Format(obj.format)

		resp := convert2Data(obj.data, format)
		
		result.Data = EProperty{
			Name: C.GoString(obj.name), 
			Format: format, 
			Property: resp,
		}
		return &result
	case EventLogMessage:
		obj := *(*C.mpv_event_log_message)(e.data)
		
		result.Data = ELogMessage{
			Prefix: C.GoString(obj.prefix), 
			Text: C.GoString(obj.text), 
			Level: LogLevel(C.int(obj.log_level)),
		}
		return &result
	case EventClientMessage:
		obj := *(*C.mpv_event_client_message)(e.data)
		length := int(C.int(obj.num_args))
		var args **C.char = obj.args

		tmp := unsafe.Slice(args, length)
		arr := make(EClientMessage, length)

		for i, v := range tmp {
			arr[i] = C.GoString(v)
		}
		result.Data = arr
		return &result
	case EventStartFile:
		obj := *(*C.mpv_event_start_file)(e.data)
		result.Data = EStartFile(obj.playlist_entry_id)
		return &result
	case EventEndFile:
		obj := *(*C.mpv_event_end_file)(e.data)
		result.Data = EEndFile{
			PlaylistEntryID:       int64(obj.playlist_entry_id),
			PlaylistInsertID:      int64(obj.playlist_insert_id),
			PlaylistInsertEntries: int(obj.playlist_insert_num_entries),
			Error:                 Error(obj.error),
			Reason:                EndFileReason(obj.reason),
		}
		return &result
	case EventHook:
		obj := *(*C.mpv_event_hook)(e.data)

		result.Data = EHook{
			Name: C.GoString(obj.name),
			ID: uint64(obj.id),
		}
		return &result
	case EventCommandReply:
		obj := *(*C.mpv_event_command)(e.data)
		resp := convert2Data(obj.result, FormatNode)
		result.Data = ECommandReply(resp.(*Node))
		return &result
	}

	return &result
}

func EventName(event EventID) string {
	return C.GoString(C.mpv_event_name(C.mpv_event_id(event)))
}
