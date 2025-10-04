export class Result<T> {
    Data: T
    Err: any

    constructor(data: T, err: string) {
        this.Data = data
        this.Err = err
    }
}