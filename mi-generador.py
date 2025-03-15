import sys

def generate_compose(output_file, num_clients):
    """Genera un archivo docker-compose con un servidor y un número específico de clientes."""
    
    with open(output_file, 'w') as f:
        f.write("""name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
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
      - CLI_LOG_LEVEL=DEBUG
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
        if num_clients <= 0:
            print("El número de clientes debe ser positivo")
            sys.exit(1)
    except ValueError as e:
        print(f"El número de clientes debe ser un entero")
        sys.exit(1)
    
    generate_compose(output_file, num_clients)