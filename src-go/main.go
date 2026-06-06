package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type ConfigInfo struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	TargetPath string `json:"targetPath"` // Active path on C:
	RepoPath   string `json:"repoPath"`   // Path in E:\Dev\projects\winrice\config\
	Exists     bool   `json:"exists"`     // Exists at active path
	InRepo     bool   `json:"inRepo"`     // Exists in repo
	IsSymlink  bool   `json:"isSymlink"`  // Active path is a symlink
	SymlinkTo  string `json:"symlinkTo"`  // Where active symlink points
	IsDir      bool   `json:"isDir"`
}

var configs = []ConfigInfo{
	{
		Key:        "wezterm",
		Name:       "WezTerm",
		TargetPath: "C:\\Users\\Timofey\\.wezterm.lua",
		RepoPath:   "config\\wezterm\\.wezterm.lua",
		IsDir:      false,
	},
	{
		Key:        "glazewm",
		Name:       "GlazeWM",
		TargetPath: "C:\\Users\\Timofey\\.glzr\\glazewm\\config.yaml",
		RepoPath:   "config\\glazewm\\config.yaml",
		IsDir:      false,
	},
	{
		Key:        "zebar",
		Name:       "Zebar",
		TargetPath: "C:\\Users\\Timofey\\.glzr\\zebar",
		RepoPath:   "config\\zebar",
		IsDir:      true,
	},
	{
		Key:        "fastfetch",
		Name:       "Fastfetch",
		TargetPath: "C:\\Users\\Timofey\\AppData\\Roaming\\fastfetch\\config.jsonc",
		RepoPath:   "config\\fastfetch\\config.jsonc",
		IsDir:      false,
	},
	{
		Key:        "sync_script",
		Name:       "Theme Syncer",
		TargetPath: "C:\\Users\\Timofey\\Theme\\sync_theme.ps1",
		RepoPath:   "config\\Theme\\sync_theme.ps1",
		IsDir:      false,
	},
	{
		Key:        "zed",
		Name:       "Zed Editor",
		TargetPath: "C:\\Users\\Timofey\\AppData\\Roaming\\Zed\\settings.json",
		RepoPath:   "config\\Zed\\settings.json",
		IsDir:      false,
	},
	{
		Key:        "vscode",
		Name:       "VS Code",
		TargetPath: "C:\\Users\\Timofey\\AppData\\Roaming\\Code\\User\\settings.json",
		RepoPath:   "config\\VSCode\\settings.json",
		IsDir:      false,
	},
	{
		Key:        "komorebi",
		Name:       "Komorebi",
		TargetPath: "C:\\Users\\Timofey\\komorebi.json",
		RepoPath:   "config\\komorebi\\komorebi.json",
		IsDir:      false,
	},
	{
		Key:        "komorebi_bar",
		Name:       "Komorebi Bar",
		TargetPath: "C:\\Users\\Timofey\\komorebi.bar.json",
		RepoPath:   "config\\komorebi\\komorebi.bar.json",
		IsDir:      false,
	},
	{
		Key:        "whkd",
		Name:       "whkd Shortcuts",
		TargetPath: "C:\\Users\\Timofey\\.config\\whkdrc",
		RepoPath:   "config\\whkd\\whkdrc",
		IsDir:      false,
	},
	{
		Key:        "gitconfig",
		Name:       "Git Config",
		TargetPath: "C:\\Users\\Timofey\\.gitconfig",
		RepoPath:   "config\\git\\.gitconfig",
		IsDir:      false,
	},
}

var projectRoot = ""

