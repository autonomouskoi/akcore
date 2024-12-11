package bus

import "github.com/extism/go-pdk"

//go:wasmimport extism:host/user wait_for_reply
func WaitForReply(busMessage uint64, timeoutMS uint64) uint64

func MarshalArg(msg *BusMessage) (pdk.Memory, error) {
	b, err := msg.MarshalVT()
	if err != nil {
		return pdk.Memory{}, err
	}
	mem := pdk.AllocateBytes(b)
	return mem, nil
}

func UnmarshalReturn(offs uint64) (*BusMessage, error) {
	mem := pdk.FindMemory(offs)
	defer mem.Free()
	msg := &BusMessage{}
	err := msg.UnmarshalVT(mem.ReadBytes())
	return msg, err
}
