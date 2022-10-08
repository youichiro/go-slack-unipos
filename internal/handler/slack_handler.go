package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/youichiro/go-slack-my-unipos/internal/util"
)

type SlackHandler struct {
	SigninSecret string
	Token        string
}

func generateModalRequest() slack.ModalViewRequest {
	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", "My App", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Please enter your name", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	firstNameText := slack.NewTextBlockObject("plain_text", "First Name", false, false)
	firstNameHint := slack.NewTextBlockObject("plain_text", "First Name Hint", false, false)
	firstNamePlaceholder := slack.NewTextBlockObject("plain_text", "Enter your first name", false, false)
	firstNameElement := slack.NewPlainTextInputBlockElement(firstNamePlaceholder, "firstName")
	// Notice that blockID is a unique identifier for a block
	firstName := slack.NewInputBlock("First Name", firstNameText, firstNameHint, firstNameElement)

	lastNameText := slack.NewTextBlockObject("plain_text", "Last Name", false, false)
	lastNameHint := slack.NewTextBlockObject("plain_text", "Last Name Hint", false, false)
	lastNamePlaceholder := slack.NewTextBlockObject("plain_text", "Enter your first name", false, false)
	lastNameElement := slack.NewPlainTextInputBlockElement(lastNamePlaceholder, "lastName")
	lastName := slack.NewInputBlock("Last Name", lastNameText, lastNameHint, lastNameElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			firstName,
			lastName,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func (h SlackHandler) HandleSlash(c *gin.Context) {
	err := util.VerifySigningSecret(c, h.SigninSecret)
	if err != nil {
		c.IndentedJSON(401, gin.H{"message": err})
	}

	s, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.IndentedJSON(500, gin.H{"message": err})
	}

	switch s.Command {
	case "/unipos":
		api := slack.New(h.Token)
		modalRequest := generateModalRequest()
		_, err = api.OpenView(s.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
	default:
		c.IndentedJSON(500, gin.H{"message": "hg"})
	}
}

func (h SlackHandler) HandleModal(c *gin.Context) {
	err := util.VerifySigningSecret(c, h.SigninSecret)
	if err != nil {
		c.IndentedJSON(401, gin.H{"message": err})
	}

	var i slack.InteractionCallback
	err = json.Unmarshal([]byte(c.Request.FormValue("payload")), &i)
	if err != nil {
		c.IndentedJSON(401, gin.H{"message": err})
	}

	// Note there might be a better way to get this info, but I figured this structure out from looking at the json response
	firstName := i.View.State.Values["First Name"]["firstName"].Value
	lastName := i.View.State.Values["Last Name"]["lastName"].Value

	msg := fmt.Sprintf("Hello %s %s, nice to meet you!", firstName, lastName)

	api := slack.New(h.Token)
	_, _, err = api.PostMessage(i.User.ID,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAttachments())
	if err != nil {
		c.IndentedJSON(401, gin.H{"message": err})
	}
}
