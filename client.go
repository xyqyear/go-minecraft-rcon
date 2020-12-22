package main

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
)

var packetIDEnum int32 = 0

func getNewPacketID() int32 {
	packetIDEnum++
	return packetIDEnum
}

// Client is a minecraft rcon client
type Client struct {
	conn       net.Conn
	isLoggedIn bool
}

// SendPacket sends rcon packet with packetID, type and payload
func (client *Client) SendPacket(packetID int32, _type int32, payload string) (err error) {
	lengthOfPacket := int32(4 + 4 + len(payload) + 2)
	buf := make([]byte, 4+lengthOfPacket)
	binary.LittleEndian.PutUint32(buf, uint32(lengthOfPacket))
	binary.LittleEndian.PutUint32(buf[4:], uint32(packetID))
	binary.LittleEndian.PutUint32(buf[8:], uint32(_type))
	copy(buf[12:], payload)
	client.conn.Write(buf)
	return
}

// SendLoginPacket is a wrapper of SendPacket where packetID = 3
func (client *Client) SendLoginPacket(packetID int32, password string) error {
	return client.SendPacket(packetID, 3, password)
}

// SendCommandPacket is a wrapper of SendPacket where packetID = 2
func (client *Client) SendCommandPacket(packetID int32, cmd string) error {
	return client.SendPacket(packetID, 2, cmd)
}

// SendPaddingPacket sends a invalid packet with a packetID different from the packetID for the actually command
// to identify the end of a response sequence
func (client *Client) SendPaddingPacket(packetID int32) error {
	return client.SendPacket(packetID, 0, "")
}

// RecvPacket receives rcon packet
func (client *Client) RecvPacket() (packetID int32, _type int32, payload string, err error) {
	var lengthOfPacket int32
	err = binary.Read(client.conn, binary.LittleEndian, &lengthOfPacket)
	if err != nil {
		return
	}
	bytesBuf := make([]byte, lengthOfPacket)
	io.ReadFull(client.conn, bytesBuf)
	packetID = int32(binary.LittleEndian.Uint32(bytesBuf[:4]))
	_type = int32(binary.LittleEndian.Uint32(bytesBuf[4:8]))
	payload = string(bytesBuf[8 : lengthOfPacket-2])
	return
}

// NewClient dials the server address and return a Client
func NewClient(serverAddress string) (client Client, err error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return
	}
	client.conn = conn
	client.isLoggedIn = false
	return
}

// Login logs the client in
func (client *Client) Login(password string) error {
	packetID := getNewPacketID()
	err := client.SendLoginPacket(packetID, password)
	if err != nil {
		return err
	}
	RecvPacketID, _, _, err := client.RecvPacket()
	if err != nil {
		return err
	}
	if RecvPacketID == -1 {
		return errors.New("Wrong password")
	} else if RecvPacketID != packetID {
		return errors.New("Unexpected packet id while login")
	}
	client.isLoggedIn = true
	return nil
}

// SendCommandNaively is a naive implementation of sending rcon command
// it sends the command and return the first packet that has the packet id sent by this function
// it only receive the first response packet
// it discards all the other packets that does not have the expected packet id
// thus, it does not support go routine
func (client *Client) SendCommandNaively(command string) (string, error) {
	packetID := getNewPacketID()
	err := client.SendCommandPacket(packetID, command)
	if err != nil {
		return "", err
	}
	// for some reason we need a delay
	time.Sleep(time.Millisecond)
	err = client.SendPaddingPacket(packetID + 2<<29)
	if err != nil {
		return "", err
	}
	var fullResponse string
	for {
		responsePacketID, _type, response, err := client.RecvPacket()
		if err != nil {
			return "", err
		}
		if _type == 0 {
			if responsePacketID == packetID {
				fullResponse += response
			} else if responsePacketID == packetID+2<<29 {
				return fullResponse, nil
			} else {
				return "", errors.New("unknown packet id")
			}
		} else {
			return "", errors.New("unexpected packet type")
		}
	}
}
