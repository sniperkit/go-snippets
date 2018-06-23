class Foo {
    private _bar: string

    get bar(): string {
        return this._bar;
    }

    set bar(newBar: string) {
        this._bar = newBar;
    }
}
