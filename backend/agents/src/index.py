import os
import subprocess

from dotenv import load_dotenv

from stock.agent import start

load_dotenv()

start()


def run_go_binary(binary_path, *args):
    try:
        result = subprocess.run([binary_path, *args], capture_output=True, text=True)
        return result.stdout, result.stderr
    except FileNotFoundError:
        return None, f"Error: Binary not found at {binary_path}"


def call_go():

    go_path = "../al-go/main"

    stdout, stderr = run_go_binary(go_path, "arg1", "arg2")

    if stdout:
        print("Output:", stdout)
    if stderr:
        print("Error:", stderr)


# call_go()
