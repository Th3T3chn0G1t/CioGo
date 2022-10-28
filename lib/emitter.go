package cionom

type Bytecode []byte
type BytecodePosition uint32
type RoutineIndex uint8

type RoutineData struct {
	Index  RoutineIndex
	Offset BytecodePosition
}

func BytecodeAppendAll(Bytecode Bytecode, Append Bytecode) Bytecode {
	for Position := BytecodePosition(0); Position < BytecodePosition(len(Append)); Position++ {
		Bytecode = append(Bytecode, Append[Position])
	}

	return Bytecode
}

func Emit(Program Program) Bytecode {
	var Header Bytecode
	var Code Bytecode

	var Routines = make(map[string]RoutineData)
	for RoutinePosition := RoutineIndex(0); RoutinePosition < RoutineIndex(len(Program.Routines)); RoutinePosition++ {
		Routines[Program.Routines[RoutinePosition].Identifier] = RoutineData{RoutineIndex(RoutinePosition), 0}
	}

	for RoutinePosition := 0; RoutinePosition < len(Program.Routines); RoutinePosition++ {
		var Routine Routine = Program.Routines[RoutinePosition]

		if Routine.External {
			continue
		}

		var Index RoutineIndex = Routines[Routine.Identifier].Index
		Routines[Routine.Identifier] = RoutineData{Index, BytecodePosition(len(Code))}

		for CallPosition := 0; CallPosition < len(Routine.Calls); CallPosition++ {
			var Call Call = Routine.Calls[CallPosition]
			Code = append(Code, 0x0)
			for ParameterPosition := 0; ParameterPosition < len(Call.Parameters); ParameterPosition++ {
				Code = append(Code, byte(0b01111111&Call.Parameters[ParameterPosition]))
			}
			Code = append(Code, byte(0b10000000|Routines[Call.Identifier].Index))
		}
		Code = append(Code, 0xFF)
	}

	Header = append(Header, byte(len(Program.Routines)&0b01111111))
	for RoutinePosition := 0; RoutinePosition < len(Program.Routines); RoutinePosition++ {
		var Routine Routine = Program.Routines[RoutinePosition]
		if Routine.External {
			Header = BytecodeAppendAll(Header, Bytecode{0xFF, 0xFF, 0xFF, 0xFF})
		} else {
			var OffsetParts Bytecode = Bytecode{
				byte((Routines[Routine.Identifier].Offset << (3 * 8)) & 0xFF),
				byte((Routines[Routine.Identifier].Offset << (2 * 8)) & 0xFF),
				byte((Routines[Routine.Identifier].Offset << (1 * 8)) & 0xFF),
				byte((Routines[Routine.Identifier].Offset << (0 * 8)) & 0xFF),
			}
			Header = BytecodeAppendAll(Header, OffsetParts)
		}
		Header = BytecodeAppendAll(Header, Bytecode(Routine.Identifier))
		Header = append(Header, 0x0)
	}

	return BytecodeAppendAll(Header, Code)
}
