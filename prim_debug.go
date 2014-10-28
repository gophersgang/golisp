// Copyright 2014 SteelSeries ApS.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package implements a basic LISP interpretor for embedding in a go program for scripting.
// This file contains the debugging primitive functions.

package golisp

import (
	"errors"
	"fmt"
	"strings"
)

var DebugCommandPrefix string = ":"

func RegisterDebugPrimitives() {
	MakePrimitiveFunction("debug-trace", -1, DebugTraceImpl)
	MakePrimitiveFunction("debug-on-error", -1, DebugOnErrorImpl)
	MakePrimitiveFunction("debug", -1, DebugImpl)
	MakePrimitiveFunction("dump", 0, DumpSymbolTableImpl)
}

func DumpSymbolTableImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	env.Dump()
	return
}

func DebugTraceImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	if Length(args) == 1 {
		DebugTrace = BooleanValue(Car(args))
	}
	return BooleanWithValue(DebugTrace), nil
}

func DebugOnErrorImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	if Length(args) == 1 {
		DebugOnError = BooleanValue(Car(args))
	}
	return BooleanWithValue(DebugOnError), nil
}

func DebugImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	fmt.Printf("Debugger\n")

	DebugRepl(env)
	return
}

func processState(tokens []string) (ok bool, state bool) {
	if len(tokens) != 2 {
		fmt.Printf("Missing on/off.\n")
		return false, false
	} else {
		switch tokens[1] {
		case "on":
			return true, true
		case "off":
			return true, false
		default:
			fmt.Printf("on/off expected.\n")
			return false, false
		}
	}
}

func DebugRepl(env *SymbolTableFrame) {
	env.DumpHeader()
	prompt := "D> "
	for true {
		defer func() {
			if x := recover(); x != nil {
				println("BANG!")
			}
		}()
		input := *ReadLine(&prompt)
		if input != "" {
			if strings.HasPrefix(input, DebugCommandPrefix) {
				cmd := strings.TrimPrefix(input, DebugCommandPrefix)
				tokens := strings.Split(cmd, " ")
				switch tokens[0] {
				case "?":
					fmt.Printf("SteelSeries/GoLisp Debugger\n")
					fmt.Printf("---------------------------\n")
					fmt.Printf(":?        - show this command summary\n")
					fmt.Printf(":b        - show the environment stack\n")
					fmt.Printf(":c        - continue, exiting the debugger\n")
					fmt.Printf(":d        - do a full of the environment stack\n")
					fmt.Printf(":e on/off - Enable/disable debug on error\n")
					fmt.Printf(":f frame# - do a full dump of a single environment frame\n")
					fmt.Printf(":n        - step to next (run to the next evaluation in this frame)\n")
					fmt.Printf(":q        - quit GoLisp\n")
					fmt.Printf(":r sexpr  - return from the current evaluation with the specified value\n")
					fmt.Printf(":s        - single step (run to the next evaluation)\n")
					fmt.Printf(":t on/off - Enable/disable tracing\n")
					fmt.Printf(":u        - continue until the enclosing environment frame is returned to\n")
					fmt.Printf("\n")
				case "b":
					env.DumpHeaders()
					fmt.Printf("\n")
				case "c":
					DebugCurrentFrame = nil
					DebugSingleStep = false
					DebugEvalInDebugRepl = false
					return
				case "d":
					env.Dump()
				case "e":
					ok, state := processState(tokens)
					if ok {
						DebugOnError = state
					}
				case "f":
					var fnum int
					if len(tokens) != 2 {
						fmt.Printf("Missing frame number.\n")
					} else {
						_, err := fmt.Sscanf(tokens[1], "%d", &fnum)
						if err != nil {
							fmt.Printf("Bad frame number: '%s'. %s\n", tokens[1], err)
						} else {
							env.DumpSingleFrame(fnum)
						}
					}
					//				case "n":

				case "q":
					QuitImpl(nil, nil)
				case "r":
					DebugEvalInDebugRepl = true
					code, err := Parse(strings.Join(tokens[1:], " "))
					d, err := Eval(code, env)
					DebugEvalInDebugRepl = false
					if err != nil {
						fmt.Printf("Error in evaluation: %s\n", err)
					} else {
						DebugReturnValue = d
						DebugCurrentFrame = nil
						DebugSingleStep = false
						DebugEvalInDebugRepl = false
						return
					}
				case "s":
					DebugSingleStep = true
					return
				case "t":
					ok, state := processState(tokens)
					if ok {
						DebugTrace = state
					}
				case "u":
					if env.Parent != nil {
						DebugCurrentFrame = env
						return
					} else {
						fmt.Printf("Already at top frame.\n")
					}
				}
			} else {
				code, err := Parse(input)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
				} else {
					DebugEvalInDebugRepl = true
					d, err := Eval(code, env)
					DebugEvalInDebugRepl = false
					if err != nil {
						fmt.Printf("Error in evaluation: %s\n", err)
					} else {
						fmt.Printf("==> %s\n", String(d))
					}
				}
			}
		}
	}
}

func ProcessError(errorMessage string, env *SymbolTableFrame) error {
	if DebugOnError && IsInteractive && !DebugEvalInDebugRepl {
		fmt.Printf("ERROR!  %s\n", errorMessage)
		DebugRepl(env)
		return nil
	} else {
		return errors.New(errorMessage)
	}
}