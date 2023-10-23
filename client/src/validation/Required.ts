import { ValidationRule } from "./ValidationRule"

export class ValidationRuleRequired extends ValidationRule {
    protected errorMessage: string = "The field is required"

    validate(text: string): string | null {
        return text.trim().length == 0 ? this.errorMessage : null
    }
}