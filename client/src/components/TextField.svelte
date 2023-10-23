<script lang="ts">
	import type { ValidationRule } from "../validation/ValidationRule";

    export let type: string = "text"
    export let placeholder: string
    export let value: string
    export let validationRules: Array<ValidationRule> = []
    let errorMessage: string | null = null

    let onfocusout = validate
    let oninput = () => { if (errorMessage) validate() }

    function validate() {
        for (let rule of validationRules) {
            errorMessage = rule.validate(value)
            if (errorMessage) break
        }
    }

    function typeAction(el: HTMLInputElement) {
        el.type = type
    }
</script>

<div class="flex flex-col w-full items-start gap-2">
    <input use:typeAction {placeholder} bind:value on:focusout={onfocusout} on:input={oninput}
        class=" p-4 w-full bg-gray-100 border-[3px] border-black border-solid rounded-none outline-none">
    {#if errorMessage}
        <span class="text-red-500">{errorMessage}</span>
    {/if}
</div>

<style>
    /* TODO: Keep animation playing when focus goes away */
    input:focus {
        animation: pulse-black 1.5s;
    }

    @keyframes pulse-black {
        0% {
            box-shadow: 0 0 0 0 rgba(0, 0, 0, 0.7);
        }
        70% {
            box-shadow: 0 0 0 10px rgba(0, 0, 0, 0);
        }
        100% {
            box-shadow: 0 0 0 0 rgba(0, 0, 0, 0);
        }
    }
</style>