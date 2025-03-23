import socket
import logging

from communication.message import read_message_with_type, ConfirmationMessage, send_message
from communication.serialization import (
    deserialize_bets,
    deserialize_end_notification,
    serialize_confirmation,
    MESSAGE_TYPE_BET_BATCH,
    MESSAGE_TYPE_END_NOTIFICATION
)
from common.utils import store_bets

SERVER_SOCKET_TIMEOUT = 1.0
MAX_NAME_LENGTH = 50
N_AGENCIES = 5

def are_valid_bets(bets):
    """Verify if the bets are valid"""
    for bet in bets:
        if bet.agency < 0 or bet.number < 0:
            return False
        if len(bet.first_name) > MAX_NAME_LENGTH or len(bet.last_name) > MAX_NAME_LENGTH:
            return False
    return True

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self._agencies_ended = set()

    def shutdown(self):
        """Gracefully shutdown the server"""
        logging.info("action: graceful_shutdown | result: in_progress")
        self._running = False
        self.__close_server_socket()
        logging.info("action: graceful_shutdown | result: success")

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self._running:
            client_sock = self.__accept_new_connection()
            if client_sock and self._running:
                self.__handle_client_connection(client_sock)

    def __close_server_socket(self):
        """Closes the server socket"""
        logging.info("action: close_connection | result: in_progress")
        if self._server_socket:
            self._server_socket.close()
            self._server_socket = None
        logging.info("action: close_connection | result: success")

    def __handle_bet_batch(self, client_sock, serialized_payload):
        """Handles a batch of bets"""
        bets = deserialize_bets(serialized_payload)
        if are_valid_bets(bets):
            store_bets(bets)
            logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
            confirmation_msg = ConfirmationMessage("success")
        else:
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets)}")
            confirmation_msg = ConfirmationMessage("fail")
        serialized_confirmation_msg = serialize_confirmation(confirmation_msg)
        send_message(client_sock, serialized_confirmation_msg)

    def __handle_end_notification(self, _client_sock, serialized_payload):
        """Handles an end notification"""
        agency = deserialize_end_notification(serialized_payload)
        self._agencies_ended.add(agency)

        if len(self._agencies_ended) == N_AGENCIES:
            logging.info("action: sorteo | result: success")

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            message_type, serialized_payload = read_message_with_type(client_sock)
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')

            if message_type == MESSAGE_TYPE_BET_BATCH:
                self.__handle_bet_batch(client_sock, serialized_payload)
            elif message_type == MESSAGE_TYPE_END_NOTIFICATION:
                self.__handle_end_notification(client_sock, serialized_payload)

        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: unexpected_error | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')

        try:
            self._server_socket.settimeout(SERVER_SOCKET_TIMEOUT)
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except socket.timeout:
            return None
        except OSError as e:
            if self._running:
                logging.error(f"action: accept_connections | result: fail | error: {e}")
            return None
        
