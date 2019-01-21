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

type MultiEnv struct {
	DataSource DataSource `profile:"database"`
	Eureka     Eureka
	Logging    map[string]interface{} `profile:"logging.level" profileDefault:"{\"github.com/flyleft/consul-iris\":\"debug\"}"`
	Users      []interface{}          `profile:"users" profileDefault:"[\"admin\",\"test\",\"root\"]"`
}

func TestProfileNoException(t *testing.T) {
	env, err := Profile(&MultiEnv{}, "test-multi-profile.yml", true)
	if err != nil {
		t.Error("Profile execute error", err)
	}
	fmt.Printf("Application active env: %+v\n", env)
}

func TestProfileActiveProfile(t *testing.T) {
	os.Setenv("PROFILES_ACTIVE", "production")
	env, err := Profile(&MultiEnv{}, "test-multi-profile.yml", true)
	if err != nil {
		t.Error("Profile execute error", err)
	}
	trueEnv := env.(*MultiEnv)
	fmt.Printf("Application active env: %+v\n", trueEnv)
	if trueEnv.DataSource.Username != "production" {
		t.Error("Active profile failed")
	}

}

func TestProfileEnv(t *testing.T) {
	eurekaZone := "http://192.168.1.10:8000/eureka/"
	os.Setenv("DEV_EUREKA_ZONE", eurekaZone)
	os.Setenv("PROFILES_ACTIVE", "dev")
	env, err := Profile(&MultiEnv{}, "test-multi-profile.yml", true)
	if err != nil {
		t.Error("Profile execute error", err)
	}
	trueEnv := env.(*MultiEnv)
	fmt.Printf("Application active env: %+v\n", trueEnv)
	if trueEnv.Eureka.Zone != eurekaZone {
		t.Error("Set value by env failed")
	}
}

type SingleEnv struct {
	Eureka  SingleEureka
	Logging map[string]interface{} `profile:"logging.level" profileDefault:"{\"github.com/flyleft/consul-iris\":\"debug\"}"`
	Skip    string                 `profile:"_"`
	Test    []string               `profile:"test"`
	Names   []string               `profile:"names" profileDefault:"[\"aaa\",\"bb\",\"你好哦\"]"`
	//SkipEureka *SingleEureka          `profile:"_"`
}

type SingleEureka struct {
	PreferIpAddress                  bool   `profile:"instance.preferIpAddress"`
	LeaseRenewalIntervalInSeconds    int32  `profile:"instance.leaseRenewalIntervalInSeconds"`
	LeaseExpirationDurationInSeconds uint   `profile:"instance.leaseExpirationDurationInSeconds"`
	ServerDefaultZone                string `profile:"client.serviceUrl.defaultZone" profileDefault:"http://localhost:8000/eureka/"`
	RegistryFetchIntervalSeconds     byte   `profile:"client.registryFetchIntervalSeconds"`
}

func TestProfileSingleEnv(t *testing.T) {
	os.Setenv("EUREKA_INSTANCE_LEASERENEWALINTERVALINSECONDS", "99")
	env, err := Profile(&SingleEnv{}, "test-single-profile.yml", true)
	if err != nil {
		t.Error("Profile execute error", err)
	}
	trueEnv := env.(*SingleEnv)
	fmt.Printf("Application active env: %+v\n", trueEnv)
	if trueEnv.Eureka.LeaseRenewalIntervalInSeconds != 99 || trueEnv.Skip != "" {
		t.Error("Set value by env failed")
	}
}
