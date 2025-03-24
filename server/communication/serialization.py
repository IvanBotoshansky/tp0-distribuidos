from common.utils import Bet

class Separator:
    BET_FIELDS = ","
    BETS = ";"
    WINNERS = ","

class MessageType:
    BET_BATCH = 1
    END_NOTIFICATION = 2
    WINNERS_REQUEST = 3
    CONFIRMATION = 11
    WINNERS_LIST = 12

def serialize_with_type(message_type, data):
    """Serializes data with a message type"""
    data_bytes = data.encode('utf-8')
    length = 1 + len(data_bytes)
    length_bytes = length.to_bytes(4, byteorder='big')
    type_bytes = message_type.to_bytes(1, byteorder='big')
    return length_bytes + type_bytes + data_bytes

def serialize_confirmation(confirmation_message):
    """Serializes the confirmation message"""
    return serialize_with_type(MessageType.CONFIRMATION, confirmation_message.status)

def serialize_winners_list(winners_list_message):
    """Serializes the winners list"""
    return serialize_with_type(MessageType.WINNERS_LIST, Separator.WINNERS.join(winners_list_message.winners))

def deserialize_bets(data):
    """Deserializes the data into a list of bets"""
    bets = []
    data = data.decode('utf-8')
    bets_data = data.split(Separator.BETS)
    for bet_data in bets_data:
        fields = bet_data.split(Separator.BET_FIELDS)
        bets.append(Bet(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5]))
    return bets

def deserialize_end_notification(data):
    """Deserializes the end notification"""
    return data.decode('utf-8')

def deserialize_winners_request(data):
    """Deserializes the winners request"""
    agency = data.decode('utf-8')
    return int(agency)
