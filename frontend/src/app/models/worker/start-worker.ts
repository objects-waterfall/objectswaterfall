export class StartWorkerData {
    host: string = ""
    authModel: AuthModel = new AuthModel()
}

export class AuthModel {
    authUrl: string = ""
    model: string = ""
}