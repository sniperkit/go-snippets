int vuln (char *userinput){
	char buffer[256];

	// force an overflow
	memcpy(buffer, input, 1024);

	return 1;
}