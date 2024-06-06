package networking

import "fmt"

type SimpleIPManager struct {
	availableIPs []string
	usedIPs      map[string]bool
}

func NewSimpleIPManager() IPManager {
	return &SimpleIPManager{
		availableIPs: []string{"192.168.1.100", "192.168.1.101", "192.168.1.102"},
		usedIPs:      make(map[string]bool),
	}
}

func (manager *SimpleIPManager) AssignIP() (string, error) {
	if len(manager.availableIPs) == 0 {
		return "", fmt.Errorf("no available IPs")
	}
	ip := manager.availableIPs[0]
	manager.availableIPs = manager.availableIPs[1:]
	manager.usedIPs[ip] = true
	return ip, nil
}

func (manager *SimpleIPManager) ReleaseIP(ip string) error {
	if !manager.usedIPs[ip] {
		return fmt.Errorf("IP %s is not in use", ip)
	}
	manager.availableIPs = append(manager.availableIPs, ip)
	delete(manager.usedIPs, ip)
	return nil
}