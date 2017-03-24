/*
 * Copyright (C) 2017 Meng Shi
 */

package types

type Channel struct {
    *Cycle
    *File

     upstream   chan string
     downstream chan string

     channal    ChannelIf
}

type ChannelIf interface {
    push(name string, data string) int
    pull(name string) int
}

func NewChannel() *Channel {
    return &Channel{}
}

func (c *Channel) push(name string, data string) int {
    c.upstream <- data

    return Ok
}

func (c *Channel) pull(name string) string {
    data := <-c.downstream

    return data
}
