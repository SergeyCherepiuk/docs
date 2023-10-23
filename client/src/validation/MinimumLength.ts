import { ValidationRule } from "./ValidationRule"

export class ValidationRuleMinimumLength extends ValidationRule {
    private minimumLength: number
    protected errorMessage: string

    constructor(minimumLength: number) {
        super()
        this.minimumLength = minimumLength
        this.errorMessage = `Must be at least ${minimumLength} character(s) long`
    }

    validate(text: string): string | null {
        return text.length < this.minimumLength ? this.errorMessage : null
    }
} 