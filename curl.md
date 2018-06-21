# Use GET method and save the cookie into liwei
curl -v localhost:8080/set -c liwei

# Use GET method and read cookie from file liwei
 curl -v localhost:8080/read -b liwei

