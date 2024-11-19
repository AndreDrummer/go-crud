package dbhandler

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

type DB struct{}

func openFileWithPerm(flags int) (*os.File, error) {
	file, err := os.OpenFile("server/db/db.txt", flags|os.O_CREATE, 0644)

	if err != nil {
		slog.Error("error opening DB file", "error", err)
		return nil, err
	}

	return file, err
}

func fileContent() ([]string, error) {
	file, err := openFileWithPerm(os.O_RDONLY)

	if err != nil {
		slog.Error(fmt.Sprintf("Error reading FileContent: %v", err))
		return nil, err
	}

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

		return content, nil
	}

	return []string{}, nil
}

func overrideFileContent(file *os.File, newcontent []string) {
	file.Truncate(0)
	file.Seek(0, 0)
	for _, v := range newcontent {
		file.WriteString(fmt.Sprintf("%s\n", v))
	}
}

// {"id":"a952983d-9857-4096-a376-feed20a20168","first_name":"André","last_name":"Silva","biography":"Lindão"}
func GetEntryID(entry string) string {
	re := regexp.MustCompile(`"([a-f0-9-]{36})"`)
	match := re.FindStringSubmatch(entry)

	if len(match) > 1 {
		return match[1]
	} else {
		return ""
	}
}

func OpenDB() *DB {
	return &DB{}
}

func (d *DB) FindAll() (string, error) {
	content, err := fileContent()

	if err != nil {
		return "", err
	}

	var stringBuilder strings.Builder

	for _, v := range content {
		if v == "" {
			continue
		} else {
			stringBuilder.WriteString(fmt.Sprintf("%s\n", v))
		}
	}

	return stringBuilder.String(), nil
}

func (d *DB) FindByID(ID string) (string, error) {
	content, err := fileContent()

	if err != nil {
		return "", err
	}

	for _, v := range content {
		if v == "" {
			continue
		} else {
			vID := GetEntryID(v)
			if vID == ID {
				return v, nil
			}
		}
	}

	return "", fmt.Errorf("user of ID %s was not found", ID)
}

func (d *DB) Insert(entry string) error {
	file, err := openFileWithPerm(os.O_APPEND | os.O_WRONLY)

	if err != nil {
		slog.Error(fmt.Sprintf("Error inserting in DB: %v", err))
		return err
	}

	if file != nil {
		defer file.Close()
		file.WriteString(fmt.Sprintf("\n%s", entry))
	}

	return nil
}

func (d *DB) Update(entry string) error {
	file, err := openFileWithPerm(os.O_RDWR)

	if err != nil {
		slog.Error(fmt.Sprintf("Error updating from DB: %v", err))
		return err
	}

	if file != nil {
		defer file.Close()

		fileContent, err := fileContent()

		if err != nil {
			return err
		}

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

	return nil
}

func (d *DB) Delete(entry string) error {
	file, err := openFileWithPerm(os.O_RDWR)

	if err != nil {
		slog.Error(fmt.Sprintf("Error deleting from DB: %v", err))
		return err
	}

	if file != nil {
		defer file.Close()
		var newContent []string

		content, err := fileContent()

		if err != nil {
			return err
		}

		for _, v := range content {
			if v == "" || v == entry {
				continue
			} else {
				newContent = append(newContent, v)
			}
		}

		overrideFileContent(file, newContent)
	}
	return nil
}

func (d *DB) Clear() error {
	file, err := openFileWithPerm(os.O_TRUNC)

	if err != nil {
		slog.Error(fmt.Sprintf("Error clearing DB: %v", err))
		return err
	}

	if file != nil {
		defer file.Close()
		file.Truncate(0)
	}

	return nil
}
