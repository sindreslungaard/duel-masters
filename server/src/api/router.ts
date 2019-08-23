import { Router } from "express"
import auth from "./routes/auth"
import legacy from "./routes/legacy"

const router = Router()

router.use("/auth", auth)
router.use("/", legacy)

export default router