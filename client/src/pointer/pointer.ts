// TODO: [0, 0] should be center of a screen
type Position = {
    x: number,
    y: number,
}

function randomColor(): string {
    return "#" + Math.floor(Math.random()*16777215).toString(16);
}

export class Pointer {
    color: string

    constructor (
        public id: string,
        color: string | undefined,
        public position: Position,
    ) {
        this.color = color ? color : randomColor()
    }    
}