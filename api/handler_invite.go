package api

import (
	"fmt"
	"net/http"
)

func (api *API) inviteHandler(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`
	<html>
		<head>
			<title>Redirecting you..</title>
			<meta property="og:type" content="website" />
			<meta name="og:title" property="og:title" content="Duel invite">
			<meta name="og:description" property="og:description" content="You have been invited to a duel">
			<meta name="og:image" property="og:image" content="https://i.imgur.com/8PlN43q.png">
		</head>
		<body style="background: #36393F">
			<p>Please wait while we redirect you.. Make sure javascript is enabled.</p>
			<script>if(!navigator.userAgent.includes("discord")) { window.location.replace("/overview?invite=%s"); }</script>
		</body>
	</html>	
	`, r.Header.Get("id"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
