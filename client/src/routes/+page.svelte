<script lang="ts">
	import { onMount } from "svelte";
	import { writable } from "svelte/store";
	import { Pointer } from "../pointer/pointer";
	import { v1 as uuid } from "uuid";

    // TODO: [0, 0] should be center of a screen
    let id = uuid()
    let pointers = writable(new Map<string, Pointer>([]))

    onMount(() => {
        let socket = new WebSocket("ws://localhost:3000/api/mouse") 
        socket.onmessage = e => {
            let pointer = JSON.parse(e.data) as Pointer
            pointers.update(map => map.set(pointer.id, pointer))
        }

        let plain = document.getElementById("plain") as HTMLDivElement
        plain.onmousemove = e => {
            let pointer = new Pointer(
                id,
                $pointers.get(id)?.color,
                { x: e.clientX, y: e.clientY },
            )

            pointers.update(map => map.set(pointer.id, pointer))
            send(socket, pointer)
        }
    })    

    function send(socket: WebSocket, value: any) {
        if (socket.readyState == socket.OPEN) {
            socket.send(JSON.stringify(value))
        }
    }

    
</script>

<div id="plain" class="w-full h-full bg-gray-300 text-center">
    {#each $pointers.values() as pointer}
        {#if pointer.id != id}
            <svg class="absolute w-full h-full">
                <circle cx={pointer.position.x} cy={pointer.position.y} r=15 fill={pointer.color} />
            </svg>
        {/if}
    {/each}
</div>