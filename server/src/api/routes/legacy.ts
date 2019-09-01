import { Router, Request, Response } from "express"

const router = Router()

router.post("/register", (req: Request, res: Response) => res.redirect(308, `/api/auth/register`))
router.post("/auth", (req: Request, res: Response) => res.redirect(308, `/api/auth/login`))

export default router