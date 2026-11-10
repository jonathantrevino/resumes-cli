package vault

import (
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
func Add(src, company, role string) (string, error) {

	// get users vault location
	root, err := home()
	if err != nil {
		return "", err
	}

	// build the full directory path where file will be placed
	dir := filepath.Join(root, company, role)

	// create location file will be stored
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("error creating %s, %w", dir, err)
	}

	// append date to file name
	name := time.Now().Format("01-02-2006") + "-" + filepath.Base(src)

	// store file in location
	dest, err := copyNew(src, filepath.Join(dir, name)) // add appended name to destination

	// return written destination back to caller
	return dest, nil
}

// copyNew returns destination copy was written to empty string if error
func copyNew(src, dst string) (string, error) {

	// open source file to be copied
	in, err := os.Open(src)
	if err != nil {
		return "", err
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
			return "", err
		}

		// destination file is opened and file does not already exist
		// copy source in to destination out
		if _, err := io.Copy(out, in); err != nil {

			// close destination file
			// remove the destination file
			// return error to caller since copy failed
			out.Close()
			os.Remove(dst)
			return "", fmt.Errorf("copying %s: %w", src, err)
		}

		// copy successful return destination copy was written to
		// Close() returns an error send it back to caller
		return dst, out.Close()
	}
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
