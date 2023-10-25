<script lang="ts">
	import TextField from "./../../components/TextField.svelte";
	import Button from "../../components/Button.svelte";
    import { ValidationRuleRequired } from "../../validation/Required";
	import { ValidationRuleMinimumLength } from "../../validation/MinimumLength";
	import { ValidationRuleUpperCaseLetters } from "../../validation/UpperCaseLetters";
	import { goto } from "$app/navigation";
    
    let errorMessage: string | null = null
    let username: string = ""
    let password: string = ""

    let usernameValidationRules = [new ValidationRuleRequired(), new ValidationRuleMinimumLength(3)]
    let passwordValidationRules = [new ValidationRuleRequired(), new ValidationRuleMinimumLength(8), new ValidationRuleUpperCaseLetters(3)]

    let signupPromise: Promise<Response> | undefined = undefined
    function signup() {
        // TODO: Validate all TextFields

        signupPromise = fetch("http://localhost:3000/api/v1/auth/signup", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "username": username,
                "password": password,
            }),
        })

        signupPromise.then(response => {
            if (response.ok) {
                goto("/home")
            } else {
                response.json().then(body => errorMessage = body.message)
            }
        })
    }

    function wait(ms: number): Promise<any> {
        return new Promise((res) => setTimeout(res , ms))
    }
</script>

<div class="flex flex-col gap-4 w-full h-full items-center justify-center content-center">
    <div id="inputs" class="flex flex-col items-center gap-8 w-[90%] min-[440px]:w-[80%] min-[550px]:w-[70%] min-[660px]:w-[60%] min-[770px]:w-[50%] min-[990px]:w-[40%] min-[1210px]:w-[30%] min-[1430px]:w-[20%]">
        {#if errorMessage}
            <div class="text-red-500">{errorMessage}</div>
        {/if}
        <TextField bind:value={username} placeholder="Username" validationRules={usernameValidationRules} />
        <TextField bind:value={password} type="password" placeholder="Password" validationRules={passwordValidationRules} />
        <Button type="submit" onClick={signup}>
            <div class="flex flex-row gap-2 items-center">
                Sign up
                {#await signupPromise}
                    {#await wait(500)}<!--Do nothing-->{:then}
                        <img src="/images/spinner.svg" alt="Loading spinner" />    
                    {/await}
                {/await}
            </div>
        </Button>
    </div>
</div>