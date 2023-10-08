export type Position = { x: number, y: number }

export class Pointer {
    constructor (
        public id: string,
        public position: Position,
        public scroll: Position,
    ) {}
}