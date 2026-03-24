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

def main():
    if not sys.argv[CLIENTS_NUM_ARGP].isdigit():
        print(f'Error: Argument {sys.argv[CLIENTS_NUM_ARGP]} is not a valid integer')
        return

    try: 
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
                     "      - batch=54\n"
                     "    networks:\n",
                     "      - testing_net\n",
                     "    depends_on:\n",
                     "      - server\n",
                     "    volumes:\n",
                     "      - type: bind\n",
                     "        source: ./client/config.yaml\n",
                     "        target: /config.yaml\n"
                ]
                f.writelines(client_config)

            f.writelines(DOCKER_NETWORKS_CONFIG)
    
    except PermissionError:
        print(f"Error: You don't have permission to write to '{sys.argv[FILE_NAME_ARGP]}'.")
    except OSError as e:
        print(f'error: {e}')
    # agregar de nombre invalido tmb

main()
