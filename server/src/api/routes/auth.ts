import { createError } from './../error'
import { Router, Request, Response } from "express"
import User from "../../models/user"
import { hash } from "bcryptjs"
import jwt from "jsonwebtoken"
import uniqid from "uniqid"

const router = Router()

router.post("/register", async (req: Request, res: Response) => {

    // TODO: Captcha and or max users per ip
    
    if(!req.body.username || !req.body.username || !req.body.email) {
        return createError(res, 400, "Username, email and password is required")
    }

    if(req.body.username.length < 3 || req.body.username.length > 20) {
        return createError(res, 400, "Username must be between 3 to 20 characters long")
    }

    if(req.body.password.length < 6) {
        return createError(res, 400, "Password must be at least 6 characters long")
    }

    if(!/^\w+-?\w+(?!-)$/.test(req.body.username)) {
        return createError(res, 400, "The username can only contain letters, numbers and individual dashes")
    }

    try {

        let userByName = await User.findOne({'username': { $regex : new RegExp(req.body.username, "i") }})

        if(userByName) {
            return createError(res, 400, "The username has already been taken")
        }

        let userByEmail = await User.findOne({'email': { $regex : new RegExp(req.body.email, "i") }})

		if(userByEmail) {
            return createError(res, 400, "The email has already been registered")
        }

        let hashedPassword = await hash(req.body.password, 10)

        let ttl = Math.round(Date.now() / 1000 + 31536000)

        let token = await jwt.sign({
            exp: ttl, // 1 year
            data: {
                username: req.body.username,
                email: req.body.email   
            }
        }, process.env.SECRET_KEY)

        let user = await User.create({
			uid: uniqid(),
			username: req.body.username,
			email: req.body.email,
			password: hashedPassword,
			sessions: [{
				token,
				ip: req.connection.remoteAddress,
				expires: ttl
			}],
			permissions: []
        })
        
        return res.json({
			uid: user.uid,
			username: user.username,
			email: user.email,
			permissions: user.permissions,
			token: token
		})

    }
    catch(error) {
        return createError(res, 500, "An error occured during registration")
    }

})

export default router