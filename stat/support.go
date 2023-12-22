package stat

func GetDescription() (int, int) {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	return Statistics.diaps, Statistics.counts
}
func GetCountValues() []int {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	ret := make([]int, 0)
	for i := 0; i < Statistics.counts; i++ {
		r, is := Statistics.chanels[i]
		if !is {
			continue
		}
		ret = append(ret, r.CountValues.Value...)
	}
	return ret
}
func GetOcupaeValues() []int {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	ret := make([]int, 0)
	for i := 0; i < Statistics.counts; i++ {
		r, is := Statistics.chanels[i]
		if !is {
			continue
		}
		ret = append(ret, r.OcupaeValues.Value...)
	}
	return ret
}

func GetSpeedValues() []int {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	ret := make([]int, 0)
	for i := 0; i < Statistics.counts; i++ {
		r, is := Statistics.chanels[i]
		if !is {
			continue
		}
		ret = append(ret, r.SpeedValues.Value...)
	}
	return ret
}
func ClearCountValues() {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	for i := 0; i < Statistics.counts; i++ {
		r, is := Statistics.chanels[i]
		if !is {
			continue
		}
		r.CountValues.clear(Statistics.diaps)
	}
}
func ClearSpeedValues() {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	for i := 0; i < Statistics.counts; i++ {
		r, is := Statistics.chanels[i]
		if !is {
			continue
		}
		r.SpeedValues.clear(Statistics.diaps)
	}
}
