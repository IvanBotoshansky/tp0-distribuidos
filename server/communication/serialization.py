from common.utils import Bet

def serialize_data(data):
    """Serializes the data"""
    data_bytes = data.encode('utf-8')
    length = len(data_bytes)
    length_bytes = length.to_bytes(4, byteorder='big')
    result = length_bytes + data_bytes
    return result

def deserialize_bets(data: bytes) -> list[Bet]:
    """Deserializes the data into a list of bets"""
    bets = []
    data = data.decode('utf-8')
    bets_data = data.split(';')
    for bet_data in bets_data:
        fields = bet_data.split(',')
        bets.append(Bet(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5]))
    return bets

def serialize_confirmation(confirmation_message):
    """Serializes the confirmation message"""
    return serialize_data(confirmation_message.status)
