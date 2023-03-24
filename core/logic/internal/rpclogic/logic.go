package rpclogic

import (
	"bytes"
	"fmt"

	"github.com/licat233/genzero/config"
	"github.com/licat233/genzero/core/utils"
	"github.com/licat233/genzero/sql"
	"github.com/licat233/genzero/tools"
)

type Logic struct {
	CamelName      string
	LowerCamelName string
	SnakeName      string
	ModelName      string
	RpcGoPkgName   string
	PluralizedName string
	Dir            string
	Multiple       bool

	ConveFields string //注意：每个方法的数据不一样，会变
	HasUuid     bool
	HasName     bool

	Table *sql.Table
}

type LogicCollection []*Logic

func NewLogic(t *sql.Table) *Logic {
	return &Logic{
		CamelName:      tools.ToCamel(t.Name),
		LowerCamelName: tools.ToLowerCamel(t.Name),
		SnakeName:      tools.ToSnake(t.Name),
		ModelName:      tools.ToCamel(t.Name) + "Model",
		RpcGoPkgName:   config.C.Pb.GoPackage,
		PluralizedName: tools.PluralizedName(tools.ToCamel(t.Name)),
		Dir:            config.C.Logic.Rpc.Dir,
		Multiple:       config.C.Logic.Rpc.Multiple,
		ConveFields:    "",
		HasUuid:        utils.HasUuid(t.Fields),
		HasName:        utils.HasName(t.Fields),
		Table:          t,
	}
}

func (l *Logic) Run() (err error) {
	if err = l.Get(); err != nil {
		return err
	}
	if err = l.Add(); err != nil {
		return err
	}
	if err = l.Put(); err != nil {
		return err
	}
	if err = l.Del(); err != nil {
		return err
	}
	if err = l.List(); err != nil {
		return err
	}
	if err = l.Enums(); err != nil {
		return err
	}
	return
}

func (l *Logic) Get() (err error) {
	filename := l.getLogicFilename("get", "")
	ok, err := l.fileValidator(filename)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	logicContentTpl := `res, err := l.svcCtx.{{.ModelName}}.FindOne(l.ctx, in.Id)
	if err != nil && err != model.ErrNotFound {
		l.Logger.Error(err)
		return nil, errorx.IntRpcErr(err)
	}
	if err == model.ErrNotFound {
		return nil, nil
	}`

	logicContent, err := tools.ParserTpl(logicContentTpl, l)
	if err != nil {
		return err
	}

	returnDataTpl := `{{.CamelName}}: dataconv.Md{{.CamelName}}ToPb{{.CamelName}}(res),`
	returnData, err := tools.ParserTpl(returnDataTpl, l)
	if err != nil {
		return err
	}

	err = modifyLogicFileContent(filename, logicContent, returnData)
	if err != nil {
		return err
	}

	return
}

func (l *Logic) Add() (err error) {
	filename := l.getLogicFilename("add", "")
	ok, err := l.fileValidator(filename)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	var conveFieldsBuf bytes.Buffer
	for _, field := range l.AddFields() {
		if field.Type == "time.Time" {
			conveFieldsBuf.WriteString(fmt.Sprintf("%s: time.Unix(in.%s, 0).Local(),\n", field.UpperCamelCaseName, field.UpperCamelCaseName))
			continue
		}
		conveFieldsBuf.WriteString(fmt.Sprintf("%s: in.%s,\n", field.UpperCamelCaseName, field.UpperCamelCaseName))
	}
	l.ConveFields = conveFieldsBuf.String()

	logicContentTpl := `data := &model.{{.CamelName}}{
		{{.ConveFields}}
	}
	result, err := l.svcCtx.{{.ModelName}}.Insert(l.ctx, data)
	if err != nil {
		l.Logger.Error(err)
		return nil, errorx.IntRpcErr(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		l.Logger.Error(err)
		return nil, errorx.IntRpcErr(err)
	}

	res, err := l.svcCtx.{{.CamelName}}Model.FindOne(l.ctx, id)
	if err != nil && err != model.ErrNotFound {
		l.Logger.Error(err)
		return nil, errorx.IntRpcErr(err)
	}
	if err == model.ErrNotFound {
		l.Logger.Error(err)
		return nil, errorx.IntRpcErr(err)
	}`

	logicContent, err := tools.ParserTpl(logicContentTpl, l)
	if err != nil {
		return err
	}

	returnDataTpl := `{{.CamelName}}: dataconv.Md{{.CamelName}}ToPb{{.CamelName}}(res),`
	returnData, err := tools.ParserTpl(returnDataTpl, l)
	if err != nil {
		return err
	}

	err = modifyLogicFileContent(filename, logicContent, returnData)

	return err
}

