package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// BackupDatabase triggers a database backup based on the driver type
func BackupDatabase(driver, host, port, user, pass, name, outputFile string) error {
	switch driver {
	case "sqlite", "sqlite3":
		// For SQLite, we just copy the file
		return copyFile(name, outputFile)
	case "mysql", "mariadb":
		// mysqldump -u[user] -p[pass] -h[host] -P[port] [dbname] > [outputFile]
		cmd := exec.Command("mysqldump",
			fmt.Sprintf("-u%s", user),
			fmt.Sprintf("-p%s", pass),
			fmt.Sprintf("-h%s", host),
			fmt.Sprintf("-P%s", port),
			name,
		)
		outfile, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer outfile.Close()
		cmd.Stdout = outfile
		return runCommandWithStderr(cmd)
	case "postgres":
		// PGPASSWORD=pass pg_dump -U [user] -h [host] -p [port] [dbname] > [outputFile]
		cmd := exec.Command("pg_dump",
			"-U", user,
			"-h", host,
			"-p", port,
			name,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", pass))
		outfile, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer outfile.Close()
		cmd.Stdout = outfile
		return runCommandWithStderr(cmd)
	case "sqlserver":
		return errors.New("SQL Server backup enabled only via T-SQL to local server disk, not implemented for remote stream yet")
	case "oracle":
		// exp userid=user/pass@host:port/service file=outputFile
		// Connection string: user/pass@//host:port/service
		connStr := fmt.Sprintf("%s/%s@//%s:%s/%s", user, pass, host, port, name)
		cmd := exec.Command("exp",
			fmt.Sprintf("userid=%s", connStr),
			fmt.Sprintf("file=%s", outputFile),
		)
		return runCommandWithStderr(cmd)
	default:
		return fmt.Errorf("unsupported driver for backup: %s", driver)
	}
}

// RestoreDatabase triggers a database restore based on the driver type
func RestoreDatabase(driver, host, port, user, pass, name, inputFile string) error {
	switch driver {
	case "sqlite", "sqlite3":
		// For SQLite, we copy the input file to the db location
		return copyFile(inputFile, name)
	case "mysql", "mariadb":
		// mysql -u[user] -p[pass] -h[host] -P[port] [dbname] < [inputFile]
		cmd := exec.Command("mysql",
			fmt.Sprintf("-u%s", user),
			fmt.Sprintf("-p%s", pass),
			fmt.Sprintf("-h%s", host),
			fmt.Sprintf("-P%s", port),
			name,
		)
		infile, err := os.Open(inputFile)
		if err != nil {
			return err
		}
		defer infile.Close()
		cmd.Stdin = infile
		return runCommandWithStderr(cmd)
	case "postgres":
		// PGPASSWORD=pass psql -U [user] -h [host] -p [port] [dbname] < [inputFile]
		cmd := exec.Command("psql",
			"-U", user,
			"-h", host,
			"-p", port,
			name,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", pass))
		infile, err := os.Open(inputFile)
		if err != nil {
			return err
		}
		defer infile.Close()
		cmd.Stdin = infile
		return runCommandWithStderr(cmd)
	case "sqlserver":
		return errors.New("SQL Server restore enabled only via T-SQL from local server disk, not implemented for remote stream yet")
	case "oracle":
		// imp userid=user/pass@host:port/service file=inputFile full=y
		connStr := fmt.Sprintf("%s/%s@//%s:%s/%s", user, pass, host, port, name)
		cmd := exec.Command("imp",
			fmt.Sprintf("userid=%s", connStr),
			fmt.Sprintf("file=%s", inputFile),
			"full=y",
		)
		return runCommandWithStderr(cmd)
	default:
		return fmt.Errorf("unsupported driver for restore: %s", driver)
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func runCommandWithStderr(cmd *exec.Cmd) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr.String())
	}
	return nil
}
