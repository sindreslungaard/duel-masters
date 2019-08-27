declare namespace Express {
    export interface Request {
        user?: import("./src/models/user").IUser
    }
 }