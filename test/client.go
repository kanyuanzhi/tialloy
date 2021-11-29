package test

import (
	"github.com/kanyuanzhi/tialloy/tinet"
	"io"
	"log"
	"net"
	"time"
)

type Client struct {
	ClientID uint32
	MsgID    uint32

	Conn net.Conn
}

func (c *Client) Start() {
	log.Println("Client start")

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	c.Conn = conn
	if err != nil {
		log.Println(err.Error())
		return
	}
	message := tinet.NewMessage(c.MsgID, []byte("hello111111"))
	dp := tinet.NewDataPack()
	go func() {
		for {
			binaryMessage, _ := dp.Pack(message)
			_, err = conn.Write(binaryMessage)
			if err != nil {
				log.Println(err.Error())
				return
			}
			message.SetMsgID(c.MsgID)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			dataHeadBuf := make([]byte, dp.GetHeadLen())
			_, err = io.ReadFull(conn, dataHeadBuf)
			if err != nil {
				log.Println("read head err", err)
				return
			}

			recvMessage, err := dp.Unpack(dataHeadBuf)
			if err != nil {
				log.Println("unpack head err", err)
				return
			}

			var dataBuf []byte
			if recvMessage.GetDataLen() > 0 {
				dataBuf = make([]byte, recvMessage.GetDataLen())
				if _, err := io.ReadFull(conn, dataBuf); err != nil {
					log.Println("read message data err", err.Error())
					return
				}
			}
			recvMessage.SetData(dataBuf)

			log.Printf("received message id=%d, data=%s", recvMessage.GetMsgID(), recvMessage.GetData())
		}
	}()

	time.Sleep(10 * time.Second)
	conn.Close()

	for {
		time.Sleep(100 * time.Second)
	}
}

func NewClient(clientID uint32, msgID uint32) *Client {
	return &Client{
		ClientID: clientID,
		MsgID:    msgID,
	}
}
