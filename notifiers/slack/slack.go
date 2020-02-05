package slack

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/mitchellh/mapstructure"
	slackApi "github.com/nlopes/slack"
	"github.com/pkg/errors"
	"statusbay/notifiers/common"
	"statusbay/serverutil"
	watcherCommon "statusbay/watcher/kubernetes/common"
	"strings"
	"time"
)

var (
	NoTokenErr = errors.New("slack token is required")
)

var (
	//UpdateSlackUserInterval interval for update slack user list
	UpdateSlackUserInterval = time.Hour
)

// NewSlack returns a slack notifier
func NewSlack(urlBase string) common.Notifier {
	return &Manager{
		config:  Config{MessageTemplates: defaultMessageConfig},
		urlBase: urlBase,
	}
}

// LoadConfig maps a generic notifier config (map[string]interface{}) to a concrete type
func (sl *Manager) LoadConfig(notifierConfig common.NotifierConfig) (err error) {
	sl.config = Config{}
	if err = mapstructure.Decode(notifierConfig, &sl.config); err != nil {
		return
	}

	// validate config
	if sl.config.Token == "" {
		return NoTokenErr
	}

	// init slack client
	sl.client = slackApi.New(sl.config.Token)

	return
}

// sendToAll sends the provided message to all valid recipients
func (sl *Manager) sendToAll(stage ReportStage, message watcherCommon.DeploymentReport, color MessageColor) {
	var (
		deployBy string
		err      error
	)

	if deployBy, err = sl.getUserIdByEmail(message.DeployBy); err != nil {
		deployBy = message.DeployBy
	} else {
		deployBy = fmt.Sprintf("by <@%s>", deployBy)
	}

	status := strings.ToUpper(string(message.Status))
	link := fmt.Sprintf("%s/%s", sl.urlBase, message.URI)

	for _, to := range distinct(append(message.To, sl.config.DefaultChannels...)) {
		if to == "" {
			continue
		}
		toChannel, err := sl.GetChannelId(to)
		if err == nil {
			attachment := slackApi.Attachment{
				Title:   replacePlaceholders(sl.config.MessageTemplates[stage].Title, status, link, deployBy),
				Pretext: replacePlaceholders(sl.config.MessageTemplates[stage].Pretext, status, link, deployBy),
				Text:    replacePlaceholders(sl.config.MessageTemplates[stage].Text, status, link, deployBy),
				Color:   string(color),
				// TODO:: add cluster + namespace name
				Fields: []slackApi.AttachmentField{
					{
						Title: "Application Name:",
						Value: message.Name,
						Short: false,
					},
				},
			}
			sl.send(toChannel, attachment)

		} else {
			log.WithFields(log.Fields{
				"to": to,
			}).Debug("Slack id not found")
		}

	}
}

// ReportStarted sends a deployment start report
func (sl *Manager) ReportStarted(message watcherCommon.DeploymentReport) {
	sl.sendToAll(started, message, blue)
}

// ReportDeleted sends a deployment deleted report
func (sl *Manager) ReportDeleted(message watcherCommon.DeploymentReport) {
	sl.sendToAll(deleted, message, red)
}

// ReportEnded sends a deployment end report
func (sl *Manager) ReportEnded(message watcherCommon.DeploymentReport) {
	color := green

	switch message.Status {
	case watcherCommon.DeploymentSuccessful:
		color = green
	case watcherCommon.DeploymentCanceled:
		color = yellow
	case watcherCommon.DeploymentStatusFailed:
		color = red
	}

	sl.sendToAll(ended, message, color)
}

// Serve will periodically check slack for a change in the list of existing users
func (sl *Manager) Serve() serverutil.StopFunc {
	sl.updateUsers()

	ctx, cancelFn := context.WithCancel(context.Background())
	stopped := make(chan bool)
	go func() {
		for {
			select {
			case <-time.After(UpdateSlackUserInterval):
				sl.updateUsers()
			case <-ctx.Done():
				log.Warn("Slack Loop has been shut down")
				stopped <- true
				return
			}
		}
	}()
	return func() {
		cancelFn()
		<-stopped
	}
}

// updateUsers updates the list of users available in slack
func (sl *Manager) updateUsers() {
	currentUsers := map[string]string{}

	users, err := sl.client.GetUsers()
	if err != nil {
		log.WithError(err).Error("unable to update user list")
		return
	}

	for _, user := range users {
		if !user.Deleted && user.Profile.Email != "" {
			currentUsers[user.Profile.Email] = user.ID
		}
	}
	if len(currentUsers) != len(sl.emailToUser) {
		sl.emailToUser = currentUsers
		log.Info(fmt.Sprintf("Found %d slack users", len(sl.emailToUser)))
	}
}

// getUserIdByEmail return slack user by email
func (sl *Manager) getUserIdByEmail(email string) (string, error) {
	if userId, ok := sl.emailToUser[email]; !ok {
		//log.WithField("email", email).Warn("Slack user by email not found")
		return "", errors.New("slack user by email not found")
	} else {
		return userId, nil
	}
}

// send sends a slack notification to user
func (sl *Manager) send(channelID string, attachment slackApi.Attachment) {
	_, _, err := sl.client.PostMessage(channelID, slackApi.MsgOptionAttachments(attachment), slackApi.MsgOptionAsUser(true))
	if err != nil {

		log.WithError(err).WithFields(log.Fields{
			"channel_id": channelID,
		}).Error("Error when trying to send post message")
	}
	log.WithFields(log.Fields{
		"channel_id": channelID,
	}).Debug("Slack message was sent")
}

// GetChannelId returns the channel id. if is it email, search the user channel id by his email
func (sl *Manager) GetChannelId(to string) (string, error) {
	if strings.HasPrefix(to, "#") {
		return to, nil
	}
	return sl.getUserIdByEmail(to)
}

// distinct de-duplicates a slice
func distinct(inputSlice []string) []string {
	keys := make(map[string]struct{})
	var list []string
	for _, entry := range inputSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}

// replacePlaceholders replaces Status, Link and DeployedBy placeholders from the templates with the actual values
func replacePlaceholders(input, status, link, deployedBy string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(input, common.StatusPlaceholder, status),
			common.LinkPlaceholder, link),
		common.DeployedByPlaceholder, deployedBy)
}