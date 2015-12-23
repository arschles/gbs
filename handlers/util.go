package handlers

func buildStatusURL(containerID string) string {
	return "/status/" + containerID
}
