import { Router } from "express"
import auth from "./routes/auth"
import legacy from "./routes/legacy"
import match from "./routes/match"

const router = Router()

router.use("/auth", auth)
router.use("/match", match)
router.use("/", legacy)

export default router