<script lang="ts">
	import { onMount } from "svelte";
	import { get, writable, type Writable } from "svelte/store";
	import { Pointer, type Position } from "../pointer/pointer";
	import { stringify, v1 as uuid } from "uuid";
	import type { Selection } from "../selection/selection";

    let content: Writable<{ data: string, toSend: boolean }> = writable({ data: "", toSend: false })

    function updateContent(data: string, toSend: boolean) {
        content.set({ data: data, toSend: toSend })
    }

    let id = uuid()

    let pointer = writable(new Pointer(id, { x: 0, y: 0 }, { x: 0, y: 0 }))
    let pointers = writable(new Map<string, Pointer>([]))

    function updatePointerPosition(pointer: Writable<Pointer>, position: Position) {
        pointer.update(p => {
            p.position = { x: position.x, y: position.y }
            return p
        })
    }

    function updatePointerScroll(pointer: Writable<Pointer>, scroll: Position) {
        pointer.update(p => {
            p.scroll = { x: scroll.x, y: scroll.y }
            return p
        })
    }

    let selection = writable({ id: id, start: 0, end: 0 })
    let selections = writable(new Map<string, Selection>([]))

    function updateSelection(s: { start: number, end: number }) {
        selection.update(old => {
            old.start = s.start,
            old.end = s.end
            return old
        })
    }

    onMount(() => {
        let area = document.getElementById("area") as HTMLTextAreaElement

        let pointerSock = new WebSocket("ws://localhost:3000/api/pointer") 
        pointerSock.onmessage = e => {
            let pointer = JSON.parse(e.data) as Pointer
            pointers.update(map => map.set(pointer.id, pointer))

            let image = document.getElementById(`pointer-${pointer.id}`) as HTMLImageElement
            if (image) {
                image.style.left = pointer.position.x + pointer.scroll.x + "px"
                image.style.top = pointer.position.y + pointer.scroll.y + "px"
            }
        }

        area.onmousemove = e => updatePointerPosition(pointer, { x: e.clientX, y: e.clientY })
        window.onscroll = () => updatePointerScroll(pointer, { x: window.scrollX, y: window.scrollY })

        let contentSock = new WebSocket("ws://localhost:3000/api/content") 
        contentSock.onmessage = e => updateContent(e.data, false) 

        area.oninput = e => {
            let value = (e.target as HTMLTextAreaElement).value;
            updateContent(value, true)
        }

        let selectionSock = new WebSocket("ws://localhost:3000/api/selection") 
        selectionSock.onmessage = e => {
            let selection = JSON.parse(e.data) as Selection
            console.log(selection)
        }

        area.onselect = e => {
            let target = (e.target as HTMLTextAreaElement)
            updateSelection({ start: target.selectionStart, end: target.selectionEnd })
        }

        pointer.subscribe(p => send(pointerSock, p))
        content.subscribe(c => {
            if (c.toSend) send(contentSock, c.data)
        })
        selection.subscribe(s => send(selectionSock, s))
    })

    function send(socket: WebSocket, value: any) {
        if (socket.readyState != socket.OPEN) {
            return
        }

        if (typeof value == "string") {
            socket.send(value)
        } else {
            socket.send(JSON.stringify(value))
        }
    }
</script>

<div class="flex w-full h-full justify-center bg-gray-300">
    {#each $pointers.values() as pointer}
        <img id={`pointer-${pointer.id}`} src="/pointer.svg" alt="Pointer" class="absolute pointer-events-none" />
    {/each}
    <textarea id="area" class="w-[768px] h-[2000px] resize-none p-8" value={$content.data} />
</div>