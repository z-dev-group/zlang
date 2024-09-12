package object

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB
var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of argments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got=%s", args[0].Type())
			}
		}},
	},
	{
		"puts",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				if arg.Inspect() == "\\n" {
					fmt.Println()
				} else {
					fmt.Print(arg.Inspect())
				}
			}
			return nil
		}},
	},
	{
		"first",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}
			arr := args[0].(*Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return nil
		}},
	},
	{
		"last",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got=%s", args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return nil
		}},
	},
	{
		"rest",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got=%s", args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &Array{Elements: newElements}
			}
			return nil
		}},
	},
	{
		"push",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got=%s", args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			newElements := make([]Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &Array{Elements: newElements}
		}},
	},
	{
		"execute",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != STRING_OBJ {
				return newError("argument to `push` must be ARRAY, got=%s", args[0].Type())
			}
			str := args[0].(*String).Value

			commands := strings.Fields(str)
			cmd := exec.Command(commands[0], commands[1:]...)
			stdout, err := cmd.Output()

			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
			return &String{Value: string(stdout)}
		}},
	},
	{
		"mysql_init",
		&Builtin{Fn: func(args ...Object) Object {
			if args[0].Type() != STRING_OBJ {
				return newError("argument 1 to `mysql_init` must be String, got=%s", args[0].Type())
			}
			server := args[0].(*String).Value
			if args[1].Type() != STRING_OBJ {
				return newError("argument 2 to `mysql_init` must be String, got=%s", args[0].Type())
			}
			user := args[1].(*String).Value
			if args[2].Type() != STRING_OBJ {
				return newError("argument 3 to `mysql_init` must be String, got=%s", args[0].Type())
			}
			password := args[2].(*String).Value
			if args[3].Type() != STRING_OBJ {
				return newError("argument 4 to `mysql_init` must be String, got=%s", args[0].Type())
			}
			database := args[3].(*String).Value
			cfg := mysql.Config{
				User:   user,
				Passwd: password,
				Net:    "tcp",
				Addr:   server,
				DBName: database,
			}
			dsn := cfg.FormatDSN()
			dsn = strings.Replace(dsn, "allowNativePasswords=false", "allowNativePasswords=true", 1)
			db, _ = sql.Open("mysql", dsn)
			pingErr := db.Ping()
			if pingErr != nil {
				log.Fatal(pingErr)
			}
			return nil
		}},
	},
	{
		"mysql_query",
		&Builtin{Fn: func(args ...Object) Object {
			sql := args[0].(*String).Value
			rows, _ := db.Query(sql)
			result := Array{}
			columns, _ := rows.Columns()
			count := len(columns)
			values := make([]interface{}, count)
			valuePtrs := make([]interface{}, count)
			rowNames := Array{}
			for _, name := range columns {
				rowName := String{}
				rowName.Value = name
				rowNames.Elements = append(rowNames.Elements, Object(&rowName))
			}
			result.Elements = append(result.Elements, Object(&rowNames))
			for rows.Next() {
				for i := range columns {
					valuePtrs[i] = &values[i]
				}
				err := rows.Scan(valuePtrs...)
				if err != nil {
					panic(err)
				}
				row := Array{}
				for i := range columns {
					val := values[i]

					b, ok := val.([]byte)
					var v string
					if ok {
						v = string(b)
					}
					ui, ok := val.([]uint8)
					if ok {
						v = B2S(ui)
					}
					i64, ok := val.(int64)
					if ok {
						v = strconv.FormatInt(i64, 10)
					}
					value := String{}
					value.Value = v
					row.Elements = append(row.Elements, Object(&value))
				}
				result.Elements = append(result.Elements, Object(&row))
			}
			return &result
		}},
	},
	{
		"typeof",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			objectType := args[0].Type()
			return &String{Value: strings.ToLower(string(objectType))}
		}},
	},
	{
		"fetch",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError("wrong number of arguments. need more than one, got=%d", len(args))
			}
			urlObject := args[0]
			url, _ := urlObject.(*String)

			option := &Hash{}
			option.Pairs = make(map[HashKey]HashPair)
			if len(args) == 2 {
				option, _ = args[1].(*Hash)
			}
			client := &http.Client{
				Timeout: time.Second * 10, // 设置超时时间
			}
			methodObject := String{Value: "method"}
			methodHashPair, ok := option.Pairs[methodObject.HashKey()]
			method := "GET"
			if ok {
				method = methodHashPair.Value.Inspect()
			}
			bodyObject := String{Value: "body"}
			bodyHashParir, ok := option.Pairs[bodyObject.HashKey()]
			body := ""
			if ok {
				body = bodyHashParir.Value.Json()
			}
			method = strings.ToUpper(method)
			req, err := http.NewRequest(method, url.Value, bytes.NewReader([]byte(body)))
			if err != nil {
				return &Error{Message: err.Error()}
			}
			headersObject := String{Value: "headers"}
			headersPair, ok := option.Pairs[headersObject.HashKey()]
			headers := make(map[string]string)
			if ok {
				headersMap, ok := headersPair.Value.(*Hash)
				if ok {
					for _, value := range headersMap.Pairs {
						headers[value.Key.Inspect()] = value.Value.Inspect()
					}
				}
			}

			for key, value := range headers {
				req.Header.Set(key, value)
			}

			resp, err := client.Do(req)
			if err != nil {
				return &Error{Message: err.Error()}
			}
			defer resp.Body.Close()
			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return &Error{Message: err.Error()}
			}
			return &String{Value: string(responseBody)}
		}},
	},
}

func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
