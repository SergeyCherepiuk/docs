<script lang="ts">
	import TextField from "./../../components/TextField.svelte";
	import Button from "../../components/Button.svelte";
    import { ValidationRuleRequired } from "../../validation/Required";
	import { ValidationRuleMinimumLength } from "../../validation/MinimumLength";
	import { ValidationRuleUpperCaseLetters } from "../../validation/UpperCaseLetters";
    
    let username: string = ""
    let password: string = ""

    let usernameValidationRules = [new ValidationRuleRequired(), new ValidationRuleMinimumLength(3)]
    let passwordValidationRules = [new ValidationRuleRequired(), new ValidationRuleMinimumLength(8), new ValidationRuleUpperCaseLetters(3)]
</script>

<div class="flex w-full h-full items-center justify-center content-center">
    <div id="inputs" class="flex flex-col items-center gap-8 w-[90%] min-[440px]:w-[80%] min-[550px]:w-[70%] min-[660px]:w-[60%] min-[770px]:w-[50%] min-[990px]:w-[40%] min-[1210px]:w-[30%] min-[1430px]:w-[20%]">
        <TextField bind:value={username} placeholder="Username" validationRules={usernameValidationRules} />
        <TextField bind:value={password} type="password" placeholder="Password" validationRules={passwordValidationRules} />
        <Button type="submit" onClick={() => { /* TODO: Login */ }}>Login</Button>
    </div>
</div>