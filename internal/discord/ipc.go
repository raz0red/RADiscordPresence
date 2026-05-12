// Package discord implements the Discord Rich Presence IPC protocol.
package discord

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	opHandshake uint32 = 0
	opFrame     uint32 = 1
)

// AppID is the Discord application client ID for RAPresence.
const AppID = "1503853960786874402"

// Activity is the Rich Presence payload sent to Discord.
// Assets must be nested under "assets" and image/text fields must NOT appear
// at the top level — Discord silently ignores top-level image fields.
type Activity struct {
	Name       string              `json:"name,omitempty"`
	Type       int                 `json:"type"`
	Details    string              `json:"details,omitempty"`
	State      string              `json:"state,omitempty"`
	Assets     *ActivityAssets     `json:"assets,omitempty"`
	Timestamps *ActivityTimestamps `json:"timestamps,omitempty"`
	Buttons    []Button            `json:"buttons,omitempty"`
}

// ActivityAssets holds the image/text fields for a Rich Presence activity.
type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

// ActivityTimestamps holds the start/end epoch seconds for an activity.
type ActivityTimestamps struct {
	Start int64 `json:"start,omitempty"`
}

// Button is a clickable link shown on the Discord presence card.
type Button struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Client manages a Discord IPC connection.
type Client struct {
	conn  net.Conn
	nonce int
}

// Connect opens a Discord IPC connection and performs the handshake.
func Connect(clientID string) (*Client, error) {
	conn, err := openConn()
	if err != nil {
		return nil, fmt.Errorf("discord IPC: %w", err)
	}
	c := &Client{conn: conn}
	if err := c.handshake(clientID); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("discord handshake: %w", err)
	}
	return c, nil
}

// SetActivity updates the Discord Rich Presence. Pass nil to clear.
func (c *Client) SetActivity(a *Activity) error {
	c.nonce++
	type setArgs struct {
		PID      int       `json:"pid"`
		Activity *Activity `json:"activity"`
	}
	type frame struct {
		Cmd   string  `json:"cmd"`
		Args  setArgs `json:"args"`
		Nonce string  `json:"nonce"`
	}
	req := frame{
		Cmd:   "SET_ACTIVITY",
		Args:  setArgs{PID: os.Getpid(), Activity: a},
		Nonce: strconv.Itoa(c.nonce),
	}
	if err := c.writeFrame(opFrame, req); err != nil {
		return err
	}
	_, _, err := c.readFrame()
	return err
}

// Close closes the Discord IPC connection.
func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
}

func (c *Client) handshake(clientID string) error {
	type hs struct {
		V        int    `json:"v"`
		ClientID string `json:"client_id"`
	}
	if err := c.writeFrame(opHandshake, hs{V: 1, ClientID: clientID}); err != nil {
		return err
	}
	op, data, err := c.readFrame()
	if err != nil {
		return err
	}
	if op != opFrame {
		return fmt.Errorf("expected opcode %d, got %d", opFrame, op)
	}
	var resp struct {
		Evt string `json:"evt"`
	}
	_ = json.Unmarshal(data, &resp)
	if resp.Evt != "READY" {
		return fmt.Errorf("expected READY event, got %q", resp.Evt)
	}
	return nil
}

func (c *Client) writeFrame(op uint32, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	hdr := make([]byte, 8)
	binary.LittleEndian.PutUint32(hdr[0:4], op)
	binary.LittleEndian.PutUint32(hdr[4:8], uint32(len(data)))
	_ = c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err = c.conn.Write(append(hdr, data...))
	return err
}

func (c *Client) readFrame() (uint32, []byte, error) {
	_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	hdr := make([]byte, 8)
	if _, err := io.ReadFull(c.conn, hdr); err != nil {
		return 0, nil, fmt.Errorf("read header: %w", err)
	}
	op := binary.LittleEndian.Uint32(hdr[0:4])
	n := binary.LittleEndian.Uint32(hdr[4:8])
	if n > 1<<20 {
		return 0, nil, fmt.Errorf("frame too large (%d bytes)", n)
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(c.conn, buf); err != nil {
		return 0, nil, fmt.Errorf("read payload: %w", err)
	}
	return op, buf, nil
}
