<script lang="ts">
	import { onMount } from "svelte";
	import { writable, type Writable } from "svelte/store";
	import { Pointer, type Position } from "../pointer/pointer";
	import { v1 as uuid } from "uuid";
	import type { Message } from "../ws/message";

    let content = writable("")

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
        let socket = new WebSocket("ws://localhost:3000/api/listen") 
        socket.onmessage = e => {
            let message = JSON.parse(e.data) as Message
            if (message.type == "pointer-movement") {
                let pointer = message.payload as Pointer
                pointers.update(map => map.set(pointer.id, pointer))

                let image = document.getElementById(`pointer-${pointer.id}`) as HTMLImageElement
                if (image) {
                    image.style.left = pointer.position.x + pointer.scroll.x + "px"
                    image.style.top = pointer.position.y + pointer.scroll.y + "px"
                }
            } else if (message.type == "content-change") {
                content.update(() => message.payload as string)
            }
        }

        let area = document.getElementById("area") as HTMLTextAreaElement
        area.onmousemove = e => updatePointerPosition(pointer, { x: e.clientX, y: e.clientY })
        window.onscroll = () => updatePointerScroll(pointer, { x: window.scrollX, y: window.scrollY })

        pointer.subscribe(p => sendPointer(socket, p))
        content.subscribe(c => sendContent(socket, c))
    })

    function sendContent(socket: WebSocket, content: string) {
        let message = { type: "content-change", payload: content}
        send(socket, message)
    }

    function sendPointer(socket: WebSocket, pointer: Pointer) {
        let message = { type: "pointer-movement", payload: pointer }
        send(socket, message)
    }

    function send(socket: WebSocket, value: any) {
        if (socket.readyState == socket.OPEN) {
            socket.send(JSON.stringify(value))
        }
    }
</script>

{#each $pointers.values() as pointer}
    <img id={`pointer-${pointer.id}`} src="/pointer.svg" alt="Pointer" class="absolute pointer-events-none" />
{/each}
<textarea id="area" class="w-[768px] h-[2000px] resize-none p-8" bind:value={$content} />
