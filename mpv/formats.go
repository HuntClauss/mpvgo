package mpv

// #include "utils.h"
// #include <stdlib.h>
import "C"
import (
	"encoding/binary"
	"unsafe"
)


func convert2Pointer(data interface{}, format Format) unsafe.Pointer {
	switch format {
	case FormatNone:
		break
	case FormatString, FormatOsdString:
		if data == nil {
			var result *C.char
			return unsafe.Pointer(&result)
		}
		return unsafe.Pointer(C.CString(data.(string)))
	case FormatFlag:
		if data == nil {
			var result C.int
			return unsafe.Pointer(&result)
		}
		result := C.int(0)
		if data.(bool) {
			result = 1
		}
		return unsafe.Pointer(&result)
	case FormatInt64:
		if data == nil {
			var result C.int64_t
			return unsafe.Pointer(&result)
		}
		result := C.int64_t(data.(int64))
		return unsafe.Pointer(&result)
	case FormatDouble:
		if data == nil {
			var result C.long
			return unsafe.Pointer(&result)
		}
		result := C.long(data.(float64))
		return unsafe.Pointer(&result)
	case FormatNode:
		if data == nil {
			var result C.mpv_node
			return unsafe.Pointer(&result)
		}
		// TODO
	case FormatNodeArray:
		if data == nil {
			var result C.mpv_node_list
			return unsafe.Pointer(&result)
		}

		result := C.mpv_node_list{}
		arr := data.(NodeList)

		result.num = C.int(len(arr))
		result.values = C.makeNodeList(result.num)
		result.keys = nil

		for i, v := range arr {
			if v.Data != nil {
				C.setNodeListElement(result.values, C.int(i), *v.CNode())
			}
		}

		return unsafe.Pointer(&result)
	case FormatNodeMap:
		if data == nil {
			var result C.mpv_node_list
			return unsafe.Pointer(&result)
		}

		result := C.mpv_node_list{}
		arr := data.(NodeMap)

		result.num = C.int(len(arr))
		result.values = C.makeNodeList(result.num)
		result.keys = C.makeStringArray(result.num)

		index := 0
		for k, v := range arr {
			C.setString(result.keys, C.int(index), C.CString(k))
			if v.Data != nil {
				C.setNodeListElement(result.values, C.int(index), *v.CNode())
			}
			index += 1
		}

		return unsafe.Pointer(&result)
	case FormatByteArray:
		if data == nil {
			var result C.mpv_byte_array
			return unsafe.Pointer(&result)
		}
	}
	return nil
}

func convert2Data(data interface{}, format Format) interface{} {
	switch format {
	case FormatNone:
		return nil
	case FormatString:
		if val, ok := data.(*C.char); ok {
			return C.GoString(val)
		} else
		if val, ok := data.(unsafe.Pointer); ok {
			return C.GoString((*C.char)(val))
		}
		val := binary.LittleEndian.Uint64(data.([]byte))
		return C.GoString((*C.char)(unsafe.Pointer(uintptr(val))))
	case FormatFlag:
		if val, ok := data.(C.int); ok {
			return val == 1
		} else
		if val, ok := data.(unsafe.Pointer); ok {
			return *(*C.int)(val) == 1
		}
		return data.(C.int) == 1
	case FormatInt64:
		if val, ok := data.(C.int64_t); ok {
			return int64(val)
		} else
		if val, ok := data.(unsafe.Pointer); ok { // Idk if this could even be unsafe.Pointer in real usage.
			return int64(*(*C.long)(val))
		}
		return int64(binary.LittleEndian.Uint64(data.([]byte)))
	case FormatDouble:
		if val, ok := data.(C.double); ok {
			return float64(val)
		} else
		if val, ok := data.(unsafe.Pointer); ok {
			return float64(*(*C.double)(val))
		}
		return nil
	case FormatNode:
		if val, ok := data.(C.mpv_node); ok {
			content := convert2Data(val.u[:], Format(val.format))
			return &Node{Data: content, Format: Format(val.format)}
		} else if val, ok := data.(*C.mpv_node); ok {
			content := convert2Data(val.u[:], Format(val.format))
			return &Node{Data: content, Format: Format(val.format)}
		} else if val, ok := data.(unsafe.Pointer); ok {
			cnode := *(*C.mpv_node)(val)
			content := convert2Data(cnode.u[:], Format(cnode.format))
			return &Node{Data: content, Format: Format(cnode.format)}
		}
		return nil
	case FormatNodeArray:
		ptr := binary.LittleEndian.Uint64(data.([]byte))
		arr := *(*C.mpv_node_list)(unsafe.Pointer(uintptr(ptr)))

		var cvalues *C.mpv_node = arr.values

		values := unsafe.Slice(cvalues, int(arr.num))
		result := make(NodeList, arr.num)

		for i, v := range values {
			val := convert2Data(v.u[:], Format(v.format))

			node := Node{Data: val, Format: Format(v.format)}
			result[i] = node
		}
		return result
	case FormatNodeMap:
		ptr := binary.LittleEndian.Uint64(data.([]byte))
		nodeMap := *(*C.mpv_node_list)(unsafe.Pointer(uintptr(ptr)))

		var ckeys **C.char = nodeMap.keys
		var cvalues *C.mpv_node = nodeMap.values

		keys := unsafe.Slice(ckeys, int(nodeMap.num))
		values := unsafe.Slice(cvalues, int(nodeMap.num))
		result := make(NodeMap, nodeMap.num)

		for i, v := range values {
			val := convert2Data(v.u[:], Format(v.format))

			node := Node{Data: val, Format: Format(v.format)}
			key := C.GoString((*C.char)(keys[C.int(i)]))
			result[key] = node
		}
		return result
	case FormatByteArray:
		// This is always address to value
		val := binary.LittleEndian.Uint64(data.([]byte))
		ptr := unsafe.Pointer(uintptr(val))
		tmp := *(*C.mpv_byte_array)(ptr)
		// In libmpv 'size' is int64, but C.GoBytes accept only int
		// Could be a problem in the future
		return C.GoBytes(tmp.data, C.int(tmp.size))
	}

	return nil
}