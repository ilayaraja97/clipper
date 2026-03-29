package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	exec_color    = "#ffa657"
	config_color  = "#ffffff"
	chat_color    = "#66b3ff"
	help_color    = "#aaaaaa"
	error_color   = "#cc3333"
	warning_color = "#ffcc00"
	success_color = "#46b946"
)

type Renderer struct {
	contentRenderer *glamour.TermRenderer
	successRenderer lipgloss.Style
	warningRenderer lipgloss.Style
	errorRenderer   lipgloss.Style
	helpRenderer    lipgloss.Style
}

func NewRenderer(options ...glamour.TermRendererOption) *Renderer {
	contentRenderer, err := glamour.NewTermRenderer(options...)
	if err != nil {
		return nil
	}

	successRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(success_color))
	warningRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(warning_color))
	errorRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(error_color))
	helpRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(help_color)).Italic(true)

	return &Renderer{
		contentRenderer: contentRenderer,
		successRenderer: successRenderer,
		warningRenderer: warningRenderer,
		errorRenderer:   errorRenderer,
		helpRenderer:    helpRenderer,
	}
}

func (r *Renderer) RenderContent(in string) string {
	out, _ := r.contentRenderer.Render(in)

	return out
}

func (r *Renderer) RenderSuccess(in string) string {
	return r.successRenderer.Render(in)
}

func (r *Renderer) RenderWarning(in string) string {
	return r.warningRenderer.Render(in)
}

func (r *Renderer) RenderError(in string) string {
	return r.errorRenderer.Render(in)
}

func (r *Renderer) RenderHelp(in string) string {
	return r.helpRenderer.Render(in)
}

func (r *Renderer) RenderConfigMessage(currentLabel string, currentDescription string, progress []string) string {
	welcome := "Welcome!  \n\n"
	welcome += "I cannot find a configuration file, so let's set it up interactively.  \n"
	welcome += "Press `enter` to save the current field, and leave a field blank to accept its default.  \n"
	welcome += "Use `up` and `down` to move between fields.  \n"
	welcome += "You can use cloud AI providers like OpenRouter, OpenAI, Gemini, or a local OpenAI-compatible server.  \n\n"
	welcome += fmt.Sprintf("**Current field:** %s  \n", currentLabel)
	welcome += fmt.Sprintf("%s  \n\n", currentDescription)
	welcome += "**Setup progress**  \n"
	welcome += strings.Join(progress, "  \n")

	return welcome
}

func (r *Renderer) RenderHelpMessage() string {
	help := "**Help**\n"
	help += "- `↑`/`↓` : navigate in history\n"
	help += "- `tab`   : switch between `🚀 exec` and `💬 chat` prompt modes\n"
	help += "- `ctrl+h`: show help\n"
	help += "- `ctrl+s`: edit settings\n"
	help += "- `ctrl+r`: clear terminal and reset discussion history\n"
	help += "- `ctrl+l`: clear terminal but keep discussion history\n"
	help += "- `ctrl+c`: exit or interrupt command execution\n"

	return help
}
