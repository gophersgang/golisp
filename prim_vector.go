// Copyright 2014 SteelSeries ApS.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package implements a basic LISP interpretor for embedding in a go program for scripting.
// This file contains the vector primitive functions.

package golisp

import (
	"fmt"
	"math"
)

func RegisterVectorPrimitives() {
	MakePrimitiveFunction("make-vector", "1|2", MakeVectorImpl)
	MakePrimitiveFunction("vector", ">=1", VectorImpl)
	MakePrimitiveFunction("vector-copy", "1", VectorCopyImpl)
	MakePrimitiveFunction("list->vector", "1", ListToVectorImpl)
	MakePrimitiveFunction("vector->list", "1", VectorToListImpl)
	MakePrimitiveFunction("make-initialized-vector", "2", MakeInitializedVectorImpl)
	MakePrimitiveFunction("vector-grow", "2", VectorGrowImpl)
	MakePrimitiveFunction("vector-map", ">=2", VectorMapImpl)          // also called by map
	MakePrimitiveFunction("vector-for-each", ">=2", VectorForEachImpl) // also called by for-each
	MakePrimitiveFunction("vector-reduce", "3", VectorReduceImpl)      // also called by reduce
	MakePrimitiveFunction("vector-filter", "2", VectorFilterImpl)      // also called by filter
	MakePrimitiveFunction("vector-remove", "2", VectorRemoveImpl)      // also clled by remove
	MakePrimitiveFunction("vector?", "1", VectorPImpl)
	MakePrimitiveFunction("vector-length", "1", VectorLengthImpl)   // also called by length
	MakePrimitiveFunction("vector-ref", "2", VectorRefImpl)         // also called by nth
	MakePrimitiveFunction("vector-set!", "3", VectorSetImpl)        // also called by set-nth!
	MakePrimitiveFunction("vector-first", "1", VectorFirstImpl)     // also called by first
	MakePrimitiveFunction("vector-second", "1", VectorSecondImpl)   // also called by second
	MakePrimitiveFunction("vector-third", "1", VectorThirdImpl)     // also called by third
	MakePrimitiveFunction("vector-fourth", "1", VectorFourthImpl)   // also called by fourth
	MakePrimitiveFunction("vector-fifth", "1", VectorFifthImpl)     // also called by fifth
	MakePrimitiveFunction("vector-sixth", "1", VectorSixthImpl)     // also called by sixth
	MakePrimitiveFunction("vector-seventh", "1", VectorSeventhImpl) // also called by seventh
	MakePrimitiveFunction("vector-eighth", "1", VectorEighthImpl)   // also called by eighth
	MakePrimitiveFunction("vector-ninth", "1", VectorNinthImpl)     // also called by ninth
	MakePrimitiveFunction("vector-tenth", "1", VectorTenthImpl)     // also called by tenth
	MakePrimitiveFunction("vector-binary-search", "4", VectorBinarySearchImpl)
	MakePrimitiveFunction("vector-find", "2", VectorFindImpl) // also called by find
	MakePrimitiveFunction("subvector", "3", SubVectorImpl)
	MakePrimitiveFunction("vector-head", "2", VectorHeadImpl) // also called by take
	MakePrimitiveFunction("vector-tail", "2", VectorTailImpl) // also called by drop
	MakePrimitiveFunction("vector-fill!", "2", VectorFillImpl)
	MakePrimitiveFunction("subvector-fill!", "4", SubVectorFillImpl)
	MakePrimitiveFunction("subvector-move-left!", "5", SubVectorMoveLeftImpl)
	MakePrimitiveFunction("subvector-move-right!", "5", SubVectorMoveRightImpl)
	MakePrimitiveFunction("vector-sort", "2", VectorSortImpl) // also called by sort
	MakePrimitiveFunction("vector-sort!", "2", VectorSortInPlaceImpl)
}

func MakeVectorImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	k := First(args)
	if !IntegerP(k) {
		err = ProcessError(fmt.Sprintf("make-vector needs an integer as its first argument, but got %s.", String(k)), env)
		return
	}

	size := IntegerValue(k)

	var value *Data = nil
	if Length(args) == 2 {
		value = Second(args)
	}

	vals := make([]*Data, size)
	for i := int64(0); i < size; i++ {
		vals[i] = value
	}

	result = VectorWithValue(vals)
	return
}

func VectorImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	result = VectorWithValue(ToArray(args))
	return
}

func VectorCopyImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	vect := First(args)
	if !VectorP(vect) {
		err = ProcessError(fmt.Sprintf("vector-copy needs a vector as its argument, but got %s.", String(vect)), env)
		return
	}

	v := VectorValue(vect)
	newV := make([]*Data, 0, len(v))
	for _, e := range v {
		newV = append(newV, e)
	}
	result = VectorWithValue(newV)
	return
}

func ListToVectorImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	l := First(args)
	if !ListP(l) {
		err = ProcessError(fmt.Sprintf("list->vector needs a list as its argument, but got %s.", String(l)), env)
		return
	}

	result = VectorWithValue(ToArray(l))
	return
}

func VectorToListImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector->list needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	result = ArrayToList(VectorValue(v))
	return
}

func MakeInitializedVectorImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	k := First(args)
	if !IntegerP(k) {
		err = ProcessError(fmt.Sprintf("make-initialized-vector needs an integer as its first argument, but got %s.", String(k)), env)
		return
	}

	size := IntegerValue(k)

	f := Second(args)
	if !FunctionOrPrimitiveP(f) {
		err = ProcessError(fmt.Sprintf("make-initialized-vector needs a function as its second argument, but got %s.", String(f)), env)
		return
	}

	vals := make([]*Data, size)
	for i := int64(0); i < size; i++ {
		vals[i], err = Apply(f, InternalMakeList(IntegerWithValue(i)), env)
		if err != nil {
			return
		}
	}

	result = VectorWithValue(vals)
	return
}

func VectorGrowImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-grow needs a vector as its first argument, but got %s.", String(v)), env)
		return
	}
	originalValues := VectorValue(v)

	k := Second(args)
	if !IntegerP(k) {
		err = ProcessError(fmt.Sprintf("vector-grow needs an integer as its second argument, but got %s.", String(k)), env)
		return
	}

	size := IntegerValue(k)

	if int(size) <= len(originalValues) {
		err = ProcessError(fmt.Sprintf("vector-grow needs a new size that is larger than the size of its vector argument (%d), but got %s.", len(originalValues), size), env)
		return
	}

	vals := make([]*Data, size)
	for i, val := range originalValues {
		vals[i] = val
	}

	result = VectorWithValue(vals)
	return
}

func VectorMapImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	f := First(args)
	if !FunctionOrPrimitiveP(f) {
		err = ProcessError(fmt.Sprintf("vector-map needs a function as its first argument, but got %s.", String(f)), env)
		return
	}

	var collections [][]*Data = make([][]*Data, 0, Length(args)-1)
	var loopCount int64 = math.MaxInt64
	var col *Data
	for a := Cdr(args); NotNilP(a); a = Cdr(a) {
		col = Car(a)
		if !VectorP(col) {
			err = ProcessError(fmt.Sprintf("vector-map needs vectors as its other arguments, but got %s.", String(col)), env)
			return
		}
		if NilP(col) || col == nil {
			return
		}
		collections = append(collections, VectorValue(col))
		loopCount = intMin(loopCount, int64(Length(col)))
	}

	if loopCount == math.MaxInt64 {
		return
	}

	var vals []*Data = make([]*Data, loopCount)
	var v *Data
	var a *Data
	for index := 0; index < int(loopCount); index++ {
		mapArgs := make([]*Data, 0, len(collections))
		for _, mapArgCollection := range collections {
			a = mapArgCollection[index]
			mapArgs = append(mapArgs, a)
		}
		v, err = ApplyWithoutEval(f, ArrayToList(mapArgs), env)
		if err != nil {
			return
		}
		vals[index] = v
	}

	result = VectorWithValue(vals)
	return
}

func VectorForEachImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	f := First(args)
	if !FunctionOrPrimitiveP(f) {
		err = ProcessError(fmt.Sprintf("vector-for-each needs a function as its first argument, but got %s.", String(f)), env)
		return
	}

	var collections [][]*Data = make([][]*Data, 0, Length(args)-1)
	var loopCount int64 = math.MaxInt64
	var col *Data
	for a := Cdr(args); NotNilP(a); a = Cdr(a) {
		col = Car(a)
		if !VectorP(col) {
			err = ProcessError(fmt.Sprintf("vector-for-each needs vectors as its other arguments, but got %s.", String(col)), env)
			return
		}
		if NilP(col) || col == nil {
			return
		}
		collections = append(collections, VectorValue(col))
		loopCount = intMin(loopCount, int64(Length(col)))
	}

	if loopCount == math.MaxInt64 {
		return
	}

	var a *Data
	for index := 0; index < int(loopCount); index++ {
		mapArgs := make([]*Data, 0, len(collections))
		for _, mapArgCollection := range collections {
			a = mapArgCollection[index]
			mapArgs = append(mapArgs, a)
		}
		_, err = ApplyWithoutEval(f, ArrayToList(mapArgs), env)
		if err != nil {
			return
		}
	}

	return
}

func VectorReduceImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	f := First(args)
	if !FunctionOrPrimitiveP(f) {
		err = ProcessError("vector-reduce needs a function as its first argument", env)
		return
	}

	initial := Second(args)
	col := Third(args)

	if !VectorP(col) {
		err = ProcessError("vector-reduce needs a vector as its third argument", env)
		return
	}

	v := VectorValue(col)

	if Length(col) == 0 {
		return initial, nil
	}

	if Length(col) == 1 {
		return v[0], nil
	}

	result = v[0]
	for _, val := range v[1:] {
		result, err = ApplyWithoutEval(f, InternalMakeList(result, val), env)
		if err != nil {
			return
		}
	}

	return
}

func VectorFilterImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	f := First(args)
	if !FunctionOrPrimitiveP(f) {
		err = ProcessError(fmt.Sprintf("vector-filter needs a function as its first argument, but got %s.", String(f)), env)
		return
	}

	v := Second(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-filter needs a vector as its second argument, but got %s.", String(v)), env)
		return
	}
	vals := VectorValue(v)

	var d []*Data = make([]*Data, 0, len(vals))
	for _, val := range vals {
		v, err = ApplyWithoutEval(f, InternalMakeList(val), env)
		if err != nil {
			return
		}
		if !BooleanP(v) {
			err = ProcessError("vector-filter needs a predicate function as its first argument.", env)
			return
		}

		if BooleanValue(v) {
			d = append(d, val)
		}
	}

	return VectorWithValue(d), nil
}

func VectorRemoveImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	f := First(args)
	if !FunctionOrPrimitiveP(f) {
		err = ProcessError(fmt.Sprintf("vector-remove needs a function as its first argument, but got %s.", String(f)), env)
		return
	}

	v := Second(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-remove needs a vector as its second argument, but got %s.", String(v)), env)
		return
	}
	vals := VectorValue(v)

	var d []*Data = make([]*Data, 0, len(vals))
	for _, val := range vals {
		v, err = ApplyWithoutEval(f, InternalMakeList(val), env)
		if err != nil {
			return
		}
		if !BooleanP(v) {
			err = ProcessError("vector-remove needs a predicate function as its first argument.", env)
			return
		}

		if !BooleanValue(v) {
			d = append(d, val)
		}
	}

	return VectorWithValue(d), nil
}

func VectorPImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	return BooleanWithValue(VectorP(v)), nil
}

func VectorLengthImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-length needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	result = IntegerWithValue(int64(Length(v)))
	return
}

func VectorRefImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-ref needs a vector as its first argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	k := Second(args)
	if !IntegerP(k) {
		err = ProcessError(fmt.Sprintf("vector-ref needs an integer as its second argument, but got %s.", String(k)), env)
		return
	}
	index := int(IntegerValue(k))

	if index >= len(values) {
		err = ProcessError(fmt.Sprintf("vector-ref needs an index less than the vector length, but got %d.", index), env)
		return
	}

	result = values[index]

	return
}

func VectorSetImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-set! needs a vector as its first argument, but got %s.", String(v)), env)
		return
	}

	k := Second(args)
	if !IntegerP(k) {
		err = ProcessError(fmt.Sprintf("vector-set! needs an integer as its second argument, but got %s.", String(k)), env)
		return
	}

	values := VectorValue(v)
	kval := int(IntegerValue(k))
	newValue := Third(args)

	values[kval] = newValue
	result = StringWithValue("OK")
	return
}

func VectorFirstImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-first needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 0 {
		result = values[0]
	} else {
		err = ProcessError(fmt.Sprintf("vector-first needs a vector with length of at least 1, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorSecondImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-second needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 1 {
		result = values[1]
	} else {
		err = ProcessError(fmt.Sprintf("vector-second needs a vector with length of at least 2, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorThirdImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-third needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 2 {
		result = values[2]
	} else {
		err = ProcessError(fmt.Sprintf("vector-third needs a vector with length of at least 3, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorFourthImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-fourth needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 3 {
		result = values[3]
	} else {
		err = ProcessError(fmt.Sprintf("vector-fourth needs a vector with length of at least 4, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorFifthImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-fifth needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 4 {
		result = values[4]
	} else {
		err = ProcessError(fmt.Sprintf("vector-fifth needs a vector with length of at least 5, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorSixthImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-sixth needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 5 {
		result = values[5]
	} else {
		err = ProcessError(fmt.Sprintf("vector-sixth needs a vector with length of at least 6, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorSeventhImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-seventh needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 6 {
		result = values[6]
	} else {
		err = ProcessError(fmt.Sprintf("vector-seventh needs a vector with length of at least 7, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorEighthImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-eighth needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 7 {
		result = values[7]
	} else {
		err = ProcessError(fmt.Sprintf("vector-eigth needs a vector with length of at least 8, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorNinthImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-ninth needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 8 {
		result = values[8]
	} else {
		err = ProcessError(fmt.Sprintf("vector-ninth needs a vector with length of at least 9, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorTenthImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-tenth needs a vector as its argument, but got %s.", String(v)), env)
		return
	}

	values := VectorValue(v)
	if len(values) > 9 {
		result = values[9]
	} else {
		err = ProcessError(fmt.Sprintf("vector-tenth needs a vector with length of at least 10, but got %d.", len(values)), env)
		return
	}

	return
}

func VectorBinarySearchImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	return
}

func VectorFindImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	f := First(args)
	if !FunctionOrPrimitiveP(f) {
		err = ProcessError("vector-find needs a function as its first argument", env)
		return
	}

	v := Second(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-find needs a vector as its second argument, but got %s.", String(v)), env)
		return
	}

	var found *Data
	for _, val := range VectorValue(v) {
		found, err = ApplyWithoutEval(f, InternalMakeList(val), env)
		if !BooleanP(found) {
			err = ProcessError("vector-find needs a predicate function as its first argument.", env)
			return
		}
		if BooleanValue(found) {
			return val, nil
		}
	}

	return LispFalse, nil
}

func SubVectorImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("subvector needs a vector as its argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	start := Second(args)
	if !IntegerP(start) {
		err = ProcessError(fmt.Sprintf("subvector needs an integer as its starting index, but got %s.", String(start)), env)
		return
	}
	startIndex := int(IntegerValue(start))

	if startIndex < 0 || startIndex >= len(values) {
		err = ProcessError(fmt.Sprintf("subvector starting index is out of bounds (0-%d), got %d.", len(values)-1, startIndex), env)
		return
	}

	end := Third(args)
	if !IntegerP(end) {
		err = ProcessError(fmt.Sprintf("subvector needs an integer as its ending index, but got %s.", String(end)), env)
		return
	}
	endIndex := int(IntegerValue(end))

	if endIndex < startIndex || endIndex > len(values) {
		err = ProcessError(fmt.Sprintf("subvector ending index is out of bounds (%d-%d), got %d.", startIndex, len(values), startIndex), env)
		return
	}

	result = VectorWithValue(values[startIndex:endIndex])
	return
}

func VectorHeadImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-head needs a vector as its argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	end := Second(args)
	if !IntegerP(end) {
		err = ProcessError(fmt.Sprintf("vector-head needs an integer as its ending index, but got %s.", String(end)), env)
		return
	}
	endIndex := int(IntegerValue(end))

	if endIndex < 0 || endIndex > len(values) {
		err = ProcessError(fmt.Sprintf("vector-head ending index is out of bounds (0-%d), got %d.", len(values), endIndex), env)
		return
	}

	result = VectorWithValue(values[:endIndex])
	return
}

func VectorTailImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-tail needs a vector as its argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	start := Second(args)
	if !IntegerP(start) {
		err = ProcessError(fmt.Sprintf("vector-tail needs an integer as its starting index, but got %s.", String(start)), env)
		return
	}
	startIndex := int(IntegerValue(start))

	if startIndex < 0 || startIndex > len(values) {
		err = ProcessError(fmt.Sprintf("vector-tail starting index is out of bounds (0-%d), got %d.", len(values), startIndex), env)
		return
	}

	result = VectorWithValue(values[startIndex:])
	return
}

func VectorFillImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-fill! needs a vector as its argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	newValue := Second(args)

	for i, _ := range values {
		values[i] = newValue
	}

	return
}

func SubVectorFillImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("subvector-fill! needs a vector as its argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	start := Second(args)
	if !IntegerP(start) {
		err = ProcessError(fmt.Sprintf("subvector-fill! needs an integer as its starting index, but got %s.", String(start)), env)
		return
	}
	startIndex := int(IntegerValue(start))

	if startIndex < 0 || startIndex >= len(values) {
		err = ProcessError(fmt.Sprintf("subvector-fill! starting index is out of bounds (0-%d), got %d.", len(values)-1, startIndex), env)
		return
	}

	end := Third(args)
	if !IntegerP(end) {
		err = ProcessError(fmt.Sprintf("subvector-fill! needs an integer as its ending index, but got %s.", String(end)), env)
		return
	}
	endIndex := int(IntegerValue(end))

	if endIndex < startIndex || endIndex > len(values) {
		err = ProcessError(fmt.Sprintf("subvector-fill! ending index is out of bounds (%d-%d), got %d.", startIndex, len(values), startIndex), env)
		return
	}

	newValue := Fourth(args)

	for i := startIndex; i < endIndex; i = i + 1 {
		values[i] = newValue
	}
	return
}

func SubVectorMoveLeftImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("subvector-move-left! needs a vector as its first vector argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	start := Second(args)
	if !IntegerP(start) {
		err = ProcessError(fmt.Sprintf("subvector-move-left! needs an integer as its starting index, but got %s.", String(start)), env)
		return
	}
	startIndex := int(IntegerValue(start))

	if startIndex < 0 || startIndex >= len(values) {
		err = ProcessError(fmt.Sprintf("subvector-move-left! starting index is out of bounds (0-%d), got %d.", len(values)-1, startIndex), env)
		return
	}

	end := Third(args)
	if !IntegerP(end) {
		err = ProcessError(fmt.Sprintf("subvector-move-left! needs an integer as its ending index, but got %s.", String(end)), env)
		return
	}
	endIndex := int(IntegerValue(end))

	if endIndex < startIndex || endIndex > len(values) {
		err = ProcessError(fmt.Sprintf("subvector-move-left! ending index is out of bounds (%d-%d), got %d.", startIndex, len(values), startIndex), env)
		return
	}

	v2 := Fourth(args)
	if !VectorP(v2) {
		err = ProcessError(fmt.Sprintf("subvector-move-left! needs a vector as its second vector argument, but got %s.", String(v2)), env)
		return
	}
	values2 := VectorValue(v2)

	start2 := Fifth(args)
	if !IntegerP(start2) {
		err = ProcessError(fmt.Sprintf("subvector-move-left! needs an integer as its second starting index, but got %s.", String(start2)), env)
		return
	}
	startIndex2 := int(IntegerValue(start2))

	if startIndex < 0 || startIndex >= len(values2) {
		err = ProcessError(fmt.Sprintf("subvector-move-left! starting index is out of bounds (0-%d), got %d.", len(values)-1, startIndex), env)
		return
	}

	sourceLength := endIndex - startIndex
	tailSize2 := len(values2) - startIndex2
	if sourceLength > tailSize2 {
		err = ProcessError(fmt.Sprintf("subvector-move-left! source subvector is longer than the available space in the destination (0-%d), got %d.", tailSize2, sourceLength), env)
		return
	}

	for i, i2 := startIndex, startIndex2; i < endIndex; i, i2 = i+1, i2+1 {
		values2[i2] = values[i]
	}

	result = v2
	return
}

func SubVectorMoveRightImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("subvector-move-right! needs a vector as its first vector argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	start := Second(args)
	if !IntegerP(start) {
		err = ProcessError(fmt.Sprintf("subvector-move-right! needs an integer as its starting index, but got %s.", String(start)), env)
		return
	}
	startIndex := int(IntegerValue(start))

	if startIndex < 0 || startIndex >= len(values) {
		err = ProcessError(fmt.Sprintf("subvector-move-right! starting index is out of bounds (0-%d), got %d.", len(values)-1, startIndex), env)
		return
	}

	end := Third(args)
	if !IntegerP(end) {
		err = ProcessError(fmt.Sprintf("subvector-move-right! needs an integer as its ending index, but got %s.", String(end)), env)
		return
	}
	endIndex := int(IntegerValue(end))

	if endIndex < startIndex || endIndex > len(values) {
		err = ProcessError(fmt.Sprintf("subvector-move-right! ending index is out of bounds (%d-%d), got %d.", startIndex, len(values), startIndex), env)
		return
	}

	v2 := Fourth(args)
	if !VectorP(v2) {
		err = ProcessError(fmt.Sprintf("subvector-move-right! needs a vector as its second vector argument, but got %s.", String(v2)), env)
		return
	}
	values2 := VectorValue(v2)

	start2 := Fifth(args)
	if !IntegerP(start2) {
		err = ProcessError(fmt.Sprintf("subvector-move-right! needs an integer as its second starting index, but got %s.", String(start2)), env)
		return
	}
	startIndex2 := int(IntegerValue(start2))

	if startIndex < 0 || startIndex >= len(values2) {
		err = ProcessError(fmt.Sprintf("subvector-move-right! starting index is out of bounds (0-%d), got %d.", len(values)-1, startIndex), env)
		return
	}

	sourceLength := endIndex - startIndex
	tailSize2 := len(values2) - startIndex2
	if sourceLength > tailSize2 {
		err = ProcessError(fmt.Sprintf("subvector-move-right! source subvector is longer than the available space in the destination (0-%d), got %d.", tailSize2, sourceLength), env)
		return
	}

	for i, i2 := endIndex-1, startIndex2+sourceLength-1; i >= startIndex; i, i2 = i-1, i2-1 {
		values2[i2] = values[i]
	}

	result = v2
	return
}

func VectorSortImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-sort needs a vector as its argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	proc := Second(args)
	if !FunctionOrPrimitiveP(proc) {
		err = ProcessError(fmt.Sprintf("vector-sort requires a function or primitive as it's second argument, but got %s.", String(proc)), env)
		return
	}

	sorted, err := MergeSort(values, proc, env)
	if err != nil {
		return
	}

	result = VectorWithValue(sorted)
	return
}

func VectorSortInPlaceImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
	v := First(args)
	if !VectorP(v) {
		err = ProcessError(fmt.Sprintf("vector-sort! needs a vector as its argument, but got %s.", String(v)), env)
		return
	}
	values := VectorValue(v)

	proc := Second(args)
	if !FunctionOrPrimitiveP(proc) {
		err = ProcessError(fmt.Sprintf("vector-sort! requires a function or primitive as it's second argument, but got %s.", String(proc)), env)
		return
	}

	sorted, err := MergeSort(values, proc, env)
	if err != nil {
		return
	}

	for i, val := range sorted {
		values[i] = val
	}

	result = v
	return
}
