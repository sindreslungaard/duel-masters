import { Router } from "express"
import auth from "./routes/auth"

const router = Router()

router.use("/auth", auth)

export default router