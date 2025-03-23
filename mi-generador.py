import sys

def generate_compose(output_file, num_clients):
    """Genera un archivo docker-compose con un servidor y un número específico de clientes."""
    
    with open(output_file, 'w') as f:
        f.write(f"""name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - N_AGENCIES={num_clients}
    volumes:
      - ./server/config.ini:/config.ini
    networks:
      - testing_net
""")
        
        for i in range(1, num_clients + 1):
            f.write(f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={i}
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-{i}.csv:/.data/agency-{i}.csv
    networks:
      - testing_net
    depends_on:
      - server
""")
        
        f.write("""
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
""")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Cantidad de argumentos incorrecta")
        sys.exit(1)
    
    output_file = sys.argv[1]
    try:
        num_clients = int(sys.argv[2])
    except ValueError as e:
        print(f"El número de clientes debe ser un entero")
        sys.exit(1)
    
    generate_compose(output_file, num_clients)