package dbhandler

import (
	"bufio"
	"fmt"
	dberrors "go-crud/internal/server/app/db/errors"
	"log/slog"
	"os"
	"regexp"
)

type DB struct{}

func OpenDB() *DB {
	return &DB{}
}

func (d *DB) FindAll() ([]string, error) {
	content, err := fileContent()
	var output = make([]string, 0)

	if err != nil {
		return []string{}, err
	}

	for _, v := range content {
		if v == "" {
			continue
		} else {
			output = append(output, v)
		}
	}

	return output, nil
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

	return "", &dberrors.DBNotFoundError{}
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

func (d *DB) Update(entryID, newEntry string) error {
	file, err := openFileWithPerm(os.O_RDWR)
	var found bool

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

		for _, v := range fileContent {
			if v == "" {
				continue
			} else {
				vID := GetEntryID(v)
				if vID == entryID {
					found = true
					newContent = append(newContent, newEntry)
				} else {
					newContent = append(newContent, v)
				}
			}
		}

		overrideFileContent(file, newContent)
	}

	if found {
		return nil
	} else {
		return &dberrors.DBNotFoundError{}
	}
}

func (d *DB) Delete(ID string) error {
	file, err := openFileWithPerm(os.O_RDWR)
	var found bool

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
			entryID := GetEntryID(v)
			if entryID == ID {
				found = true
				continue
			} else {
				newContent = append(newContent, v)
			}
		}

		overrideFileContent(file, newContent)
	}
	if found {
		return nil
	} else {
		return &dberrors.DBNotFoundError{}
	}
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

func openFileWithPerm(flags int) (*os.File, error) {
	file, err := os.OpenFile("internal/server/app/db/db.txt", flags|os.O_CREATE, 0644)

	if err != nil {
		slog.Error("ERROR", "DB Path", "Failed to open DB file.")
		return nil, &dberrors.DBPathError{
			Err: err,
		}
	}

	return file, err
}

func fileContent() ([]string, error) {
	file, err := openFileWithPerm(os.O_RDONLY)

	if err != nil {
		slog.Error("Error reading FileContent")
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

func GetEntryID(entry string) string {
	if entry == "" {
		return ""
	}

	re := regexp.MustCompile(`"([a-f0-9-]{36})"`)
	match := re.FindStringSubmatch(entry)

	if len(match) > 1 {
		return match[1]
	} else {
		return ""
	}
}
