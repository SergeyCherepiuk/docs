<script lang="ts">
	import type { MouseEventHandler } from "svelte/elements";

    export let type: "button" | "submit" | "reset" = "button"
    export let onClick: MouseEventHandler<HTMLButtonElement>
</script>

<button {type} on:click={onClick} class="py-4 px-8 w-min bg-black text-white">
    <slot />
</button>