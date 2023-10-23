export abstract class ValidationRule {
    protected abstract errorMessage: string

    abstract validate(text: string): string | null
}