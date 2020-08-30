package service

import (
	"github.com/flosch/pongo2"
	"os"
)

type tagPublicNode struct {
	Name string
}

func (self *tagPublicNode) Execute(ctx *pongo2.ExecutionContext, buffer pongo2.TemplateWriter) *pongo2.Error {
	if os.Getenv("APP_ENV") == "production" {
		buffer.WriteString("https://"+ ctx.Public["publicPath"].(string) + "/" + self.Name)
	} else {
		buffer.WriteString("/public/" + self.Name)
	}
	return nil
}

func tagPublicParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagPublicNode{}

	if filenameToken := arguments.MatchType(pongo2.TokenString); filenameToken != nil {
		node.Name = filenameToken.Val
	}

	return node, nil
}
