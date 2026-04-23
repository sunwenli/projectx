package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/sunwenli/projectx/core"
)

type MessageType byte

const (
	// 枚举类型如果是int,只用赋值第一个，之后的会自增
	// 如果是 byte 的要全部都赋值才行
	MessageTypeTx        MessageType = 0x1
	MessageTypeBlock     MessageType = 0x2
	MessageTypeGetBlocks MessageType = 0x3
	MessageTypeStatus    MessageType = 0x4
	MessageTypeGetStatus MessageType = 0x5
)

type RPC struct {
	// From    NetAddr
	From    net.Addr
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

func (msg *Message) Byte() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type RPCProcessor interface {
	ProcessMessage(*DecodeMessage) error
}

type DecodeMessage struct {
	From net.Addr
	Data any
}
type RPCDecodeFunc func(RPC) (*DecodeMessage, error)

func DefaultDecodeRPCFunc(rpc RPC) (*DecodeMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s ：%s", rpc.From, err)
	}

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming msg")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewTxGobDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	case MessageTypeBlock:
		block := new(core.Block)
		if err := block.Decode(core.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: block,
		}, nil
	case MessageTypeStatus:
		statusmessage := new(StatusMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(statusmessage); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: statusmessage,
		}, nil
	case MessageTypeGetStatus:
		return &DecodeMessage{
			From: rpc.From,
			Data: &GetStatusMessage{},
		}, nil
	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}
}
