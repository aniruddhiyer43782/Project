package util

import "os"


const dockerEnvFile string = "/.dockerenv"


func IsRunInDocker() bool {
	_, err := os.Stat(dockerEnvFile)
	return err == nil
}
