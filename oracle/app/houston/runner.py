import signal
import threading


def terminator():
    done = threading.Event()

    def on_stop(_, __):
        done.set()

    signal.signal(signal.SIGTERM, on_stop)
    signal.signal(signal.SIGINT, on_stop)

    return done
