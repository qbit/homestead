#!/usr/bin/env python
 #  @filename   :   main.cpp
 #  @brief      :   7.5inch e-paper display demo
 #  @author     :   Yehui from Waveshare
 #
 #  Copyright (C) Waveshare     July 28 2017
 #
 # Permission is hereby granted, free of charge, to any person obtaining a copy
 # of this software and associated documnetation files (the "Software"), to deal
 # in the Software without restriction, including without limitation the rights
 # to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 # copies of the Software, and to permit persons to  whom the Software is
 # furished to do so, subject to the following conditions:
 #
 # The above copyright notice and this permission notice shall be included in
 # all copies or substantial portions of the Software.
 #
 # THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 # IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 # FITNESS OR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 # AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 # LIABILITY WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 # OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 # THE SOFTWARE.
 ##

import datetime
import epd7in5
import Image
import ImageDraw
import ImageFont
import json
import math
import os
import time
import urllib2
from subprocess import call

EPD_WIDTH = 640
EPD_HEIGHT = 384

wx_set = False

epd = epd7in5.EPD()

def init():
    epd.init()

w_key = "15d19578a1df62e09f70e6c9a0ccf439"
w_url = "http://api.openweathermap.org/data/2.5/%s?APPID=%s&zip=%s"

l_tmp = 0
l_wind_mpg = 0
l_wind_dir = ""
l_humidity = 0
l_rain = 0
l_name = ""
l_batt = ""

def get_weather_icon(code):
    return "no"

def local_weather():
    print("Updating local weather")
    op = os.popen("tail -n1 /home/pi/433.json").read()
    op = json.loads(op)
    global l_name
    global l_batt

    try:
        l_name = op['model']
        l_batt = op['battery']

        if op['message_type'] == 56:
            global l_tmp
            global l_wind_mpg
            global l_humidity
            l_tmp = op['temperature_F']
            l_wind_mpg = op['wind_speed_mph']
            l_humidity = op['humidity']

        if op['message_type'] == 49:
            global l_wind_dir
            global l_rain
            l_wind_dir = op['wind_dir']
            l_rain = op['rainfall_accumulation_inch']
    except:
        print "Found non-normal device"
        print op

def get_weather(loc=81069):
    wu = w_url % ("weather", w_key, loc)
    js = """
        {"name":"failed","main":{"temp":0},"wind":{"speed":0}, "weather":[{"main":"fail"}]}
        """
    weather = json.loads(js)

    try:
        js = urllib2.urlopen(wu).read()
    except:
        print("Can't query api.openweathermap.org", wu)

    try: 
        weather = json.loads(js)
    except:
        print("Can't read json from api.openweathermap.org")
    return weather

def ktof(k):
    return (k - 273.15) * 1.8000 + 32.00
def toc(f):
    return (f -32) * 5/9

count = 0
def main():
    global count
    try: 
        gdata = json.loads(urllib2.urlopen("https://data.bolddaemon.com/data/current/GreenHouse").read())
    except:
        print("Can't query data.bolddaemon.com!")
        gdata = {'temp': 0}
    try: 
        hdata = json.loads(urllib2.urlopen("https://data.bolddaemon.com/data/current/House").read())
    except:
        print("Can't query data.bolddaemon.com!")
        hdata = {'temp': 0}

        wx_set = True

    lwx = local_weather()
    if count % 12 == 0:
        print "Updating display"
        wx = get_weather()
        epd.init()
        gtemp = math.ceil(gdata['temp'])
        htemp = math.ceil(hdata['temp'])
        # For simplicity, the arguments are explicit numerical coordinates
        image = Image.new('1', (EPD_WIDTH, EPD_HEIGHT), 1)    # 1: clear the frame
        draw = ImageDraw.Draw(image)

        hfont = ImageFont.truetype('/usr/share/fonts/truetype/anonymous-pro/Anonymous Pro B.ttf', 24)
        font = ImageFont.truetype('/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf', 20)
        bfont = ImageFont.truetype('/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf', 60)
        cfont = ImageFont.truetype('/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf', 50)

        draw.rectangle((0, 6, 640, 40), fill = 0)

        os.environ["TZ"]="US/Mountain"
        time.tzset()
        n = datetime.datetime.now()

        draw.text((30, 10), '%s' % n.strftime("%x"), font = hfont, fill = 255)
        draw.text((200, 10), 'HomeStead Status', font = hfont, fill = 255)
        draw.text((550, 10), '%s' % n.strftime("%I:%M"), font = hfont, fill = 255)

        # Greenhouse
        draw.text((30, 50), 'GreenHouse:', font = hfont, fill = 0)
        draw.text((30, 70), '  Temp: %dF (%dC)' % (gtemp, toc(gtemp)), font = font, fill = 0)


        # House
        draw.text((300, 50), 'House:', font = hfont, fill = 0)
        draw.text((300, 70), '  Temp: %dF (%dC)' % (htemp, toc(htemp)), font = font, fill = 0)

        draw.rectangle([0,98, 640,100], fill = 0)

        print(wx)
        draw.text((30, 105), "%s:" % wx['name'], font = hfont, fill = 0)
        draw.text((30, 125), '  Temp: %dF (%dC)' % (ktof(wx['main']['temp']), wx['main']['temp'] - 273.15), font = font, fill = 0)
        draw.text((30, 145), '  Wind: %d mph' % math.ceil(wx['wind']['speed']), font = font, fill = 0)
        draw.text((30, 165), '  Sky: %s' % wx['weather'][0]['main'], font = font, fill = 0)

        draw.text((300, 105), "%s (%s):" % (l_name, l_batt), font = hfont, fill = 0)
        #draw.text((300, 125), '  Temp: %dF (%dC)' % (l_tmp, toc(l_tmp)), font = font, fill = 0)
        draw.text((300, 125), '  Wind: %d mph (%s)' % (math.ceil(l_wind_mpg), l_wind_dir), font = font, fill = 0)
        draw.text((300, 145), '  Humidity: %d%%' % l_humidity, font = font, fill = 0)
        draw.text((300, 165), '  Rain: %d in' % l_rain, font = font, fill = 0)

        draw.text((90, 220), '  %dF (%dC)' % (l_tmp, toc(l_tmp)), font = bfont, fill = 0)
        draw.text((20, 290), '  %dC' % toc(gtemp), font = cfont, fill = 0)
        draw.text((400, 290), '  %dC' % toc(htemp), font = cfont, fill = 0)

        epd.display_frame(epd.get_frame_buffer(image))
        epd.sleep()
    count = count + 1

if __name__ == '__main__':
    init()
    while 1:
        main()
        time.sleep(10) 
