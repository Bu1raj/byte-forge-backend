package executor

import (
	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

func RunSubmission(id, dir, lang string, timeout time.Duration) models.Result {
	start := time.Now()
	res := models.Result{ID: id}
	absDir, _ := filepath.Abs(dir)

	image, runCmd := DockerImageAndCmd(lang)
	if image == "" {
		res.ExitErr = "unsupported language"
		return res
	}

	args := []string{
		"run", "--rm",
		"--network", "none",
		"--pids-limit", "64",
		"--memory", "256m",
		"--cpus", "0.5",
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",
		"-v", fmt.Sprintf("%s:/work", absDir),
		"-w", "/work",
		image,
		"/bin/sh", "-c", runCmd,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.CombinedOutput()

	res.Duration = time.Since(start).String()

	if ctx.Err() == context.DeadlineExceeded {
		res.TimedOut = true
		res.ExitErr = "timeout"
		res.Stdout = string(out)
	} else if err != nil {
		res.ExitErr = err.Error()
		res.Stdout = string(out)
	} else {
		res.Stdout = string(out)
		res.Completed = true
	}
	return res
}
