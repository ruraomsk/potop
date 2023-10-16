package radar

import (
	"math/rand"
	"sync"
	"time"

	"github.com/ruraomsk/potop/setup"
	"github.com/simonvetter/modbus"
)

type handler struct {
	lock   sync.Mutex
	uptime time.Time
	dates  [16]uint16
	reg16  [4]uint16
}

// reg16 to holdig
func (h *handler) unpack() {
	k := 0
	j := 0
	for i := 0; i < setup.Set.ModbusRadar.Chanels; i++ {
		h.dates[i] = (h.reg16[k] >> j) & 0xf
		j += 4
		if j > 12 {
			j = 0
			k++
		}
	}
}
func (h *handler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	err = modbus.ErrIllegalFunction
	return
}

func (h *handler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	err = modbus.ErrIllegalFunction
	return
}

func (h *handler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	err = nil
	h.lock.Lock()
	defer h.lock.Unlock()

	if req.UnitId != uint8(setup.Set.ModbusRadar.ID) {
		err = modbus.ErrIllegalFunction
		return
	}

	if int(req.Addr)+int(req.Quantity) > len(h.reg16) {
		err = modbus.ErrIllegalDataAddress
		return
	}

	if req.IsWrite {
		h.uptime = time.Now()
	}
	if !setup.Set.ModbusRadar.Debug {
		for i := 0; i < int(req.Quantity); i++ {
			if req.IsWrite {
				h.reg16[int(req.Addr)+i] = req.Args[i]
			}
			res = append(res, h.reg16[int(req.Addr)+i])
		}
	} else {
		for i := 0; i < int(req.Quantity); i++ {
			if !req.IsWrite {
				h.reg16[int(req.Addr)+i] = 0
				for l := 0; l < 4; l++ {
					h.reg16[int(req.Addr)+i] |= uint16((rand.Intn(3) << (l * 4)))
				}
			}
			res = append(res, h.reg16[int(req.Addr)+i])
		}
	}
	return
}

func (h *handler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	err = modbus.ErrIllegalFunction
	return
}
