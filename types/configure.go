/*
 * Copyright (C) 2016 Meng Shi
 */

package types

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "unsafe"
)

var (
    ConfigOk    =  0
    ConfigError = -1
)

type Configure struct {
    *Log
    *File

     resource     string
     fileName     string
     commandType  int
     moduleType   int64
     value        interface{}
     configure    ConfigureIf
}

type ConfigureIf interface {
    Parse() int
    ReadToken() int
}

func NewConfigure(log *Log) *Configure {
    return &Configure{
        Log          : log,
        File : NewFile(log),
    }
}

func (c *Configure) SetName(file string) int {
    if file == "" {
        return Error
    }

    if c.File.SetName(file) == Error {
        return Error
    }

    return Ok
}

func (c *Configure) GetName() string {
    return c.File.GetName()
}

func (c *Configure) SetFileName(fileName string) int {
    if fileName == "" {
        return Error
    }

    c.fileName = fileName

    return Ok
}

func (c *Configure) GetFileName() string {
    return c.fileName
}

func (c *Configure) SetFileType(fileType string) int {
    if fileType == "" {
        return Error
    }

    if c.SetName(fileType) == Error {
        return Error
    }

    return Ok
}

func (c *Configure) GetFileType() string {
    if fileType := c.GetName(); fileType != "" {
        return fileType
    }

    return ""
}

func (c *Configure) SetFile(action IO) int {
    if action == nil {
        return Error
    }

    if c.File.Set(action) == Error {
        return Error
    }

    return Ok
}

func (c *Configure) GetFile() IO {
    if file := c.File.Get(); file != nil {
        return file
    }

    return nil
}

func (c *Configure) SetResource(resource string) int {
    if resource == "" {
        return Error
    }

    c.resource = resource

    return Ok
}

func (c *Configure) GetResource(resource string) string {
    return c.resource
}

func (c *Configure) Get() ConfigureIf {
    log := c.Log.Get()


    file := c.File.Get()
    if file == nil {
        file = NewFile(c.Log)
    }

    if file.Open(c.resource) == Error {
        log.Error("configure open file error")
        return nil
    }

    if file.Read() == Error {
        log.Error("configure read file error")
        goto JMP_CLOSE
        return nil
    }

    if content := file.Type().GetContent(); content != nil {
        c.content = content
    } else {
        log.Warn("not found content: %d\n", 10)
    }

JMP_CLOSE:
    if file.Close() == Error {
        log.Warn("file close error: %d\n", 10)
        return nil
    }

    return c.configure
}

func (c *Configure) Set(configre ConfigureIf) int {
    if configre == nil {
        return Error
    }

    c.configure = configre

    return Ok
}

func (c *Configure) SetModuleType(moduleType int64) int {
    if moduleType <= 0 {
        return Error
    }

    c.moduleType = moduleType

    return Ok
}

func (c *Configure) SetCommandType(commandType int) int {
    if commandType <= 0 {
        return Error
    }

    c.commandType = commandType

    return Ok
}

func (c *Configure) GetValue() interface{} {
   return c.value
}

func SetFlag(cycle *Cycle, command *Command, p *unsafe.Pointer) int {
    if cycle == nil || p == nil {
        return Error
    }

    field := (*bool)(unsafe.Pointer(uintptr(*p) + command.Offset))
    if field == nil {
        return Error
    }

    configure := cycle.GetConfigure()
    if configure == nil {
        return Error
    }

    flag := configure.GetValue()
    if flag == true {
        *field = true
    } else if flag == false {
        *field = false
    } else {
        return Error
    }

    /*
    if command.Post != nil {
        post := command.Post.(DvrConfPostType);
        post.Handler(cf, post, *p);
    }
    */

    return Ok
}

func SetString(cycle *Cycle, command *Command, p *unsafe.Pointer) int {
    if cycle == nil || p == nil {
        return Error
    }

    field := (*string)(unsafe.Pointer(uintptr(*p) + command.Offset))
    if field == nil {
        return Error
    }

    configure := cycle.GetConfigure()
    if configure == nil {
        return Error
    }

    strings := configure.GetValue()
    if strings == nil {
        return Error
    }

    *field = strings.(string)

    return Ok
}

func SetNumber(cycle *Cycle, command *Command, p *unsafe.Pointer) int {
    if cycle == nil || p == nil {
        return Error
    }

    field := (*int)(unsafe.Pointer(uintptr(*p) + command.Offset))
    if field == nil {
        return Error
    }

    configure := cycle.GetConfigure()
    if configure == nil {
        return Error
    }

    number := configure.GetValue()
    if number == nil {
        return Error
    }

    *field = number.(int)

    return Error
}

func (c *Configure) Parse(cycle *Cycle) int {
    log := c.Log.Get()

    if configure := c.Get(); configure != nil {
        if configure.Parse() == Error {
            return Error
        }

        return Ok
    }

    // TODO default process
    if c.value == nil {
        content := c.GetContent()
        if content == nil {
            log.Error("configure content: %s, filename: %s, size: %d\n",
                      content, c.GetFileName(), c.GetSize())

            return Error
        }

        error := yaml.Unmarshal(content, &c.value)
        if error != nil {
            log.Error("yanm unmarshal error: %s\n", error)
            return Error
        }
    }

    switch v := c.value.(type) {

    case []interface{} :
        for _, value := range v {
            c.value = value
            c.Parse(cycle)
        }

    case map[interface{}]interface{}:
        if c.doParse(v, cycle) == Error {
            return Error
        }

    default:
        fmt.Println("unknown")
    }

    return Ok
}

func (c *Configure) doParse(materialized map[interface{}]interface{}, cycle *Cycle) int {
    log := c.Log.Get()

    flag := Ok

    for key, value := range materialized {

        if key != nil && value != nil {
            flag = Ok
        }

        name := key.(string)
        found := false

        for m := 0; flag != Error && !found && Modules[m] != nil; m++ {
            module := Modules[m]
								    /*
            if module.Type != CONFIG_MODULE &&
               module.Type != c.moduleType {

                continue;
            }
            */

            commands := module.Commands;
            if commands == nil {
                continue;
            }

            //fmt.Printf("%s, %X, %X, %d\n", name, module.Type, c.moduleType, m)

            for i := 0; commands[i].Name.Len != 0; i++ {

                command := commands[i]

                if len(name) == command.Name.Len &&
                        name == command.Name.Data.(string) {

                				found = true;

                    if command.Type & c.commandType == 0 &&
                       command.Type & MAIN_CONFIG == 0 {

                        //flag = Error
																				    found = false
                        break
                    }

                    //log.Error("directive \"%s\" is not allowed here", name)
                    //					flag = Error
                    context := cycle.GetContext(module.Index)

                    c.value = value
																    if cycle.SetConfigure(c) == Error {
                        flag = Error
																				    break
                    }

                    command.Set(cycle, &command, context)
                }
            }
        }

        if !found {
            log.Error("unkown")

            flag = Error
            break
        }

        if flag == Error {
            break
        }
    }

    if flag == Error {
        return ConfigError
    }

    return ConfigOk
}

func (c *Configure) ReadToken() int {
    fmt.Println("configure read token")
    return Ok
}
