package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/faruk/secscan/backend/internal/scanner"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, service *scanner.Service) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		api.POST("/scan", func(c *gin.Context) {
			var req scanner.ScanRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			req.TargetURL = strings.TrimSpace(req.TargetURL)
			record, err := service.Start(req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusAccepted, gin.H{
				"scan_id": record.ID,
				"status":  record.Status,
				"target":  record.TargetURL,
			})
		})

		api.GET("/scan/:id", func(c *gin.Context) {
			record, ok := service.Get(c.Param("id"))
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "scan not found"})
				return
			}
			c.JSON(http.StatusOK, record)
		})

		api.GET("/scan/:id/stream", func(c *gin.Context) {
			id := c.Param("id")
			_, ok := service.Get(id)
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "scan not found"})
				return
			}

			events, cancel := service.Subscribe(id)
			defer cancel()

			c.Writer.Header().Set("Content-Type", "text/event-stream")
			c.Writer.Header().Set("Cache-Control", "no-cache")
			c.Writer.Header().Set("Connection", "keep-alive")
			c.Writer.Flush()

			for {
				select {
				case <-c.Request.Context().Done():
					return
				case event := <-events:
					payload, _ := json.Marshal(event)
					fmt.Fprintf(c.Writer, "data: %s\n\n", payload)
					c.Writer.Flush()
				}
			}
		})

		api.GET("/scan/:id/report.pdf", func(c *gin.Context) {
			record, ok := service.Get(c.Param("id"))
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "scan not found"})
				return
			}

			body := []byte(buildPDFLikeReport(record))
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s-report.pdf", record.ID))
			c.Data(http.StatusOK, "application/pdf", body)
		})
	}
}

func buildPDFLikeReport(record *scanner.ScanRecord) string {
	lines := []string{
		"%PDF-1.1",
		"SecScan report",
		fmt.Sprintf("ID: %s", record.ID),
		fmt.Sprintf("Target: %s", record.TargetURL),
		fmt.Sprintf("Grade: %s", record.Grade),
		fmt.Sprintf("Summary: %s", record.Summary),
	}

	for _, mod := range record.Modules {
		lines = append(lines, fmt.Sprintf("- %s [%s]: %s", mod.Name, mod.Grade, mod.Summary))
	}
	lines = append(lines, "%%EOF")
	return strings.Join(lines, "\n")
}

