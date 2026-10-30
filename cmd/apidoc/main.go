// Command apidoc 解析 Go-Zero 的 cms.api 文件，生成 OpenAPI 3.0.3 文档（openapi.json）。
// 用法：go run ./cmd/apidoc  （在项目根目录执行，输出 openapi.json）
//
// 由于 goctl 1.10.2 不提供 openapi 子命令，这里用零依赖的方式从 .api 源文件直接生成，
// 保证文档与接口定义同源、可复现（make gen-doc）。
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type field struct {
	name     string
	json     string
	goType   string
	optional bool
}

type typeDef struct {
	name   string
	fields []field
}

type route struct {
	handler string
	method  string
	path    string
	req     string
	resp    string
}

func main() {
	src := loadAPI("api/cms/v1/cms.api")
	spec := generate(src)

	b, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "marshal:", err)
		os.Exit(1)
	}
	if err := os.WriteFile("openapi.json", b, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "write:", err)
		os.Exit(1)
	}
	fmt.Println("openapi.json generated")
}

// loadAPI 读取主 api 文件及其 import 的其它 .api 文件（go-zero 支持 import 拆分定义）。
func loadAPI(mainPath string) string {
	read := func(p string) string {
		b, err := os.ReadFile(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read", p, err)
			return ""
		}
		return string(b)
	}
	base := read(mainPath)
	dir := filepath.Dir(mainPath)
	reImp := regexp.MustCompile(`import\s*\(([\s\S]*?)\)`)
	for _, m := range reImp.FindAllStringSubmatch(base, -1) {
		for _, line := range strings.Split(m[1], "\n") {
			line = strings.TrimSpace(line)
			line = strings.Trim(line, `"`)
			if line == "" || !strings.HasSuffix(line, ".api") {
				continue
			}
			base += "\n" + read(filepath.Join(dir, line))
		}
	}
	return base
}

func stripComments(s string) string {
	reBlock := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	s = reBlock.ReplaceAllString(s, "")
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if idx := strings.Index(l, "//"); idx >= 0 {
			lines[i] = l[:idx]
		}
	}
	return strings.Join(lines, "\n")
}

// findBlock 找到 token 之后的开闭括号（openB/closeB）之间的内容，可匹配多层。
func findBlock(src, token string, openB, closeB byte) []string {
	var out []string
	rest := src
	for {
		idx := strings.Index(rest, token)
		if idx < 0 {
			break
		}
		start := idx + len(token)
		k := start
		for k < len(rest) && rest[k] != openB {
			k++
		}
		if k >= len(rest) {
			break
		}
		depth := 0
		j := k
		for j < len(rest) {
			switch rest[j] {
			case openB:
				depth++
			case closeB:
				depth--
				if depth == 0 {
					out = append(out, rest[k+1:j])
					rest = rest[j+1:]
					goto done
				}
			}
			j++
		}
		break
	done:
	}
	return out
}

func splitStructs(block string) []struct {
	name string
	body string
} {
	var res []struct {
		name string
		body string
	}
	rest := block
	for {
		idx := strings.Index(rest, "{")
		if idx < 0 {
			break
		}
		before := strings.TrimSpace(rest[:idx])
		name := before[strings.LastIndex(before, " ")+1:]
		depth := 0
		j := idx
		for j < len(rest) {
			switch rest[j] {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					res = append(res, struct {
						name string
						body string
					}{name, rest[idx+1 : j]})
					rest = rest[j+1:]
					goto nxt
				}
			}
			j++
		}
		break
	nxt:
	}
	return res
}

func parseTypes(src string) map[string]typeDef {
	types := map[string]typeDef{}
	fieldRe := regexp.MustCompile(`(\w+)\s+([\w.*\[\]]+)\s+` + "`" + `([^` + "`" + `]*)` + "`")
	for _, block := range findBlock(src, "type (", '(', ')') {
		for _, st := range splitStructs(block) {
			td := typeDef{name: st.name}
			for _, line := range strings.Split(st.body, "\n") {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				m := fieldRe.FindStringSubmatch(line)
				if m == nil {
					continue
				}
				f := field{name: m[1], goType: m[2]}
				tag := m[3]
				if jm := regexp.MustCompile(`json:"([^"]*)"`).FindStringSubmatch(tag); jm != nil {
					parts := strings.Split(jm[1], ",")
					f.json = parts[0]
					for _, p := range parts[1:] {
						if p == "optional" {
							f.optional = true
						}
					}
				}
				if f.json == "" {
					f.json = f.name
				}
				if strings.Contains(tag, "optional") {
					f.optional = true
				}
				td.fields = append(td.fields, f)
			}
			types[td.name] = td
		}
	}
	return types
}