func getInfo(c ConfigInfo) ConfigInfo {
	// Resolve full repo path
	c.RepoPath = filepath.Join(projectRoot, c.RepoPath)

	// Check repo existence
	if stat, err := os.Lstat(c.RepoPath); err == nil {
		c.InRepo = true
		c.IsDir = stat.IsDir()
	}

	// Check target existence and symlink status
	targetStat, err := os.Lstat(c.TargetPath)
	if err == nil {
		c.Exists = true
		if targetStat.Mode()&os.ModeSymlink != 0 {
			c.IsSymlink = true
			if linkDest, err := os.Readlink(c.TargetPath); err == nil {
				c.SymlinkTo = linkDest
			}
		}
	}
	return c
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var result []ConfigInfo
	for _, c := range configs {
		result = append(result, getInfo(c))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleReadFile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	var found *ConfigInfo
	for i := range configs {
		if configs[i].Key == key {
			found = &configs[i]
			break
		}
	}

	if found == nil {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	info := getInfo(*found)
	if !info.InRepo {
		http.Error(w, "Config not in repository yet", http.StatusNotFound)
		return
	}

	if info.IsDir {
		http.Error(w, "Target is a directory, cannot read directly", http.StatusBadRequest)
		return
	}

	data, err := os.ReadFile(info.RepoPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}

type WritePayload struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

func handleWriteFile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload WritePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var found *ConfigInfo
	for i := range configs {
		if configs[i].Key == payload.Key {
			found = &configs[i]
			break
		}
	}

	if found == nil {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	info := getInfo(*found)
	if info.IsDir {
		http.Error(w, "Target is a directory, cannot write directly", http.StatusBadRequest)
		return
	}

	// Ensure parent directory in repo exists
	repoDir := filepath.Dir(info.RepoPath)
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create repo dir: %v", err), http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(info.RepoPath, []byte(payload.Content), 0644); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write file: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Success")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleAdopt(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	var found *ConfigInfo
	for i := range configs {
		if configs[i].Key == key {
			found = &configs[i]
			break
		}
	}

	if found == nil {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	info := getInfo(*found)
	if !info.Exists {
		http.Error(w, "Source file does not exist to adopt", http.StatusBadRequest)
		return
	}

	if info.InRepo {
		http.Error(w, "Already in repository", http.StatusBadRequest)
		return
	}

	// Create repo dir
	repoDir := filepath.Dir(info.RepoPath)
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create repo dir: %v", err), http.StatusInternalServerError)
		return
	}

	// Copy to repo
	if info.IsDir {
		if err := copyDir(info.TargetPath, info.RepoPath); err != nil {
			http.Error(w, fmt.Sprintf("Failed to copy dir: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		if err := copyFile(info.TargetPath, info.RepoPath); err != nil {
			http.Error(w, fmt.Sprintf("Failed to copy file: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Adopted successfully")
}

func handleLink(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	var found *ConfigInfo
	for i := range configs {
		if configs[i].Key == key {
			found = &configs[i]
			break
		}
	}

	if found == nil {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	info := getInfo(*found)
	if !info.InRepo {
		http.Error(w, "Config must be in repo before linking", http.StatusBadRequest)
		return
	}

	// Backup existing active target if it is not already a symlink
	if info.Exists && !info.IsSymlink {
		backupPath := info.TargetPath + ".backup"
		if err := os.Rename(info.TargetPath, backupPath); err != nil {
			http.Error(w, fmt.Sprintf("Failed to backup existing file: %v", err), http.StatusInternalServerError)
			return
		}
	} else if info.Exists && info.IsSymlink {
		// Just remove the old symlink
		if err := os.Remove(info.TargetPath); err != nil {
			http.Error(w, fmt.Sprintf("Failed to remove old symlink: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Ensure target parent dir exists
	targetParent := filepath.Dir(info.TargetPath)
	if err := os.MkdirAll(targetParent, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create target parent dir: %v", err), http.StatusInternalServerError)
		return
	}

	// Create symlink
	err := os.Symlink(info.RepoPath, info.TargetPath)
	if err != nil {
		// If native fails, try PowerShell
		var cmd *exec.Cmd
		if info.IsDir {
			cmd = exec.Command("powershell", "-Command", fmt.Sprintf("New-Item -ItemType Junction -Path %q -Target %q -Force", info.TargetPath, info.RepoPath))
		} else {
			cmd = exec.Command("powershell", "-Command", fmt.Sprintf("New-Item -ItemType SymbolicLink -Path %q -Target %q -Force", info.TargetPath, info.RepoPath))
		}
		if output, err := cmd.CombinedOutput(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create symlink: %v\nOutput: %s", err, string(output)), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Linked successfully")
}

func handleSync(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Run Theme Syncer script
	syncScript := "C:\\Users\\Timofey\\Theme\\sync_theme.ps1"
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", syncScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Sync failed: %v\nOutput: %s", err, string(output)), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func handleReloadGlazeWM(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cmd := exec.Command("glazewm", "command", "wm-reload-config")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Reload failed: %v\nOutput: %s", err, string(output)), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func main() {
	port := flag.Int("port", 54321, "Port to run the backend on")
	flag.Parse()

	// Locate project root
	var err error
	projectRoot, err = os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	http.HandleFunc("/api/status", handleStatus)
	http.HandleFunc("/api/file", handleReadFile)
	http.HandleFunc("/api/file/write", handleWriteFile)
	http.HandleFunc("/api/action/adopt", handleAdopt)
	http.HandleFunc("/api/action/link", handleLink)
	http.HandleFunc("/api/action/sync", handleSync)
	http.HandleFunc("/api/action/reload-glazewm", handleReloadGlazeWM)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	fmt.Printf("WinRice Backend starting on http://%s\n", addr)
	
	if runtime.GOOS == "windows" {
		os.MkdirAll(filepath.Join(projectRoot, "config"), 0755)
	}

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