func (l *Logic) Put() (err error) {
	filename := l.getLogicFilename("put", "")
	ok, err := l.fileValidator(filename)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	var conveFieldsBuf bytes.Buffer
	for _, field := range l.PutFields() {
		if field.Type == "time.Time" {
			conveFieldsBuf.WriteString(fmt.Sprintf("%s: in.%s.Unix(),\n", field.UpperCamelCaseName, field.UpperCamelCaseName))
			continue
		}
		conveFieldsBuf.WriteString(fmt.Sprintf("%s: in.%s,\n", field.UpperCamelCaseName, field.UpperCamelCaseName))
	}
	l.ConveFields = conveFieldsBuf.String()

	logicContentTpl := `data := &model.User{
		{{.ConveFields}}
	}

	if err := l.svcCtx.{{.ModelName}}.Update(l.ctx, data); err != nil {
		l.Logger.Error(err)
		return nil, errorx.IntRpcErr(err)
	}`

	logicContent, err := tools.ParserTpl(logicContentTpl, l)
	if err != nil {
		return err
	}

	err = modifyLogicFileContent(filename, logicContent, "")

	return err
}

func (l *Logic) Del() (err error) {
	filename := l.getLogicFilename("del", "")
	ok, err := l.fileValidator(filename)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	//分为软删除和硬删除
	var logicContentTpl string
	if l.Table.HasDeleteFiled {
		logicContentTpl = `err := l.svcCtx.{{.ModelName}}.DeleteSoft(l.ctx, in.id)
		if err != nil {
			l.Logger.Error(err)
			return nil, errorx.IntRpcErr(err)
		}`
	} else {
		logicContentTpl = `err = l.svcCtx.{{.ModelName}}.Delete(l.ctx, in.id)
		if err != nil {
			l.Logger.Error(err)
			return nil, errorx.IntRpcErr(err)
		}`
	}

	logicContent, err := tools.ParserTpl(logicContentTpl, l)
	if err != nil {
		return err
	}

	err = modifyLogicFileContent(filename, logicContent, "")

	return
}

func (l *Logic) List() (err error) {
	filename := l.getLogicFilename("get", "list")
	ok, err := l.fileValidator(filename)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	logicContentTpl := `
	var pageSize, page int64
	var keyword string
	if in.ListReq != nil {
		pageSize = in.ListReq.PageSize
		page = in.ListReq.Page
		keyword = in.ListReq.Keyword
	}
	in.ListReq = tools.NewListReq(in.ListReq)
	list, total, err := l.svcCtx.{{.ModelName}}.FindList(l.ctx, pageSize, page, keyword, Pb{{.CamelName}}ToMd{{.CamelName}}(in.{{.CamelName}}))
	if err != nil {
		l.Logger.Error(err)
		return nil, errorx.IntRpcErr(err)
	}`

	logicContent, err := tools.ParserTpl(logicContentTpl, l)
	if err != nil {
		return err
	}
	returnData := fmt.Sprintf("%s: dataconv.Md%sToPb%s(list),\nTotal:      total,", l.PluralizedName, l.PluralizedName, l.PluralizedName)

	err = modifyLogicFileContent(filename, logicContent, returnData)

	return
}

