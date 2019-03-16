package client

import (
	"bytes"
	"fmt"
	"github.com/tgo-team/tgo-talk/protocol/mqtt"
	"github.com/tgo-team/tgo-talk/tgo"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"net"
	"os"
	"time"
)

type Client struct {
	logicConn    *net.UDPConn // 逻辑连接
	conn         net.Conn
	pro          tgo.Protocol
	cmdResultMap map[uint16]chan []byte
	opts *Options
	heartTimer *time.Ticker
}

func New(opts *Options) *Client {
	c := &Client{
		opts:         opts,
		pro:          mqtt.NewMQTTCodec(),
		cmdResultMap: map[uint16]chan []byte{},
		heartTimer: time.NewTicker(opts.MaxHeartbeatInterval),
	}
	var err error
	c.logicConn, err = c.connect()
	if err != nil {
		panic(err)
	}
	go c.recvCMD()
	return c
}

func (c *Client) connect() (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", c.opts.UDPAddress)
	if err != nil {
		fmt.Println("net.ResolveUDPAddr fail.", err)
		os.Exit(1)
	}

	c.logicConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("net.DialUDP fail.", err)
		os.Exit(1)
	}
	return c.logicConn, nil
}

func (c *Client) SendCMD(cmd uint16, payload []byte) error {
	packet := packets.NewCMDPacket(cmd, payload)
	data, err := c.pro.EncodePacket(packet)
	if err != nil {
		return err
	}
	_, err = c.logicConn.Write(data)
	return err
}

func (c *Client) recvCMD() {

	data := make([]byte, 0x7fff)
	for {
		n, _, err := c.logicConn.ReadFromUDP(data)
		if err != nil {
			panic(err)
		}
		packet, err := c.pro.DecodePacket(bytes.NewBuffer(data[:n]))
		if err != nil {
			panic(err)
		}
		go c.handlePacket(packet)
	}

}

func (c *Client) handlePacket(packet packets.Packet) {

	cmdPacket, ok := packet.(*packets.CMDPacket)
	if ok {
		notifyChan := c.cmdResultMap[cmdPacket.CMD]
		if notifyChan != nil {
			notifyChan <- cmdPacket.Payload
			delete(c.cmdResultMap, cmdPacket.CMD)
		}
	}

}

// Register 注册客户端
func (c *Client) Register(uid uint64, password string) (uint16, error) {
	resultChan := make(chan []byte, 0)
	c.cmdResultMap[2] = resultChan
	var payload bytes.Buffer
	payload.Write(packets.EncodeUint64(uid))
	payload.Write(packets.EncodeString(password))
	err := c.SendCMD(1, payload.Bytes())
	if err != nil {
		return 0, err
	}
	resultBytes := <-resultChan

	code := packets.DecodeUint16(bytes.NewBuffer(resultBytes))
	return code, nil
}

// Login 登录到服务器
func (c *Client) Login(uid uint64,password string) error  {
	var err error
	c.conn,err = net.DialTimeout("tcp",c.opts.TCPAddress,time.Second*5)
	if err!=nil {
		return err
	}
	connectPacket := packets.NewConnectPacket(uid,password)
	err = c.sendPacket(connectPacket)
	recvPacket,err := c.pro.DecodePacket(c.conn)
	if err!=nil {
		return err
	}
	connackPacket := recvPacket.(*packets.ConnackPacket)
	if connackPacket.ReturnCode != packets.ConnReturnCodeSuccess {
		return fmt.Errorf("登录失败！")
	}
	go c.heartbeat()
	go c.msgLoop()
	return nil
}

func (c *Client) sendPacket(packet packets.Packet) error  {
	packetData,err := c.pro.EncodePacket(packet)
	if err!=nil {
		return err
	}
	_,err = c.conn.Write(packetData)
	return err
}

func (c *Client) SendMsg(cid uint64,payload []byte) error  {
	msgPacket := packets.NewMessagePacket(1,cid,payload)
	err := c.sendPacket(msgPacket)
	if err!=nil {
		return err
	}
	return nil
}

func (c *Client) heartbeat()   {
	for {
		select {
		case <-c.heartTimer.C:
			pingReqPacket := packets.NewPingreqPacket()
			err := c.sendPacket(pingReqPacket)
			if err!=nil {
				fmt.Println(fmt.Sprintf("发送心跳包失败！-> %v",err))
				return
			}
		}
	}
}

func (c *Client) msgLoop()  {
	for {
		packet,err := c.pro.DecodePacket(c.conn)
		if err!=nil {
			fmt.Println(fmt.Sprintf("解码消息失败！-> %v",err))
			return
		}
		msgPacket,ok := packet.(*packets.MessagePacket)
		if ok {
			fmt.Println(fmt.Sprintf("-> %v\r",string(msgPacket.Payload)))
			fmt.Print(">")
		}
	}
}