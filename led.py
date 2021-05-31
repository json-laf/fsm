import RPi.GPIO as GPIO
import time, sys

def Blinks():
    GPIO.setmode(GPIO.BCM)
    GPIO.setup(18, GPIO.OUT)
    GPIO.output(18, GPIO.LOW)
    blinks = 0
    print('开始闪烁')
    while (blinks < 5):
        GPIO.output(18, GPIO.HIGH)
        time.sleep(1.0)
        GPIO.output(18, GPIO.LOW)
        time.sleep(1.0)
        blinks += 1
    GPIO.output(18, GPIO.LOW)
    GPIO.cleanup()
    print('结束闪烁')

def Light():
    GPIO.setmode(GPIO.BCM)
    GPIO.setup(18, GPIO.OUT)
    GPIO.output(18, GPIO.LOW)
    blinks = 0
    print('开始常亮')
    while (blinks < 5):
        GPIO.output(18, GPIO.HIGH)
        time.sleep(1.0)
        blinks += 1
    GPIO.output(18, GPIO.LOW)
    GPIO.cleanup()
    print('结束常亮')

def main():
    if sys.argv[1] == "blinks":
        Blinks()
    if sys.argv[1] == "light":
        Light()

main()