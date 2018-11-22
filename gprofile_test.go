package gprofile

import (
	"fmt"
	"os"
	"testing"
)

type Eureka struct {
	Zone          string `profile:"zone"`
	FetchInterval int    `profile:"fetchInterval"`
}

type DataSource struct {
	Host     string `profile:"host" profileDefault:"localhost"`
	Username string `profile:"username"`
	Password string `profile:"password"`
}

type Env struct {
	DataSource DataSource `profile:"database"`
	Eureka     Eureka
	Logging    map[string]interface{} `profile:"logging.level" profileDefault:"{\"github.com/flyleft/consul-iris\":\"debug\"}"`
	Users      []interface{}          `profile:"users" profileDefault:"[\"admin\",\"test\",\"root\"]"`
}

func TestProfileNoException(t *testing.T) {
	env, err := Profile(&Env{}, "test-multi-profile.yml", true)
	if err != nil {
		t.Error("Profile execute error", err)
	}
	fmt.Printf("Application active env: %+v\n", env)
}

func TestProfileActiveProfile(t *testing.T) {
	os.Setenv("PROFILES_ACTIVE", "production")
	env, err := Profile(&Env{}, "test-multi-profile.yml", true)
	if err != nil {
		t.Error("Profile execute error", err)
	}
	trueEnv := env.(*Env)
	fmt.Printf("Application active env: %+v\n", trueEnv)
}

func TestProfileEnv(t *testing.T) {
	os.Setenv("DEV_BASE_TESTSTRING", "TestProfileEnv")
	os.Setenv("DEV_ENV", "AAA")
	env, err := Profile(&Env{}, "test-multi-profile.yml", true)
	if err != nil {
		t.Error("Profile execute error", err)
	}
	trueEnv := env.(*Env)
	fmt.Printf("Application active env: %+v\n", trueEnv)

}
