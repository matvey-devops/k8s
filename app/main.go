package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"
)

type SystemInfo struct {
	CurrentDirectory string            `json:"current_directory"`
	Hostname         string            `json:"hostname"`
	Username         string            `json:"username"`
	DNSServers       []string          `json:"dns_servers"`
	NetworkAdapters  []NetworkAdapter  `json:"network_adapters"`
	OS               string            `json:"os"`
	Architecture     string            `json:"architecture"`
	GoVersion        string            `json:"go_version"`
	ServerTime       string            `json:"server_time"`
	EnvironmentVars  map[string]string `json:"environment_vars"`
}

type NetworkAdapter struct {
	Name        string   `json:"name"`
	IPAddresses []string `json:"ip_addresses"`
	MACAddress  string   `json:"mac_address"`
	Flags       string   `json:"flags"`
}

type Response struct {
	ReceivedData    interface{}         `json:"received_data"`
	Method          string              `json:"method"`
	Headers         map[string][]string `json:"headers"`
	ContentLength   int64               `json:"content_length"`
	RemoteAddr      string              `json:"remote_addr"`
	RequestURI      string              `json:"request_uri"`
	SystemInfo      SystemInfo          `json:"system_info"`
	ServerUptime    string              `json:"server_uptime"`
	Timestamp       time.Time           `json:"timestamp"`
}

var startTime = time.Now()

func main() {
	http.HandleFunc("/", handleAll)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Сервер запущен на порту %s\n", port)
	fmt.Printf("📍 Обращайтесь к / для получения полной информации\n")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("❌ Ошибка запуска сервера: %v\n", err)
		os.Exit(1)
	}
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Определяем полученные данные
	var receivedData interface{}
	if len(body) > 0 {
		// Пытаемся распарсить как JSON
		if err := json.Unmarshal(body, &receivedData); err != nil {
			// Если не JSON, сохраняем как строку
			receivedData = string(body)
		}
	} else {
		receivedData = "Нет данных в теле запроса"
	}

	// Получаем системную информацию
	systemInfo := getSystemInfo()

	// Формируем полный ответ
	response := Response{
		ReceivedData:  receivedData,
		Method:        r.Method,
		Headers:       r.Header,
		ContentLength: r.ContentLength,
		RemoteAddr:    r.RemoteAddr,
		RequestURI:    r.RequestURI,
		SystemInfo:    systemInfo,
		ServerUptime:  time.Since(startTime).String(),
		Timestamp:     time.Now().UTC(),
	}

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(response); err != nil {
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
	}
}

func getSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
		ServerTime:   time.Now().Format(time.RFC3339),
	}

	if cwd, err := os.Getwd(); err == nil {
		info.CurrentDirectory = cwd
	}

	if hostname, err := os.Hostname(); err == nil {
		info.Hostname = hostname
	}

	if currentUser, err := user.Current(); err == nil {
		info.Username = currentUser.Username
	}

	info.DNSServers = getDNSServers()
	info.NetworkAdapters = getNetworkAdapters()
	info.EnvironmentVars = getFilteredEnvVars()

	return info
}

func getDNSServers() []string {
	var servers []string

	content, err := os.ReadFile("/etc/resolv.conf")
	if err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "nameserver") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					servers = append(servers, parts[1])
				}
			}
		}
	}

	return servers
}

func getNetworkAdapters() []NetworkAdapter {
	var adapters []NetworkAdapter

	interfaces, err := net.Interfaces()
	if err != nil {
		return adapters
	}

	for _, iface := range interfaces {
		adapter := NetworkAdapter{
			Name:       iface.Name,
			MACAddress: iface.HardwareAddr.String(),
			Flags:      iface.Flags.String(),
		}

		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				adapter.IPAddresses = append(adapter.IPAddresses, addr.String())
			}
		}

		adapters = append(adapters, adapter)
	}

	return adapters
}

func getFilteredEnvVars() map[string]string {
	envVars := make(map[string]string)
	
	safeVars := []string{
		"PATH", "HOME", "USER", "SHELL", "LANG", "LC_ALL",
		"HOSTNAME", "PWD", "TERM", "TZ", "PORT",
		"KUBERNETES_SERVICE_HOST", "KUBERNETES_SERVICE_PORT",
		"NODE_NAME", "POD_NAME", "POD_NAMESPACE", "POD_IP",
	}

	for _, key := range safeVars {
		if value := os.Getenv(key); value != "" {
			envVars[key] = value
		}
	}

	return envVars
}
