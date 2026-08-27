package service

type SpecimenDetailBuffer struct {
	Values []string
}

var sharedSpecimenBuffers = make(map[uint]SpecimenDetailBuffer)

func specimenPositionTokens(containerID uint, parts ...string) SpecimenDetailBuffer {
	buffer := sharedSpecimenBuffers[containerID]
	if cap(buffer.Values) < len(parts) {
		buffer.Values = make([]string, 0, len(parts))
	}
	buffer.Values = append(buffer.Values[:0], parts...)
	sharedSpecimenBuffers[containerID] = buffer
	return buffer
}
