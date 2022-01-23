package mpv

// #include "utils.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type Mpv struct {
	ctx *C.mpv_handle
}

// Create creates new mpv instance and client API handle to control the mpv instance.
func Create() (*Mpv, error) {
	handle := C.mpv_create()
	if handle == nil {
		return nil, errors.New("cannot create mpv instance. Possible reasons: *out of memory* or *LC_NUMERIC != \"C\"*")
	}
	return &Mpv{ctx: handle}, nil
}

// CreateClient creates a new client handle connected to the same player core as current client.
func (m *Mpv) CreateClient(name string) (*Mpv, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	handle := C.mpv_create_client(m.ctx, cname)
	if handle == nil {
		return nil, errors.New("cannot create mpv client instance")
	}
	return &Mpv{ctx: handle}, nil
}

// CreateWeakClient creates weak handle reference.
//
// If all handles are weak references, core is automatically destroyed.
func (m *Mpv) CreateWeakClient(name string) (*Mpv, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	handle := C.mpv_create_weak_client(m.ctx, cname)
	if handle == nil {
		return nil, errors.New("cannot create mpv client instance")
	}
	return &Mpv{ctx: handle}, nil
}

// Initialize initializes uninitialized mpv instance
func (m *Mpv) Initialize() error {
	return Error(C.mpv_initialize(m.ctx)).Err()
}

// Destroy disconnects and destroys mpv handle
func (m *Mpv) Destroy() {
	C.mpv_destroy(m.ctx)
}

// Terminate terminates the player and all clients, and waits until all of them are destroyed
func (m *Mpv) Terminate() {
	C.mpv_terminate_destroy(m.ctx)
}

// ClientName returns the name of current client handle
func (m *Mpv) ClientName() string {
	return C.GoString(C.mpv_client_name(m.ctx))
}

// ClientID returns the ID of current client handle
func (m *Mpv) ClientID() int64 {
	return int64(C.mpv_client_id(m.ctx))
}

// LoadConfig loads and parse provided file.
//
// filename should be absolute path to file
func (m *Mpv) LoadConfig(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	return Error(C.mpv_load_config_file(m.ctx, cfilename)).Err()
}

// InternalTime returns internal time in microseconds.
// This has an arbitrary start offset,
// but will never wrap or go backwards.
func (m *Mpv) InternalTime() int64 {
	return int64(C.mpv_get_time_us(m.ctx))
}

func (m *Mpv) SetOption(name string, option interface{}, format Format) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	code := C.mpv_set_option(m.ctx, cname, C.mpv_format(format), convert2Pointer(option, format))
	return Error(code).Err()
}

func (m *Mpv) SetOptionString(name, option string) error {
	cname := C.CString(name)
	coption := C.CString(option)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(coption))

	code := C.mpv_set_option_string(m.ctx, cname, coption)
	return Error(code).Err()
}

func (m *Mpv) SetProperty(name string, property interface{}, format Format) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	code := C.mpv_set_property(m.ctx, cname, C.mpv_format(format), convert2Pointer(property, format))
	return Error(code).Err()
}

func (m *Mpv) SetPropertyString(name, property string) error {
	cname := C.CString(name)
	cproperty := C.CString(property)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(cproperty))

	code := C.mpv_set_property_string(m.ctx, cname, cproperty)
	return Error(code).Err()
}

func (m *Mpv) SetPropertyAsync(name string, property interface{}, id uint64, format Format) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	code := C.mpv_set_property_async(m.ctx, C.ulong(id), cname, C.mpv_format(format), convert2Pointer(property, format))
	return Error(code).Err()
}

func (m *Mpv) GetProperty(name string, format Format) (interface{}, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	result := convert2Pointer(nil, format)

	switch format {
	case FormatString:
		defer C.free(result)
	case FormatNode:
		defer C.mpv_free_node_contents((*C.mpv_node)(result))
	}

	code := C.mpv_get_property(m.ctx, cname, C.mpv_format(format), result)
	if code != 0 {
		return nil, Error(code).Err()
	}
	if result == nil {
		return nil, fmt.Errorf("returned value is nil. Probably format type is invalid")
	}

	return convert2Data(result, format), nil
}

