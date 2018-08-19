package crosspostcontrol

import (
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/model"
	"crypto/sha1"
	"encoding/hex"
)

const RejectMessage = "Message intercepted by Crosspost Control" 

type Config struct {
	Matching               string
	Mode string
}

type Server struct {
	config Config
	api  plugin.API
	history map[string]*model.Post
}

func NewServer(api plugin.API, config Config) (*Server, error) {
	s := Server{}
	s.api = api
	s.config = config
	s.history = make(map[string]*model.Post)
	return &s, nil
}


func (s *Server) HandleMessage(post *model.Post, intercept bool) (*model.Post, string) {
	h := sha1.New()
	h.Write([]byte(post.UserId + post.Message))
	hash := hex.EncodeToString(h.Sum(nil))
	previousPost, ok := s.history[hash]

	if !ok {
		s.history[hash] = post
		return post, ""
	}

	channel, _ := s.api.GetChannel(previousPost.ChannelId)

	response := &model.Post{
		Message: "You posted a message with the same contents in '" + channel.Name + "' - Please avoid crossposting to keep things organized.",
		ChannelId: post.ChannelId,
	}
	s.api.SendEphemeralPost(post.UserId, response)
	return nil, RejectMessage
}

func (s *Server) ReloadConfig(config Config) {
	s.config = config
}
