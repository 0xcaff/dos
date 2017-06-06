package main

import (
	"fmt"
	"time"

	dosProto "github.com/0xcaff/dos/proto"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func Read(conn *websocket.Conn, message proto.Message) error {
	format, buf, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	if format != websocket.BinaryMessage {
		log.WithFields(log.Fields{
			"remoteAddr": conn.RemoteAddr(),
		}).Warning("got non binary websocket message")

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
		log.Warning("websocket ", err)
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

		err = fmt.Errorf("failed to parse protobuf envelope: %#v", err)
		log.Warning("websocket ", err)
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

		log.Warning("failed to compose protobuf message: ", err)
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		log.Warning("failed to write websocket message: ", err)
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
