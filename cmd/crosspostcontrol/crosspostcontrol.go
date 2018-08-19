package main

import (
	"github.com/mattermost/mattermost-server/model"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/plugin"
	"fmt"
	"errors"
	"github.com/DSchalla/MM-CrosspostControl/crosspostcontrol"
)

type CrosspostControlPlugin struct {
	plugin.MattermostPlugin
	server *crosspostcontrol.Server
	config *crosspostcontrol.Config
}

func (c *CrosspostControlPlugin) OnActivate() error {
	mlog.Debug("[CROSSPOSTCONTROL-PLUGIN] OnActivate Hook Start")
	var err error
	c.readConfig()
	c.server, err = crosspostcontrol.NewServer(c.API, *c.config)
	if err != nil {
		mlog.Debug(fmt.Sprintf("[CROSSPOSTCONTROL-PLUGIN] NewBotServer returned error: %s", err))
	}
	mlog.Debug("[CROSSPOSTCONTROL-PLUGIN] OnActivate Hook End")

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (c *CrosspostControlPlugin) OnConfigurationChange() error {
	err := c.readConfig()
	if err != nil {
		return err
	}
	return c.reloadConfig()
}

func (c *CrosspostControlPlugin) MessageWillBePosted(context *plugin.Context, post *model.Post) (*model.Post, string) {
	if post.Props["from_crosspostcontrol"] != nil && post.Props["from_crosspostcontrol"].(bool) == true {
		return post, ""
	}
	c.API.LogDebug("[CROSSPOSTCONTROL-PLUGIN] MessageWillBePosted Hook Start")
	post, rejectMessage := c.server.HandleMessage(post, true)
	c.API.LogDebug("[CROSSPOSTCONTROL-PLUGIN] MessageWillBePosted Hook End")
	return post, rejectMessage
}

func (c *CrosspostControlPlugin) readConfig() error {
	c.config = &crosspostcontrol.Config{}
	err := c.API.LoadPluginConfiguration(c.config)
	return err
}

func (c *CrosspostControlPlugin) reloadConfig() error {
	if c.server != nil {
		c.server.ReloadConfig(*c.config)
	}

	return nil
}

func main() {
	plugin.ClientMain(&CrosspostControlPlugin{})
}
