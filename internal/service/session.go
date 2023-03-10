package service

type sessionData struct {
	Username string
}

type Session struct {
	Data map[string]*sessionData
}

func NewSession() *Session {
	s := new(Session)
	s.Data = make(map[string]*sessionData)
	return s
}

func (s *Session) Init(username string) string {
	sessionId := GenerateId()
	data := &sessionData{Username: username}
	s.Data[sessionId] = data
	return sessionId
}

func (s *Session) Get(sessionId string) string {
	data := s.Data[sessionId]
	if data == nil {
		return ""
	}
	return data.Username
}
