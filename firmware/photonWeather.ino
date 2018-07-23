#include "homestead_store.h"
#include "SparkFun_Photon_Weather_Shield_Library.h"

const int port = 80;
const char server[] = "data";

HomeStead homestead(server); // Create a Homestead object

//char sensorName[] = "outside";
char sensorName[] = "GreenHouse";
float humidity = 0;
float tempf = 0;
float pascals = 0;
float baroTemp = 0;

long lastPrint = 0;

int upload = 300000;
//int upload = 10000;

Weather sensor;

void setup_leds()
{
	#define RGB_NOTIFICATIONS_CONNECTING_ONLY 1
	RGB.control(true);
	RGB.brightness(0);
}

STARTUP(setup_leds());

//---------------------------------------------------------------
void setup()
{
	Serial.begin(9600);   // open serial over USB at 9600 baud

	//WiFi.selectAntenna(ANT_EXTERNAL);

	// Make sure your Serial Terminal app is closed before powering your device
	// Now open your Serial Terminal, and hit any key to continue!
	//Serial.println("Press any key to begin");
	//This line pauses the Serial port until a key is pressed
	//while(!Serial.available()) Spark.process();

	//Initialize the I2C sensors and ping them
	sensor.begin();

	/*You can only receive accurate barometric readings or accurate altitude
	  readings at a given time, not both at the same time. The following two lines
	  tell the sensor what mode to use. You could easily write a function that
	  takes a reading in one made and then switches to the other mode to grab that
	  reading, resulting in data that contains both accurate altitude and barometric
	  readings. For this example, we will only be using the barometer mode. Be sure
	  to only uncomment one line at a time. */
	sensor.setModeBarometer();//Set to Barometer Mode
	//baro.setModeAltimeter();//Set to altimeter Mode

	//These are additional MPL3115A2 functions that MUST be called for the sensor to work.
	sensor.setOversampleRate(7); // Set Oversample rate
	//Call with a rate from 0 to 7. See page 33 for table of ratios.
	//Sets the over sample rate. Datasheet calls for 128 but you can set it
	//from 1 to 128 samples. The higher the oversample rate the greater
	//the time between data samples.


	sensor.enableEventFlags(); //Necessary register calls to enble temp, baro and alt

        postToHomestead();
}
//---------------------------------------------------------------
void loop()
{
      //Get readings from all sensors
      getWeather();

      // This math looks at the current time vs the last time a publish happened
      if(millis() - lastPrint > upload)
      {
        // Record when you published
        lastPrint = millis();

        // Use the printInfo() function to print data out to Serial
        printInfo();
        postToHomestead();
      }
}
//---------------------------------------------------------------
void getWeather()
{
  // Measure Relative Humidity from the HTU21D or Si7021
  humidity = sensor.getRH();

  // Measure Temperature from the HTU21D or Si7021
  tempf = sensor.getTempF();
  // Temperature is measured every time RH is requested.
  // It is faster, therefore, to read it from previous RH
  // measurement with getTemp() instead with readTemp()

  //Measure the Barometer temperature in F from the MPL3115A2
  baroTemp = sensor.readBaroTempF();

  //Measure Pressure from the MPL3115A2
  pascals = sensor.readPressure();

  //If in altitude mode, you can get a reading in feet with this line:
  //float altf = sensor.readAltitudeFt();
}
//---------------------------------------------------------------

int postToHomestead()
{

	// baro_temp humidity pres_hpa pres_inhg temp
	homestead.add("sensor", sensorName);
	homestead.add("baro_temp", baroTemp);
	homestead.add("humidity", humidity);
	homestead.add("pres_hpa", pascals/100);
	homestead.add("pres_inhg", (pascals/100) * 0.0295300);
	homestead.add("temp", tempf);

	TCPClient client;
	char response[512];
	int i = 0;
	int retVal = 0;

	Serial.println("connecting");
	if (client.connect("10.0.1.5", port)) // Connect to the server
	{
		// Post message to indicate connect success
		Serial.println("Posting!");

		// homestead.post() will return a string formatted as an HTTP POST.
		// It'll include all of the field/data values we added before.
		// Use client.print() to send that string to the server.

		client.print(homestead.post());
		delay(1000);
		// Now we'll do some simple checking to see what (if any) response
		// the server gives us.
		while (client.available())
		{
			char c = client.read();
			Serial.print(c);	// Print the response for debugging help.
			if (i < 512)
				response[i++] = c; // Add character to response string
		}
		// Search the response string for "200 OK", if that's found the post
		// succeeded.
		if (strstr(response, "200 OK"))
		{
			Serial.println("Post success!");

			retVal = 1;
		}
		else if (strstr(response, "400 Bad Request"))
		{	// "400 Bad Request" means the Homestead POST was formatted incorrectly.
			// This most commonly ocurrs because a field is either missing,
			// duplicated, or misspelled.
			Serial.println("Bad request");

			retVal = -1;
		}
		else
		{
			// Otherwise we got a response we weren't looking for.
			retVal = -2;
		}
	}
	else
	{	// If the connection failed, print a message:
		Serial.println("connection failed");

		retVal = -3;
	}

	client.stop();	// Close the connection to server.
	return retVal;	// Return error (or success) code.
}

void printInfo()
{
	//This function prints the weather data out to the default Serial Port
	Serial.print("Temp:");
	Serial.print(tempf);
	Serial.print("F, ");

	Serial.print("Humidity:");
	Serial.print(humidity);
	Serial.print("%, ");

	Serial.print("Baro_Temp:");
	Serial.print(baroTemp);
	Serial.print("F, ");

	Serial.print("Pressure:");
	Serial.print(pascals/100);
	Serial.print("hPa, ");
	Serial.print((pascals/100) * 0.0295300);
	Serial.println("in.Hg");
	//The MPL3115A2 outputs the pressure in Pascals. However, most weather stations
	//report pressure in hectopascals or millibars. Divide by 100 to get a reading
	//more closely resembling what online weather reports may say in hPa or mb.
	//Another common unit for pressure is Inches of Mercury (in.Hg). To convert
	//from mb to in.Hg, use the following formula. P(inHg) = 0.0295300 * P(mb)
	//More info on conversion can be found here:
	//www.srh.noaa.gov/images/epz/wxcalc/pressureConversion.pdf

	//If in altitude mode, print with these lines
	//Serial.print("Altitude:");
	//Serial.print(altf);
	//Serial.println("ft.");

}

