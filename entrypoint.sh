echo "${TUNNEL_SSH_KEY}" > ./ssh_key
chmod 400 ./ssh_key

ssh -i ./ssh_key "${TUNNEL_USER}@${TUNNEL_IP}" \
    -N -D localhost:5000 \
    -o StrictHostKeyChecking=no

export HTTP_PROXY=socks5://localhost:5000

/app/tunnel-http-socks5
