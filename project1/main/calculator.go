package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

// TODO: golang을 설치했다는 환경에서 실행해야 하므로, 이를 검증하기 위해서는
// 		 코드를 테스트 할 수 있는 스크립트를 작성하여 환경 구성을 위해
//       최적화 해놓을 것

// TODO: 괜찮다면, Dockerfile을 만들어서 도커에서 테스트 할 수 있도록 한다.

// TODO: 동시성 프로그래밍을 지원하기 위해서, 각 단계를 각각의 객체로 구성한 후에
//	     메모리 락을 할 수 있도록 한다.

// TODO: 모든 레지스터의 연산 구현은 CPU의 API를 사용하는 방식으로 구현한다.
// 		 논리 레벨에서 구현하고 있으면 정신 건강에 안좋아질 수 있으니, 단순하게 구현하도록 한다.

// TODO: 모든 레지스터는 가능한 어셈블리 레지스터를 다루는 것처럼 사용할 수 있도록 한다.
type generalRegisters struct {
	r0 [4]byte
	r1 [4]byte
	r2 [4]byte
	r3 [4]byte
	r4 [4]byte
	r5 [4]byte
	r6 [4]byte
	r7 [4]byte
	r8 [4]byte
	r9 [4]byte
}
type specificRegisters struct {
	instructionReg    string
	bufferReg         [4]byte
	flagReg           bool
	operatorReg       string
	isTargetImmediate bool
	sourceReg         [4]byte
	targetReg         [4]byte
}

func main() {
	startTime := time.Now()
	gr := generalRegisters{
		r0: [4]byte{},
		r1: [4]byte{},
		r2: [4]byte{},
		r3: [4]byte{},
		r4: [4]byte{},
		r5: [4]byte{},
		r6: [4]byte{},
		r7: [4]byte{},
		r8: [4]byte{},
		r9: [4]byte{},
	}
	sr := specificRegisters{
		instructionReg:    "",
		bufferReg:         [4]byte{},
		flagReg:           false,
		operatorReg:       "",
		isTargetImmediate: false,
		sourceReg:         [4]byte{},
		targetReg:         [4]byte{},
	}
	//init configuration for calculator
	if len(os.Args) <= 1 {
		fmt.Println("Error : No input file specified")
		return
	}
	fmt.Println("The input file name : ", os.Args[1])

	//read input text file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error : did you specified the input text file?")
		panic(err)
	}
	defer file.Close()

	//start core state machine
	// TODO: how to implement registers of virtual machine?
	// we have to check the difference between actual computer and this virtual calculator

	// hand over the file to scanner
	// calculator Single Cycle up and running
	scanner := bufio.NewScanner(file)

	// 초기 구성까지 몇 초가 걸렸는지 체크한다.
	fmt.Printf("Log : Calculator finished within %d\n", time.Since(startTime))
	// 각 명령어 사이클 파트를 실행하는데 각각 몇 ms가 걸리는지 파악할 것
	for scanner.Scan() {
		pointTime := time.Now()

		// fetch instruction
		sr.instructionReg = fetchInstruction(scanner)
		fetchTime := time.Since(pointTime)
		fmt.Printf("Log: fetched instruction within %d\nLog: %s\n", fetchTime, sr.instructionReg)

		// Decode instruction
		err := decodeInstruction(&gr, &sr)
		if err != nil {
			// TODO: handle the error for decoding level
			// 		 such as malformed Instruction
			//
			panic("Unforeseen Error Occurred")
		}
		// execute instruction
		err = executeInstruction(&gr, &sr)
		if err != nil {
			panic("Execute Error : ")
		}
		// Write Back instruction

		// TODO: remove this comment before submit
		// No Implementation in this part, just left over this procedure for following CPU Cycle convention respectively
		// Maybe, If i have to fix the decode and execute level for correctness,
		// i would gonna pissed off seriously.
		fmt.Printf("Log : Temporary Result, %s\n\n\n", hex.EncodeToString(gr.r0[:]))
	}

	//halt calculator
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}

func (g *generalRegisters) getRegisterReferenceFromInteger(index int) *[4]byte {
	switch index {
	case 0:
		return &g.r0
	case 1:
		return &g.r1
	case 2:
		return &g.r2
	case 3:
		return &g.r3
	case 4:
		return &g.r4
	case 5:
		return &g.r5
	case 6:
		return &g.r6
	case 7:
		return &g.r7
	case 8:
		return &g.r8
	case 9:
		return &g.r9
	default:
		return nil
	}
}

func (g *generalRegisters) getRegisterReferenceFromUint32(index uint32) *[4]byte {
	switch index {
	case 0:
		return &g.r0
	case 1:
		return &g.r1
	case 2:
		return &g.r2
	case 3:
		return &g.r3
	case 4:
		return &g.r4
	case 5:
		return &g.r5
	case 6:
		return &g.r6
	case 7:
		return &g.r7
	case 8:
		return &g.r8
	case 9:
		return &g.r9
	default:
		return nil
	}
}
