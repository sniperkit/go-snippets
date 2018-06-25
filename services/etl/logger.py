import logging

logger = logging.getLogger('pedsnet.etlconv')
logger.setLevel(logging.INFO)

handler = logging.StreamHandler()
handler.setLevel(logging.DEBUG)
logger.addHandler(handler)
