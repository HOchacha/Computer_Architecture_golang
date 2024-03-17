package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// this instruction fetches
func fetchInstruction(scanner *bufio.Scanner) string {
	return scanner.Text()
}

// 해당 함수에서는 tokenizing 및 레지스터 및 상수 값 구별을 할 수 있도록 한다.
// 값을 읽어들인 후에, 상수 혹은 레지스터가 판별되면 해당 바이트 스트림에 데이터를 입력하도록 한다.
// 반환은 에러 여부만 반환하도록 한다.
// 사실상의 동작이 추상화되어 있으며, 엄밀히 동작이 컴포넌트 별로 구분된 것은 아니다.
// 가령, ALU의 연산을 진행하기 위해서는 제어 신호를 받아들이는 부분이 있어야 한다.
// 하지만,
func decodeInstruction(
	gr *generalRegisters, sr *specificRegisters) error {
	// 데이터 토큰화 하기
	// 토큰화 이후 반환할 값 혹은 레지스터 참조 반환
	fields := strings.Fields(sr.instructionReg)

	writeInstructionFieldOnLog(fields)

	sr.operatorReg = fields[0]
	// 상수 간의 연산을 진행하는 경우에는, R0에 입력한다.
	// 그러지 않는 경우, 첫번째 Operand를 target으로, 두번째 Operand를 target으로 작성한다.
	switch fields[0] {
	case "+", "-", "*", "/", "M":
		// TODO : 일단 값을 문자열 해석 없이 넣는 방식으로 구현하였음,
		//        우선 계산기 내에 버퍼 역할을 수행하는 레지스터를 따로 구현할 생각은 없음
		// 		  특수 목적 레지스터 내에 값을 작성하여 연산을 수행하도록 한다면
		// 		  MIPS와 설계가 다소 멀어진다. 주소 체계로 접근할 수 있도록 하는 방법을 취하고 싶은데
		// 		  상수 값을 받아들일 때가 문제가 된다.
		//		  -> 상수 값을 받아들일 때는 특수 목적 레지스터 내에 상수 전용 레지스터를 사용하도록 하는 것이 적절해보임
		// 		  그러면, targetReg, SourceReg에 현재 단계에서 값을 입력을 하고, 값을 처리하는 execution 단계에서는
		//  	  단순히 레지스터에 참조를 전달하는 방식을 사용하도록 하자
		//        명령어의 구조는 [Operator] [SourceReg] [TargetReg]로 지정하도록 한다.
		sr.operatorReg = fields[0]

		if strings.HasPrefix(fields[0], "0x") {
			integer, _ := strconv.Atoi(strings.TrimPrefix(fields[0], "0x"))
			writeValueInRegister(uint32(integer), &sr.sourceReg)
		} else {
			integer, _ := strconv.Atoi(strings.TrimPrefix(fields[0], "R"))
			register := gr.getRegisterReferenceFromInteger(integer)
			sr.sourceReg = *register
		}

		if strings.HasPrefix(fields[1], "0x") {
			integer, _ := strconv.Atoi(strings.TrimPrefix(fields[1], "0x"))
			writeValueInRegister(uint32(integer), &sr.targetReg)
			sr.isTargetImmediate = true
		} else {
			integer, _ := strconv.Atoi(strings.TrimPrefix(fields[1], "R"))
			register := gr.getRegisterReferenceFromInteger(integer)
			sr.targetReg = *register
			sr.isTargetImmediate = false
		}
	case "B":
		//나중에 구현될 예정
	// Branch 연산의 경우, 스택을 구현하여 플로우 복귀까지 구현해야 할 수 있다.
	// 특히나 One-by-One으로 실행하는 경우에는 스택을 구현할 필요는 없다.
	case "C":
		// 나중에 구현될 예정
	}

	fmt.Printf("Log : Operator [%s], target register [%s], source register [%s]\n",
		sr.operatorReg, hex.EncodeToString(sr.sourceReg[:]), hex.EncodeToString(sr.targetReg[:]))

	//
	return nil
}

func executeInstruction(gr *generalRegisters, sr *specificRegisters) error {
	// 명령어 종류에 따라서 명령어를 실행시킬 수 있도록 한다.
	endian := systemEndian()
	switch sr.operatorReg {
	case "+":
		addByteArray(endian, gr, sr)
	case "-":
		subByteArray(endian, gr, sr)
	case "*":
		multiplyByteArray(endian, gr, sr)
	case "/":
		// TODO : 에러 처리가 필요하다 (Divided by Zero)
		divByteArray(endian, gr, sr)
	case "M":
		// 이 경우, 레지스터에 값을 넣는 작업을 해야하기 때문에, 불편하지만 다시 레지스터를 얻는 과정을 추가해야 한다.
		// 주소 지정하는 방법에 대해서는 다소 복잡한 방식을 취하고 있음
		moveByteArray(endian, gr, sr)
	}

	return nil
}

// 레지스터에 직접 접근하여 데이터를 입출력하도록 한다.
func WriteBackinMemory() {
	//Dismissed
}

func writeInstructionFieldOnLog(fields []string) {
	fmt.Printf("Log : ")
	for index, token := range fields {
		if index == 0 {
			fmt.Printf("Operator: %s\n", token)
		} else {
			fmt.Printf("Operand: %s\n", token)
		}
	}
}

func writeValueInRegister(value uint32, register *[4]byte) {
	systemType := systemEndian()
	if systemType == binary.BigEndian {
		binary.BigEndian.PutUint32(register[:], value)
	} else {
		binary.LittleEndian.PutUint32(register[:], value)
	}
}

func systemEndian() binary.ByteOrder {
	var i int32 = 0x01020304
	firstByte := (*byte)(unsafe.Pointer(&i))
	if *firstByte == 0x04 {
		return binary.LittleEndian
	}
	return binary.BigEndian
}
