import { Response } from "express"

export const createError = (res: Response, code: number, message: string) => {
    res.status(code).json({
        code,
        message
    })
}