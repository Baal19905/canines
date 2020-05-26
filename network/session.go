package network

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/Baal19905/canines/interfaces"
	"github.com/Baal19905/canines/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
	"time"
)

type Session struct {
	conn      net.Conn
	id        string
	server    interfaces.IServer
	sendqueue chan *Response
	trycount  uint32
	sendexit  chan bool
	exitflag  bool
	lock      sync.Mutex
}

func NewSession(conn net.Conn, server interfaces.IServer) (*Session, error) {
	s := &Session{
		conn:     conn,
		id:       "",
		server:   server,
		exitflag: false,
	}
	id, err := newSessionId()
	if err != nil {
		return nil, err
	}
	s.id = id
	s.server.GetSessionMgr().Add(s)
	s.sendqueue = make(chan *Response, utils.ConfInfo.MaxSendQueue)
	s.trycount = 0
	return s, nil
}

func newSessionId() (string, error) {
	b := make([]byte, 64)
	n, err := rand.Read(b)
	if err != nil || n != len(b) {
		return "", fmt.Errorf("can't genorate sessionid")
	}
	return hex.EncodeToString(b), nil
}

func (s *Session) Start() {
	s.server.CallOnConnect(s.id)
	go s.startReader()
	go s.startWriter()
}

func (s *Session) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.exitflag {
		return
	}
	s.server.CallOnDisConnect(s.GetSessionID())
	s.exitflag = true
	s.sendexit <- true
	s.conn.Close()
	s.server.GetSessionMgr().Remove(s.GetSessionID())
	close(s.sendqueue)
}

func (s *Session) GetSessionID() string {
	return s.id
}

func (s *Session) PushMsg(opcode uint32, msg []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.exitflag {
		log.Warn("session is closing, msg ignored, opcode: ", opcode)
		return
	}
	response := Response{opcode, msg}
	s.sendqueue <- &response
}

func (s *Session) GetRemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *Session) startReader() {
	defer s.Stop()
	headObj := utils.HeadTemplate
	headBytes := make([]byte, headObj.GetHeadLen())
	for {
		if _, err := io.ReadFull(s.conn, headBytes); err != nil {
			if err == io.EOF {
				log.Info("client disconnected, addr: ", s.GetRemoteAddr(), ", sessionid: ", s.GetSessionID())
			} else {
				log.Error("read err, addr: ", s.GetRemoteAddr(), "sessionid: ", s.GetSessionID(), ", err: ", err.Error())
			}
			return
		}
		if err := headObj.UnMarshal(headBytes); err != nil {
			log.Error("unmarshal head failed, addr: ", s.GetRemoteAddr(), "sessionid: ", s.GetSessionID(), ", err: ", err.Error())
			return
		}
		if err := headObj.Check(); err != nil {
			log.Error("invalid head from addr: ", s.GetRemoteAddr(), "sessionid: ", s.GetSessionID(), ", err: ", err.Error())
			return
		}
		tmpBytes := make([]byte, headObj.GetBodyLen())
		if _, err := io.ReadFull(s.conn, tmpBytes); err != nil {
			if err == io.EOF {
				log.Info("client disconnected, addr: ", s.GetRemoteAddr(), ", sessionid: ", s.GetSessionID())
			} else {
				log.Error("read err, addr: ", s.GetRemoteAddr(), "sessionid: ", s.GetSessionID(), ", err: ", err.Error())
			}
			return
		}
		bodyBytes, err := headObj.PreHandle(tmpBytes)
		if err != nil {
			log.Error("PreHandle failed, addr: ", s.GetRemoteAddr(), "sessionid: ", s.GetSessionID(), ", opcode: ", headObj.GetOpcode(), "err: ", err.Error())
			continue
		}
		request := Request{
			sessionid: s.id,
			opcode:    headObj.GetOpcode(),
			msg:       bodyBytes,
		}
		s.server.GetHandleMgr().SendToHandler(&request)
	}
}

func (s *Session) startWriter() {
	defer s.Stop()
	headObj := utils.HeadTemplate
	for {
		select {
		case <-s.sendexit:
			return;
		case response := <-s.sendqueue:
			sndBody, err := headObj.PreSend(response.Msg)
			if err != nil {
				log.Error("PreSend failed, addr: ", s.GetRemoteAddr(), "sessionid: ", s.GetSessionID(), ", opcode: ", headObj.GetOpcode(), "err: ", err.Error())
				continue
			}
			headObj.SetOpcode(response.Opcode)
			headObj.SetBodyLen(uint32(len(sndBody)))
			sndStart := 0
			sndLen := headObj.GetHeadLen() + headObj.GetBodyLen()
			sndBytes := headObj.Marshal()
			sndBytes = append(sndBytes, sndBody...)
			for {
				if s.trycount >= 3 {
					log.Warn("write failed 3 times, end session, addr: ", s.GetRemoteAddr(), "sessionid: ", s.GetSessionID())
					return
				}
				s.conn.SetWriteDeadline(time.Now().Add(time.Duration(utils.ConfInfo.SendTimeOut) * time.Second))
				n, err := s.conn.Write(sndBytes[sndStart:sndLen])
				if err != nil {
					log.Error("Write failed, addr: ", s.GetRemoteAddr(), "sessionid: ", s.GetSessionID(), ", opcode: ", headObj.GetOpcode(), "err: ", err.Error())
					s.trycount++
					sndStart = n
					continue
				}
				break
				s.trycount = 0
			}
		}
	}
}
