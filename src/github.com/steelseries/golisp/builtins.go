// Copyright 2013 SteelSeries ApS. All rights reserved.
// No license is given for the use of this source code.

// This package impliments a basic LISP interpretor for embedding in a go program for scripting.
// This file pre-loads primitive bui8ltin functions

// Basic arithmetic is implimented, with skeletons for the rest of the "special" symbols.
// Flesh out as required. Remember to add tests to builtins_test.go

package golisp

import (
    "container/list"
    "errors"
)

func init() {
    symbolTable = &SymbolTable{list.New()}
    PushLocalBindings()
    InitBuiltins()
}

func InitBuiltins() {
    // MakePrimitiveFunction(<symbol>, <required # args, -1 means >= 1>, <function>)
    MakePrimitiveFunction("+", -1, Add)
    MakePrimitiveFunction("-", -1, Subtract)
    MakePrimitiveFunction("*", -1, Multiply)
    MakePrimitiveFunction("/", -1, Quotient)
    // MakePrimitiveFunction("%", 2, Remainder)
    // MakePrimitiveFunction("<", -1, LessThan)
    // MakePrimitiveFunction(">", -1, GreaterThan)
    // MakePrimitiveFunction("==", 2, EqualTo)
    // MakePrimitiveFunction("!", 1, Not)
    // MakePrimitiveFunction("!=", 2, NotEqual)
    // MakePrimitiveFunction("<=", -1, LessThanOrEqualTo)
    // MakePrimitiveFunction(">=", -1, GreaterThanOrEqualTo)
    MakePrimitiveFunction("if", -1, If)
    // MakePrimitiveFunction("var", 2, Var)
}

func Add(args *Data) (result *Data, err error) {
    var acc int = 0
    for c := args; NotNilP(c); c = Cdr(c) {
        if TypeOf(Car(c)) != NumberType {
            err = errors.New("Number expected")
            return
        }
        acc += IntValue(Car(c))
    }
    return NumberWithValue(acc), nil
}

func Subtract(args *Data) (result *Data, err error) {
    if TypeOf(Car(args)) != NumberType {
        err = errors.New("Number expected")
        return
    }
    var acc int = IntValue(Car(args))
    if Length(args) == 1 { //negation
        acc = -acc
    } else {
        for c := Cdr(args); NotNilP(c); c = Cdr(c) {
            if TypeOf(Car(c)) != NumberType {
                err = errors.New("Number expected")
                return
            }
            acc -= IntValue(Car(c))
        }
    }
    return NumberWithValue(acc), nil
}

func Multiply(args *Data) (result *Data, err error) {
    var acc int = 1
    for c := args; NotNilP(c); c = Cdr(c) {
        if TypeOf(Car(c)) != NumberType {
            err = errors.New("Number expected")
            return
        }
        acc *= IntValue(Car(c))
    }
    return NumberWithValue(acc), nil
}

func Quotient(args *Data) (result *Data, err error) {
    if TypeOf(Car(args)) != NumberType {
        err = errors.New("Number expected")
        return
    }
    var acc int = IntValue(Car(args))
    for c := Cdr(args); NotNilP(c); c = Cdr(c) {
        if TypeOf(Car(c)) != NumberType {
            err = errors.New("Number expected")
            return
        }
        acc /= IntValue(Car(c))
    }
    return NumberWithValue(acc), nil
}

// func Remainder(args *Data) (result *Data, err error) {
// }

// func LessThan(args *Data) (result *Data, err error) {
// }

// func GreaterThan(args *Data) (result *Data, err error) {
// }

// func EqualTo(args *Data) (result *Data, err error) {
// }

// func Not(args *Data) (result *Data, err error) {
// }

// func NotEqual(args *Data) (result *Data, err error) {
// }

// func LessThanOrEqualTo(args *Data) (result *Data, err error) {
// }

// func GreaterThanOrEqualTo(args *Data) (result *Data, err error) {
// }

func If(args *Data) (result *Data, err error) {
    condition := BooleanValue(Eval(Car(args)))
    thenClause := Cadr(args)
    elseClause := Caddr(args)
    if Cdddr(args) != nil {
        err = errors.New("Too many arguments to IF")
        return
    }
    if condition {
        return Eval(thenClause), nil
    } else {
        return Eval(elseClause), nil
    }
}

// func Var(args *Data) (result *Data, err error) {
// }
