### Free radius config

COnfig options:

1. Add latency to free radius responses

cd /etc/freeradius/3.0/
mkdir -p scripts
touch scripts/random_delay.sh

cat <<EOF > scripts/random_delay.sh
#!/bin/bash
sleep $((RANDOM % 2))
exit 0
EOF

chmod +x scripts/random_delay.sh

vim mods-available/exec

exec random_delay {
    wait = yes
    program = "/etc/freeradius/3.0/scripts/random_delay.sh"
}

ln -s /etc/freeradius/3.0/mods-available/exec /etc/freeradius/3.0/mods-enabled/exec 2>/dev/null || true

vim sites-enabled/default

under post-auth{}

random_delay

radiusd -X