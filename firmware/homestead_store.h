/**
 * Original Author: Todd Treece <todd@sparkfun.com>
 * Edited for the Spark by: Jim Lindblom <jim@sparkfun.com>
 * Edited for HomeStead by: Aaron Bieber <aaron@bolddaemon.com>
 *
 * Copyright (c) 2014 SparkFun Electronics.
 * Licensed under the GPL v3 license.
 *
 */

#ifndef HomeStead_h
#define HomeStead_h

#include "application.h"

class HomeStead {

  public:
    HomeStead(String host);
    void add(String field, String data);
    void add(String field, char data);
    void add(String field, int data);
    void add(String field, byte data);
    void add(String field, long data);
    void add(String field, unsigned int data);
    void add(String field, unsigned long data);
    void add(String field, float data);
    void add(String field, double data);

    String queryString();
    String url();
    String post();

  private:
    String _host;
    String _params;
};

#endif
