package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/okta/okta-cli-client/iostream"
)

var stdErrWriter = survey.WithStdio(iostream.Input, iostream.Messages, iostream.Messages)

var Icons = survey.WithIcons(func(icons *survey.IconSet) {
	icons.Question.Text = ""
})

func Ask(inputs []*survey.Question, response interface{}) error {
	return survey.Ask(inputs, response, stdErrWriter, Icons)
}

func AskOne(input *survey.Question, response interface{}) error {
	return survey.Ask([]*survey.Question{input}, response, stdErrWriter, Icons)
}

func TextInput(name string, message string, help string, defaultValue string, required bool) *survey.Question {
	input := &survey.Question{
		Name:   name,
		Prompt: &survey.Input{Message: message, Help: help, Default: defaultValue},
	}

	if required {
		input.Validate = survey.Required
	}

	return input
}
