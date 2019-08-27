import { Request, Response, NextFunction } from 'express'
import User from "../../models/user"
import { createError } from './../error'

const authorized = async (req: Request, res: Response, next: NextFunction) => {

    const token = req.header('authorization')

    let user = await User.findOne({'sessions': { $elemMatch: { token: token }}})
    
    if(!user) {
		return createError(res, 403, "Unauthorized")
    }
    
    req.user = user

    next()

}

export default authorized