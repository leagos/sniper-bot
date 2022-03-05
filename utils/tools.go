package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetCpuId() string {
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorID")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	cpuid := string(out)
	cpuid = cpuid[12 : len(cpuid)-2]
	cpuid = strings.ReplaceAll(cpuid, "\n", "")
	return cpuid
}
