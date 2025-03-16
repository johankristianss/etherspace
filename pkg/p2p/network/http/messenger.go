package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	p2pnet "github.com/johankristianss/evrium/pkg/p2p/network"
	log "github.com/sirupsen/logrus"
)

type HTTPMessenger struct {
	Address string
}

func NewHTTPMessenger(address string) *HTTPMessenger {
	gin.SetMode(gin.ReleaseMode)
	return &HTTPMessenger{
		Address: address,
	}
}

func (m *HTTPMessenger) Send(msg p2pnet.Message, ctx context.Context) error {
	log.WithFields(log.Fields{"To": msg.To.Addr, "From": msg.From.Addr, "Type": msg.Type}).Info("Sending message")

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+msg.To.Addr+"/message",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.WithFields(log.Fields{"To": msg.To.Addr, "From": msg.From.Addr, "Type": msg.Type, "Error": err}).Error("Failed to send message")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{"To": msg.To.Addr, "From": msg.From.Addr, "Type": msg.Type}).Error("Status code", resp.Status, " is not OK")
		return fmt.Errorf("failed to send message: %s", resp.Status)
	}

	return nil
}

func (m *HTTPMessenger) ListenForever(msgChan chan p2pnet.Message, ctx context.Context) error {
	router := gin.Default()

	router.POST("/message", func(c *gin.Context) {
		var msg p2pnet.Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Failed to bind JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.WithFields(log.Fields{"To": msg.To.Addr, "From": msg.From.Addr, "Type": msg.Type}).Info("Received message")
		msgChan <- msg

		c.JSON(http.StatusOK, gin.H{"status": "received"})
	})

	server := &http.Server{
		Addr:    m.Address,
		Handler: router,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	log.WithFields(log.Fields{"Address": m.Address}).Info("Starting HTTP server")
	return server.ListenAndServe()
}
