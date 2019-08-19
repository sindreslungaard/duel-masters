import { createServer } from "http"
import logger from "./utils/logger"
import api from "./api"

export const bootstrap = async () => {

    logger.info("Starting up...")

    const web = createServer(api)

    web.listen(process.env.PORT || 3000, () => {
        logger.info(`Listening on port ${process.env.PORT || 3000}`)
    })

}