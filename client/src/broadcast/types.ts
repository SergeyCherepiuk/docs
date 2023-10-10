export type Message = {
    messageType: string
    rawMessage: any
}

export type User = {
    id: string
    pointer: Pointer
    selection: Selection
}

export type Pointer = {
    position: Position
    scroll: Position
}

export type Position = {
    x: number
    y: number
}

export type Selection = {
    start: number
    end: number
}