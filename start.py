import subprocess

def start(address: str):
    proc = subprocess.Popen(['', 'bar.py', None], shell=True)

start("cd C:/Users/1/Desktop/program/InforShowServer1/qq/3101522606 && go-cqhttp_windows_amd64.exe")