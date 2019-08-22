import { createServer } from "http"
import logger from "./utils/logger"
import api from "./api"
import * as server from "./net/server"

export const bootstrap = async () => {

    logger.info("Starting up...")

    const web = createServer(api)

    server.connect(web)

    web.listen(process.env.PORT || 3000, () => {
        logger.info(`Listening on port ${process.env.PORT || 3000}`)
    })

} 