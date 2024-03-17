package main

import (
	"encoding/binary"
)

func addByteArray(endian binary.ByteOrder, gr *generalRegisters, sr *specificRegisters) {

	if endian == binary.BigEndian {
		source := binary.BigEndian.Uint32(sr.sourceReg[:])
		target := binary.BigEndian.Uint32(sr.targetReg[:])

		binary.BigEndian.PutUint32(gr.r0[:], source+target)
	} else {
		source := binary.LittleEndian.Uint32(sr.sourceReg[:])
		target := binary.LittleEndian.Uint32(sr.targetReg[:])

		binary.LittleEndian.PutUint32(gr.r0[:], source+target)
	}
}

func subByteArray(endian binary.ByteOrder, gr *generalRegisters, sr *specificRegisters) {
	if endian == binary.BigEndian {
		source := binary.BigEndian.Uint32(sr.sourceReg[:])
		target := binary.BigEndian.Uint32(sr.targetReg[:])

		binary.BigEndian.PutUint32(gr.r0[:], source+target)
	} else {
		source := binary.LittleEndian.Uint32(sr.sourceReg[:])
		target := binary.LittleEndian.Uint32(sr.targetReg[:])

		binary.LittleEndian.PutUint32(gr.r0[:], source-target)
	}
}

func multiplyByteArray(endian binary.ByteOrder, gr *generalRegisters, sr *specificRegisters) {
	if endian == binary.BigEndian {
		source := binary.BigEndian.Uint32(sr.sourceReg[:])
		target := binary.BigEndian.Uint32(sr.targetReg[:])

		binary.BigEndian.PutUint32(gr.r0[:], source*target)
	} else {
		source := binary.LittleEndian.Uint32(sr.sourceReg[:])
		target := binary.LittleEndian.Uint32(sr.targetReg[:])

		binary.LittleEndian.PutUint32(gr.r0[:], source*target)
	}
}

func divByteArray(endian binary.ByteOrder, gr *generalRegisters, sr *specificRegisters) {
	if endian == binary.BigEndian {
		source := binary.BigEndian.Uint32(sr.sourceReg[:])
		target := binary.BigEndian.Uint32(sr.targetReg[:])

		binary.BigEndian.PutUint32(gr.r0[:], source*target)
	} else {
		source := binary.LittleEndian.Uint32(sr.sourceReg[:])
		target := binary.LittleEndian.Uint32(sr.targetReg[:])

		binary.LittleEndian.PutUint32(gr.r0[:], source*target)
	}
}

func moveByteArray(endian binary.ByteOrder, gr *generalRegisters, sr *specificRegisters) {
	var sourceIndex uint32
	var targetIndex uint32
	if endian == binary.BigEndian {
		sourceIndex = binary.BigEndian.Uint32(sr.sourceReg[:])
		targetIndex = binary.BigEndian.Uint32(sr.sourceReg[:])
	} else {
		sourceIndex = binary.LittleEndian.Uint32(sr.sourceReg[:])
		targetIndex = binary.LittleEndian.Uint32(sr.sourceReg[:])
	}

	sourceRegisterPointer := gr.getRegisterReferenceFromUint32(sourceIndex)
	if sr.isTargetImmediate == true {
		*sourceRegisterPointer = sr.targetReg
	} else {
		targetRegisterPointer := gr.getRegisterReferenceFromUint32(targetIndex)
		*sourceRegisterPointer = *targetRegisterPointer
	}
}
