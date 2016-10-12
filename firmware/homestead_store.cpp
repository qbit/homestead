/**
 * Original Author: Todd Treece <todd@sparkfun.com>
 * Edited for the Spark by: Jim Lindblom <jim@sparkfun.com>
 * Edited for HomeStead by: Aaron bieber <aaron@bolddaemon.com>
 *
 * Copyright (c) 2014 SparkFun Electronics.
 * Licensed under the GPL v3 license.
 *
 */

#include "homestead_store.h"
#include <stdlib.h>

HomeStead::HomeStead(String host) {
	_host = host;
	_params = "";
}

void HomeStead::add(String field, String data) {

	_params += "&" + field + "=" + data;

}

void HomeStead::add(String field, char data) {

	_params += "&" + field + "=" + String(data);

}


void HomeStead::add(String field, int data) {

  _params += "&" + field + "=" + String(data);

}


void HomeStead::add(String field, byte data) {

  _params += "&" + field + "=" + String(data);

}


void HomeStead::add(String field, long data) {

  _params += "&" + field + "=" + String(data);

}

void HomeStead::add(String field, unsigned int data) {

  _params += "&" + field + "=" + String(data);

}

void HomeStead::add(String field, unsigned long data) {

  _params += "&" + field + "=" + String(data);

}

void HomeStead::add(String field, double data) {

  char tmp[30];

  //dtostrf(data, 1, 4, tmp);
  sprintf(tmp, "%f", data);

  _params += "&" + field + "=" + String(tmp);

}

void HomeStead::add(String field, float data) {

  char tmp[30];

  //dtostrf(data, 1, 4, tmp);
  sprintf(tmp, "%f", data);

  _params += "&" + field + "=" + String(tmp);

}

String HomeStead::queryString() {
  return String(_params);
}

String HomeStead::post() {

	String params = _params.substring(1);
	String result = "POST /data/store HTTP/1.1\n";
	result += "Host: " + _host + "\n";
	result += "Connection: close\n";
	result += "Content-Type: application/x-www-form-urlencoded\n";
	result += "Content-Length: " + String(params.length()) + "\n\n";
	result += params;

	_params = "";
	return result;

}

