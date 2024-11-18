package dbhandler

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

type DB struct{}

func openFileWithPerm(flags int) *os.File {
	file, err := os.OpenFile("server/db/db.txt", flags|os.O_CREATE, 0644)

	if err != nil {
		slog.Error("error opening DB file", "error", err)
		return nil
	}

	return file
}

func fileContent() []string {
	file := openFileWithPerm(os.O_RDONLY)
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	var content []string

	if file != nil {
		defer file.Close()
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				content = append(content, line)
			}
		}

		return content
	}

	return []string{}
}

func overrideFileContent(file *os.File, newcontent []string) {
	file.Truncate(0)
	file.Seek(0, 0)
	for _, v := range newcontent {
		file.WriteString(fmt.Sprintf("%s\n", v))
	}
}

func AnyToString(a any) string {
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)

	var builder strings.Builder

	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i)
		builder.WriteString(fmt.Sprintf("%s ", fieldValue))
	}

	return builder.String()
}

func OpenDB() *DB {
	return &DB{}
}

func (d *DB) Insert(entry string) {
	file := openFileWithPerm(os.O_APPEND | os.O_WRONLY)

	if file != nil {
		defer file.Close()
		file.WriteString(fmt.Sprintf("\n%s", entry))
	}
}

func (d *DB) Update(entry string) {
	file := openFileWithPerm(os.O_RDWR)

	if file != nil {
		defer file.Close()

		fileContent := fileContent()

		var newContent []string
		entryID := strings.Split(entry, " ")[0]

		for _, v := range fileContent {
			if v == "" {
				continue
			} else {
				vID := strings.Split(v, " ")[0]
				if vID != "" && vID == entryID {
					newContent = append(newContent, entry)
				} else {
					newContent = append(newContent, v)
				}
			}
		}

		overrideFileContent(file, newContent)
	}
}

func (d *DB) GetByID(ID string) string {
	content := fileContent()

	for _, v := range content {
		if v == "" {
			continue
		} else {
			vID := strings.Split(v, " ")[0]
			if vID == ID {
				return v
			}
		}
	}

	return ""
}

func (d *DB) GetAll() string {
	content := fileContent()
	var stringBuilder strings.Builder

	for _, v := range content {
		if v == "" {
			continue
		} else {
			stringBuilder.WriteString(fmt.Sprintf("%s\n", v))
		}
	}

	return stringBuilder.String()
}

func (d *DB) Delete(entry string) {
	file := openFileWithPerm(os.O_RDWR)

	if file != nil {
		defer file.Close()
		content := fileContent()
		var newContent []string

		for _, v := range content {
			if v == "" || v == entry {
				continue
			} else {
				newContent = append(newContent, v)
			}
		}

		overrideFileContent(file, newContent)
	}
}

func (d *DB) Clear() {
	file := openFileWithPerm(os.O_TRUNC)

	if file != nil {
		defer file.Close()
		file.Truncate(0)
	}
}
