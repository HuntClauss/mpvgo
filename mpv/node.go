package mpv
import "C"

// #cgo LDFLAGS: -lmpv
// #include <mpv/client.h>
// #include <stdlib.h>
/*

//mpv_node* makeNodeList(int length) {
//	return calloc(sizeof(mpv_node), length);
//}
//
//void setNodeListElement(mpv_node* values, int index, mpv_node value) {
//	values[index] = value;
//}

 */
import "C"
import (
	"encoding/binary"
)

type Node struct {
	Data interface{}
	Format Format
}


func (n *Node) CNode() *C.mpv_node {
	ptr := convert2Pointer(n.Data, n.Format)
	if ptr == nil {
		return nil
	}

	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(uintptr(ptr)))

	result := &C.mpv_node{}
	result.u = buf
	result.format = C.mpv_format(n.Format)
	return result
}

type NodeList []Node
type NodeMap map[string]Node
