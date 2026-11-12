package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type Resume struct {
	Company          string    `json:"company"`
	Role             string    `json:"role"`
	OriginalFileName string    `json:"orginal_name"`
	FiledFileName    string    `json:"filed_name"`
	AddedAt          time.Time `json:"added_at"`
}

// home returns vault folder: $RESUMES_HOME if set, ~/resumes if empty
func home() (string, error) {

	// non empty RESUMES_HOME value sent to caller
	if dir := os.Getenv("RESUMES_HOME"); dir != "" {
		return dir, nil
	}

	// get users home directory
	usersHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error finding home directory %w\n", err)
	}

	// RESUMES_HOME not set so return default ~/resumes location
	return filepath.Join(usersHome, "resumes"), nil

}

// Add returns string location file was added to, empty string if error
func Add(src, company, role string) (*Resume, error) {

	// get users vault location
	root, err := home()
	if err != nil {
		return nil, err
	}

	// build the full directory path where file will be placed
	dir := filepath.Join(root, company, role)

	// create location file will be stored
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("error creating %s, %w", dir, err)
	}

	// append date to file name
	name := time.Now().Format("01-02-2006") + "-" + filepath.Base(src)

	// store file in location
	resume, err := copyNew(src, filepath.Join(dir, name), company, role) // add appended name to destination

	// return written destination back to caller
	return resume, nil
}

// copyNew returns destination copy was written to empty string if error
func copyNew(src, dst, company, role string) (*Resume, error) {

	fmt.Println("destination in copyNew: ", dst)
	// open source file to be copied
	in, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	// always close source file on function finish
	defer in.Close()

	// seperate the path and file extension
	// used to insert incremented value in case of existing files with matching name
	extension := filepath.Ext(dst)
	base := strings.TrimSuffix(dst, extension)

	// increment until incremented value no longer matches file name
	for i := 2; ; i++ {

		// open destination file and check if file already exists
		out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if errors.Is(err, fs.ErrExist) {

			// format name to include incremented value
			// continue to next for loop iteration to try to open file with new incremented name
			dst = fmt.Sprintf("%s-%d%s", base, i, extension)
			continue
		}

		// return error to caller if OpenFile fails
		if err != nil {
			return nil, err
		}

		// destination file is opened and file does not already exist
		// copy source in to destination out
		if _, err := io.Copy(out, in); err != nil {

			// close destination file
			// remove the destination file
			// return error to caller since copy failed
			out.Close()
			os.Remove(dst)
			return nil, fmt.Errorf("copying %s: %w", src, err)
		}

		resume, err := createJSON(src, dst, company, role)
		if err != nil {
			// TOOD: delete copied file if json file not created
			return nil, fmt.Errorf("creating json %w", err)
		}
		fmt.Println(resume)

		// copy successful return destination copy was written to
		// Close() returns an error send it back to caller
		return resume, out.Close()
	}
}

// TODO: deleteNew should return true when copyNew copy is deleted false if error
func deleteNew(dst string) (bool, error) {
	return true, nil
}

// createJson returns resume struct to caller error if false
func createJSON(original_name, filed_name, company, role string) (*Resume, error) {

	// create struct object
	resume := Resume{Company: company, Role: role, OriginalFileName: original_name, FiledFileName: filed_name, AddedAt: time.Now()}

	json_filed_name := filed_name + ".json"

	// OpenFile() opens new file for .json
	// WriteFile() can also be used to handle creating the file and write to file
	jsonFile, err := os.OpenFile(json_filed_name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open json file %w", err)
	}

	// create and use encoder to write json to file, SetIndent() formats / indents json in file
	enc := json.NewEncoder(jsonFile)
	enc.SetIndent("", "  ")
	enc.Encode(resume)

	defer jsonFile.Close()

	return &resume, nil
}

// Slug returns formatted string, empty if error
func Slug(s string) (string, error) {

	s = strings.ToLower(s)

	// for every rune only return characters or numbers
	trimmed := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// take slice of characters and numbers and join them
	joined := strings.Join(trimmed, "-")
	s = joined

	// return empty string if s is empty and error to be handled independently
	if s == "" {
		return "", errors.New("String is empty")
	}

	return s, nil
}
