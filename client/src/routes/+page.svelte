<script lang="ts">
	import { onMount } from "svelte";
	import { writable } from "svelte/store";
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
            } else if (message.type == "content-change") {
                content.update(() => message.payload as string)
            }
        }

        let plain = document.getElementById("plain") as HTMLDivElement
        plain.onmousemove = e => {
            let pointer = new Pointer(
                id, $pointers.get(id)?.color, { x: e.clientX, y: e.clientY },
            )
            pointers.update(map => map.set(pointer.id, pointer))

            let message = { type: "pointer-movement", payload: pointer }
            send(socket, message)
        }

        content.subscribe(c => {
            console.log(c)
            let message = { type: "content-change", payload: c}
            send(socket, message)
        })
    })    

    function send(socket: WebSocket, value: any) {
        if (socket.readyState == socket.OPEN) {
            socket.send(JSON.stringify(value))
        }
    }
</script>

<div id="plain" class="w-full h-full bg-gray-300 text-center">
    <svg class="absolute w-full h-full pointer-events-none">
        {#each $pointers.values() as pointer}
            {#if pointer.id != id}
                <circle cx={pointer.position.x} cy={pointer.position.y} r=10 fill={pointer.color} />
            {/if}
        {/each}
    </svg>
    <textarea id="area" class="w-[768px] h-full resize-none p-8" bind:value={$content} />
</div>