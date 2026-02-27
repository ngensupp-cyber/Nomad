package payload

import (
	"fmt"
	"os"
	"os/exec"
)

func BuildGoAgent(osType, arch, c2Addr, outputPath string) error {
	// Prepare build command
	// Note: We need to use the absolute or correct relative path to agents/main.go
	cmd := exec.Command("go", "build", "-ldflags", fmt.Sprintf("-X main.serverAddr=%s", c2Addr), "-o", outputPath, "agents/main.go")
	
	// Set environment variables for cross-compilation
	env := os.Environ()
	env = append(env, "GOOS="+osType)
	env = append(env, "GOARCH="+arch)
	if osType == "windows" {
		env = append(env, "CGO_ENABLED=1")
		env = append(env, "CC=x86_64-w64-mingw32-gcc")
	} else {
		env = append(env, "CGO_ENABLED=0")
	}
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %v\n%s", err, string(output))
	}

	return nil
}
