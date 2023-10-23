import { ValidationRule } from "./ValidationRule";

export class ValidationRuleUpperCaseLetters extends ValidationRule {
    private minimumNumber: number
    protected errorMessage: string

    constructor(minimumNumber: number) {
        super()
        this.minimumNumber = minimumNumber;
        this.errorMessage = `Must contain at least ${minimumNumber} uppercase letters`
    }

    validate(text: string): string | null {
        let number = Array.from(text).filter(ch => ch === ch.toUpperCase()).length 
        return number < this.minimumNumber ? this.errorMessage : null 
    }
}