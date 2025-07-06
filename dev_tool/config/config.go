package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port         int    `yaml:"port"`
		FrontendPort int    `yaml:"frontend_port"`
		Mode         string `yaml:"mode"`
	} `yaml:"server"`

	Cors struct {
		AllowOrigins     []string `yaml:"allow_origins"`
		AllowMethods     []string `yaml:"allow_methods"`
		AllowHeaders     []string `yaml:"allow_headers"`
		AllowCredentials bool     `yaml:"allow_credentials"`
		MaxAge           int      `yaml:"max_age"`
	} `yaml:"cors"`

	Database struct {
		SQLite struct {
			Path string `yaml:"path"`
		} `yaml:"sqlite"`
	} `yaml:"database"`
}

var GlobalConfig Config

// executableDir stores the directory containing the executable
var executableDir string

func Init() {
	var configPath string
	var configBaseDir string // Directory from which config was successfully loaded
	var data []byte
	var err error

	// --- Try loading config relative to executable (for compiled binary) ---
	exePath, errExe := os.Executable()
	if errExe == nil {
		executableDir := filepath.Dir(exePath)
		pathTry := filepath.Join(executableDir, "config", "config.yaml")
		log.Printf("[Attempt 1] Reading config relative to executable: %s", pathTry)
		data, err = ioutil.ReadFile(pathTry)
		if err == nil {
			configPath = pathTry
			configBaseDir = executableDir
			log.Printf("Successfully read config relative to executable.")
		} else if !os.IsNotExist(err) {
			// Other error reading file relative to executable
			log.Fatalf("读取配置文件失败 (relative to exe %s): %v", pathTry, err)
		}
	} else {
		log.Printf("Warning: Could not get executable path: %v. Skipping check relative to executable.", errExe)
	}

	// --- If not found or executable path failed, try relative to CWD (for go run) ---
	if configPath == "" { // Only try CWD if not already loaded successfully
		pathTry := filepath.Join(".", "config", "config.yaml") // "." is Current Working Directory
		log.Printf("[Attempt 2] Reading config relative to current working directory: %s", pathTry)
		data, err = ioutil.ReadFile(pathTry)
		if err == nil {
			configPath = pathTry
			configBaseDir, _ = os.Getwd() // Use CWD as base
			log.Printf("Successfully read config relative to CWD.")
		} else {
			// Failed both ways
			log.Fatalf("读取配置文件失败: Tried relative to executable (if possible) and relative to CWD (%s): %v", pathTry, err)
		}
	}

	// --- Parse the config data ---
	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		log.Fatalf("解析配置文件失败 (%s): %v", configPath, err)
	}

	// --- Adjust database path based on where config was found ---
	if !filepath.IsAbs(GlobalConfig.Database.SQLite.Path) {
		// configBaseDir is either executableDir or CWD
		GlobalConfig.Database.SQLite.Path = filepath.Join(configBaseDir, GlobalConfig.Database.SQLite.Path)
		log.Printf("Updated database path relative to config location: %s", GlobalConfig.Database.SQLite.Path)
	} else {
		log.Printf("Database path is absolute: %s", GlobalConfig.Database.SQLite.Path)
	}

	// --- Ensure database directory exists ---
	dbDir := filepath.Dir(GlobalConfig.Database.SQLite.Path)
	log.Printf("Ensuring database directory exists: %s", dbDir)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("创建数据库目录失败: %v", err)
	}

	log.Println("Configuration initialized successfully.")
}
