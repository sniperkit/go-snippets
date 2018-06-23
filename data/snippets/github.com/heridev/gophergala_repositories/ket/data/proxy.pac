function FindProxyForURL(url, host) {
	if (host === "ket" || url.indexOf("http:") == 0)
		return "PROXY localhost:8080; DIRECT";
	return "DIRECT";
}