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
	"regexp"
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
	c.RepoPath = filepath.Join(projectRoot, c.RepoPath)

	if stat, err := os.Lstat(c.RepoPath); err == nil {
		c.InRepo = true
		c.IsDir = stat.IsDir()
	}

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

	repoDir := filepath.Dir(info.RepoPath)
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create repo dir: %v", err), http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(info.RepoPath, []byte(payload.Content), 0644); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write file: %v", err), http.StatusInternalServerError)
		return
	}

	// Fallback copy: if symlinking is not active, also write to the target on C:
	if !info.IsSymlink && info.Exists {
		_ = os.WriteFile(info.TargetPath, []byte(payload.Content), 0644)
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

func adoptConfigSync(c ConfigInfo) error {
	info := getInfo(c)
	if !info.Exists {
		return fmt.Errorf("source does not exist")
	}
	if info.InRepo {
		return nil
	}

	repoDir := filepath.Dir(info.RepoPath)
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return err
	}

	if info.IsDir {
		return copyDir(info.TargetPath, info.RepoPath)
	}
	return copyFile(info.TargetPath, info.RepoPath)
}

func linkConfigSync(c ConfigInfo) error {
	info := getInfo(c)
	if !info.InRepo {
		return fmt.Errorf("must be in repo before linking")
	}

	// Try creating symlink/junction natively or via PowerShell
	var linkErr error
	if info.IsDir {
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("New-Item -ItemType Junction -Path %q -Target %q -Force", info.TargetPath, info.RepoPath))
		_, linkErr = cmd.CombinedOutput()
	} else {
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("New-Item -ItemType SymbolicLink -Path %q -Target %q -Force", info.TargetPath, info.RepoPath))
		_, linkErr = cmd.CombinedOutput()
	}

	// Fallback copy mode if symlinking requires elevation
	if linkErr != nil {
		fmt.Printf("Link failed for %s (likely elevation required). Falling back to direct sync-copy.\n", c.Name)
		if info.IsDir {
			os.RemoveAll(info.TargetPath)
			return copyDir(info.RepoPath, info.TargetPath)
		} else {
			return copyFile(info.RepoPath, info.TargetPath)
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

	if err := adoptConfigSync(*found); err != nil {
		http.Error(w, fmt.Sprintf("Adopt failed: %v", err), http.StatusInternalServerError)
		return
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

	if err := linkConfigSync(*found); err != nil {
		http.Error(w, fmt.Sprintf("Link failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Linked successfully")
}

func performAutoAdoptAndLink() {
	fmt.Println("Starting automatic adopt and link process...")
	for _, c := range configs {
		info := getInfo(c)
		// 1. Adopt if exists on C: but not in repo
		if info.Exists && !info.InRepo {
			fmt.Printf("Auto-Adopting configuration: %s...\n", c.Name)
			if err := adoptConfigSync(c); err != nil {
				fmt.Printf("Warning: failed to auto-adopt %s: %v\n", c.Name, err)
			}
		}

		// Refresh info after possible adopt
		info = getInfo(c)

		// 2. Link if in repo but not symlinked
		if info.InRepo && !info.IsSymlink {
			fmt.Printf("Auto-Linking configuration: %s...\n", c.Name)
			if err := linkConfigSync(c); err != nil {
				fmt.Printf("Warning: failed to auto-link %s: %v\n", c.Name, err)
			}
		}
	}
	fmt.Println("Auto-adopt and link process completed.")
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

type ThemeInfo struct {
	FontFamily      string  `json:"font_family"`
	FontSize        float64 `json:"font_size"`
	Opacity         float64 `json:"opacity"`
	BorderFocused   string  `json:"border_focused"`
	BorderUnfocused string  `json:"border_unfocused"`
	BgColor         string  `json:"bg_color"`
	FgColor         string  `json:"fg_color"`
	AccentColor     string  `json:"accent_color"`
	Lavender        string  `json:"lavender"`
	Lilac           string  `json:"lilac"`
	LavenderGrey    string  `json:"lavender_grey"`
	PineBlue        string  `json:"pine_blue"`
	JungleTeal      string  `json:"jungle_teal"`
}

func getStringVal(luaBlock, key string, defaultVal string) string {
	re := regexp.MustCompile(key + `\s*=\s*['"]([^'"]+)['"]`)
	m := re.FindStringSubmatch(luaBlock)
	if len(m) > 1 {
		return m[1]
	}
	return defaultVal
}

func getFloatVal(luaBlock, key string, defaultVal float64) float64 {
	re := regexp.MustCompile(key + `\s*=\s*([\d\.]+)`)
	m := re.FindStringSubmatch(luaBlock)
	if len(m) > 1 {
		var v float64
		fmt.Sscanf(m[1], "%f", &v)
		return v
	}
	return defaultVal
}

func handleGetTheme(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	weztermConfig := getInfo(configs[0])
	defaultTheme := ThemeInfo{
		FontFamily:      "JetBrains Mono",
		FontSize:        14.0,
		Opacity:         0.85,
		BorderFocused:   "#2b6f7c",
		BorderUnfocused: "#b2a9a5",
		BgColor:         "#f0edec",
		FgColor:         "#2c363c",
		AccentColor:     "#cbd9e3",
		Lavender:        "#dfd9e2",
		Lilac:           "#c3acce",
		LavenderGrey:    "#89909f",
		PineBlue:        "#538083",
		JungleTeal:      "#2a7f62",
	}

	if !weztermConfig.InRepo {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(defaultTheme)
		return
	}

	luaBytes, err := os.ReadFile(weztermConfig.RepoPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read WezTerm lua: %v", err), http.StatusInternalServerError)
		return
	}

	luaContent := string(luaBytes)
	re := regexp.MustCompile(`(?s)--\s*@theme\r?\n(.*?)\r?\n--\s*@theme-end`)
	m := re.FindStringSubmatch(luaContent)
	if len(m) < 2 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(defaultTheme)
		return
	}

	themeBlock := m[1]
	theme := ThemeInfo{
		FontFamily:      getStringVal(themeBlock, "font_family", defaultTheme.FontFamily),
		FontSize:        getFloatVal(themeBlock, "font_size", defaultTheme.FontSize),
		Opacity:         getFloatVal(themeBlock, "opacity", defaultTheme.Opacity),
		BorderFocused:   getStringVal(themeBlock, "border_focused", defaultTheme.BorderFocused),
		BorderUnfocused: getStringVal(themeBlock, "border_unfocused", defaultTheme.BorderUnfocused),
		BgColor:         getStringVal(themeBlock, "bg_color", defaultTheme.BgColor),
		FgColor:         getStringVal(themeBlock, "fg_color", defaultTheme.FgColor),
		AccentColor:     getStringVal(themeBlock, "accent_color", defaultTheme.AccentColor),
		Lavender:        getStringVal(themeBlock, "lavender", defaultTheme.Lavender),
		Lilac:           getStringVal(themeBlock, "lilac", defaultTheme.Lilac),
		LavenderGrey:    getStringVal(themeBlock, "lavender_grey", defaultTheme.LavenderGrey),
		PineBlue:        getStringVal(themeBlock, "pine_blue", defaultTheme.PineBlue),
		JungleTeal:      getStringVal(themeBlock, "jungle_teal", defaultTheme.JungleTeal),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(theme)
}

func handleSaveTheme(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var theme ThemeInfo
	if err := json.NewDecoder(r.Body).Decode(&theme); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	weztermConfig := getInfo(configs[0])
	if !weztermConfig.InRepo {
		http.Error(w, "WezTerm config must be adopted first", http.StatusBadRequest)
		return
	}

	luaBytes, err := os.ReadFile(weztermConfig.RepoPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read WezTerm config: %v", err), http.StatusInternalServerError)
		return
	}

	luaContent := string(luaBytes)

	newThemeBlock := fmt.Sprintf(`-- @theme
local theme = {
  color_scheme = 'zenbones',
  font_family = '%s',
  font_size = %.1f,
  opacity = %.2f,
  border_focused = '%s',
  border_unfocused = '%s',
  bg_color = '%s',
  fg_color = '%s',
  accent_color = '%s',
  lavender = '%s',
  lilac = '%s',
  lavender_grey = '%s',
  pine_blue = '%s',
  jungle_teal = '%s',
}
-- @theme-end`,
		theme.FontFamily,
		theme.FontSize,
		theme.Opacity,
		theme.BorderFocused,
		theme.BorderUnfocused,
		theme.BgColor,
		theme.FgColor,
		theme.AccentColor,
		theme.Lavender,
		theme.Lilac,
		theme.LavenderGrey,
		theme.PineBlue,
		theme.JungleTeal,
	)

	re := regexp.MustCompile(`(?s)--\s*@theme\r?\n.*?\r?\n--\s*@theme-end`)
	updatedLua := re.ReplaceAllString(luaContent, newThemeBlock)

	if err := os.WriteFile(weztermConfig.RepoPath, []byte(updatedLua), 0644); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write WezTerm lua: %v", err), http.StatusInternalServerError)
		return
	}

	// Fallback copy for WezTerm active path if it is not linked as a symlink
	if !weztermConfig.IsSymlink && weztermConfig.Exists {
		_ = os.WriteFile(weztermConfig.TargetPath, []byte(updatedLua), 0644)
	}

	syncScript := "C:\\Users\\Timofey\\Theme\\sync_theme.ps1"
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", syncScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to sync theme: %v\nOutput: %s", err, string(output)), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func main() {
	port := flag.Int("port", 54321, "Port to run the backend on")
	flag.Parse()

	var err error
	projectRoot, err = os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	if filepath.Base(projectRoot) == "src-go" {
		projectRoot = filepath.Dir(projectRoot)
	}

	if runtime.GOOS == "windows" {
		os.MkdirAll(filepath.Join(projectRoot, "config"), 0755)
	}

	go performAutoAdoptAndLink()

	http.HandleFunc("/api/status", handleStatus)
	http.HandleFunc("/api/file", handleReadFile)
	http.HandleFunc("/api/file/write", handleWriteFile)
	http.HandleFunc("/api/action/adopt", handleAdopt)
	http.HandleFunc("/api/action/link", handleLink)
	http.HandleFunc("/api/action/sync", handleSync)
	http.HandleFunc("/api/action/reload-glazewm", handleReloadGlazeWM)
	http.HandleFunc("/api/theme", handleGetTheme)
	http.HandleFunc("/api/theme/save", handleSaveTheme)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	fmt.Printf("WinRice Backend starting on http://%s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