func (l *Logic) Enums() (err error) {
	//前提是要存在name字段
	if !l.HasName {
		return nil
	}

	filename := l.getLogicFilename("get", "enums")
	ok, err := l.fileValidator(filename)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	logicContentTpl := `list, _, err := l.svcCtx.{{.ModelName}}.FindList(l.ctx, -1, -1, "", nil)
	if err != nil && err != model.ErrNotFound {
		l.Logger.Error(err)
		return nil, errorx.IntRpcErr(err)
	}
	enums := []*{{.RpcGoPkgName}}.Enum{}
	for _, item := range list {
		enums = append(enums, &{{.RpcGoPkgName}}.Enum{
			Label: item.Name,
			Value: item.Id,
		})
	}`

	logicContent, err := tools.ParserTpl(logicContentTpl, l)
	if err != nil {
		return err
	}

	err = modifyLogicFileContent(filename, logicContent, "Enums: enums,")

	return
}

// 生成pb结构体转md结构体的方法
func (l *Logic) PbToMd() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\nfunc Pb%sToMd%s(pb *%s.%s) *model.%s {\n", l.CamelName, l.CamelName, l.RpcGoPkgName, l.CamelName, l.CamelName))
	buf.WriteString(`if pb == nil {
		return nil
	}
	`)
	buf.WriteString(fmt.Sprintf("return &model.%s{\n", l.CamelName))
	for _, field := range l.PbFields() {
		if field.Type == "time.Time" {
			buf.WriteString(fmt.Sprintf("%s: time.Unix(pb.%s, 0).Local(),\n", field.UpperCamelCaseName, field.UpperCamelCaseName))
			continue
		}
		buf.WriteString(fmt.Sprintf("%s: pb.%s,\n", field.UpperCamelCaseName, field.UpperCamelCaseName))
	}
	buf.WriteString("}\n")
	buf.WriteString("}\n")
	return buf.String()
}

// 生成md结构体转pb结构体的方法
func (l *Logic) MdToPb() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\nfunc Md%sToPb%s(md *model.%s) *%s.%s {\n", l.CamelName, l.CamelName, l.CamelName, l.RpcGoPkgName, l.CamelName))
	buf.WriteString(`if md == nil {
		return nil
	}
	`)
	buf.WriteString(fmt.Sprintf("return &%s.%s{\n", l.RpcGoPkgName, l.CamelName))
	for _, field := range l.PbFields() {
		if field.Type == "time.Time" {
			buf.WriteString(fmt.Sprintf("%s: md.%s.Unix(),\n", field.UpperCamelCaseName, field.UpperCamelCaseName))
			continue
		}
		buf.WriteString(fmt.Sprintf("%s: md.%s,\n", field.UpperCamelCaseName, field.UpperCamelCaseName))
	}
	buf.WriteString("}\n")
	buf.WriteString("}\n")
	return buf.String()
}

func (l *Logic) PbList2MdList() (string, error) {
	tpl := `
	func Pb{{.PluralizedName}}2Md{{.PluralizedName}}(pbList []*{{.RpcGoPkgName}}.{{.CamelName}}) []*model.{{.CamelName}} {
		if pbList == nil {
			return nil
		}
		data := make([]*model.{{.CamelName}}, 0)
		for _, v := range pbList {
			data = append(data, Pb{{.CamelName}}ToMd{{.CamelName}}(v))
		}
		return data
	}
	`

	return tools.ParserTpl(tpl, l)
}

func (l *Logic) MdList2PbList() (string, error) {
	tpl := `
	func Md{{.PluralizedName}}2Pb{{.PluralizedName}}(mdList []*model.{{.CamelName}}) []*{{.RpcGoPkgName}}.{{.CamelName}} {
		if mdList == nil {
			return nil
		}
		data := make([]*{{.RpcGoPkgName}}.{{.CamelName}}, 0)
		for _, v := range mdList {
			data = append(data, Md{{.CamelName}}ToPb{{.CamelName}}(v))
		}
		return data
	}
	`

	return tools.ParserTpl(tpl, l)
}
