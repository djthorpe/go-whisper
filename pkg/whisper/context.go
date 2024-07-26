package whisper

import (
	"path/filepath"

	// Packages
	"github.com/mutablelogic/go-whisper/sys/whisper"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Context struct {
	*whisper.Context

	// Parameters for the model
	params *whisper.Params

	// Temporary parameters
	temp *whisper.Params
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewContextWithModel(base string, model *Model) *Context {
	path := filepath.Join(base, model.Path)
	ctx := whisper.Whisper_init(path)
	if ctx == nil {
		return nil
	}
	params := whisper.NewParams(whisper.SAMPLING_GREEDY)
	if params == nil {
		ctx.Whisper_free()
		return nil
	}

	return &Context{
		Context: ctx,
		params:  params,
		temp:    nil,
	}
}

func (c *Context) Close() error {
	var result error

	if c.Context != nil {
		c.Context.Whisper_free()
	}
	if c.params != nil {
		c.params.Close()
	}

	// Return any error
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (c *Context) Params() *whisper.Params {
	if c.temp != nil {
		return c.temp
	}
	return c.params
}

func (c *Context) SetLanguage(v string) {
	if c.temp == nil {
		t := *c.params
		c.temp = &t
	}
	if v == "" || v == "auto" {
		c.temp.SetLanguage(-1)
		return
	}
	lang := c.Whisper_lang_id(v)
	if lang != -1 {
		c.temp.SetLanguage(lang)
	}
}

func (c *Context) SetTranslate(v bool) {
	if c.temp == nil {
		t := *c.params
		c.temp = &t
	}
	c.temp.SetTranslate(v)
}