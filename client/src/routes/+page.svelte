<script lang="ts">
	import { onMount } from "svelte";
	import { get, writable } from "svelte/store";
	import { Pointer } from "../pointer/pointer";
	import { v1 as uuid } from "uuid";
	import type { Message } from "../ws/message";

    let content = writable("")

    let id = uuid()
    let pointers = writable(new Map<string, Pointer>([]))

    onMount(() => {
        let socket = new WebSocket("ws://localhost:3000/api/listen") 
        socket.onmessage = e => {
            let message = JSON.parse(e.data) as Message
            if (message.type == "pointer-movement") {
                let pointer = message.payload as Pointer
                pointers.update(map => map.set(pointer.id, pointer))

                let image = document.getElementById(`pointer-${pointer.id}`) as HTMLImageElement
                if (image) {
                    image.style.left = pointer.position.x + "px"
                    image.style.top = pointer.position.y + "px"
                }
            } else if (message.type == "content-change") {
                content.update(() => message.payload as string)
            }
        }

        let area = document.getElementById("area") as HTMLTextAreaElement
        area.onmousemove = e => {
            let pointer = new Pointer(id, {
                x: e.clientX + window.scrollX,
                y: e.clientY + window.scrollY
            })
            sendPointer(socket, pointer)
        }
        window.onscroll = () => {
            let pointer = get(pointers).get(id)
            if (pointer) {
                sendPointer(socket, pointer)
            }
        }

        content.subscribe(c => {
            let message = { type: "content-change", payload: c}
            send(socket, message)
        })
    })

    function sendPointer(socket: WebSocket, pointer: Pointer) {
        pointers.update(map => map.set(id, pointer))
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
    {#if pointer.id != id}
        <img id={`pointer-${pointer.id}`} src="/pointer.svg" alt="Pointer" class="absolute pointer-events-none" />
    {/if}
{/each}
<textarea id="area" class="w-[768px] h-[2000px] resize-none p-8" bind:value={$content} />