func parseRoutes(src string) []route {
	var routes []route
	re := regexp.MustCompile(`^(get|post|put|delete|patch|head|options)\s+(\S+)\s*(?:\((\w+)\))?(?:\s+returns\s*\((\w+)\))?`)
	for _, block := range findBlock(src, "service", '{', '}') {
		lines := strings.Split(block, "\n")
		var cur route
		inRoute := false
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "@handler") {
				cur = route{handler: strings.TrimSpace(strings.TrimPrefix(line, "@handler"))}
				inRoute = true
				continue
			}
			if inRoute {
				if m := re.FindStringSubmatch(line); m != nil {
					cur.method = strings.ToUpper(m[1])
					cur.path = m[2]
					cur.req = m[3]
					cur.resp = m[4]
					routes = append(routes, cur)
					inRoute = false
				}
			}
		}
	}
	return routes
}

func toSchema(t string) map[string]interface{} {
	t = strings.TrimSpace(t)
	nullable := false
	if strings.HasPrefix(t, "*") {
		nullable = true
		t = t[1:]
	}
	if strings.HasPrefix(t, "[]") {
		item := toSchema(t[2:])
		return map[string]interface{}{"type": "array", "items": item}
	}
	if strings.HasPrefix(t, "map[string]") {
		val := toSchema(t[len("map[string]"):])
		return map[string]interface{}{"type": "object", "additionalProperties": val}
	}
	var schema map[string]interface{}
	switch t {
	case "string":
		schema = map[string]interface{}{"type": "string"}
	case "int64", "int", "int32", "int16", "uint", "uint64", "uint32":
		schema = map[string]interface{}{"type": "integer"}
	case "float64", "float32", "float":
		schema = map[string]interface{}{"type": "number"}
	case "bool":
		schema = map[string]interface{}{"type": "boolean"}
	case "time.Time":
		schema = map[string]interface{}{"type": "string", "format": "date-time"}
	default:
		schema = map[string]interface{}{"type": "string"}
	}
	if nullable {
		schema["nullable"] = true
	}
	return schema
}

func buildSchema(td typeDef) map[string]interface{} {
	props := map[string]interface{}{}
	required := []string{}
	for _, f := range td.fields {
		props[f.json] = toSchema(f.goType)
		if !f.optional {
			required = append(required, f.json)
		}
	}
	schema := map[string]interface{}{"type": "object", "properties": props}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func generate(src string) map[string]interface{} {
	src = stripComments(src)
	types := parseTypes(src)
	routes := parseRoutes(src)

	paths := map[string]interface{}{}
	for _, r := range routes {
		p := strings.ReplaceAll(r.path, ":id", "{id}")
		p = strings.ReplaceAll(p, ":Id", "{id}")
		method := strings.ToLower(r.method)
		op := map[string]interface{}{
			"summary":     r.handler,
			"operationId": r.handler,
		}
		if r.req != "" {
			op["requestBody"] = map[string]interface{}{
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{"$ref": "#/components/schemas/" + r.req},
					},
				},
			}
		}
		respContent := map[string]interface{}{}
		if r.resp != "" {
			respContent["schema"] = map[string]interface{}{"$ref": "#/components/schemas/" + r.resp}
		} else {
			respContent["schema"] = map[string]interface{}{"type": "object"}
		}
		op["responses"] = map[string]interface{}{
			"200": map[string]interface{}{
				"description": "成功",
				"content": map[string]interface{}{
					"application/json": respContent,
				},
			},
		}
		pm, ok := paths[p].(map[string]interface{})
		if !ok {
			pm = map[string]interface{}{}
			paths[p] = pm
		}
		pm[method] = op
	}

	schemas := map[string]interface{}{}
	for name, td := range types {
		schemas[name] = buildSchema(td)
	}

	return map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":   "Go Meteor CMS API",
			"version": "v1",
			"description": "由 cms.api 自动生成的 OpenAPI 文档（make gen-doc）。" +
				"所有接口统一响应结构为 { code:int, msg:string, data:any }，其中 data 即下方各接口标注的响应模型。",
		},
		"servers": []interface{}{map[string]interface{}{"url": "/"}},
		"paths":   paths,
		"components": map[string]interface{}{
			"schemas": schemas,
		},
	}
}
