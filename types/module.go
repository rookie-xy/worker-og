/*
 * Copyright (C) 2016 Meng Shi
 */

package types

import "unsafe"

var (
    Ok     =  0
    Error  = -1
    Again  = -2
    Ignore = -3
)

const (
    MODULE_V1  = 0
    CONTEXT_V1 = 0

    SYSTEM_MODULE = 0x10000000
    CONFIG_MODULE = 0xF0000000
)

type InitFunc func(cycle *Cycle) int
type MainFunc func(cycle *Cycle) int

type Module struct {
    CtxIndex   uint
    Index      uint
    Context    unsafe.Pointer
    Commands   []Command
    Type       int64
    Init       InitFunc
    Main       MainFunc
}

var Modules []*Module

func Load(modules []*Module, module *Module) []*Module {
    if modules == nil && module == nil {
        return nil
    }

    modules = append(modules, module)

    return modules
}
