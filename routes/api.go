package routes

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	// tools stores the vote counts for each tool
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

// normalizeToolName converts a tool name to lowercase for case-insensitive lookup
// but returns the original capitalized name if it exists in the tools map.
func normalizeToolName(name string) string {
	mutex.RLock()
	defer mutex.RUnlock()

	// Convert input to lowercase for case-insensitive lookup
	lowerName := strings.ToLower(name)

	// Iterate through the tools map to find a match
	for tool := range tools {
		if strings.ToLower(tool) == lowerName {
			return tool // Return the original capitalized name
		}
	}

	return "" // Tool not found
}

// GetTools returns the list of tools and their votes
func GetTools(c *gin.Context) {
	mutex.RLock() // Lock for reading
	defer mutex.RUnlock()

	// Return the tools map as JSON
	c.JSON(http.StatusOK, tools)
}

// GetToolByName returns the details of a specific tool
func GetToolByName(c *gin.Context) {
	name := c.Param("name")

	// Normalize tool name and get the original capitalized name
	normalizedName := normalizeToolName(name)
	if normalizedName == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Tool '%s' not found", name)})
		return
	}

	mutex.RLock()
	defer mutex.RUnlock()

	// Return the tool details
	c.JSON(http.StatusOK, gin.H{"tool": normalizedName, "votes": tools[normalizedName]})
}

// VoteForTool increments the vote count for a specific tool
func VoteForTool(c *gin.Context) {
	name := c.Param("tool")

	// Normalize tool name and get the original capitalized name
	normalizedName := normalizeToolName(name)
	if normalizedName == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Tool '%s' not found", name)})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Increment vote count
	tools[normalizedName]++
	log.Printf("Vote received for tool: %s (new vote count: %d)", normalizedName, tools[normalizedName])

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s vote count updated", normalizedName), "votes": tools[normalizedName]})
}
