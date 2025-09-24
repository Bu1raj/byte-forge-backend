package executor

func FilenameForLang(lang string) string {
	switch lang {
	case "python":
		return "main.py"
	case "c":
		return "main.c"
	case "go":
		return "main.go"
	case "javascript", "node":
		return "main.js"
	default:
		return ""
	}
}

func DockerImageAndCmd(lang string) (image, cmd string) {
	switch lang {
	case "python":
		return "python:3.11-alpine", "python main.py"
	case "c":
		return "gcc:latest", "gcc -O2 -std=gnu11 -o main main.c && ./main"
	case "go":
		return "golang:1.20", "go run main.go"
	case "javascript", "node":
		return "node:20-alpine", "node main.js"
	default:
		return "", ""
	}
}
