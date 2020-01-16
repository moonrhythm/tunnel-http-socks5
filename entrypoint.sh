echo "${TUNNEL_SSH_KEY}" | base64 -d > ./ssh_key
chmod 400 ./ssh_key

ssh -i ./ssh_key "${TUNNEL_USER}@${TUNNEL_IP}" \
    -N -D 127.0.0.1:5000 \
    -o StrictHostKeyChecking=no

export HTTP_PROXY=socks5://127.0.0.1:5000

/app/tunnel-http-socks5
