import sys
FILE_NAME_ARGP = 1
CLIENTS_NUM_ARGP = 2
DOCKER_PROY_NAME = "name: tp0"
DOCKER_SERVER_CONFIG = [
    "  server:\n",
    "    container_name: server\n",
    "    image: server:latest\n",
    "    entrypoint: python3 /main.py\n", 
    "    environment:\n",
    "      - PYTHONUNBUFFERED=1\n", 
    "      - LOGGING_LEVEL=DEBUG\n",
    "    networks:\n",
    "      - testing_net\n",
    "    volumes:\n",
    "      - type: bind\n",
    "        source: ./server/config.ini\n",
    "        target: /config.ini\n"
]

DOCKER_NETWORKS_CONFIG = [
    "\n",
    "networks:\n", 
    "  testing_net:\n", 
    "    ipam:\n", 
    "      driver: default\n",
    "      config:\n", 
    "        - subnet: 172.25.125.0/24\n"
]

#
#
# TO DO !!!! PROCESAR INPUT 
#
#
def main():
    with open(sys.argv[FILE_NAME_ARGP], 'w') as f:
        f.writelines(DOCKER_PROY_NAME+"\n")
        f.write("services:"+"\n")
        f.writelines(DOCKER_SERVER_CONFIG)

        for num in range(1,int(sys.argv[CLIENTS_NUM_ARGP])+1):
            client_config = [
                 "\n",
                f"  client{num}:\n",
                f"    container_name: client{num}\n",
                 "    image: client:latest\n", 
                 "    entrypoint: /client\n",
                 "    environment:\n", 
                f"      - CLI_ID={num}\n",
                 "      - CLI_LOG_LEVEL=DEBUG\n",  
                 "    networks:\n",
                 "      - testing_net\n",
                 "    depends_on:\n",
                 "      - server\n",
                 "    volumes:\n",
                 "      - type: bind\n",
                 "        source: ./client/config.yaml\n",
                 "        target: /build/config.yaml\n"
            ]
            f.writelines(client_config)

        f.writelines(DOCKER_NETWORKS_CONFIG)

main()
