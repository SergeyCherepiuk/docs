<script lang="ts">
	import { onMount } from "svelte";
	import { get, writable, type Writable } from "svelte/store";
	import { Pointer, type Position } from "../pointer/pointer";
	import { v1 as uuid } from "uuid";

    let content: Writable<{ data: string, toSend: boolean }> = writable({ data: "", toSend: false })

    function updateContent(data: string, toSend: boolean) {
        content.set({ data: data, toSend: toSend })
    }

    let pointer = writable(new Pointer(uuid(), { x: 0, y: 0 }, { x: 0, y: 0 }))
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

    onMount(() => {
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

        let contentSock = new WebSocket("ws://localhost:3000/api/content") 
        contentSock.onmessage = e => updateContent(e.data, false)

        let area = document.getElementById("area") as HTMLTextAreaElement
        area.onmousemove = e => updatePointerPosition(pointer, { x: e.clientX, y: e.clientY })
        window.onscroll = () => updatePointerScroll(pointer, { x: window.scrollX, y: window.scrollY })

        area.addEventListener("input", event => {
            let value = (event.target as HTMLTextAreaElement).value;
            updateContent(value, true)
        });

        pointer.subscribe(p => send(pointerSock, p))
        content.subscribe(c => {
            if (c.toSend) send(contentSock, c.data)
        })
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

{#each $pointers.values() as pointer}
    <img id={`pointer-${pointer.id}`} src="/pointer.svg" alt="Pointer" class="absolute pointer-events-none" />
{/each}
<textarea id="area" class="w-[768px] h-[2000px] resize-none p-8" value={$content.data} />