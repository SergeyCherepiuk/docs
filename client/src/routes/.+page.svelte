<script lang="ts">
	import { writable, type Writable } from "svelte/store";
	import { v1 as uuid } from "uuid";
	import { onMount } from "svelte";
	import type { User, Message, Position, Selection } from "../broadcast/types";

    let content: Writable<{ data: string, toSend: boolean }> = writable({ data: "", toSend: false }) 

    let id = uuid()

    let user: Writable<User> = writable({
        id: id,
        pointer: {
            position: { x: 0, y: 0 }, 
            scroll: { x: 0, y: 0 }
        },
        selection: { start: 0, end: 0 },
    })
    let users = writable(new Map<string, User>([])) 

    onMount(() => {
        let area = document.getElementById("area") as HTMLTextAreaElement

        let socket = new WebSocket("ws://localhost:3000/api/v1/connect")
        socket.onopen = () => send({ messageType: "user", rawMessage: $user })
        socket.onmessage = e => {
            // TODO: Sanitize corresponding of the data on disconnect message

            let message = JSON.parse(e.data) as Message
            if (message.messageType == "pointer") {
                let user = message.rawMessage as User
                users.update(map => map.set(user.id, user))

                let image = document.getElementById(`pointer-${user.id}`) as HTMLImageElement
                if (image) {
                    image.style.left = user.pointer.position.x + user.pointer.scroll.x + "px"
                    image.style.top = user.pointer.position.y + user.pointer.scroll.y + "px"
                }
            } else if (message.messageType == "selection") {
                console.log(message.rawMessage) // NOTE: rawMessage of type User
            } else if (message.messageType == "content") {
                content.set({ data: message.rawMessage as string, toSend: false })
            } else if (message.messageType == "user") {
                let user = message.rawMessage as User
                users.update(map => map.set(user.id, user))
            }
        }

        area.onmousemove = e => updatePointerPosition({ x: e.clientX, y: e.clientY })
        window.onscroll = () => updatePointerScroll({ x: window.scrollX, y: window.scrollY }) 

        area.onselect = e => {
            let target = (e.target as HTMLTextAreaElement)
            updateSelection({ start: target.selectionStart, end: target.selectionEnd })
        }
        
        area.oninput = e => {
            let value = (e.target as HTMLTextAreaElement).value;
            updateContent(value, true)
        }

        function updatePointerPosition(position: Position) {
            user.update(u => {
                u.pointer.position = position
                send({ messageType: "pointer", rawMessage: u })
                return u
            })
        }

        function updatePointerScroll(scroll: Position) {
            user.update(u => {
                u.pointer.scroll = scroll
                send({ messageType: "pointer", rawMessage: u })
                return u
            })
        }
        
        function updateSelection(selection: Selection) {
            user.update(u => {
                u.selection = selection
                send({ messageType: "selection", rawMessage: u })
                return u
            })
        }

        function updateContent(data: string, toSend: boolean) {
            content.set({ data: data, toSend: toSend })
            if (toSend) send({ messageType: "content", rawMessage: data}) 
        }

        function send(message: Message) {
            if (socket.readyState == socket.OPEN) {
                socket.send(JSON.stringify(message))
            }
        }
    })
</script>

<!-- <div class="flex w-full min-w-fit h-full justify-center bg-gray-300"> -->
    {#each $users.keys() as id}
        <img id={`pointer-${id}`} src="/images/pointer.svg" alt="Pointer" class="absolute pointer-events-none" />
    {/each}
    <textarea id="area" class="w-[768px] h-[2000px] resize-none p-8" value={$content.data} />
<!-- </div> -->