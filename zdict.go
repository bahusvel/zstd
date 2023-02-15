package zstd

/*
#define ZDICT_STATIC_LINKING_ONLY
#include "zdict.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

func TrainFromBuffer(samples [][]byte) ([]byte, error) {
	nbSamples := C.unsigned(len(samples))
	sampleSizes := make([]C.size_t, len(samples))
	var dictBufferCapacity C.size_t = 0
	for i, sample := range samples {
		sampleSizes[i] = C.size_t(len(sample))
		dictBufferCapacity += C.size_t(len(sample))
	}
	samplesBuffer := make([]byte, dictBufferCapacity)
	offset := 0
	for _, sample := range samples {
		copy(samplesBuffer[offset:offset+len(sample)], sample)
		offset += len(sample)
	}
	dictBuffer := make([]byte, dictBufferCapacity)

	dictSize := C.ZDICT_trainFromBuffer(unsafe.Pointer(&dictBuffer[0]), dictBufferCapacity, unsafe.Pointer(&samplesBuffer[0]), (*C.size_t)(unsafe.Pointer(&sampleSizes[0])), nbSamples)
	if C.ZDICT_isError(dictSize) != 0 {
		return []byte{}, errors.New(C.GoString(C.ZDICT_getErrorName(dictSize)))
	}

	return dictBuffer[:dictSize], nil
}
