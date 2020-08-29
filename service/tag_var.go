package service

import (
    "github.com/flosch/pongo2"
)

type tagVarNode struct {
    Name string
}

func (self *tagVarNode) Execute(ctx *pongo2.ExecutionContext, buffer pongo2.TemplateWriter) *pongo2.Error {
    buffer.WriteString("{{" + self.Name + "}}")
    return nil
}

func tagVarParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
    node := &tagVarNode{}
    if filenameToken := arguments.MatchType(pongo2.TokenString); filenameToken != nil {
        node.Name = filenameToken.Val
    }
    return node, nil
}