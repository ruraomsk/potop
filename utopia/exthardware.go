package utopia

import "github.com/ruraomsk/potop/hardware"

func GetStatusDirs() []uint8 {
	StateHardware := hardware.GetStateHard()
	result := make([]uint8, 0)
	var b uint8
	for _, v := range StateHardware.StatusDirs {
		switch v {
		case 0:
			//   OFF = 0, //все сигналы выключены
			b = 0xE
		case 1:
			//   DEACTIV_YELLOW=1, //направление перешло в неактивное состояние, желтый после зеленого
			b = 0x1
		case 2:
			//   DEACTIV_RED=2, //направление перешло в неактивное состояние, красный
			b = 0x0
		case 3:
			//   ACTIV_RED=3, //направление перешло в активное состояние, красный
			b = 0x0
		case 4:
			//   ACTIV_REDYELLOW=4, //направление перешло в активное состояние, красный c желтым
			b = 0x2
		case 5:
			//   ACTIV_GREEN=5, //направление перешло в активное состояние, зеленый
			b = 0x8
		case 6:
			//   UNCHANGE_GREEN=6, //направление не меняло свое состояние, зеленый
			b = 0x8
		case 7:
			//   UNCHANGE_RED=7, //направление не меняло свое состояние, красный
			b = 0x0
		case 8:
			//   GREEN_BLINK=8, //зеленый мигающий сигнал
			b = 0xA
		case 9:
			//   ZM_YELLOW_BLINK=9, //желтый мигающий в режиме ЖМ
			b = 0x9
		case 10:
			//   OS_OFF=10,	//сигналы выключены в режиме ОС
			b = 0xe
		case 11:
			//   UNUSED=11 //неиспользуемое направление
			b = 0xf
		default:
			b = 0xe
		}
		result = append(result, b)
	}
	return result
}

func GetStatusUtopia() byte {
	StateHardware := hardware.GetStateHard()
	if !StateHardware.Connect {
		return 6
	}
	if StateHardware.Dark {
		return 6
	}
	if StateHardware.Flashing {
		return 3
	}
	if StateHardware.AllRed {
		return 4
	}
	return byte(ctrl.status)
}
func GetDiagnosticUtopia() byte {
	StateHardware := hardware.GetStateHard()
	result := byte(0)
	//bit01 (0x02): “communication error”
	if !StateHardware.Connect {
		result |= 2
		return result
	}
	if StateHardware.WatchDog == 0 {
		result |= 1
	}
	if StateHardware.Status[0] == 0 {
		return result
	}
	// bit02 (0x04): “conflicting signal group command”
	if StateHardware.Status[0] == 10 || StateHardware.Status[0] == 11 {
		result |= 0x04
	}
	//  bit03 (0x08): “centralisation inhibited”
	//  bit04 (0x10): “inter-green violation”
	if StateHardware.Status[0] == 2 {
		result |= 0x10
	}
	//  bit05 (0x20): “lamp fault”
	if StateHardware.Status[0] == 1 {
		result |= 0x20
	}
	//  bit06 (0x40): [not used]
	//  bit07 (0x80): “extended diagnostics update”
	if StateHardware.Status[0] >= 3 && StateHardware.Status[0] <= 9 {
		result |= 0x80
	}
	return result
}
