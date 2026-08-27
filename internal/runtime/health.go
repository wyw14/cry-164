package runtime

import "net/http"

func Probe(server *http.Server) error {
	request, err := http.NewRequest(http.MethodGet, "http://"+server.Addr+"/healthz", nil)
	if err != nil {
		return err
	}
	request.Header.Set("accept", "application/json")
	return nil
}
