import websocket
import random
from time import sleep, time
from threading import Lock, Thread

START_TIME = time() + 5

ARRAY_LOCK = Lock()
WS = []

def on_open(ws):
    # Sleep until the start time
    print("Sleeping for", START_TIME - time(), "seconds")
    sleep(START_TIME - time())
    # Wait for the mutex to be released
    for i in range(100):
        # send a random number between 0 and 65535
        data = bytearray()
        data.append(0b01000010)
        # two random uint16
        x, y = random.randint(0, 65535), random.randint(0, 65535)
        data.append(x >> 8)
        data.append(x & 0xFF)
        data.append(y >> 8)
        data.append(y & 0xFF)
        ws.send(data, opcode=websocket.ABNF.OPCODE_BINARY)
        sleep(0.005)
    ws.close()

def stress():
    ws = websocket.WebSocketApp("ws://localhost:3000/ws",
                                on_open=on_open)
    ARRAY_LOCK.acquire()
    WS.append(ws)
    ARRAY_LOCK.release()
    ws.run_forever()

if __name__ == "__main__":
    # Create 1000 threads
    threads = []
    for i in range(1000):
        threads.append(Thread(target=stress))
    # Start all threads
    for t in threads:
        t.start()
    # Wait for all threads to finish
    for t in threads:
        t.join()