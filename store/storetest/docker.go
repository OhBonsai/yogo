package storetest

import (
	"fmt"
	"os/exec"

	"github.com/OhBonsai/yogo/mlog"
	"github.com/OhBonsai/yogo/model"
	"strings"
	"encoding/json"
	"time"
	"net"
	"io"
)

type Container struct {
	Id string
	NetworkSettings struct{
		Ports map[string][]struct{
			HostPost string
		}
	}
}


type RunningContainer struct {
	Container
}

func (c *RunningContainer) Stop() error {
	mlog.Info(fmt.Sprintf("Removing container: %v", c.Id))
	return exec.Command("docker", "rm", "-f", c.Id).Run()
}


func NewPostgreSQLContainer() (*RunningContainer, *model.SqlSettings, error) {
	container, err := runContainer([]string{
		"-e", "POSTGRES_USER=bonsai",
		"-e", "POSTGRES_PASSWORD=123456",
		"--tmpfs", "/var/lib/postgresql/data",
		"postgres:9.4",
	})

	if err != nil {
		return nil, nil, err
	}

	mlog.Info("Waiting for postgres connectivity")
	port := container.NetworkSettings.Ports["5432/tcp"][0].HostPost
	if err := waitForPort(port); err != nil {
		container.Stop()
		return nil, nil, err
	}
	return container, databaseSettings("postgres", "postgres://bonsai:123456@127.0.0.1:"+port+"?sslmode=disable"), nil
}

func databaseSettings(driver, dataSource string) *model.SqlSettings {
	settings := &model.SqlSettings{
		DriverName:               &driver,
		DataSource:               &dataSource,
		DataSourceReplicas:       []string{},
		DataSourceSearchReplicas: []string{},
		MaxIdleConns:             new(int),
		MaxOpenConns:             new(int),
		Trace:                    false,
		AtRestEncryptKey:         model.NewRandomString(32),
		QueryTimeout:             new(int),
	}
	*settings.MaxIdleConns = 10
	*settings.MaxOpenConns = 100
	*settings.QueryTimeout = 10
	return settings
}



func runContainer(args []string) (*RunningContainer, error) {
	name := "yogo-storetest-" + model.NewId()
	dockerArgs := append([]string{"run", "-d", "-P", "--name", name}, args...)
	out, err := exec.Command("docker", dockerArgs...).Output()

	if err != nil {
		return nil, err
	}

	id := strings.TrimSpace(string(out))
	out, err = exec.Command("docker", "inspect", id).Output()
	if err != nil {
		exec.Command("docker", "rm", "-f", id).Run()
		return nil, err
	}

	var containers []Container
	if err := json.Unmarshal(out, &containers); err != nil {
		exec.Command("docker", "rm", "-f", id).Run()
		return nil, err
	}

	mlog.Info(fmt.Sprintf("Running container: %v", id))
	return &RunningContainer{containers[0]}, nil
}

func waitForPort(port string) error {
	deadline := time.Now().Add(time.Minute * 10)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, time.Minute)
		if err != nil {
			return err
		}
		if err = conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
			return err
		}
		_, err = conn.Read(make([]byte, 1))
		conn.Close()
		if err == nil {
			return nil
		}
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return nil
		}
		if err != io.EOF {
			return err
		}
		time.Sleep(time.Millisecond * 200)
	}
	return fmt.Errorf("timeout waiting for port %v", port)
}