package executor

// not using this, commented out for now

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"io/fs"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	"github.com/Bu1raj/byte-forge-backend/internal/models"
// )

// // RunSubmission executes the code submission in a Docker container with resource limits.
// // It returns a Result struct containing the execution details.
// func RunSubmission(id string, req models.SubmitReq) models.Result {
// 	start := time.Now()
// 	res := models.Result{ID: id}
// 	timeoutSec := req.TimeoutSecond

// 	image, runCmd := DockerImageAndCmd(req.Language)
// 	if image == "" {
// 		res.ExitErr = "unsupported language"
// 		return res
// 	}

// 	tmpDir, err := os.MkdirTemp("", "submission-*") // can use the id as prefix?
// 	if err != nil {
// 		return models.Result{ID: id, ExitErr: "temp dir error: " + err.Error()}
// 	}
// 	defer os.RemoveAll(tmpDir) // clean up

// 	filename := FilenameForLang(req.Language)
// 	if filename == "" {
// 		return models.Result{ID: id, ExitErr: "unsupported language"}
// 	}

// 	err = os.WriteFile(filepath.Join(tmpDir, filename), []byte(req.Code), fs.FileMode(0644))
// 	if err != nil {
// 		return models.Result{ID: id, ExitErr: "write error: " + err.Error()}
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
// 	defer cancel()

// 	args := []string{
// 		"run", "--rm",
// 		"--network", "none",
// 		"--pids-limit", "64",
// 		"--memory", "256m",
// 		"--cpus", "0.5",
// 		"--cap-drop", "ALL",
// 		"--security-opt", "no-new-privileges",
// 		"-v", fmt.Sprintf("%s:/work", tmpDir),
// 		"-w", "/work",
// 		image,
// 		"/bin/sh", "-c", runCmd,
// 	}

// 	cmd := exec.CommandContext(ctx, "docker", args...)

// 	cmd.Stdin = strings.NewReader(req.Code)

// 	var stderr, stdout bytes.Buffer
// 	cmd.Stderr = &stderr
// 	cmd.Stdout = &stdout
	
// 	err = cmd.Run()

// 	res.Duration = time.Since(start).String()
// 	res.Stdout = stdout.String()
// 	res.Stderr = stderr.String()

// 	if ctx.Err() == context.DeadlineExceeded {
// 		res.TimedOut = true
// 		res.ExitErr = "timeout"
// 		return res
// 	}
// 	if err != nil {
// 		if exitErr, ok := err.(*exec.ExitError); ok {
// 			res.ExitErr = fmt.Sprintf("exit code %d", exitErr.ExitCode())
// 		} else {
// 			res.ExitErr = err.Error()
// 		}
// 		return res
// 	}

// 	res.Completed = true
// 	return res
// }
