package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	// We can do this natively just as easily, but this framework makes
	// the examples a bit more clear.
	"github.com/gin-gonic/gin"
)

var (
	beats, port int
	heartbeat   *time.Ticker
	router      *gin.Engine
)

// Do all of this stuff first.
func init() {
	// In this example, we will hard code the port.  Later the environment
	// will dictate.
	port = 7718
	// Set up the heartbeat ticker.
	heartbeat = time.NewTicker(10 * time.Second)

	// Setup the service router.
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	// These are the services we will be listening for.
	router.GET("/beats", GetHeartbeatCount)
	router.GET("/ping", PingTheAPI)

} // func

// GetHeartbeatCount sends the number of times the heartbeat ticker has
// fired since the program started.
func GetHeartbeatCount(c *gin.Context) {
	content := gin.H{"payload": beats}
	c.JSON(http.StatusOK, content)
}

// PingTheAPI lets the caller know we are alive.
func PingTheAPI(c *gin.Context) {
	content := gin.H{"payload": "pong"}
	c.JSON(http.StatusOK, content)
}

// Manage the processes.
func main() {
	// Dispatch a process into the background.
	go func() {
		// Now run it forever.
		for {
			// Watch for stuff to happen.
			select {
			// When the Heartbeat ticker is fired, execute this.
			case <-heartbeat.C:
				beats++
				fmt.Printf("bump,Bump... @ %s\n", time.Now().UTC())
			} // select
		} // for
	}() // go func

	// After the 'go func' is dispatched, start the server and listen on the
	// specified port.
	fmt.Printf("ready on port %d\n", port)
	router.Run(":" + strconv.Itoa(port))
} // func
