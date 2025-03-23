LENGTH_HEADER_SIZE = 4

class ConfirmationMessage:
    def __init__(self, status):
        """Creates a new ConfirmationMessage instance"""
        self.status = status

def send_message(socket, data):
    """Sends a message through the socket"""
    total_bytes_sent = 0
    while total_bytes_sent < len(data):
        bytes_sent = socket.send(data[total_bytes_sent:])
        if bytes_sent == 0:
            raise ConnectionError
        total_bytes_sent += bytes_sent
    return total_bytes_sent

def read_message(socket, length):
    """Reads a message from the socket"""
    data = bytearray()
    total_bytes_received = 0
    while total_bytes_received < length:
        bytesReceived = socket.recv(length - total_bytes_received)
        if not bytesReceived:
            raise ConnectionError
        data.extend(bytesReceived)
        total_bytes_received += len(bytesReceived)
    return data

def receive_message(socket):
    """Receives the payload of a message and its type from the socket"""
    length_bytes = read_message(socket, LENGTH_HEADER_SIZE)
    message_length = int.from_bytes(length_bytes, byteorder='big')
    full_message = read_message(socket, message_length)
    message_type = full_message[0]
    payload = full_message[1:]
    return message_type, payload
