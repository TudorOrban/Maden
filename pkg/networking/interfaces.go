package networking

type IPManager interface {
	AssignIP() (string, error)
	ReleaseIP(ip string) error
}