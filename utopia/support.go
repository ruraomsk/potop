package utopia

func bcc8(b byte, check uint16) uint16 {
	if b != 0 {
		var rbit uint16
		check = check ^ uint16(b)
		for i := 0; i < 8; i++ {
			rbit = check & 0x01
			check = check >> 1
			if rbit != 0 {
				check = check | 0x8000
				check = check ^ 0x2001
			}
		}
	}
	return check
}
func crc16_calc(data []byte) uint16 {
	var crc uint16
	crc = 0
	for _, v := range data {
		crc = bcc8(v, crc)
	}
	return crc
}
