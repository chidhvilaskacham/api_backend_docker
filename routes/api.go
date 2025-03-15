package routes

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"github.com/gin-gonic/gin"
)

var (
	tools = map[string]int{
		"Ansible":       0,
		"Visual_studio": 0,
		"Docker":        0,
		"Prometheus":    0,
		"Git":           0,
		"Jenkins":       0,
	}
	mutex sync.RWMutex // Use RWMutex to allow concurrent reads
)

// GetTools returns the list of tools and their votes
func GetTools(c *gin.Context) {
	mutex.RLock() // Lock for reading
	defer mutex.RUnlock()

	c.JSON(http.StatusOK, tools)
}

// GetToolByName returns the details of a specific tool
func GetToolByName(c *gin.Context) {
	name := c.Param("name")

	mutex.RLock()
	defer mutex.RUnlock()

	// Case-insensitive lookup
	for tool, votes := range tools {
		if strings.EqualFold(tool, name) {
			c.JSON(http.StatusOK, gin.H{"tool": tool, "votes": votes})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Tool not found"})
}

// VoteForTool increments the vote count for a specific tool
func VoteForTool(c *gin.Context) {
	name := c.Param("tool")

	mutex.Lock()
	defer mutex.Unlock()

	// Case-insensitive lookup
	for tool := range tools {
		if strings.EqualFold(tool, name) {
			tools[tool]++
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s vote count updated", tool), "votes": tools[tool]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Tool not found"})
}
