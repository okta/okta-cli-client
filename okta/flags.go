package okta

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/okta/okta-cli-client/prompt"
	"github.com/spf13/cobra"
)

type commandInput interface {
	GetName() string
	GetLabel() string
	GetHelp() string
	GetIsRequired() bool
}

func isInputRequired(i commandInput) bool {
	return i.GetIsRequired()
}

func shouldAsk(cmd *cobra.Command, f *Flag) bool {
	return shouldPrompt(cmd, f)
}

func shouldPrompt(cmd *cobra.Command, flag *Flag) bool {
	return canPrompt(cmd) && !flag.IsSet(cmd)
}

func askFlag(cmd *cobra.Command, f *Flag, value interface{}, defaultValue *string) error {
	if shouldAsk(cmd, f) {
		return ask(f, value, defaultValue)
	}

	return nil
}

func ask(i commandInput, value interface{}, defaultValue *string) error {
	var val string
	if defaultValue != nil {
		val = *defaultValue
	}
	isRequired := isInputRequired(i)
	input := prompt.TextInput("", i.GetLabel(), i.GetHelp(), val, isRequired)

	if err := prompt.AskOne(input, value); err != nil {
		return handleInputError(err)
	}

	return nil
}

type Flag struct {
	Name         string
	LongForm     string
	ShortForm    string
	Help         string
	IsRequired   bool
	AlwaysPrompt bool
}

func (f Flag) GetName() string {
	return f.Name
}

func (f Flag) GetLabel() string {
	return inputLabel(f.Name)
}

func (f Flag) GetHelp() string {
	return f.Help
}

func (f Flag) GetIsRequired() bool {
	return f.IsRequired
}

func (f *Flag) IsSet(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(f.LongForm)
}

func inputLabel(name string) string {
	return fmt.Sprintf("%s:", name)
}

func (f *Flag) Ask(cmd *cobra.Command, value interface{}, defaultValue *string) error {
	return askFlag(cmd, f, value, defaultValue)
}

func (f *Flag) RegisterString(cmd *cobra.Command, value *string, defaultValue string) {
	registerString(cmd, f, value, defaultValue)
}

func (f *Flag) RegisterBool(cmd *cobra.Command, value *bool, defaultValue bool) {
	registerBool(cmd, f, value, defaultValue)
}

func registerString(cmd *cobra.Command, f *Flag, value *string, defaultValue string) {
	cmd.Flags().StringVarP(value, f.LongForm, f.ShortForm, defaultValue, f.Help)
	if err := markFlagRequired(cmd, f); err != nil {
		panic(fmt.Errorf("failed to register string flag %v", err))
	}
}

func registerBool(cmd *cobra.Command, f *Flag, value *bool, defaultValue bool) {
	cmd.Flags().BoolVarP(value, f.LongForm, f.ShortForm, defaultValue, f.Help)
	if err := markFlagRequired(cmd, f); err != nil {
		panic(fmt.Errorf("failed to register bool flag %v", err))
	}
}

func markFlagRequired(cmd *cobra.Command, f *Flag) error {
	if f.IsRequired {
		return cmd.MarkFlagRequired(f.LongForm)
	}

	return nil
}

func unexpectedError(err error) error {
	return fmt.Errorf("An unexpected error occurred: %w", err)
}

func handleInputError(err error) error {
	if err == terminal.InterruptErr {
		os.Exit(0)
	}

	return unexpectedError(err)
}
