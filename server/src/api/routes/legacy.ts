import { Router, Request, Response } from "express"

const router = Router()

router.get("/register", (req: Request, res: Response) => res.redirect(301, `/auth/register`))
router.get("/login", (req: Request, res: Response) => res.redirect(301, `/auth/login`))

export default router