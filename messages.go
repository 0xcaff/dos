package main

import (
	"fmt"
	"log"
	"time"

	dosProto "github.com/caffinatedmonkey/dos/proto"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

func Read(conn *websocket.Conn, message proto.Message) error {
	format, buf, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	if format != websocket.BinaryMessage {
		log.Println("[websocket] got non binary message from", conn.RemoteAddr())
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
			time.Now().Add(time.Second),
		)
		conn.Close()

		return fmt.Errorf("dos: got non binary message")
	}

	if message == nil {
		return nil
	}

	err = proto.Unmarshal(buf, message)
	if err != nil {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
			time.Now().Add(time.Second),
		)
		conn.Close()

		return fmt.Errorf("[protobuf] failed to parse message: %#v", err)
	}

	return nil
}

func ReadMessage(conn *websocket.Conn, typ dosProto.MessageType, message proto.Message) error {
	envelope := dosProto.Envelope{}
	err := Read(conn, &envelope)
	if err != nil {
		return err
	}

	if envelope.Type != typ {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
			time.Now().Add(time.Second),
		)
		conn.Close()

		err = fmt.Errorf("Received type %s instead of type %s", envelope.Type.String(), typ.String())
		log.Println("[websocket]", err)
		return err
	}

	if message == nil {
		return nil
	}

	err = proto.Unmarshal(envelope.Contents, message)
	if err != nil {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
			time.Now().Add(time.Second),
		)
		conn.Close()

		err = fmt.Errorf("[protobuf] failed to parse envelope: %#v", err)
		log.Println("[websocket]", err)
		return err
	}

	return nil
}

func WriteMessage(conn *websocket.Conn, typ dosProto.MessageType, message proto.Message) error {
	buf, err := ZipMessage(typ, message)
	if err != nil {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, ""),
			time.Now().Add(time.Second),
		)
		conn.Close()

		log.Println("[protobuf] failed to compose message:", err)
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		log.Println("[websocket] failed to write message:", err)
		return err
	}

	return nil
}

func ZipMessage(typ dosProto.MessageType, message proto.Message) ([]byte, error) {
	var buf []byte
	var err error

	if message != nil {
		buf, err = proto.Marshal(message)

		if err != nil {
			return nil, err
		}
	}

	envelope := dosProto.Envelope{}
	envelope.Type = typ
	envelope.Contents = buf

	return proto.Marshal(&envelope)
}