func (m *Mpv) GetPropertyString(name string) (string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cresult := C.mpv_get_property_string(m.ctx, cname)

	if cresult == nil {
		return "", fmt.Errorf("cannot find property with provided name ('%s')", name)
	}
	defer C.mpv_free(unsafe.Pointer(cresult))
	return C.GoString(cresult), nil
}

func (m *Mpv) GetPropertyOsdString(name string) (string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cresult := C.mpv_get_property_osd_string(m.ctx, cname)

	if cresult == nil {
		return "", fmt.Errorf("cannot find property with provided name ('%s')", name)
	}
	defer C.mpv_free(unsafe.Pointer(cresult))
	return C.GoString(cresult), nil
}

func (m *Mpv) GetPropertyAsync(name string, id uint64, format Format) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	code := C.mpv_get_property_async(m.ctx, C.ulong(id), cname, C.mpv_format(format))
	return Error(code).Err()
}

func (m *Mpv) Command(args []string) error {
	array := C.makeStringArray(C.int(len(args)))

	for i, v := range args {
		cvalue := C.CString(v)
		defer C.free(unsafe.Pointer(cvalue))

		C.setString(array, C.int(i), cvalue)
	}

	return Error(C.mpv_command(m.ctx, array)).Err()
}

func (m *Mpv) CommandString(command string) error {
	ccmd := C.CString(command)
	defer C.free(unsafe.Pointer(ccmd))

	return Error(C.mpv_command_string(m.ctx, ccmd)).Err()
}

func (m *Mpv) CommandAsync(args []string, id uint64) error {
	array := C.makeStringArray(C.int(len(args)))

	for i, v := range args {
		cvalue := C.CString(v)
		defer C.free(unsafe.Pointer(cvalue))

		C.setString(array, C.int(i), cvalue)
	}

	code := C.mpv_command_async(m.ctx, C.ulong(id), array)
	return Error(code).Err()
}

func (m *Mpv) CommandNode(args *Node) (*Node, error) {
	cnode := args.CNode()

	var cresult *C.mpv_node = &C.mpv_node{}
	if code := C.mpv_command_node(m.ctx, cnode, cresult); code != 0 {
		return nil, Error(code).Err()
	}

	if cresult != nil {
		defer C.mpv_free_node_contents(cresult)
		node := convert2Data(cresult, FormatNode)
		return node.(*Node), nil
	}
	return nil, nil
}

func (m *Mpv) CommandAsyncNode(args *Node, id uint64) error {
	cnode := args.CNode()

	code := C.mpv_command_node_async(m.ctx, C.ulong(id), cnode)
	return Error(code).Err()
}

func (m *Mpv) CommandReturn(args []string) (*Node, error) {
	array := C.makeStringArray(C.int(len(args)))

	for i, v := range args {
		cvalue := C.CString(v)
		defer C.free(unsafe.Pointer(cvalue))

		C.setString(array, C.int(i), cvalue)
	}

	var cresult *C.mpv_node = &C.mpv_node{}
	code := C.mpv_command_ret(m.ctx, array, cresult)
	if code != 0 {
		return nil, Error(code).Err()
	}

	if cresult != nil {
		defer C.mpv_free_node_contents(cresult)
		node := convert2Data(cresult, FormatNode)
		return node.(*Node), nil
	}
	return nil, nil
}

func (m *Mpv) AbortAsyncCommand(id uint64) {
	C.mpv_abort_async_command(m.ctx, C.ulong(id))
}

func (m *Mpv) ObserveProperty(name string, id uint64, format Format) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	code := C.mpv_observe_property(m.ctx, C.ulong(id), cname, C.mpv_format(format))
	return Error(code).Err()
}

func (m *Mpv) UnObserveProperty(id uint64) (int, error) {
	num := C.mpv_unobserve_property(m.ctx, C.ulong(id))
	if num < 0 {
		return 0, Error(num).Err()
	}
	return int(num), nil
}

// ClientApiVersion returns version of compiled mpv
func ClientApiVersion() uint64 {
	return uint64(C.mpv_client_api_version())
}