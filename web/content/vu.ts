type ValueSubscriber<T> = (value: T) => void;

class ValueUpdater<T> {
    private _subs: ValueSubscriber<T>[] = [];
    private _last: T;

    constructor(initial: T) {
        this._last = initial;
    }

    subscribe(vs: ValueSubscriber<T>): () => void {
        this._subs.push(vs);
        return () => {
            this._subs = this._subs.filter((v) => v !== vs);
        };
    }

    get last(): T {
        return this._last;
    }

    update(v: T) {
        this._last = v;
        this._subs.forEach((vs) => vs(v));
    }
}

export { ValueSubscriber, ValueUpdater }