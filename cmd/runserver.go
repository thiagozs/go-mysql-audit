package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mercadobitcoin/go-proxy-audit/cfg"
	"github.com/mercadobitcoin/go-proxy-audit/proxy"
	"github.com/spf13/cobra"
)

var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Run proxy for mysql",
	Run:   runserver,
}

var (
	mysqlPort  string
	proxyPort  string
	proxyDebug bool
	Version    string
	Build      string
)

func init() {
	runserverCmd.PersistentFlags().StringVar(&mysqlPort, "mysql", "3306", "port for mysql server")
	runserverCmd.PersistentFlags().StringVar(&proxyPort, "proxy", "33060", "port for proxy server")
	runserverCmd.PersistentFlags().BoolVar(&proxyDebug, "debug", false, "set debug log in bytes stdout")
	rootCmd.AddCommand(runserverCmd)
}

func runserver(cmd *cobra.Command, args []string) {

	if cmd.Flags().Lookup("proxy").Changed {
		cfg.Config.Proxy = proxyPort
	}

	if cmd.Flags().Lookup("mysql").Changed {
		cfg.Config.Mysql = mysqlPort
	}

	if cmd.Flags().Lookup("debug").Changed {
		cfg.Config.Debug = proxyDebug
	}

	log.Printf("[SERVER] Proxy MySQL Audit\n")
	log.Printf("[SERVER] Version: %s\n", Version)
	log.Printf("[SERVER] Build  : %s\n", Build)
	log.Printf("[SERVER] proxy server on host=%s, mysql server listening host=%s...\n",
		cfg.Config.Proxy, cfg.Config.Mysql)

	if strings.Contains(mysqlPort, ":") {
		splt := strings.Split(mysqlPort, ":")
		if splt[0] == "" {
			mysqlPort = fmt.Sprintf("%s:%s", "", splt[1])
		}
	} else if !strings.Contains(proxyPort, ":") {
		mysqlPort = fmt.Sprintf("%s:%s", "", mysqlPort)
	}

	if strings.Contains(proxyPort, ":") {
		splt := strings.Split(proxyPort, ":")
		if splt[0] == "" {
			proxyPort = fmt.Sprintf("%s:%s", "", splt[1])
		}
	} else if !strings.Contains(proxyPort, ":") {
		proxyPort = fmt.Sprintf("%s:%s", "", proxyPort)
	}

	server, _ := net.Listen("tcp", proxyPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	conns := clientConns(server)

	for {
		select {
		case <-quit:
			log.Printf("[SERVER] Shutdown proxy server...\n")
			server.Close()
			os.Exit(0)
		case conn := <-conns:

			go func(conn net.Conn, p *proxy.ProxyConn) {
				p.NewClientConn(conn)
				p.NewMysqlConn(mysqlPort)
				err := p.Handshake()
				if err != nil {
					p.Close()
				}

				go p.PipeMysql2Client()
				go p.PipeClient2Mysql()

				for {
					if !p.IsClientClose() {
						continue
					}
					p.NewClientConn(conn)
					err := p.FakeHandshake()
					if err != nil {
						p.CloseClient()
					}
					go p.PipeClient2Mysql()
				}

			}(conn, proxy.New(cfg.Config.Debug))
		}
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			client, _ := listener.Accept()
			if client == nil {
				continue
			}
			ch <- client
		}
	}()
	return ch
}
