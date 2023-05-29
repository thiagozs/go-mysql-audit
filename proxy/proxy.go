package proxy

import (
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/thiagozs/go-proxy-audit/utils"
)

type ProxyConn struct {
	MysqlConn             *net.TCPConn
	ClientConn            *net.TCPConn
	InitHandshakePacket   Packet
	FinishHandshakePacket Packet
	clientClose           bool
	uuid                  string
	username              string
	logger                zerolog.Logger
}

func New(debug bool) *ProxyConn {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	return &ProxyConn{logger: logger}
}

func (p *ProxyConn) NewMysqlConn(url string) {
	mysqlConn, _ := net.Dial("tcp", url)
	p.MysqlConn = mysqlConn.(*net.TCPConn)
	p.MysqlConn.SetNoDelay(true)
	p.MysqlConn.SetKeepAlive(true)
}

func (p *ProxyConn) NewClientConn(conn net.Conn) {
	p.ClientConn = conn.(*net.TCPConn)
	p.ClientConn.SetNoDelay(true)
	p.clientClose = false

}

func (p ProxyConn) IsClientClose() bool {
	return p.clientClose
}

func (p *ProxyConn) CloseClient() {
	p.clientClose = true
	p.ClientConn.Close()
}

func (p *ProxyConn) Close() {
	p.clientClose = true
	p.ClientConn.Close()
	p.MysqlConn.Close()
}

func (p ProxyConn) ReadMysql(whoSent string) (Packet, error) {
	packet, err := ReadPacket(p.MysqlConn, whoSent, p.logger)
	return packet, err
}

func (p ProxyConn) ReadClient(whoSent string) (Packet, error) {
	packet, err := ReadPacket(p.ClientConn, whoSent, p.logger)
	return packet, err
}

func (p ProxyConn) SendMysql(packet Packet) error {
	_, err := p.MysqlConn.Write(packet.Raw())
	return err
}

func (p ProxyConn) SendClient(packet Packet) error {
	_, err := p.ClientConn.Write(packet.Raw())
	return err
}

func (p *ProxyConn) Handshake() error {

	p.uuid = uuid.New().String()

	packet, err := p.ReadMysql("handshake-readmysql")
	if err != nil {
		return err
	}

	if err = p.SendClient(packet); err != nil {
		return err
	}

	p.InitHandshakePacket = packet
	packet, err = p.ReadClient("handshake-readclient")
	if err != nil {
		return err
	}

	// start position username
	offset := 36
	username, dlen := utils.ParseStringNUL(packet.Raw()[offset:])
	if dlen > -1 {
		p.logger.Info().
			Str("uuid", p.uuid).
			Str("user", username).
			Msg("handshake")
		p.username = string(username)
	}
	offset += dlen

	authData, dlen := utils.ParseLenEncUint(packet.Raw()[offset:])
	if dlen > -1 {
		p.logger.Info().
			Str("uuid", p.uuid).
			Str("user", p.username).
			Int("lenAuthData", int(authData)).
			Msg("handshake")
	}

	offset += dlen
	//offset += int(authData)

	info, dlen := utils.ParseStringNULEOF(packet.Raw()[offset:])
	if dlen > -1 {
		p.logger.Info().
			Str("uuid", p.uuid).
			Str("user", p.username).
			Str("info", info).
			Msg("handshake")
	}

	if err = p.SendMysql(packet); err != nil {
		return err
	}
	packet, err = p.ReadMysql("handshake-readmysql")
	if err != nil {
		return err
	}

	if err = p.SendClient(packet); err != nil {
		return err
	}
	p.FinishHandshakePacket = packet

	return nil
}

func (p ProxyConn) FakeHandshake() error {
	if err := p.SendClient(p.InitHandshakePacket); err != nil {
		return err
	}
	if _, err := p.ReadClient("handshake-readclient"); err != nil {
		return err
	}
	if err := p.SendClient(p.FinishHandshakePacket); err != nil {
		return err
	}
	return nil
}

func (p *ProxyConn) PipeClient2Mysql() {
	for {
		packet, err := p.ReadClient("pipe-client")
		if err != nil {
			p.CloseClient()
			break
		}
		if len(packet.Data()) == 1 && packet.Data()[0] == 1 {
		} else {
			p.dispatch("client", packet.Data())
			p.SendMysql(packet)
		}
	}
}

func (p *ProxyConn) PipeMysql2Client() {
	for {
		packet, err := p.ReadMysql("pipe-mysql")
		if err != nil {
			p.Close()
			break
		}
		p.dispatch("server", packet.Data())
		p.SendClient(packet)
	}
}

func (p *ProxyConn) dispatch(from string, data []byte) {
	offset := 1
	cmd := data[0]

	switch cmd {
	case COM_QUERY:
		p.logger.Info().Str("message_from", from).
			Str("uuid", p.uuid).Str("user", p.username).
			Str("query", string(data[offset:])).
			Msg("COM_QUERY")
	case COM_INIT_DB:
		p.logger.Info().Str("message_from", from).
			Str("uuid", p.uuid).Str("user", p.username).
			Str("cmd", string(data[offset:])).
			Msg("COM_INIT_DB")
	default:
		//TODO: implement others status
	}
}
