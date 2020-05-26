package network

import (
	log "github.com/sirupsen/logrus"
	"github.com/Baal19905/canines/handler"
	"github.com/Baal19905/canines/interfaces"
	"github.com/Baal19905/canines/logfile"
	"github.com/Baal19905/canines/utils"
	"net"
)

type TcpServer struct {
	Addr         string
	SessionMgr   interfaces.ISessionMgr
	HandlerMgr   interfaces.IHandleMgr
	OnConnect    interfaces.NotifyCallback
	OnDisConnect interfaces.NotifyCallback
	exit         chan bool
}

func NewTcpServer(cfg string) interfaces.IServer {
	utils.ConfInfo.Load(cfg)
	log.SetOutput(logfile.SetLogfile(utils.ConfInfo.LogPath))
	log.SetLevel(log.Level(utils.ConfInfo.LogLevel))
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000000",
		PrettyPrint:     true,
	})
	tcp := &TcpServer{
		Addr:         utils.ConfInfo.Addr,
		SessionMgr:   NewSessionManage(),
		OnConnect:    nil,
		OnDisConnect: nil,
	}
	tcp.HandlerMgr = handler.NewHandlerMgr(tcp.SessionMgr)
	return tcp
}

func (s *TcpServer) start() {
	s.HandlerMgr.StartPool()
	go func() {
		ln, err := net.Listen("tcp", s.Addr)
		if err != nil {
			log.Error("Listen on ", s.Addr, "failed, err: ", err.Error())
			return
		}
		log.Info("Listening on ", s.Addr, " ...")
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Error("Accept failed, err: ", err.Error())
				return
			}
			if utils.ConfInfo.MaxCli <= s.SessionMgr.Len() {
				conn.Close()
				log.Warn("MaxCli: ", utils.ConfInfo.MaxCli, ", rejected")
				continue
			}
			session, err := NewSession(conn, s)
			if err != nil {
				log.Error("create session failed, err: ", err.Error())
				continue
			}
			session.Start()
			log.Info("accepted client, addr: ", session.GetRemoteAddr(), ", sessionid: ", session.GetSessionID(), ", cliNum: ", s.SessionMgr.Len())
		}
	}()
}

func (s *TcpServer) Stop() {
	s.SessionMgr.Stop()
	s.HandlerMgr.StopPool()
	s.exit <- true
}

func (s *TcpServer) Serve() {
	s.start()
	select {
	case <-s.exit:
		return
	}
	log.Info("server stopped")
}

func (s *TcpServer) RegisterRouter(id uint32, router interfaces.RouterCallback) {
	s.HandlerMgr.AddRouter(id, router)
}

func (s *TcpServer) Push(sid string, opcode uint32, msg []byte) {
	session := s.SessionMgr.GetSession(sid)
	if session == nil {
		return
	}
	session.PushMsg(opcode, msg)
}

func (s *TcpServer) RegisterOnConnect(connect interfaces.NotifyCallback) {
	s.OnConnect = connect
}

func (s *TcpServer) RegisterOnDisConnect(disconnect interfaces.NotifyCallback) {
	s.OnDisConnect = disconnect
}

func (s *TcpServer) CallOnConnect(sid string) {
	if s.OnConnect != nil {
		session := s.GetSessionMgr().GetSession(sid)
		if session == nil {
			log.Warn("nil OnConnect")
			return
		}
		s.OnConnect(session)
	}
}

func (s *TcpServer) CallOnDisConnect(sid string) {
	if s.OnDisConnect != nil {
		session := s.GetSessionMgr().GetSession(sid)
		if session == nil {
			log.Warn("nil OnDisConnect")
			return
		}
		s.OnDisConnect(session)
	}
}

func (s *TcpServer) GetHandleMgr() interfaces.IHandleMgr {
	return s.HandlerMgr
}

func (s *TcpServer) GetSessionMgr() interfaces.ISessionMgr {
	return s.SessionMgr
}
