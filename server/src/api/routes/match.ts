import { Router, Request, Response } from "express"
import { createError } from './../error'
import authorized from "../middleware/authorized"
import { createMatch } from "../../game/match"

const router = Router()

router.post("/", authorized, (req: Request, res: Response) => {

    // TODO: Rate limit

    if(!req.body.name) {
        return createError(res, 400, "Missing required property \"name\"")
    }

    if(req.body.visibility !== "public" && req.body.visibility !== "private") {
		return createError(res, 400, "Invalid value for property \"visibility\"")
    }
    
    if(!req.user) {
        createError(res, 500, "Request not linked to a user")
    }

    let match = createMatch(req.user.uid, req.body.name, req.body.description || "")

    return res.json({
        matchUid: match.id,
        inviteId: match.inviteId
    })
    
})

export default router